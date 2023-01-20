package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	awsRules "github.com/terraform-linters/tflint-ruleset-aws/rules"
)

func Test_EnsureDefaultTagsRule(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "Succeeds_WithProviderTags",
			Content: `
			provider "aws" {
				region = "eu-west-1"
				default_tags {
					tags = {
						team = "platform-engineering",
					}
				}
			}
			
			resource "aws_instance" "ec2_instance" {
				region = "eu-west-1"
			}`,
			Config: `
			rule "ensure_default_tags" {
			  enabled   = true
			  tags		= ["team"]
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Fails_WithMissingProviderTags",
			Content: `
			provider "aws" {
				region = "eu-west-1"
				default_tags {
					tags = {
						team = "platform-engineering",
					}
				}
			}
			
			resource "aws_instance" "ec2_instance" {
				region = "eu-west-1"
			}`,
			Config: `
			rule "ensure_default_tags" {
			  enabled   = true
			  tags		= ["team", "application"]
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewEnsureDefaultTagsRule(),
					Message: "The provider is missing the following tags: \"application\".",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 13},
						End:      hcl.Pos{Line: 7, Column: 7},
					},
				},
			},
		},
		{
			Name: "Succeeds_WithResourceTagsPresent",
			Content: `
			provider "aws" {
				region = "eu-west-1"
			}
			
			resource "aws_instance" "ec2_instance" {
				region = "eu-west-1"
				tags = {
					team = "platform-engineering"
				}
			}`,
			Config: `
			rule "ensure_default_tags" {
			  enabled   = true
			  tags		= ["team"]
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Fails_WithResourceTagsMissing",
			Content: `
			provider "aws" {
				region = "eu-west-1"
			}
			
			resource "aws_instance" "ec2_instance" {
				region = "eu-west-1"
			}
			  `,
			Config: `
			rule "ensure_default_tags" {
			  enabled   = true
			  tags		= ["team"]
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewEnsureDefaultTagsRule(),
					Message: "default_tags is missing",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
				{
					Rule:    awsRules.NewAwsResourceMissingTagsRule(),
					Message: "The resource is missing the following tags: \"team\".",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 4},
						End:      hcl.Pos{Line: 6, Column: 42},
					},
				},
			},
		},
	}

	rule := NewEnsureDefaultTagsRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content, ".tflint.hcl": test.Config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
