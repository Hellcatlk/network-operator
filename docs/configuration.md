# Configuration Provider Specification

## Overview

Configuration represents a configuration of a port, and provides available port types and specific configuration information.

## Date Types

1. Must define an API type for “configuration” resources. The type:
    1. Must be implemented as a CustomResourceDefinition
    2. Must be namespace-scoped
    3. Must have the standard Kubernetes “type metadata” and “object metadata”
