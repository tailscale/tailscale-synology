# Tailscale package for Synology NAS
![CI](https://github.com/tailscale/tailscale-synology/workflows/CI/badge.svg)

Synology NAS package for Tailscale based on precompiled static binaries.

## Disclaimer

You use everything here at your own risk. Make sure you have other network
paths to your NAS before installing this, in case something goes wrong.

## Installation

Check the [releases](https://github.com/tailscale/tailscale-synology/releases)
page for SPKs for your platform. If there is no SPK you have to compile
it yourself using the instructions below.

1.  In the Synology DSM web admin UI, open the Package Center.
2.  Press the *Manual install* button and provide the SPK file.
3.  Follow the wizard until done.
4.  At this point `tailscaled` should be up and running.
5.  SSH into the  machine, and run `sudo tailscale up` so you can authenticate.

## Compatibility

The current package is confirmed to be working in different Synology models and architectures.

The package is created based on Tailscale [static binaries](https://pkgs.tailscale.com/stable/#static), and if your NAS has any of the supported architectures (x86, x86_64, arm, arm64) it should _just_ work.

If in doubt, check the [synology model list](docs/platforms.md) for the matching architecture.

## Making packages

This project builds Synology packages "by hand", based on pre-compiled tailscale static binaries.

You can build the packages using `make`
```bash
git clone https://github.com/tailscale/tailscale-synology.git
cd tailscale-synology/
make
```
If everything worked you should have a directory called `spks` that contains your SPK files.

## Credits and References

- Thanks to [@nirev](https://github.com/nirev) for creating this project and transferring it to Tailscale's GitHub org.
- https://haugene.github.io/docker-transmission-openvpn/synology-nas/ for the /dev/net/tun thing
- Package structure: [Synology Package Developer Guide](https://help.synology.com/developer-guide/index.html)
- Official Package building tool: [pkgscripts-ng](https://github.com/SynologyOpenSource/pkgscripts-ng)
- The package building process was originally based on [synology-wireguard](https://github.com/runfalk/synology-wireguard) \
If you need to _**compile**_ a synology package, check it out.
