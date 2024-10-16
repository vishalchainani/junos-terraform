terraform {
  required_providers {
    junos-vsrx = {
      source = "juniper/providers/junos-vsrx"
      version = "21.31.108"
    }
  }
}

provider "junos-vsrx" {
    host = "10.52.180.39"
    port = 830
    username = "regress"
    password = "MaRtInI"
    sshkey = ""
}

module "vsrx_1" {
  source = "./vsrx_1"

  providers = {junos-vsrx = junos-vsrx}

  depends_on = [junos-vsrx_JunosDestroyCommit.commit-main]
}

resource "junos-vsrx_JunosDeviceCommit" "commit-main" {
  resource_name = "commit"
  depends_on = [module.vsrx_1]
}

resource "junos-vsrx_JunosDestroyCommit" "commit-main" {
  resource_name = "destroycommit"
}
