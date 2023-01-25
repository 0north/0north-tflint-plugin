package utils

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func NewIssue(rule tflint.Rule, message string, issueRange hcl.Range) helper.Issue {
	return helper.Issue{Rule: rule, Message: message, Range: issueRange}
}
