package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// EnsureDefaultTagsRule checks whether ...
type EnsureDefaultTagsRule struct {
	tflint.DefaultRule
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
	return true
}

// Severity returns the rule severity
func (r *EnsureDefaultTagsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *EnsureDefaultTagsRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *EnsureDefaultTagsRule) Check(runner tflint.Runner) error {
	// This rule is an example to get a top-level resource attribute.
	resources, err := runner.GetProviderContent("aws", &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "default_tags",
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type: "tags",
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	// Put a log that can be output with `TFLINT_LOG=debug`
	logger.Debug(fmt.Sprintf("Get %d instances", len(resources.Blocks)))

	for _, resource := range resources.Blocks {
		defaultTagsBlock := resource.Body.Blocks.OfType("default_tags")

		if len(defaultTagsBlock) == 0 {
			return runner.EmitIssue(
				r,
				"default_tags is missing",
				resource.DefRange,
			)
		}
	}

	return nil
}
