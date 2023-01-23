package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/tags"
	"golang.org/x/exp/slices"
)

// ValidateTagsRule definition
type ValidateTagsRule struct {
	tflint.DefaultRule
}

// ValidateTagsRuleConfig is a config of ValidateTagsRule
type ValidateTagsRuleConfig struct {
	Tags []struct {
		Tag           string   `cty:"tag"`
		AllowedValues []string `cty:"allowed_values"`
	} `hclext:"tags"`
	Exclude []string `hclext:"exclude,optional"`
}

// NewValidateTagsRule returns a new rule
func NewValidateTagsRule() *ValidateTagsRule {
	return &ValidateTagsRule{}
}

// Name returns the rule name
func (r *ValidateTagsRule) Name() string {
	return "validate_tags"
}

// Enabled returns whether the rule is enabled by default
func (r *ValidateTagsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *ValidateTagsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *ValidateTagsRule) Link() string {
	return "https://github.com/0north/tflint-ruleset-0north-plugin/blob/main/docs/rules/validate_tags_rule.md"
}

// Checks the rule
func (r *ValidateTagsRule) Check(runner tflint.Runner) error {
	config := &ValidateTagsRuleConfig{}
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
	for _, provider := range providers.Blocks {
		// Get the default_tags block on the provider
		defaultTagsBlocks := provider.Body.Blocks.OfType("default_tags")

		// Check for allowed tags
		for _, defaultTagsBlock := range defaultTagsBlocks {
			err := r.verifyValidTag(runner, config, defaultTagsBlock)
			if err != nil {
				return err
			}
		}
	}

	// Go through all resources
	for _, resourceType := range tags.Resources {
		// Skip this resource if its type is excluded in the configuration
		if slices.Contains(config.Exclude, resourceType) {
			continue
		}

		resources, err := runner.GetResourceContent(resourceType, &hclext.BodySchema{
			Attributes: []hclext.AttributeSchema{{Name: "tags"}},
		}, nil)
		if err != nil {
			return err
		}

		// Go through all resources and check for allowed tag values
		for _, resource := range resources.Blocks {
			err := r.verifyValidTag(runner, config, resource)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Takes a Terraform block and verifies that if the team tag is present it has one of the valid values
func (r *ValidateTagsRule) verifyValidTag(runner tflint.Runner, config *ValidateTagsRuleConfig, block *hclext.Block) error {
	attribute, exists := block.Body.Attributes["tags"]
	if !exists {
		return nil
	}

	var tags map[string]string
	err := runner.EvaluateExpr(attribute.Expr, &tags, nil)
	if err != nil {
		return err
	}

	err = runner.EnsureNoError(err, func() error {
		for tag := range tags {
			for _, validatedTag := range config.Tags {
				if tag == validatedTag.Tag {
					if !slices.Contains(validatedTag.AllowedValues, tags[tag]) {
						err := runner.EmitIssue(
							r,
							fmt.Sprintf("Tag value %s is not allowed for tag %s (valid values are %s)", tags[tag], tag, strings.Join(validatedTag.AllowedValues, ", ")),
							attribute.Range,
						)
						if err != nil {
							return err
						}
					}
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
