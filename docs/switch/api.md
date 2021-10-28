# API and Resource Definitions

## Switch

The **Switch** describes a switch.

### Switch spec

The `Switch's` spec defines the desire state of the switch.

#### Provider

A reference to a `SwitchProvider` contains the login information and back-end method of the switch.

#### Ports

Ports is a map whose key is the name of a SwitchPort CR and value is some info and limit of the port that SwitchPort CR represents.

* Port -- Indicates the specific restriction on the port.
  * physicalPortName -- The real port name in the switch.
  * disabled -- True if this port is not available, false otherwise.
  * vlanRange -- Indicates the range of VLANs allowed by this port in the switch.
  * trunkDisable -- True if this port can be used as a trunk port, false otherwise.

### Switch status

 The `Switch's` status which represents the switch's current state.

 #### State

 The current configuration status of the switch.

 #### Provider

 The reference of switch provider.

 #### Ports

 Restricted ports in the switch.

 #### Error

The error message of the port.

Example Switch:

``` yaml
apiVersion: metal3.io/v1alpha1
kind: Switch
metadata:
  creationTimestamp: "2021-10-26T02:36:49Z"
  finalizers:
  - metal3.io
  generation: 1
  name: switch-example
  namespace: default
  resourceVersion: "2835135"
  uid: b2d17d19-a74d-4c9d-9a30-be20312bd207
spec:
  ports:
    switchport-example:
      physicalPortName: br-test
      vlanRange: 1-100
  provider:
    kind: AnsibleSwitch
    name: ansible-example
    namespace: default
status:
  ports:
    switchport-example:
      physicalPortName: br-test
      vlanRange: 1-100
  provider:
    kind: AnsibleSwitch
    name: ansible-example
    namespace: default
  state: Running
```

## AnsibleSwitch

Use ansible as the backend to connect to the configuration switch.

#### os

The `os` is operator system of switch.

#### ip

The `ip` is ipv4 address of the switch.

#### bridge

Indicates the bridge where the port to be configured is located.Only for ovs switch.

#### port

The `port` is which port we can `ssh` to the switch.

#### credentialsSecret

The `credentialsSecret` is a secret resource contains username and password for the switch.

Example AnsibleSwitch:

```yaml
apiVersion: metal3.io/v1alpha1
kind: AnsibleSwitch
metadata:
  creationTimestamp: "2021-10-22T05:28:44Z"
  generation: 3
  name: ansible-example
  namespace: default
  resourceVersion: "2835123"
  uid: 04c971ab-2f46-4d8f-8597-c3a71f436c0e
spec:
  bridge: br-test
  credentials:
    name: switch-example-secret
    namespace: default
  host: 192.168.0.1
  os: openvswitch

```

## SwitchPort

**SwitchPort** CR represents a specific port of a network device, including port information,
the performance of the network device to which it belongs, and the performance of
the connected network interface card.

### SwitchPort Spec

The *SwitchPort Spec* defines the port on which network device and what configuration should be configured.

#### configuration

The reference of PortConfiguration CR.

### SwitchPort status

 The `SwitchPort's` status which represents the switchPort's current state.

#### physicalPortName

The `physicalPortName` field is the port's name on network device.

#### configuration

A reference to define which configuration should be configured.

#### state

The `state` shows the progress of configuring the network.

* *\<empty string\>* -- Indicates the status of the Port CR when it was first created.
* *Idle* -- Indicates waiting for spec.configurationRef to be assigned.
* *Configuring* -- Indicates that the port is being configured.
* *Active* -- Indicates that the port configuration is complete.
* *Deconfiguring* -- Indicates that the port configuration is being cleared.
* *Deletingted* -- Indicates that the port configuration has been cleared.

#### error

The error message of the port.

#### deviceRef

A reference to define this port on which network device.

Example SwitchPort:

```yaml
apiVersion: metal3.io/v1alpha1
kind: SwitchPort
metadata:
  creationTimestamp: "2021-10-26T02:36:49Z"
  finalizers:
  - metal3.io
  generation: 4
  name: switchport-example
  namespace: default
  ownerReferences:
  - apiVersion: metal3.io/v1alpha1
    blockOwnerDeletion: true
    kind: Switch
    name: switch-example
    uid: b2d17d19-a74d-4c9d-9a30-be20312bd207
  resourceVersion: "3791497"
  uid: d8e68088-8da6-436c-a683-2f6959ae84b6
spec:
  configuration:
    name: switchportconfiguration-example
    namespace: default
status:
  configuration:
    untaggedVLAN: 11
  portName: br-test
  state: Configuring
```

## SwitchPortConfiguration

The **SwitchPortConfiguration** is a kind configure for port on switch.

### SwitchPortConfiguration Spec

The *SwitchPortConfiguration Spec* defines details of configuration.

#### acl

The `acl` defines access control list of switch's port.

The sub-fields are

* *ipVersion* --
* *action* --
* *protocol* --
* *sourceIP* --
* *sourcePortRange* --
* *destinationIP* --
* *destinationPortRange* --

#### untaggedVLAN

Indicates which VLAN this port should be placed in.

#### taggedVLANRange

The range of tagged vlans.

#### disable

Disable port if true.

Example SwitchPort:

```yaml
apiVersion: metal3.io/v1alpha1
kind: SwitchPortConfiguration
metadata:
  creationTimestamp: "2021-10-22T05:28:44Z"
  generation: 1
  name: switchportconfiguration-example
  namespace: default
  resourceVersion: "1005794"
  uid: b37f856d-0faa-408d-ba0c-cde9b3718382
spec:
  untaggedVLAN: 11
```

## SwitchResourceLimit

`SwitchResourceLimit` represents information about the resources currently
available for the tenant in the switch.
It is created by the `SwitchResource` controller according to the
administrator's setting in the `SwitchResource`.

### SwitchResourceLimit status

The `SwitchResourceLimit's` status which represents the SwitchResourceLimit's current state.

#### vlanRange

Indicates the range of VLANs allowed by the user.

#### switchResourceRef

A reference to a switch.

#### usedVLAN

Indicates the vlan that the user has used.


Example SwitchResourceLimit:

```yaml
apiVersion: metal3.io/v1alpha1
kind: SwitchResourceLimit
metadata:
  creationTimestamp: "2021-10-26T02:36:49Z"
  generation: 1
  name: user-limit
  namespace: default
  resourceVersion: "3796423"
  uid: 2942c30a-ce5e-4233-8f3d-af90f1f29cf1
spec: {}
status:
  switchResourceRef:
    name: switchresource-example
    namespace: default
  usedVLAN: "11"
  vlanRange: 1-100
```

## SwitchResource

`SwitchResource` represents the resource in the switch.
The administrator writes the initial available resources
in the `spec` according to the actual situation, and then the controller
updates the real-time available resource to the `status`.
The administrator writes the user's restrictions into tenantLimits field.

### SwitchResource Spec

#### vlanRange

Indicates the initial allocatable vlan range.

#### tenantLimits

Indicates the resource limit for the tenant.

The sub-fields are
  * namespace -- The namespace where the restricted tenant is located.
  * vlanRange -- The range of VLANs allowed to be used.

### SwitchResource Status

#### availableVLAN

Indicates the vlan range that the administrator can assign to the user currently.

Example SwitchResource:

```yaml
apiVersion: metal3.io/v1alpha1
kind: SwitchResource
metadata:
  creationTimestamp: "2021-10-26T02:36:49Z"
  finalizers:
  - metal3.io
  generation: 1
  name: switchresource-example
  namespace: default
  resourceVersion: "2835132"
  uid: 1a6cd362-bbe4-4e49-a47f-63b6d430793b
spec:
  tenantLimits:
    user-1:
      namespace: default
      vlanRange: 1-100
  vlanRange: 1-1000
status:
  availableVLAN: 101-1000
  state: Running
  tenantLimits:
    user-1:
      namespace: default
      vlanRange: 1-100
```