package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_EnsureDefaultTagsRule_Fails_WithMissingAttribute(t *testing.T) {
	test := struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		Name: "issue found",
		Content: `
		provider "aws" {
			region = "eu-west-1"
		}`,
		Expected: helper.Issues{
			{
				Rule:    NewEnsureDefaultTagsRule(),
				Message: "default_tags is missing",
				Range: hcl.Range{
					Filename: "resource.tf",
					Start:    hcl.Pos{Line: 2, Column: 3},
					End:      hcl.Pos{Line: 2, Column: 17},
				},
			},
		},
	}

	rule := NewEnsureDefaultTagsRule()

	t.Run(test.Name, func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, test.Expected, runner.Issues)
	})
}

func Test_EnsureDefaultTagsRule_Succeeds_WithPresentAttribute(t *testing.T) {
	test := struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		Name: "issue found",
		Content: `
		provider "aws" {
			region = "eu-west-1"
			default_tags {
				tags = {
					team = "platform-engineering",
				}
			}
		}`,
		Expected: helper.Issues{},
	}

	rule := NewEnsureDefaultTagsRule()

	t.Run(test.Name, func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, test.Expected, runner.Issues)
	})
}
