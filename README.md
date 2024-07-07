# SHA Exporter



## Overview

The SHA Exporter is intended to provide a means to ensure files and groups conform to predefined sha256 hashes.  Many organizations use tools, such as [Ansible](https://ansible.com), to deploy files to multiple hosts or to define the user membership of a group, such as the [wheel group](https://en.wikipedia.org/wiki/Wheel_(computing)), for granting privilege escalation.  This exporter provides a means to ensure these files and groups remain consistent with their expected content.

## Getting Started
* Clone the sha_exporter [repository](https://github.com/crooks/sha_exporter) and compile it using the standard `go build` command on your platform of choice.
* Copy the binary to somewhere sane on your system. A good choice is `/usr/local/bin/sha_exporter`.
* Create a YAML config file (usually something like `/etc/prometheus/sha_exporter.yml` on a Linux system).
* Populate the config file. There's an [example in the repo](https://github.com/crooks/sha_exporter/blob/main/examples/sha_exporter.yml) and details in the next section.
* Ensure your host(s) is reachable from [Prometheus](https://prometheus.io/) on the port you've defined in the config (default is TCP/9773).
* Start the binary and ensure it works.  On Linix you can do this with: `curl http://localhost:9773/metrics`
* You might want to create a dedicated user account for running Exporters: `useradd -m -r -U -d /opt/exporter -c "Prometheus Exporters" exporter`
* Create a systemd unit file for the exporter.  There's an [example in the repo](https://github.com/crooks/sha_exporter/blob/main/examples/sha_exporter.service).
* Start and enable the service: `systemctl enable --now sha_exporter.service`

## Configuration Options
The configuration file should be created in yaml format.  All the options are shown below.  Many of these have sane defaults if they're not specified in the config file.
```
groups:
  wheel:
    hash: 4cb48f0f2d2401f605cc9ba480fa238901224b017d5f89755d7fc9b0cc90f0f8

files:
  sshdcfg:
    path: /etc/ssh/sshd_config
    hash: 3055540351c2f0e6ce1c4bf1a63315dda30a646bbcc9e6c12ab9cddbaecf59de
  sha_exporter_unit:
    path: /etc/systemd/system/sha_exporter.service
    hash: 12b9dfefc0256e99716a1bb23f21bd3c8ce30970b880948de4d83a41bb475db3

groupfile: /etc/group
scrape_interval: 120

exporter:
  address: 0.0.0.0
  port: 9773

logging:
  journal: false
  level: debug
```

## Generating SHA-256 hashes
Each file or group defined in the config file requires an associated [SHA-256](https://en.wikipedia.org/wiki/SHA-2) hash.  There are many ways to generate such a hash depending on your operating system.  On Linux, the `sha256sum` is a convenient option.
* To generate a file hash: `sha256sum <filename>`
* To generate the hash of a string (like the users in a group): `echo -n "user1,user2,user3" | sha256sum`

**Note**: The ordering of users in a group is not important.  The code will sort them before generating a hash.  For example, `user1,user2,user3` will generate the same hash as `user2,user3,user1`.  When generating the hash for the configuration, it's important to first place them in alphanumeric order so as to generate the correct hash.  Using the previous example, you would use `echo -n "user1,user2,user3" | sha256sum`, **not** `echo -n "user2,user3,user1" | sha256sum`.

