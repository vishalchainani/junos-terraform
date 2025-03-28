terraform {
  required_providers {
    junos-vsrx = {
      source = "juniper/providers/junos-vsrx"
      version = "21.31.108"
    }
  }
}

resource "junos-vsrx_Interfaces" "vsrx_2" {
  resource_name = "example_resource"
  interface = [
    {
    name = "ge-0/0/3"
    description = "Main_Ethernet"
    mtu = 9192
    vlan_tagging = true
    unit = [
      {
      name = 0
      description = "unit_description"
      vlan_id = 100
      family = [{
        inet = [{
          address = [
            {
            name = "192.168.103.1/24"
            },
            {
            name = "193.168.103.1/24"
            }

          ]
        }]
        inet6 = [{
          address = [
            {
            name = "2001:db8:85a3::8a2e:370:7334/64"
            },
            {
            name = "2001:db9:85a3::8a2e:370:7334/64"
            }
          ]
        }]
      }]
    },
    {
      name = 1
      description = "unit_description"
      vlan_id = 200
      family = [{
        inet = [{
          address = [
            {
            name = "192.169.103.1/24"
            },
            {
            name = "193.169.103.1/24"
            }
          ]
        }]
        inet6 = [{
          address = [
            {
            name = "2003:db8:85a3::8a2e:370:7334/64"
            },
            {
            name = "2003:db9:85a3::8a2e:370:7334/64"
            }
          ]
        }]
      }]
    }
    ]
  },
  {
    name = "ge-0/0/4"
    description = "Main_Ethernet"
    mtu = 9192
    vlan_tagging = true
    unit = [
      {
      name = 0
      description = "unit_description"
      vlan_id = 100
      family = [{
        inet = [{
          address = [
            {
            name = "192.168.104.1/24"
            },
            {
            name = "193.168.104.1/24"
            }
          ]
        }]
        inet6 = [{
          address = [
            {
            name = "2002:db8:85a3::8a2e:370:7334/64"
            },
            {
            name = "2002:db9:85a3::8a2e:370:7334/64"
            }
          ]
        }]
      }]
    },
    {
      name = 1
      description = "unit_description"
      vlan_id = 200
      family = [{
        inet = [{
          address = [
            {
            name = "192.169.104.1/24"
            },
            {
            name = "193.169.104.1/24"
            }
          ]
        }]
        inet6 = [{
          address = [
            {
            name = "2004:db8:85a3::8a2e:370:7334/64"
            },
            {
            name = "2004:db9:85a3::8a2e:370:7334/64"
            }
          ]
        }]
      }]
    }
    ]
  }
  ]
}
