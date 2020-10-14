variable "do_token" {}

variable "droplet_count" {
  default = 2
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_kubernetes_cluster" "yeye" {
  name = "testcluster"
  region = "nyc3"
  version = "1.18.8-do.1"
  # auto_upgrade = true

  node_pool {
    name       = "worker-pool"
    size       = "s-2vcpu-2gb"
    node_count = 3
    auto_scale = true
    min_nodes  = 3
    max_nodes  = 5
    tags = ["yeye-node"]
  }
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

  droplet_tag = "yeye-node"
}

resource "digitalocean_project" "tutorials" {
  name        = "testprjc"
  resources   = flatten(["do:kubernetes:${digitalocean_kubernetes_cluster.yeye.id}", digitalocean_loadbalancer.loadb.urn]) # have to construct urn manually
}

output "lb_output" {
  value = digitalocean_loadbalancer.loadb.ip
}

output "kube_info" {
  value = digitalocean_kubernetes_cluster.yeye.node_pool.0.nodes
}
