# go-icinga2-api

go-icinga2-api is a Go client library for configuring Icinga2 server via the [Icinga2 API](http://docs.icinga.org/icinga2/latest/doc/module/icinga2/chapter/icinga2-api)

[![Build Status](https://travis-ci.org/lrsmith/go-icinga2-api.svg?branch=master)](https://travis-ci.org/lrsmith/go-icinga2-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/lrsmith/go-icinga2-api)](https://goreportcard.com/report/github.com/lrsmith/go-icinga2-api)
[![codecov](https://codecov.io/gh/lrsmith/go-icinga2-api/branch/mas-aptet/graph/badge.svg)](https://codecov.io/gh/lrsmith/go-icinga2-api)


## Motivation

This library is being written to learn Go and also to provide the framework for a [Terraform](https://www.terraform.io/) Icinga2 Provider. An [initial implementation](https://github.com/lrsmith/terraform-provider-icinga2) was done but was not portable, so this project was started to provide a more general client library for Go which could be leveraged for refactoring the Terraform
providers.

## License

This software is licensed under the [Mozilla Public License 2.0](https://www.mozilla.org/en-US/MPL/2.0/)

## Contributing

This is a work in progress both for learning Go and getting some needed tooling. Any constructive feedback
or comments will be taken. Also contributions via Pull Requests will be accepted. Ideally any code contributions
should include or extend the existing tests.

# To Do
* Extend CreateHost to allow setting the remaining items in HostAttrs.
* Refactor DeleteHost so Cascade is a configurable option.

golang
