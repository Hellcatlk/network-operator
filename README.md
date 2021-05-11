# NetworkConfiguration Operator
[![Continuous Integration](https://github.com/Hellcatlk/networkconfiguration-operator/workflows/Continuous%20Integration/badge.svg)](https://github.com/Hellcatlk/networkconfiguration-operator/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Hellcatlk/networkconfiguration-operator)](https://goreportcard.com/report/github.com/Hellcatlk/networkconfiguration-operator)

## Quick start

1. Run `make install`
2. Run `make run`

## NOTE

Use like `kubebuilder create api --group metal3.io --version v1alpha1 --kind NetworkConfiguration` add API

## Supported backend

|os|protocol|url format|
|:-|:-|:-|
|openvswitch|ssh|ssh://\<host-ip>|
