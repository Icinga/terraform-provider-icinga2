---
layout: "icinga2"
page_title: "Icinga2: user"
sidebar_current: "docs-icinga2-resource-user"
description: |-
  Configures an user resource. This allows users to be configured, updated and deleted.
---

# icinga2\_user

Configures an Icinga2 user resource. This allows users to be configured, updated,
and deleted.

## Example Usage

```hcl
# Configure a new user

provider "icinga2" {
  api_url = "https://192.168.33.5:5665/v1"
}

resource "icinga2_user" "terraform" {
  name  = "terraform"
  email = "terraform@dev.null" 
}
```

## Argument Reference

The following arguments are supported:

* `name`  - (Required) The user.
* `email` - (Optional) An email string for this user. Useful for notification commands.
