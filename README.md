# tftools

tftools is a cli utility to work with Terraform state and plan json files.


## Plan

tftools offers commands to help analyze and summarize plan files. These can
be used to either provide better plan summaries for CI/CD or help complicated
TF state surgery.

### plan summary

`tfplan plan summary`

The plan summary processes a plan file and provides summary information about the file.
The specific information that is presented can be selected by use of the following flags.
If no flag is provided the --basic flag is assumed. Any combintation of flags are supported.

|option|type|description|
|------|----|-----------|
|plan-file|string|The path to the plan file to analyze|
|basic|flag|Provides the standard basic TF summary line.|
|resource-types|flag|Provides a summary of actions by resource type.|
|resource-address|flag|Provides a summary of actions by resource address.|

#### --basic
```shell
tfplan plan summary --plan-file <plan-file-path> --basic
Plan: 0 to import, 0 to add, 0 to change, 52 to destroy
```

#### --resource-types
```shell
Destroying    4 google_compute_instance_group
Adding       36 google_tags_location_tag_binding
Updating      6 google_compute_disk
Importing     6 google_compute_instance
```

#### --resource-address
```shell
- google_compute_instance.vm01
- google_compute_instance.vm02
~ google_compute_instance.vm03
> google_compute_instance.vm04
```

| short | action|
|-------|-------|
| -     | Destroy |
| +     | Add |
| ~     | Change |
| \>    | Import |



## State

tftools offers commands to help work with Terraform state. Specifically it offers
utilities to help generate commands or statements to make state changes.

### state import generate

`state import generate`

The generate command can generate Terraform import statements for all resources
scoped by a base address. This be useful when trying to import state that exists in
another state file. For example needing to split out some resources into a new state file.

| option              | type   | description                                                                                             |
|---------------------|--------|---------------------------------------------------------------------------------------------------------|
| state-file          | string | The path to the state file in json format.                                                              |
| resource-prefix     | string | An optional state prefix to scope the import statement creation.                                        |
| new-resource-prefix | string | An optional prefix that will replace the scoped prefix in the event the resources are moving locations. |

#### Selection Process

All resources with an address that starts with "resource-prefix" will be collected. New import statements will be generated for
each resource pointing the the ID found in state. The address in the import statement will be same as the address found but with
the resource-prefix portion replaced by the new-resource-prefix if new-resource-prefix is provided.

```shell
tftools state import generate --state-file <state-file-path> --resource-prefix module.old --new-resource-prefix module.new

#module.old.google_compute_instance.compute_instance[3]
import {
    to = module.new.google_compute_instance.compute_instance[3]
    id = "projects/project/zones/us-central1-f/instances/instance1"
}

```


