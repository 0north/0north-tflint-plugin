package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_ValidateTagsRule(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "Succeeds_ForProvider_WithValidTeamName_FromString",
			Content: `
			provider "aws" {
				region = "eu-west-1"
				default_tags {
					tags = {
						team = "platform-engineering",
					}
				}
			}`,
			Config: `
			rule "validate_tags" {
				enabled = true
				tags	= [
					{
						tag = "team",
						allowed_values = ["platform-engineering", "voyage-optimization"]
					}
				]
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Succeeds_ForProvider_WithValidTeamName_FromVariable",
			Content: `
			provider "aws" {
				region = "eu-west-1"
				default_tags {
					tags = {
						team = var.team,
					}
				}
			}

			variable "team" {
				description = "Name/email of the owning team [e.g. site-reliability, data-and-integration, core-team]"
				type        = string
				default     = "platform-engineering"
			}`,
			Config: `
			rule "validate_tags" {
				enabled = true
				tags	= [
					{
						tag = "team",
						allowed_values = ["platform-engineering", "voyage-optimization"]
					}
				]
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Fails_ForProvider_WithInvalidTeamName_FromString",
			Content: `
			provider "aws" {
				region = "eu-west-1"
				default_tags {
					tags = {
						team = "cloud-crew",
					}
				}
			}`,
			Config: `
			rule "validate_tags" {
				enabled = true
				tags	= [
					{
						tag = "team",
						allowed_values = ["platform-engineering", "voyage-optimization"]
					}
				]
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewValidateTagsRule(),
					Message: "Tag value cloud-crew is not allowed for tag team (valid values are platform-engineering, voyage-optimization)",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 6},
						End:      hcl.Pos{Line: 7, Column: 7},
					},
				},
			},
		},
		{
			Name: "Fails_ForProvider_WithInvalidTeamName_FromVariable",
			Content: `
			provider "aws" {
				region = "eu-west-1"
				default_tags {
					tags = {
						team = var.team,
					}
				}
			}

			variable "team" {
				description = "Name/email of the owning team [e.g. site-reliability, data-and-integration, core-team]"
				type        = string
				default     = "cloud-crew"
			}`,
			Config: `
			rule "validate_tags" {
				enabled = true
				tags	= [
					{
						tag = "team",
						allowed_values = ["platform-engineering", "voyage-optimization"]
					}
				]
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewValidateTagsRule(),
					Message: "Tag value cloud-crew is not allowed for tag team (valid values are platform-engineering, voyage-optimization)",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 6},
						End:      hcl.Pos{Line: 7, Column: 7},
					},
				},
			},
		},
		{
			Name: "Succeeds_ForResource_WithValidTeamName_FromString",
			Content: `
			resource "aws_instance" "ec2_instance" {
				region = "eu-west-1"
				tags = {
					team = "platform-engineering"
				}
			}`,
			Config: `
			rule "validate_tags" {
				enabled = true
				tags	= [
					{
						tag = "team",
						allowed_values = ["platform-engineering", "voyage-optimization"]
					}
				]
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Succeeds_ForResource_WithValidTeamName_FromVariable",
			Content: `
			resource "aws_instance" "ec2_instance" {
				region = "eu-west-1"
				tags = {
					team = var.team
				}
			}

			variable "team" {
				description = "Name/email of the owning team [e.g. site-reliability, data-and-integration, core-team]"
				type        = string
				default     = "platform-engineering"
			}`,
			Config: `
			rule "validate_tags" {
				enabled = true
				tags	= [
					{
						tag = "team",
						allowed_values = ["platform-engineering", "voyage-optimization"]
					}
				]
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Fails_ForResource_WithInvalidTeamName_FromString",
			Content: `
			resource "aws_instance" "ec2_instance" {
				region = "eu-west-1"
				tags = {
					team = "cloud-crew"
				}
			}`,
			Config: `
			rule "validate_tags" {
				enabled = true
				tags	= [
					{
						tag = "team",
						allowed_values = ["platform-engineering", "voyage-optimization"]
					}
				]
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewValidateTagsRule(),
					Message: "Tag value cloud-crew is not allowed for tag team (valid values are platform-engineering, voyage-optimization)",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 6, Column: 6},
					},
				},
			},
		},
		{
			Name: "Fails_ForResource_WithInvalidTeamName_FromVariable",
			Content: `
			resource "aws_instance" "ec2_instance" {
				region = "eu-west-1"
				tags = {
					team = var.team
				}
			}

			variable "team" {
				description = "Name/email of the owning team [e.g. site-reliability, data-and-integration, core-team]"
				type        = string
				default     = "cloud-crew"
			}`,
			Config: `
			rule "validate_tags" {
				enabled = true
				tags	= [
					{
						tag = "team",
						allowed_values = ["platform-engineering", "voyage-optimization"]
					}
				]
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewValidateTagsRule(),
					Message: "Tag value cloud-crew is not allowed for tag team (valid values are platform-engineering, voyage-optimization)",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 6, Column: 6},
					},
				},
			},
		},
	}

	rule := NewValidateTagsRule()

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
