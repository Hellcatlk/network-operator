# Network Operator

[![Continuous Integration](https://github.com/Hellcatlk/networkconfiguration-operator/workflows/Continuous%20Integration/badge.svg)](https://github.com/Hellcatlk/networkconfiguration-operator/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Hellcatlk/networkconfiguration-operator)](https://goreportcard.com/report/github.com/Hellcatlk/networkconfiguration-operator)

## Quick start

1. Run `make install`
2. Run `make run`(run the controller locally) or `make deploy`(run the controller as a deployment)
3. Apply CRs, there are some CRs' examples under [examples](./examples)

**NOTE**: If you execute `make deploy` in a multi-node cluster environment, you need to upload the image to the image repository.

## Supported backend

|provider switch|os|protocol|
|:-|:-|:-|
|OVSSwitch|openvswitch|ssh|

## Notes

1. Run `./tools/install_kubebuilder.sh` to install kubebuilder
2. Run `kubebuilder create api --group metal3.io --version v1alpha1 --kind <ResourceKind>` add API
