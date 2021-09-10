# Tailscale package for Synology NAS
![CI](https://github.com/tailscale/tailscale-synology/workflows/CI/badge.svg)

Synology NAS package for Tailscale based on precompiled static binaries.

## Disclaimer

You use everything here at your own risk. Make sure you have other network
paths to your NAS before installing this, in case something goes wrong.

## Issue Tracker

File issues at: https://github.com/tailscale/tailscale/issues

This repo's issue tracker is disabled. (And all historical issues have been moved so the old URLs redirect)

## Installation

1.  Download precompiled [releases](https://github.com/tailscale/tailscale-synology/releases) from the page for SPKs for your platform. 
2.  In the Synology DSM web admin UI, open the Package Center.
3.  Press the *Manual install* button and provide the SPK file.
4.  Follow the wizard until done.
5.  At this point `tailscaled` should be up and running.
6.  SSH into the  machine, and run `sudo tailscale up` so you can authenticate.

> **_NOTE:_** If there is no SPK for your platform, you have to compile it yourself using the instructions [below](https://github.com/tailscale/tailscale-synology#making-packages).

## Upgrading

If upgrading to version v1.10.0, you may end up with duplicate installations of Tailscale. This is a [known](https://github.com/tailscale/tailscale/issues/2266#issuecomment-869792505) side effect of some metadata changes that were made in v1.10.0 in preparation of the installation package to be listed in the Synology Package Center. It is recommended to uninstall the old Tailscale package first before upgrading to v1.10.0. Please note that your devices Tailscale IP may change when v1.10.0 is installed.

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

> **_NOTE:_** For building on macOS the GNU core utilites are required. Homebrew users can run `brew install coreutils` and set the `PATH` variable accordingly.

## Credits and References

- Thanks to [@nirev](https://github.com/nirev) for creating this project and transferring it to Tailscale's GitHub org.
- https://haugene.github.io/docker-transmission-openvpn/synology-nas/ for the /dev/net/tun thing
- Package structure: [Synology Package Developer Guide](https://help.synology.com/developer-guide/index.html)
- Official Package building tool: [pkgscripts-ng](https://github.com/SynologyOpenSource/pkgscripts-ng)
- The package building process was originally based on [synology-wireguard](https://github.com/runfalk/synology-wireguard) \
If you need to _**compile**_ a synology package, check it out.
