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

Register it in terraform at `~/.terraformrc`:

```
providers {
    vultr = "$GOPATH/bin/terraform-provider-vultr"
}
```

Try the example:

    export VULTR_API_KEY=TODO_SET_TO_YOUR_API_KEY
    cd example
    export TF_LOG=DEBUG
    export TF_LOG_PATH=$PWD/terraform.log
    terraform init
    terraform plan
    terraform apply
    terraform show
    terraform output default_password
    ssh root@$(terraform output ipv4_address)
    terraform destroy


# Using the vultr CLI

List available regions:

    vultr regions

At time of writing these were the available regions:

```
DCID    NAME            CONTINENT       COUNTRY     STATE   STORAGE     CODE
40      Singapore       Asia            SG                  false       SGP
25      Tokyo           Asia            JP                  false       NRT
19      Sydney          Australia       AU                  false       SYD
7       Amsterdam       Europe          NL                  false       AMS
9       Frankfurt       Europe          DE                  false       FRA
8       London          Europe          GB                  false       LHR
24      Paris           Europe          FR                  false       CDG
6       Atlanta         North America   US          GA      false       ATL
2       Chicago         North America   US          IL      false       ORD
3       Dallas          North America   US          TX      false       DFW
5       Los Angeles     North America   US          CA      false       LAX
39      Miami           North America   US          FL      false       MIA
1       New Jersey      North America   US          NJ      true        EWR
4       Seattle         North America   US          WA      false       SEA
12      Silicon Valley  North America   US          CA      false       SJC
```

List the plans available on the `New Jersey` region:

    vultr plans --region 1

At the time of writing these were the available plans:

```
VPSPLANID   NAME                                    VCPU    RAM     DISK    BANDWIDTH   PRICE
200         512 MB RAM,20 GB SSD,0.50 TB BW         1       512     20      0.50        2.50
201         1024 MB RAM,25 GB SSD,1.00 TB BW        1       1024    25      1.00        5.00
202         2048 MB RAM,40 GB SSD,2.00 TB BW        1       2048    40      2.00        10.00
203         4096 MB RAM,60 GB SSD,3.00 TB BW        2       4096    60      3.00        20.00
204         8192 MB RAM,100 GB SSD,4.00 TB BW       4       8192    100     4.00        40.00
115         8192 MB RAM,110 GB SSD,10.00 TB BW      2       8192    110     10.00       60.00
116         16384 MB RAM,110 GB SSD,20.00 TB BW     4       16384   110     20.00       120.00
117         24576 MB RAM,110 GB SSD,30.00 TB BW     6       24576   110     30.00       180.00
118         32768 MB RAM,110 GB SSD,40.00 TB BW     8       32768   110     40.00       240.00
```

**NB** I've manually added the `200` plan into the list, and it seems to work.

List the available Operating Systems:

    vultr os

At the time of writing these were the available Operating Systems:

```
OSID    NAME                    ARCH    FAMILY      WINDOWS
186     Application             x64     application false
180     Backup                  x64     backup      false
147     CentOS 6 i386           i386    centos      false
127     CentOS 6 x64            x64     centos      false
167     CentOS 7 x64            x64     centos      false
179     CoreOS Stable           x64     coreos      false
159     Custom                  x64     iso         false
152     Debian 7 i386 (wheezy)  i386    debian      false
139     Debian 7 x64 (wheezy)   x64     debian      false
194     Debian 8 i386 (jessie)  i386    debian      false
193     Debian 8 x64 (jessie)   x64     debian      false
244     Debian 9 x64 (stretch)  x64     debian      false
245     Fedora 26 x64           x64     fedora      false
254     Fedora 27 x64           x64     fedora      false
140     FreeBSD 10 x64          x64     freebsd     false
230     FreeBSD 11 x64          x64     freebsd     false
234     OpenBSD 6 x64           x64     openbsd     false
164     Snapshot                x64     snapshot    false
161     Ubuntu 14.04 i386       i386    ubuntu      false
160     Ubuntu 14.04 x64        x64     ubuntu      false
216     Ubuntu 16.04 i386       i386    ubuntu      false
215     Ubuntu 16.04 x64        x64     ubuntu      false
253     Ubuntu 17.10 i386       i386    ubuntu      false
252     Ubuntu 17.10 x64        x64     ubuntu      false
124     Windows 2012 R2 x64     x64     windows     true
240     Windows 2016 x64        x64     windows     true
```

Finally, create a new Server with `Debian 9 x64 (stretch)`:

```bash
vultr server create \
    --name test \
    --hostname test.example.com \
    --region 1 \
    --plan 200 \
    --os 244
```

Which should return something like:

```
Virtual machine created

SUBID       NAME    DCID    VPSPLANID   OSID
13699733    test    1       200         244
```

See its status:

    vultr server show 13699733

Which should return something like:

**NB** While creating it might not have an IP, but eventually it will have an IP and change to the `active` state.

```
Id (SUBID):         13699733
Name:               test
Operating system:   Debian 9 x64 (stretch)
Status:             active
Power status:       running
Server state:       ok
Location:           New Jersey
Region (DCID):      1
VCPU count:         1
RAM:                512 MB
Disk:               Virtual 20 GB
Allowed bandwidth:  500
Current bandwidth:  0
Cost per month:     2.50
Pending charges:    0.01
Plan (VPSPLANID):   200
IP:                 104.156.250.252
Netmask:            255.255.254.0
Gateway:            104.156.250.1
Created date:       2018-02-25 07:32:31
Default password:   p6U}zc}+SK3dEDt)
Auto backups:       no
KVM URL:            https://my.vultr.com/subs/vps/novnc/api.php?data=KJVWUNJZLFTEG..
```

The server is only ready when its `status` is `active` and `power_status` is `running`:

```
Status:             active
Power status:       running
IP:                 104.156.250.252
```

**NB** As of 2018-02-25 there is no way to known if the server is shutdown. That is,
if you do a `poweroff` you don't have a way to known that... but if you stop the
server from the control panel the power status will change to `stopped`.

You can now try to access it:

    vultr ssh 13699733

You can also use a regular ssh client (use the password given in `Default password` attribute):

    ssh root@104.156.250.252

When you are finished playing with the server, You can now delete it:

    vultr server delete 13699733

Which should return:

```
Virtual machine deleted
```

From this point onwards, you can no longer access the server from the API.
