# Tailscale package for Synology NAS

Synology NAS package for Tailscale.

## Issue Tracker

File issues at: https://github.com/tailscale/tailscale/issues

## Installation

See the [Synology installation
guide](https://tailscale.com/kb/1131/synology/) on the Tailscale
website.

## Building from source

The source code for the Synology packages is kept in [Tailscale's main
code repository](https://github.com/tailscale/tailscale). You can
build the packages from source yourself with:

```bash
git clone https://github.com/tailscale/tailscale.git
cd tailscale
go run ./cmd/dist build synology
```

If everything worked you should have a directory called `dist` that
contains SPK files for all supported NASes and DSM versions.

## Precompiled packages
Tailscale also makes precompiled packages available for DSM6 and DSM7, supporting a variety of architectures.

 - [Stable](https://pkgs.tailscale.com/stable/#spks): stable releases. If you're not sure which track to use, pick this one.
 - [Unstable](https://pkgs.tailscale.com/unstable/#spks): the bleeding edge. Pushed early and often. Expect rough edges!

## Compatibility

The package is confirmed to be working on various Synology models. For
recent models, the correct package is usually the DSM7 package for
`x86_64` or `armv8`. For older models based on 32-bit ARM, check the
[synology model list](docs/platforms.md) to find the synology platform
name.

## Credits and References

- Thanks to [@nirev](https://github.com/nirev) for creating this
  project and transferring it to Tailscale's GitHub org.
- https://haugene.github.io/docker-transmission-openvpn/synology-nas/
  for the /dev/net/tun thing
- Package structure: [Synology Package Developer
  Guide](https://help.synology.com/developer-guide/index.html)
- The package building process was originally based on
  [synology-wireguard](https://github.com/runfalk/synology-wireguard). If
  you need to compile C code for a synology package, check it out.
