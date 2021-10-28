# Network Operator

[![Continuous Integration](https://github.com/Hellcatlk/network-operator/workflows/Continuous%20Integration/badge.svg)](https://github.com/Hellcatlk/network-operator/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Hellcatlk/network-operator)](https://goreportcard.com/report/github.com/Hellcatlk/network-operator)

Network operators is a kubernetes API for managing network devices.

## Quick start

1. Run `make install`.
2. Run `make run`(run the controller locally) or `make docker && make deploy`(run the controller as a deployment).
3. Apply CRs, [here](./config/samples) are some samples.

## Resources

* [API documentation](docs/api.md)
* [Testing](docs/testing.md)

### Important Notes

If you execute `make deploy` in a multi-node cluster environment, you need to upload the image to the image repository.

In most cases you needn't to apply `SwitchPort` manually, `Switch` controller will create it.

## Supported backend

|Device|Provider|Which backend it uses|
|:-|:-|:-|
|Switch|AnsibleSwitch|ansible|
