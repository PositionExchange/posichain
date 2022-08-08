# Introduction
This document introduces the Posichain's package release using standard packaging system, RPM and Deb packages.

Standard packaging system has many benefits, like extensive tooling, documentation, portability, and complete design to handle different situation.

# Package Content
The RPM/Deb packages will install the following files/binary in your system.
* /usr/sbin/posichain
* /usr/sbin/posichain-setup.sh
* /usr/sbin/posichain-rclone.sh
* /etc/posichain/posichain.conf
* /etc/posichain/rclone.conf
* /etc/systemd/system/posichain.service
* /etc/sysctl.d/99-posichain.conf

The package will create `posichain` group and `posichain` user on your system.
The posichain process will be run as `posichain` user.
The default blockchain DBs are stored in `/home/posichain/posichain_db_?` directory.
The configuration of posichain process is in `/etc/posichain/posichain.conf`.

# Package Manager
Please take sometime to learn about the package managers used on Fedora/Debian based distributions.
There are many other package managers can be used to manage rpm/deb packages like [Apt]<https://en.wikipedia.org/wiki/APT_(software)>,
or [Yum]<https://www.redhat.com/sysadmin/how-manage-packages>

# Setup customized repo
You just need to do the setup of posichain repo once on a new host.
**TODO**: the repo in this document are for development/testing purpose only.

Official production repo will be different.

## RPM Package
RPM is for Redhat/Fedora based Linux distributions, such as Amazon Linux and CentOS.

```bash
# do the following once to add the posichain development repo
curl -LsSf http://danny-posichain-pub.s3.amazonaws.com/pub/yum/posichain-dev.repo | sudo tee -a /etc/yum.repos.d/posichain-dev.repo
sudo rpm --import https://raw.githubusercontent.com/harmony-one/harmony-open/master/harmony-release/harmony-pub.key
```

## Deb Package
Deb is supported on Debian based Linux distributions, such as Ubuntu, MX Linux.

```bash
# do the following once to add the posichain development repo
curl -LsSf https://raw.githubusercontent.com/harmony-one/harmony-open/master/harmony-release/harmony-pub.key | sudo apt-key add
echo "deb http://danny-harmony-pub.s3.amazonaws.com/pub/repo bionic main" | sudo tee -a /etc/apt/sources.list

```

# Test cases
## installation
```
# debian/ubuntu
sudo apt-get update
sudo apt-get install posichain

# fedora/amazon linux
sudo yum install posichain
```
## configure/start
```
# dpkg-reconfigure posichain (TODO)
sudo systemctl start posichain
```

## uninstall
```
# debian/ubuntu
sudo apt-get remove posichain

# fedora/amazon linux
sudo yum remove posichain
```

## upgrade
```bash
# debian/ubuntu
sudo apt-get update
sudo apt-get upgrade

# fedora/amazon linux
sudo yum update --refresh
```

## reinstall
```bash
remove and install
```

# Rclone
## install latest rclone
```bash
# debian/ubuntu
curl -LO https://downloads.rclone.org/v1.52.3/rclone-v1.52.3-linux-amd64.deb
sudo dpkg -i rclone-v1.52.3-linux-amd64.deb

# fedora/amazon linux
curl -LO https://downloads.rclone.org/v1.52.3/rclone-v1.52.3-linux-amd64.rpm
sudo rpm -ivh rclone-v1.52.3-linux-amd64.rpm
```

## do rclone
```bash
# validator runs on shard0
sudo -u posichain posichain-rclone.sh /home/posichain 0

# explorer node on shard0
sudo -u posichain posichain-rclone.sh -a /home/posichain 0
```

# Setup explorer (non-validating) node
To setup an explorer node (non-validating) node, please run the `posichain-setup.sh` at first.

```bash
sudo /usr/sbin/posichain-setup.sh -t explorer -s 0
```
to setup the node as an explorer node w/o blskey setup.

# Setup new validator
Please copy your blskey to `/home/posichain/.psc/blskeys` directory, and start the node.
The default configuration is for validators on mainnet. No need to run `posichain-setup.sh` script.

# Start/stop node
* `systemctl start posichain` to start node
* `systemctl stop posichain` to stop node
* `systemctl status posichain` to check status of node

# Change node configuration
The node configuration file is in `/etc/posichain/posichain.conf`.  Please edit the file as you needed.
```bash
sudo vim /etc/posichain/posichain.conf
```

# Support
Please open new github issues in https://github.com/PositionExchange/posichain/issues.
