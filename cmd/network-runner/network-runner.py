#!/usr/bin/python3

import json
import sys
from network_runner import api
from network_runner.models.inventory import Host, Inventory

# Parse json data
# format:
# {
#     Host: "",
#     Cert: {
#         Username: "",
#         Password: "",
#     },
#     OS: "",
#     Operator: "",
#     Port: "",
#     Vlan: 0,
#     Vlans: [],
# }
data = json.loads(sys.argv[1])

# Initial network runner
host = Host(name='network-operator',
            ansible_host=data.Host,
            ansible_user=data.Cert.Username,
            ansible_ssh_pass=data.Cert.Password,
            ansible_network_os=data.OS)
inventory = Inventory()
inventory.hosts.add(host)
network_runner = api.NetworkRunner(inventory)


# Check operator
if data.Operator == "CreateVlan":
    network_runner.create_vlan('network-operator', data.Vlan)
    # TODO: Check return value
    exit(0)

if data.Operator == "DeleteVlan":
    network_runner.delete_vlan('network-operator', data.Vlan)
    # TODO: Check return value
    exit(0)

if data.Operator == "ConfigAccessPort":
    network_runner.conf_access_port('network-operator', data.Port, data.Vlan)
    # TODO: Check return value
    exit(0)

if data.Operator == "ConfigTrunkPort":
    network_runner.conf_trunk_port(
        'network-operator', data.Port, data.Vlan, data.Vlans)
    # TODO: Check return value
    exit(0)

if data.Operator == "DeletePort":
    network_runner.delete_port('network-operator', data.Port)
    # TODO: Check return value
    exit(0)

print("invalid operator")
exit(1)
