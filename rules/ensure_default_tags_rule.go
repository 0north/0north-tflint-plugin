package rules

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/0north/tflint-ruleset-0north-plugin/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	awsRules "github.com/terraform-linters/tflint-ruleset-aws/rules"
)

var config *EnsureDefaultTagsRuleConfig

// EnsureDefaultTagsRule definition
type EnsureDefaultTagsRule struct {
	tflint.DefaultRule
}

// EnsureDefaultTagsRuleConfig is a config of EnsureDefaultTagsRule
type EnsureDefaultTagsRuleConfig struct {
	Tags    []string `hclext:"tags"`
	Exclude []string `hclext:"exclude,optional"`
}

// NewEnsureDefaultTagsRule returns a new rule
func NewEnsureDefaultTagsRule() *EnsureDefaultTagsRule {
	return &EnsureDefaultTagsRule{}
}

// Name returns the rule name
func (r *EnsureDefaultTagsRule) Name() string {
	return "ensure_default_tags"
}

// Enabled returns whether the rule is enabled by default
func (r *EnsureDefaultTagsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *EnsureDefaultTagsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *EnsureDefaultTagsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Checks the rule
func (r *EnsureDefaultTagsRule) Check(runner tflint.Runner) error {
	config = &EnsureDefaultTagsRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	// Check provider
	providers, err := runner.GetProviderContent("aws", &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "default_tags",
				Body: &hclext.BodySchema{

					Attributes: []hclext.AttributeSchema{
						{
							Name: "tags",
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	// Go through all providers
	providerDefaultTagIssues := []helper.Issue{}
	for _, provider := range providers.Blocks {
		// Check that default_tags is present on the provider
		defaultTagsBlocks := provider.Body.Blocks.OfType("default_tags")
		if len(defaultTagsBlocks) == 0 {
			providerDefaultTagIssues = append(providerDefaultTagIssues, utils.NewIssue(r, "default_tags is missing", provider.DefRange))
			continue
		}

		// Check for required tags
		for _, defaultTagsBlock := range defaultTagsBlocks {
			r.verifyRequiredTags(defaultTagsBlock, runner)
		}
	}

	// If we have default provider tag issues, we now need to check for AWS tags on all resources
	if len(providerDefaultTagIssues) != 0 {
		// Here we inject our override methods into the runner in order to capture issues from the AWS rule
		awsRunner := &AWSRunner{
			Runner: runner,
			Issues: []helper.Issue{},
		}

		// This runs the AWS Ressource Missing Tags Rule with our custom runner (so it uses our configuration and captures issues)
		// https://github.com/terraform-linters/tflint-ruleset-aws/blob/master/docs/rules/aws_resource_missing_tags.md
		err := awsRules.NewAwsResourceMissingTagsRule().Check(awsRunner)
		if err != nil {
			return err
		}

		// If there are issues with both missing default tags and missing resource tags, output all issues found
		if len(providerDefaultTagIssues) > 0 && len(awsRunner.Issues) > 0 {
			for _, issue := range providerDefaultTagIssues {
				err := runner.EmitIssue(r, issue.Message, issue.Range)
				if err != nil {
					return err
				}
			}

			for _, issue := range awsRunner.Issues {
				err := runner.EmitIssue(r, issue.Message, issue.Range)
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}

// Takes a Terraform block and verifies that the block has all the required tags
func (r *EnsureDefaultTagsRule) verifyRequiredTags(block *hclext.Block, runner tflint.Runner) {
	tagsAttribute := block.Body.Attributes["tags"]
	if tagsAttribute == nil {
		return
	}

	tagsBlock := tagsAttribute.Expr
	tagAttributes, _ := tagsBlock.Value(&hcl.EvalContext{})

	if !tagAttributes.Type().IsObjectType() {
		return
	}

	var missingTags []string = []string{}
	attributeMap := tagAttributes.AsValueMap()
	for _, requiredTag := range config.Tags {
		found := false
		for key := range attributeMap {
			if key == requiredTag {
				found = true
				break
			}
		}
		if !found {
			missingTags = append(missingTags, requiredTag)
		}
	}

	if len(missingTags) > 0 {
		err := runner.EmitIssue(
			r,
			fmt.Sprintf("The provider is missing the following tags: %s.", "\""+strings.Join(missingTags, "\", "+"\"")+"\""),
			tagsBlock.Range(),
		)
		if err != nil {
			return
		}
	}
}

// AWS runner overrides EmitIssue and DecodeRuleConfig to wrap AwsResourceMissingTagsRule. This allows us to capture the issues it finds as well as use our own configuration
type AWSRunner struct {
	tflint.Runner
	Issues []helper.Issue
}

func (r *AWSRunner) EmitIssue(rule tflint.Rule, message string, issueRange hcl.Range) error {
	r.Issues = append(r.Issues, utils.NewIssue(rule, message, issueRange))
	return nil
}

func (r *AWSRunner) DecodeRuleConfig(ruleName string, ret interface{}) error {
	v := reflect.ValueOf(ret).Elem()
	v.Set(reflect.ValueOf(struct {
		Tags    []string `hclext:"tags"`
		Exclude []string `hclext:"exclude,optional"`
	}{
		Tags:    config.Tags,
		Exclude: config.Exclude,
	}))

	return nil
}
