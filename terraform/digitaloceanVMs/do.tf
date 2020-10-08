variable "do_token" {}

variable "droplet_count" {
  default = 2
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_droplet" "drop" {
  count = var.droplet_count
  image  = "ubuntu-20-04-x64"
  name = format("testa-%02d", count.index) # "test-${count.index}"
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  ssh_keys = [20081905]
}

resource "digitalocean_loadbalancer" "loadb" {
  name = "loadb"

  region = "nyc3"

  forwarding_rule {
    entry_port = 80
    entry_protocol = "http"

    target_port = 80
    target_protocol = "http"
  }

  healthcheck {
    port = 22
    protocol = "tcp"
  }

  droplet_ids = digitalocean_droplet.drop.*.id
}

resource "digitalocean_project" "tutorials" {
  name        = "testprjc"
  resources   = flatten([digitalocean_droplet.drop.*.urn, digitalocean_loadbalancer.loadb.urn])
}

output "lb_output" {
  value = digitalocean_loadbalancer.loadb.ip
}
