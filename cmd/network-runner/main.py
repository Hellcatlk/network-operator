#!/usr/bin/python3

import json
import sys
from network_runner import api
from network_runner.models.inventory import Host, Inventory


def _get_port_conf(port):
    """Get port configuration

    :param port: port ID
    :type data: String

    :returns: None
    """

    network_runner.get_port_conf("network-operator", port, True)
    return


def _config_access_port(port, untaggedVLAN):
    """Config untagged vlan to access port
    If untaggedVLAN isn't exist, we will create it.

    :param port: port ID
    :type data: String

    :param untaggedVLAN: untagged vlan ID and name
    :type untaggedVLAN: Int

    :returns: None
    """

    # Create untagged vlan
    network_runner.create_vlan(
        "network-operator", untaggedVLAN)
    # Configure access port
    network_runner.conf_access_port(
        "network-operator", port, untaggedVLAN)
    return


def _config_trunk_port(port, untaggedVLAN, vlans):
    """Config untagged vlan and vlans to access port.
    If untaggedVLAN or vlans aren't exist, we will create them.

    :param port: port ID
    :type data: String

    :param untaggedVLAN: untagged vlan ID and name
    :type untaggedVLAN: Int

    :param vlans: list of vlan
    :type vlans: List[]

    :returns: None
    """

    # Create untagged vlan
    if untaggedVLAN != None:
        network_runner.create_vlan(
            "network-operator", untaggedVLAN)
    # Create vlans
    for vlan in vlans:
        network_runner.create_vlan(
            "network-operator", vlan)
    # Configure trunk port
    network_runner.conf_trunk_port(
        "network-operator", port, untaggedVLAN, vlans)
    return


def _delete_port(port, bridge):
    """Clear vlan configure

    :param port: port ID
    :type data: String

    :param bridge: bridge name of openvswitch
    :type bridge: String

    :returns: None
    """

    network_runner.delete_port("network-operator", port, bridge_name=bridge)
    return


if __name__ == '__main__':
    # Parse json data
    # format:
    # {
    #     "host": "192.168.0.1",
    #     "credentials": {
    #         "username": "admin",
    #         "password": "admin"
    #     },
    #     "os": "fos",
    #     "bridge": "",
    #     "operator": "getPortConf/configAccessPort/configTrunkPort/deletePort",
    #     "port": "0/32",
    #     "untaggedVLAN": 0
    #     "vlans": [1,2,3]
    # }
    data = json.loads(sys.argv[1])

    # Initial network runner
    host = Host(name="network-operator",
                ansible_host=data["host"],
                ansible_user=data["credentials"]["username"],
                ansible_ssh_pass=data["credentials"]["password"],
                ansible_network_os=data["os"])
    inventory = Inventory()
    inventory.hosts.add(host)
    network_runner = api.NetworkRunner(inventory)

    # Deal operator
    if data["operator"] == "getPortConf":
        _get_port_conf(data["port"])
    elif data["operator"] == "configAccessPort":
        _config_access_port(data["port"], data["untaggedVLAN"])
    elif data["operator"] == "configTrunkPort":
        _config_trunk_port(data["port"], data.get(
            "untaggedVLAN"), data.get("vlans"))
    elif data["operator"] == "deletePort":
        _delete_port(data["port"], data.get("bridge"))
    else:
        print("invalid operator")
        exit(1)
