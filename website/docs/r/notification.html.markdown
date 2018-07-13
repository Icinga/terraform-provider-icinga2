---
layout: "icinga2"
page_title: "Icinga2: notifiation"
sidebar_current: "docs-icinga2-resource-notifiation"
description: |-
  Configures a notifiation resource. This allows notifications to be configured, updated and deleted.
---

#TODO : Finish documentation. 
* Add resource code to create command, user and show example
* add to arguements section

# icinga2\_notification

Configures an Icinga2 notification resource. This allows notifications to be configured, updated,
and deleted.

## Example Usage

```hcl
# Configure a new host notification

provider "icinga2" {
  api_url = "https://192.168.33.5:5665/v1"
}

resource "icinga2_notification" "host-notification" {
  hostname = "docker-icinga2"
  command  = "mail-host-notification"
  users    = ["user"]
}

```

```hcl
# Configure a new service notification

provider "icinga2" {
  api_url = "https://192.168.33.5:5665/v1"
}

resource "icinga2_notification" "ping-service-notification" {
  hostname     = "docker-icinga2"
  command      = "mail-service-notification"
  users        = ["user"]
  servicename  = "ping"
}
```

## Argument Reference

The following arguments are supported:

* `hostname`    - (Required) The hostname the notification applies to.
* `command`     - (Required) Notification command to use.
* `users`       - (Required) List of users to notification.
* `servicename` - (Optional) Service to send notification for.
