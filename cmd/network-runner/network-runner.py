#!/usr/bin/python3

import json
import sys
from network_runner import api
from network_runner.models.inventory import Host, Inventory

# Parse json data
# format:
# {
#     host: "",
#     cert: {
#         username: "",
#         password: "",
#     },
#     os: "",
#     operator: "",
#     port: "",
#     untaggedVLAN: {
#         name: "",
#         id: 0,
#     },
#     vlans: [
#         {
#             Name: "",
#             ID: 0,
#         },
#     ],
# }
data = json.loads(sys.argv[1])

# Initial network runner
host = Host(name="network-operator",
            ansible_host=data["host"],
            ansible_user=data["cert"]["username"],
            ansible_ssh_pass=data["cert"]["password"],
            ansible_network_os=data["os"])
inventory = Inventory()
inventory.hosts.add(host)
network_runner = api.NetworkRunner(inventory)

if data["operator"] == "ConfigAccessPort":
    if data.get("untaggedVLAN") == None:
        print("miss required parameter(untaggedVLAN) for ConfigAccessPort")
        exit(1)

    # Create untagged vlan
    network_runner.create_vlan(
        "network-operator", data["untaggedVLAN"]["id"], data["untaggedVLAN"].get("name"))

    # Configure access port
    network_runner.conf_access_port(
        "network-operator", data["port"], data["untaggedVLAN"]["id"])
    exit(0)

if data["operator"] == "ConfigTrunkPort":
    # Create untagged vlan
    untaggedVLAN = None
    if data.get("untaggedVLAN") != None:
        network_runner.create_vlan(
            "network-operator", data["untaggedVLAN"]["id"], data["untaggedVLAN"].get("name"))
        untaggedVLAN = data["untaggedVLAN"]["id"]

    # Create tagged vlans
    vlans = []
    for vlan in data["vlans"]:
        # Create tagged vlan
        network_runner.create_vlan(
            "network-operator", vlan["id"], vlan.get("name"))
        vlans.append(vlan["id"])

    # Configure trunk port
    network_runner.conf_trunk_port(
        "network-operator", data["port"], untaggedVLAN, vlans)
    exit(0)

if data["operator"] == "DeletePort":
    network_runner.delete_port("network-operator", data["port"])
    exit(0)

print("invalid operator")
exit(1)
