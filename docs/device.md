# Device Provider Specification

## Overview

Device represents a specific network device instance and provides basic information
and login information of the device.

## Date Types

1. A device provider must define an API type for “device” resources. The type:
    1. Must be implemented as a CustomResourceDefinition
    2. Must be namespace-scoped
    3. Must have the standard Kubernetes “type metadata” and “object metadata”
    4. Must have a `spec` field with the following:
        1. Required fields:
            1. `ip`(string): IP Address of the device
            2. `secret`(struct): the username and password used to connect to the device
            3. `protocol`(string): equipment configuration protocol
            4. `port`(int): port to use for connection

## RBAC

TODO


