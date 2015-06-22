# About

This is a [terraform](https://www.terraform.io/) provider for the [Vultr](https://www.vultr.com/) cloud.

See the [example](example/example.tf).

WARNING this is a work-in-progress. do not use in production.


# Build

Setup the Go workspace:

	mkdir -p terraform-provider-vultr/src/github.com/rgl/terraform-provider-vultr
	cd terraform-provider-vultr
	git clone https://github.com/rgl/terraform-provider-vultr src/github.com/rgl/terraform-provider-vultr
	export GOPATH=$PWD
	export PATH=$PWD/bin:$PATH
	hash -r # reset bash path

Get the dependencies:

	go get github.com/hashicorp/terraform
	go get github.com/JamesClonk/vultr

Build and test:

	cd src/github.com/rgl/terraform-provider-vultr
	go build
	go test

Copy it to the final directory:

	cp terraform-provider-vultr* $GOPATH/bin

Try the example:

	export VULTR_API_KEY=TODO_SET_TO_YOUR_API_KEY
	cd example
	terraform plan
	terraform apply
	terraform show
	terraform destroy


# Using the vultr CLI

List available regions:

	vultr regions

At time of writing these were the available regions:

```
DCID    NAME            CONTINENT       COUNTRY         STATE
4       Seattle         North America   US              WA
8       London          Europe          GB
24      Paris           Europe          FR
19      Sydney          Australia       AU
5       Los Angeles     North America   US              CA
1       New Jersey      North America   US              NJ
25      Tokyo           Asia            JP
9       Frankfurt       Europe          DE
6       Atlanta         North America   US              GA
2       Chicago         North America   US              IL
39      Miami           North America   US              FL
3       Dallas          North America   US              TX
12      Silicon Valley  North America   US              CA
7       Amsterdam       Europe          NL
```

List the plans available on the `Amsterdam` region:

	vultr plans --region 7

At the time of writing these were the available plans:

```
VPSPLANID       NAME                                    VCPU    RAM     DISK    BANDWIDTH       PRICE
29              768 MB RAM,15 GB SSD,1.00 TB BW         1       768     15      1.00            5.00
87              512 MB RAM,125 GB SATA,1.00 TB BW       1       512     125     1.00            5.00
98              32768 MB RAM,600 GB SSD,10.00 TB BW     16      32768   600     10.00           256.00
90              3072 MB RAM,750 GB SATA,4.00 TB BW      2       3072    750     4.00            30.00
91              4096 MB RAM,1000 GB SATA,5.00 TB BW     2       4096    1000    5.00            40.00
93              1024 MB RAM,20 GB SSD,2.00 TB BW        1       1024    20      2.00            8.00
94              2048 MB RAM,45 GB SSD,3.00 TB BW        2       2048    45      3.00            16.00
95              4096 MB RAM,90 GB SSD,4.00 TB BW        2       4096    90      4.00            32.00
96              8192 MB RAM,150 GB SSD,5.00 TB BW       4       8192    150     5.00            64.00
89              2048 MB RAM,500 GB SATA,3.00 TB BW      1       2048    500     3.00            20.00
100             65536 MB RAM,700 GB SSD,15.00 TB BW     24      65536   700     15.00           512.00
88              1024 MB RAM,250 GB SATA,2.00 TB BW      1       1024    250     2.00            10.00
97              16384 MB RAM,300 GB SSD,6.00 TB BW      8       16384   300     6.00            128.00
```

List the available Operating Systems:

	vultr os

At the time of writing these were the available Operating Systems:

```
OSID    NAME                    ARCH    FAMILY          WINDOWS
182     Ubuntu 14.10 i386       i386    ubuntu          false
191     Ubuntu 15.04 x64        x64     ubuntu          false
124     Windows 2012 R2 x64     x64     windows         true
159     Custom                  x64     iso             false
164     Snapshot                x64     snapshot        false
186     Application             x64     application     false
128     Ubuntu 12.04 x64        x64     ubuntu          false
160     Ubuntu 14.04 x64        x64     ubuntu          false
167     CentOS 7 x64            x64     centos          false
181     Ubuntu 14.10 x64        x64     ubuntu          false
192     Ubuntu 15.04 i386       i386    ubuntu          false
152     Debian 7 i386 (wheezy)  i386    debian          false
127     CentOS 6 x64            x64     centos          false
163     CentOS 5 i386           i386    centos          false
139     Debian 7 x64 (wheezy)   x64     debian          false
193     Debian 8 x64 (jessie)   x64     debian          false
140     FreeBSD 10.1 x64        x64     freebsd         false
179     CoreOS Stable           x64     coreos          false
180     Backup                  x64     backup          false
162     CentOS 5 x64            x64     centos          false
148     Ubuntu 12.04 i386       i386    ubuntu          false
194     Debian 8 i386 (jessie)  i386    debian          false
147     CentOS 6 i386           i386    centos          false
161     Ubuntu 14.04 i386       i386    ubuntu          false
```

Finally, create a new Server with `CoreOS Stable`:

	vultr server create --name test --region 7 --plan 29 --os 179

Which should return something like:

```
Virtual machine created

SUBID           NAME    DCID    VPSPLANID       OSID
2098003         test    7       29              179
```

See its status:

	vultr server show 2098003

Which should return something like:

```
vultr server show 2098003
Id (SUBID):             2098003
Name:                   test
Operating system:       CoreOS Stable
Status:                 pending
Power status:           running
Location:               Amsterdam
Region (DCID):          7
VCPU count:             1
RAM:                    768 MB
Disk:                   Virtual 15 GB
Allowed bandwidth:      1000
Current bandwidth:      0
Cost per month:         5.00
Pending charges:        0
Plan (VPSPLANID):       29
IP:                     0
Netmask:                0.0.0.0
Gateway:                0.0.0.0
Internal IP:
IPv6 IP:
IPv6 Network:
IPv6 Network Size:
Created date:           2015-06-21 14:16:26
Default password:       juqdotno
Auto backups:           no
KVM URL:
```

Eventually it will change to the `active` state:

```
Id (SUBID):             2098003
Name:                   test
Operating system:       CoreOS Stable
Status:                 active
Power status:           running
Location:               Amsterdam
Region (DCID):          7
VCPU count:             1
RAM:                    768 MB
Disk:                   Virtual 15 GB
Allowed bandwidth:      1000
Current bandwidth:      0
Cost per month:         5.00
Pending charges:        0.01
Plan (VPSPLANID):       29
IP:                     108.61.198.179
Netmask:                255.255.254.0
Gateway:                108.61.198.1
Internal IP:
IPv6 IP:
IPv6 Network:
IPv6 Network Size:
Created date:           2015-06-21 14:16:26
Default password:       jogdagni!4
Auto backups:           no
KVM URL:                https://my.vultr.com/subs/vps/novnc/api.php?data=OYZTSWCWINCVE..
```

The server is only ready when its `status` is `active` and `power_status` is `running`:

```
Status:                 active
Power status:           running
IP:                     108.61.198.179
```

NB As of 2015-06-21 there is no way to known if the server is shutdown. That is, if you do a `poweroff` you don't have a way to known that...

You can now try to access it:

	vultr ssh 2098003

Unfortunately that fails on my machine with:

```
panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xc0000005 code=0x1 addr=0x0 pc=0x5153ce]

goroutine 1 [running]:
golang.org/x/crypto/ssh.NewClientConn(0x2c21a0, 0xc082044050, 0xc08209d380, 0x11, 0x0, 0x0, 0x0, 0x0, 0xc0820224e0, 0x0, ...)
        c:/Users/rgl/.go/src/golang.org/x/crypto/ssh/client.go:68 +0x3fe
golang.org/x/crypto/ssh.Dial(0x7accb0, 0x3, 0xc08209d380, 0x11, 0x0, 0xc08209d380, 0x0, 0x0)
        c:/Users/rgl/.go/src/golang.org/x/crypto/ssh/client.go:176 +0x100
```

Just try again with a regular ssh client (use the password given in `Default password` attribute):

	ssh root@108.61.198.179

When you are finished playing with the server, You can now delete it:

	vultr server delete 2098003

Which should return:

```
Virtual machine deleted
```

From this point onwards, you can no longer access the server from the API.
