Icinga2 Terraform Provider
==================

[![Build Status](https://github.com/Icinga/terraform-provider-icinga2/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/Icinga/terraform-provider-icinga2/actions/workflows/unit-tests.yml/badge.svg)
[![Build Status](https://github.com/Icinga/terraform-provider-icinga2/actions/workflows/acceptance-tests.yml/badge.svg)](https://github.com/Icinga/terraform-provider-icinga2/actions/workflows/acceptance-tests.yml/badge.svg)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.26 or later
-	[Go](https://golang.org/doc/install) 1.23 or later (to build the provider plugin)

Building The Provider
---------------------

```sh
$ git clone git@github.com:Icinga/terraform-provider-icinga2
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/Icinga/terraform-provider-icinga2
$ make build
```

Using the provider
----------------------
The documentation for this provider is at the [Terraform Icinga2 provider docs](https://www.terraform.io/docs/providers/icinga2/)

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.23+ is *required*).

To compile the provider, run `make build`.

```sh
$ make build
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

For this provider you will need access to an Icinga2 server to run the acceptance tests.

```sh
$ make docker_start
$ make testacc
```
