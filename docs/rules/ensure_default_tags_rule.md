# ensure_default_tags_rule

Require specific tags for all AWS providers and AWS resource types that support them. This rule will find an issue if some providers are missing default_tags or if default_tags is not used and some resources are missing required tags.

## Configuration

```hcl
rule "ensure_default_tags_rule" {
  enabled = true
  tags = ["Foo", "Bar"]
  exclude = ["aws_autoscaling_group"] # (Optional) Exclude some resource types from tag checks
}
```

## Examples

The AWS provider and most AWS resources use the `tags` attribute with simple `key`=`value` pairs:

```hcl
provider "aws" {
  region = "eu-west-1"
  default_tags {
    tags = {
      foo = "Bar"
      bar = "Baz"
    }
  }
}

```

```
$ tflint
1 issue(s) found:

Notice: The provider is missing the following tags: "Bar", "Foo". (ensure_default_tags)

  on test.tf line 4:
   4:   tags = {
   5:     foo = "Bar"
   6:     bar = "Baz"
   7:   }
```

Iterators in `dynamic` blocks cannot be expanded, so the tags in the following example will not be detected.

```hcl
locals {
  tags = [
    {
      key   = "Name",
      value = "SomeName",
    },
    {
      key   = "env",
      value = "SomeEnv",
    },
  ]
}
resource "aws_autoscaling_group" "this" {
  dynamic "tag" {
    for_each = local.tags

    content {
      key                 = tag.key
      value               = tag.value
      propagate_at_launch = true
    }
  }
}
```

## Why

You want to set a standardized set of tags for your AWS resources.

## How To Fix

For each resource type that supports tags, ensure that each missing tag is present. Alternatively make sure your provider defines all the required tags as default_tags.
