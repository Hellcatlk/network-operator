# Port Provider Specification

## Overview

Port represents a specific port on the device, and provides information such
as port information and configuration.

## Date Types

1. A device provider must define an API type for “port” resources. The type:
    1. Must be implemented as a CustomResourceDefinition
    2. Must be namespace-scoped
    3. Must have the standard Kubernetes “type metadata” and “object metadata”
    4. Must have a `spec` field with the following:
        1. Required fields:
            1. `configurationRef`(ObjectReference): reference to the configuration CR that needs to be configured on this port. The configuration CR is considered immutable.
    5. Must have a `status` field with the following:
        1. Required fields:
            1. `state`(string): indicates the actual configuration status of the port
            2. `configurationRef`(ObjectReference): reference to the configuration information currently applied to the port
