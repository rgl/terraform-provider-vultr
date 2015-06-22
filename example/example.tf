#provider "vultr" {
#	api_key = "TODO_SET_TO_YOUR_API_KEY__OR__THE_VULTR_API_KEY_ENV_VARIABLE"
#}

resource "vultr_server" "example" {
	name = "created from terraform"

	# set the region. 7 is Amsterdam.
	# get the list of regions with the command: vultr regions
	region_id = 7

	# set the plain. 29 is 768 MB RAM,15 GB SSD,1.00 TB BW.
	# get the list of plans with the command: vultr plans --region 7
	plan_id = 29

	# set the OS image. 179 is CoreOS Stable.
	# get the list of OSs with the command: vultr os
	os_id = 179
	
	# enable IPv6.
	ipv6 = true

	# enable private networking.
	private_networking = true

	# enable one or more ssh keys on the root account.
	# get the list of keys with the command: vultr sshkeys
	ssh_key_ids = ["55586a7173780"]

	# execute a command on the local machine.
	provisioner "local-exec" {
        command = "echo local-exec ${vultr_server.example.ipv4_address}"
    }

	# execute commands on the remote machine.
	provisioner "remote-exec" {
        inline = [
			"env >> test.txt",
			"ip addr >> test.txt",
		]
    }
}
