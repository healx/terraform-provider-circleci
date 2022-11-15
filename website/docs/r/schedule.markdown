---
layout: "circleci"
page_title: "CircleCI: circleci_schedule"
sidebar_current: "docs-resource-circleci-schedule"
description: |-
  Manages CircleCI schedules.
---

# circleci_schedule

A CircleCI schedule is a pipeline configuration that allows project workflows to be triggered periodically on a schedule.

## Example Usage

Basic usage:

```hcl
resource "circleci_schedule" "schedule" {
  organization          = "organization"
  project               = "repository"
  name                  = "schedule"
  description           = "Terraform generated schedule."
  per_hour              = 1
  hours_of_day          = [9,23]
  days_of_week          = ["MON", "TUES"]
  use_scheduling_system = false
  parameters            = jsonencode({ mycoolparam = false })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the schedule.
* `organization` - (Optional) The organization where the schedule will be created.
* `project` - (Required) The name of the CircleCI project to create the schedule in.
* `description` - (Optional) The description of the schedule.
* `per_hour` - (Required) How often per hour to trigger a pipeline.
* `hours_of_day` - (Required) Which hours of the day to trigger a pipeline.
* `days_of_week` - (Required) Which days of the week (\"MON\" .. \"SUN\") to trigger a pipeline on.
* `use_scheduling_system` - (Required) Use the scheduled system actor for attribution.
* `parameters_json` - (Optional) JSON encoded pipeline parameters to pass to created pipelines.

## Attributes Reference

* `id` - The ID of the schedule.

## Import

Schedules can be imported using the schedule ID.

```shell
# id
terraform import circleci_schedule.schedule 6d87b798-5edb-4d99-b424-ce73b43affb9
```
