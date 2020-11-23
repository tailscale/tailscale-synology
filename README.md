# Tailscale package for Synology NAS
![CI](https://github.com/nirev/synology-tailscale/workflows/CI/badge.svg)

Synology NAS package for Tailscale based on precompiled static binaries.

## Disclaimer

You use everything here at your own risk. I am not responsible if this
breaks your NAS. Realistically it should not result in data loss or make
your NAS unaccessible, but one never knows.

## Installation

Check the [releases](https://github.com/nirev/synology-tailscale/releases)
page for SPKs for your platform. If there is no SPK you have to compile
it yourself using the instructions below.

1.  In the Synology DSM web admin UI, open the Package Center.
2.  Press the *Manual install* button and provide the SPK file.
3.  Follow the wizard until done.
4.  At this point `tailscaled` should be up and running.
5.  SSH into the  machine, and run `sudo tailscale up` so you can authenticate.

## Compatibility list

All models marked as *working* have been confirmed by users to work. If
your model has the same platform as one of the working ones, chances are
it will work for you too.

| Model     | Platform   | DSM version | arch | Working? |
| --------- | ---------- | ----------- | ----- | -------- |
| DS115j    | armv7l     | 6.2         | arm   | Yes      |
| DS212j    | armv5tel   | 6.2         | arm   | Yes      |
| DS213j    | armada370  | 6.2         | arm   | Yes      |
| DS215j    | armada375  | 6.2         | arm   | Yes      |
| DS214+    | armadaxp   | 6.2         | arm   | Yes      |
| DS216play | monaco     | 6.2         | arm   | Yes      |
| DS216se   | armada370  | 6.2         | arm   | Yes      |
| DS218+    | apollolake | 6.2         | x64   | Yes      |
| DS218j    | armada38x  | 6.2         | arm   | Yes      |
| DS220+    | geminilake | 6.2         | x64   | Yes      |
| DS220j    | rtd1296    | 6.2         | arm64 | Yes      |
| DS413j    | armv5tel   | 6.2         | arm   | Yes      |
| DS415+    | avoton     | 6.2         | x64   | Yes      |
| DS420+    | geminilake | 6.2         | x64  | Yes      |
| DS720+    | geminilake | 6.2         | x64   | Yes      |
| DS916+    | braswell   | 6.2         | x64   | Yes      |
| DS918+    | apollolake | 6.2         | x64   | Yes      |
| DS920+    | geminilake | 6.2         | x64   | Yes      |
| DS1812+   | cedarview  | 6.2         | x64   | Yes      |
| DS1815+   | avoton     | 6.2         | x64   | Yes      |
| DS2015xs  | alpine     | 6.2         | arm   | Yes      |

Please note that the package is currently being generated based on
Tailscale [static binaries](https://pkgs.tailscale.com/stable/#static), so
if your NAS has any of the supported architectures (x86, x86_64, arm, arm64)
it should theoretically work.

## Making packages

This project builds Synology packages "by hand", based on pre-compiled tailscale static binaries.

You can build the packages using `make`
```bash
git clone https://github.com/nirev/synology-tailscale.git
cd synology-tailscale/
make
```
If everything worked you should have a directory called `spks` that
contains your SPK files.

## Credits and References

- [Tailscale](https://github.com/tailscale) for Tailscale
- https://haugene.github.io/docker-transmission-openvpn/synology-nas/ for the /dev/net/tun thing
- Package structure: [Synology Package Developer Guide](https://help.synology.com/developer-guide/index.html)
- Official Package building tool: [pkgscripts-ng](https://github.com/SynologyOpenSource/pkgscripts-ng)
- The package building process was originally based on [synology-wireguard](https://github.com/runfalk/synology-wireguard) \
If you need to _**compile**_ a synology package, check it out.
