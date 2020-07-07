## 0.5.0 (Unreleased)
## 0.4.0 (July 07, 2020)

ENHANCEMENTS:

* Moving to the latest version of the Icinga2 Go API
* Service templates are now supported in the Service resource

## 0.3.0 (March 19, 2020)

NOTES:

This is the first release in a while that hopefully should mark the start of a more regular release cycle.

ENHANCEMENTS:

* Move to the Terraform Plugin SDK
* Updated README to be specific to this provider


BUG FIXES:

* Fixed test for variable host.
* Fixed mistake in the example code.


## 0.2.0 (December 04, 2018)

FEATURES:

* **New Resource:** `icinga2_user` 
* **New Resource:** `icinga2_notification` 

ENHANCEMENTS:

 * provider: Allow specifying API document root.
 * resource/resource_icinga2_host: Add optional parameter to support declaring groups when creating hosts.
 * resource/resource_icinga2_service: Add optional parameter to support declaring service variables. ([[#2](https://github.com/terraform-providers/terraform-provider-icinga2/issues/2)](https://github.com/terraform-providers/terraform-provider-icinga2/issues/2))

BUGS:
 * resource/resource_icinga2_host: govend latest go-icinga2-api with optionaly declaring groups when creating a host. ([[#1](https://github.com/terraform-providers/terraform-provider-icinga2/issues/1)](https://github.com/terraform-providers/terraform-provider-icinga2/issues/1))

## 0.1.1 (August 04, 2017)

 * Configure "templates" via the icinga2_host resource ([#3](https://github.com/terraform-providers/terraform-provider-icinga2/issues/3))
 
## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
