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
#     UntaggedVLAN: {
#         Name: "",
#         ID: 0,
#     },
#     VLANs: [
#         {
#             Name: "",
#             ID: 0,
#         },
#     ],
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

if data.Operator == "ConfigAccessPort":
    # Create vlan
    network_runner.create_vlan(
        'network-operator', data.VLAN.ID, data.VLAN.Name)

    # Configure access port
    network_runner.conf_access_port(
        'network-operator', data.Port, data.VLAN.ID)
    # TODO: Check return value
    exit(0)

if data.Operator == "ConfigTrunkPort":
    vlans = []
    for VLAN in data.VLANs:
        # Create vlan
        network_runner.create_vlan(
            'network-operator', VLAN.ID, VLAN.Name)
        vlans.append(VLAN.ID)

    # Configure trunk port
    network_runner.conf_trunk_port(
        'network-operator', data.Port, data.VLAN.ID, vlans)
    # TODO: Check return value
    exit(0)

if data.Operator == "DeletePort":
    network_runner.delete_port('network-operator', data.Port)
    # TODO: Check return value
    exit(0)

print("invalid operator")
exit(1)
