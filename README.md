# ZeroNorth TFLint Ruleset Plugin

[![Build Status](https://github.com/terraform-linters/tflint-ruleset-template/workflows/build/badge.svg?branch=main)](https://github.com/terraform-linters/tflint-ruleset-template/actions)

This is the repository for ZeroNorth's custom TFLint ruleset.

## Requirements

- TFLint v0.40+
- Go v1.19

## Installation

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "0north-plugin" {
  enabled = true
  version = "0.4.0"
  source = "github.com/0north/tflint-ruleset-0north-plugin"
}
```

## Rules

| Name                | Description                                                     | Severity | Enabled | Link |
| ------------------- | --------------------------------------------------------------- | -------- | ------- | ---- |
| ensure_default_tags | Rule for linting AWS tags according to ZeroNorth specifications | WARNING  | ✔       |      |

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

You can run the built plugin like the following:

```
$ cat << EOS > .tflint.hcl
plugin "0north-plugin" {
  enabled = true
}
EOS
$ tflint
```
