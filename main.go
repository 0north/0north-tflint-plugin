package main

import (
	"github.com/0north/tflint-ruleset-0north-plugin/project"
	"github.com/0north/tflint-ruleset-0north-plugin/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "tflint-ruleset-0north-plugin",
			Version: project.Version,
			Rules: []tflint.Rule{
				rules.NewEnsureDefaultTagsRule(),
				rules.NewValidateTagsRule(),
			},
		},
	})
}
