# Terratest Example

## Overview:
The goal of this project create an example terraform goal with automated testing using Terratest.

## Setup
#Todo

## Tools Used:
[testify](https://github.com/stretchr/testify) handles testing suites in golang.
[terratest](https://github.com/gruntwork-io/terratest) is used as a library to execute Terraform commands from golang

## Known issues:

* I identified an bug with testify that means it is possible that if there is a panic() during SetupSuite then the resources created would not be destroyed.  
(**NOTE**: I created a [pull request](https://github.com/stretchr/testify/pull/850) against testify to fix this, which is scheduled for the next major release.  We could use [my fork](https://github.com/AaronNBrock/testify/tree/move-teardown-defers) for the moment if this is a major issue. )
* Running tests in parallel is not supported by testify yet (and doesn't seem to be a high priority).  See: https://github.com/stretchr/testify/issues/187