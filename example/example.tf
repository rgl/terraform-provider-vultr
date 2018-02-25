#provider "vultr" {
#  api_key = "TODO_SET_TO_YOUR_API_KEY__OR__THE_VULTR_API_KEY_ENV_VARIABLE"
#}

output "ipv4_address" {
  value = "${vultr_server.example.ipv4_address}"
}

output "default_password" {
  sensitive = true
  value = "${vultr_server.example.default_password}"
}

resource "vultr_ssh_key" "example" {
  name = "example created from terraform"

  # get the public key from a local file.
  #
  # create the example_rsa.pub file with:
  #
  #	ssh-keygen -t rsa -b 4096 -C 'terraform example' -f example_rsa -N ''
  public_key = "${file("example_rsa.pub")}"
}

resource "vultr_server" "example" {
  name = "example created from terraform"

  tag = "example tag"

  hostname = "test.example.com"

  # set the region. 1 is New Jersey.
  # get the list of regions with the command: vultr regions
  region_id = 1

  # set the plan. 200 is 512 MB RAM,20 GB SSD,0.50 TB BW.
  # get the list of plans with the command: vultr plans --region 1
  plan_id = 200

  # set the OS image. 244 is Debian 9 x64 (stretch).
  # get the list of OSs with the command: vultr os
  os_id = 244

  # enable IPv6.
  ipv6 = true

  # enable private networking.
  private_networking = true

  # enable one or more ssh keys on the root account.
  ssh_key_ids = ["${vultr_ssh_key.example.id}"]

  # execute a command on the local machine.
  provisioner "local-exec" {
    command = "echo local-exec ${vultr_server.example.ipv4_address}"
  }

  # execute commands on the remote machine.
  provisioner "remote-exec" {
    inline = [
      "env",
      "cat /etc/network/interfaces",
      "ip addr",
      "uname -a",
      "df -h",
    ]
  }
}
