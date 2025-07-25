terraform {
  required_providers {
    junos-vqfx = {
      source = "junos-vqfx"
    }
  }
}

provider "junos-vqfx" {
    host = "10.52.53.119"
    port = 22
    username = "regress"
    password = "MaRtInI"
    sshkey = ""
    alias = "dc2-spine1"
}

resource "junos-vqfx_Apply_Groups" "dc2-spine1" {
  resource_name = "JTAF_dc2-spine1"
  provider = junos-vqfx.dc2-spine1
  interfaces = [
    {
      interface = [
        {
          name = "xe-0/0/0"
          description = "*** to wan-pe2 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.32.2.1/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "xe-0/0/1"
          vlan_tagging = ""
          unit = [
            {
              name = 1
              vlan_id = 1
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.95.1.1/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "xe-0/0/2"
          vlan_tagging = ""
          unit = [
            {
              name = 1
              vlan_id = 1
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.94.1.1/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "xe-0/0/3"
          ether_options = [
            {
              ieee_802_3ad = [
                {
                  bundle = "ae0"
                }
              ]
            }
          ]
        },
        {
          name = "xe-0/0/4"
          ether_options = [
            {
              ieee_802_3ad = [
                {
                  bundle = "ae1"
                }
              ]
            }
          ]
        },
        {
          name = "xe-0/0/5"
          description = "*** to dc2-spine2 ***"
          unit = [
            {
              name = 0
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.189.1/30"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "ae0"
          esi = [
            {
              identifier = "00:00:00:00:00:00:00:01:01:00"
              all_active = ""
            }
          ]
          aggregated_ether_options = [
            {
              lacp = [
                {
                  active = ""
                  periodic = "fast"
                  system_id = "00:00:00:01:01:00"
                }
              ]
            }
          ]
          unit = [
            {
              name = 0
              family = [
                {
                  ethernet_switching = [
                    {
                      vlan = [
                        {
                          members = 1002
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "ae1"
          esi = [
            {
              identifier = "00:00:00:00:00:00:00:01:02:00"
              all_active = ""
            }
          ]
          aggregated_ether_options = [
            {
              lacp = [
                {
                  active = ""
                  periodic = "fast"
                  system_id = "00:00:00:01:02:00"
                }
              ]
            }
          ]
          unit = [
            {
              name = 0
              family = [
                {
                  ethernet_switching = [
                    {
                      vlan = [
                        {
                          members = 1002
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          name = "irb"
          unit = [
            {
              name = 1002
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.1.2.1/24"
                        }
                      ]
                    }
                  ]
                }
              ]
              mac = "02:0a:01:02:01:18"
            }
          ]
        },
        {
          name = "lo0"
          unit = [
            {
              name = 0
              description = "*** loopback ***"
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.30.100.8/32"
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              name = 10001
              description = "Loopback for VXLAN control packets for VRF_10001"
              family = [
                {
                  inet = [
                    {
                      address = [
                        {
                          name = "10.40.101.1/32"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
  snmp = [
    {
      location = "JCL Labs"
      contact = "aburston@juniper.net"
      community = [
        {
          name = "public"
          authorization = "read-only"
        }
      ]
    }
  ]
  forwarding_options = [
    {
      storm_control_profiles = [
        {
          name = "default"
          all = [
            {

            }
          ]
        }
      ]
    }
  ]
  routing_options = [
    {
      static = [
        {
          route = [
            {
              name = "0.0.0.0/0"
              next_hop = "100.123.0.1"
            }
          ]
        }
      ]
      router_id = "10.30.100.8"
      forwarding_table = [
        {
          export = "PFE-LB"
          ecmp_fast_reroute = ""
          chained_composite_next_hop = [
            {
              ingress = [
                {
                  evpn = ""
                }
              ]
            }
          ]
        }
      ]
    }
  ]
  protocols = [
    {
      bgp = [
        {
          group = [
            {
              name = "WAN_OVERLAY_eBGP"
              type = "external"
              multihop = [
                {
                  no_nexthop_change = ""
                }
              ]
              local_address = "10.30.100.8"
              family = [
                {
                  evpn = [
                    {
                      signaling = [
                        {
                          delay_route_advertisements = [
                            {
                              minimum_delay = [
                                {
                                  routing_uptime = 480
                                }
                              ]
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
              local_as = [
                {
                  as_number = 65201
                }
              ]
              multipath = [
                {
                  multiple_as = ""
                }
              ]
              neighbor = [
                {
                  name = "10.30.100.1"
                  description = "DCI EBGP peering to 10.30.100.1"
                  peer_as = 65200
                },
                {
                  name = "10.30.100.2"
                  description = "DCI EBGP peering to 10.30.100.2"
                  peer_as = 65200
                }
              ]
            },
            {
              name = "EVPN_iBGP"
              type = "internal"
              local_address = "10.30.100.8"
              family = [
                {
                  evpn = [
                    {
                      signaling = [
                        {

                        }
                      ]
                    }
                  ]
                }
              ]
              cluster = "10.30.100.8"
              local_as = [
                {
                  as_number = 65201
                }
              ]
              multipath = [
                {

                }
              ]
              neighbor = [
                {
                  name = "10.30.100.9"
                }
              ]
            },
            {
              name = "IPCLOS_eBGP"
              type = "external"
              mtu_discovery = ""
              import = "IPCLOS_BGP_IMP"
              export = "IPCLOS_BGP_EXP"
              vpn_apply_export = ""
              local_as = [
                {
                  as_number = 65520
                }
              ]
              multipath = [
                {
                  multiple_as = ""
                }
              ]
              bfd_liveness_detection = [
                {
                  minimum_interval = 1000
                  multiplier = 3
                }
              ]
              neighbor = [
                {
                  name = "10.30.189.2"
                  description = "EBGP peering to 10.30.189.2"
                  peer_as = 65521
                },
                {
                  name = "10.32.2.2"
                  description = "EBGP peering to 10.32.2.2"
                  peer_as = 65401
                }
              ]
            }
          ]
        }
      ]
      evpn = [
        {
          encapsulation = "vxlan"
          multicast_mode = "ingress-replication"
          default_gateway = "do-not-advertise"
          extended_vni_list = "all"
          no_core_isolation = ""
        }
      ]
      lldp = [
        {
          interface = [
            {
              name = "all"
            }
          ]
        }
      ]
      igmp_snooping = [
        {
          vlan = [
            {
              name = "default"
            }
          ]
        }
      ]
    }
  ]
  policy_options = [
    {
      policy_statement = [
        {
          name = "EVPN_T5_EXPORT"
          term = [
            {
              name = "fm_direct"
              from = [
                {
                  protocol = "direct"
                }
              ]
              then = [
                {
                  accept = ""
                }
              ]
            },
            {
              name = "fm_static"
              from = [
                {
                  protocol = "static"
                }
              ]
              then = [
                {
                  accept = ""
                }
              ]
            },
            {
              name = "fm_v4_default"
              from = [
                {
                  protocol = "evpn"
                },
                {
                  protocol = "ospf"
                },
                {
                  route_filter = [
                    {
                      address = "0.0.0.0/0"
                      exact = ""
                    }
                  ]
                }
              ]
              then = [
                {
                  accept = ""
                }
              ]
            },
            {
              name = "fm_v4_host"
              from = [
                {
                  protocol = "evpn"
                  route_filter = [
                    {
                      address = "0.0.0.0/0"
                      prefix_length_range = "/32-/32"
                    }
                  ]
                }
              ]
              then = [
                {
                  accept = ""
                }
              ]
            },
            {
              name = "fm_v6_host"
              from = [
                {
                  protocol = "evpn"
                  route_filter = [
                    {
                      address = "0::0/0"
                      prefix_length_range = "/128-/128"
                    }
                  ]
                }
              ]
              then = [
                {
                  accept = ""
                }
              ]
            }
          ]
        },
        {
          name = "IPCLOS_BGP_EXP"
          term = [
            {
              name = "loopback"
              from = [
                {
                  protocol = "direct"
                },
                {
                  protocol = "bgp"
                }
              ]
              then = [
                {
                  community = [
                    {
                      add = ""
                      community_name = "dc2-spine1"
                    }
                  ]
                  accept = ""
                }
              ]
            },
            {
              name = "default"
              then = [
                {
                  reject = ""
                }
              ]
            }
          ]
        },
        {
          name = "IPCLOS_BGP_IMP"
          term = [
            {
              name = "loopback"
              from = [
                {
                  protocol = "bgp"
                },
                {
                  protocol = "direct"
                }
              ]
              then = [
                {
                  accept = ""
                }
              ]
            },
            {
              name = "default"
              then = [
                {
                  reject = ""
                }
              ]
            }
          ]
        },
        {
          name = "PFE-LB"
          then = [
            {
              load_balance = [
                {
                  per_packet = ""
                }
              ]
            }
          ]
        },
        {
          name = "to-ospf"
          term = [
            {
              name = 10
              from = [
                {
                  protocol = "evpn"
                  route_filter = [
                    {
                      address = "10.1.2.0/24"
                      orlonger = ""
                    }
                  ]
                }
              ]
              then = [
                {
                  accept = ""
                }
              ]
            },
            {
              name = 100
              then = [
                {
                  reject = ""
                }
              ]
            }
          ]
        }
      ]
      community = [
        {
          name = "dc2-spine1"
          members = "65520:1"
        }
      ]
    }
  ]
  routing_instances = [
    {
      instance = [
        {
          name = "VRF_10001"
          instance_type = "vrf"
          interface = [
            {
              name = "xe-0/0/1.1"
            },
            {
              name = "xe-0/0/2.1"
            },
            {
              name = "irb.1002"
            },
            {
              name = "lo0.10001"
            }
          ]
          route_distinguisher = [
            {
              rd_type = "10.40.101.1:10001"
            }
          ]
          vrf_target = [
            {
              community = "target:1:10001"
            }
          ]
          vrf_table_label = [
            {

            }
          ]
          routing_options = [
            {
              auto_export = [
                {

                }
              ]
            }
          ]
          protocols = [
            {
              ospf = [
                {
                  export = "to-ospf"
                  area = [
                    {
                      name = "0.0.0.0"
                      interface = [
                        {
                          name = "xe-0/0/1.1"
                          metric = 100
                        },
                        {
                          name = "xe-0/0/2.1"
                          metric = 200
                        }
                      ]
                    }
                  ]
                }
              ]
              evpn = [
                {
                  ip_prefix_routes = [
                    {
                      advertise = "direct-nexthop"
                      encapsulation = "vxlan"
                      vni = 10001
                      export = "EVPN_T5_EXPORT"
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
  switch_options = [
    {
      vtep_source_interface = [
        {
          interface_name = "lo0.0"
        }
      ]
      route_distinguisher = [
        {
          rd_type = "10.30.100.8:9999"
        }
      ]
      vrf_target = [
        {
          community = "target:9999:9999"
          auto = [
            {

            }
          ]
        }
      ]
    }
  ]
  vlans = [
    {
      vlan = [
        {
          name = "vlan_1002"
          vlan_id = 1002
          l3_interface = "irb.1002"
          vxlan = [
            {
              vni = 1002
            }
          ]
        }
      ]
    }
  ]
}