# validate_tags_rule

Validate tag values for all AWS providers and AWS resource types that support them.

## Configuration

```hcl
rule "validate_tags" {
  enabled = true
  tags = [
    {
        tag = "foo"
        allowed_values = ["bar, baz"]
    }
  ]
  exclude = ["aws_autoscaling_group"] # (Optional) Exclude some resource types from tag checks
}
```

## Examples

This rule ensures that a tag can only be set to one of the allowed values:

```hcl
provider "aws" {
  region = "eu-west-1"
  default_tags {
    tags = {
      foo = "x"
    }
  }
}

```

```
$ tflint
1 issue(s) found:

Notice: Tag value x is not allowed for tag foo (valid values are bar, baz) (validate_tags_rule)

  on test.tf line 4:
   4:   tags = {
   5:     foo = "x"
   6:   }
```

## Why

You want to standardize tag values for your AWS resources.

## How To Fix

For each resource or provider with invalid tags, ensure that each tag has a valid value.
