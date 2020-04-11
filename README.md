# Tailscale package for Synology NAS

Synology NAS package for Tailscale based on precompiled static binaries.

## Disclaimer

You use everything here at your own risk. I am not responsible if this
breaks your NAS. Realistically it should not result in data loss or make
your NAS unaccessible, but one never knows.

Compatibility list
------------------

All models marked as *Is working* have been confirmed by users to work. If
your model has the same platform as one of the working ones, chances are
it will work for you too.

| Model  | Platform   | DSM version | arch | Working? |
| ------ | ---------- | ----------- | ---- | -------- |
| DS218+ | apollolake | 6.2         | x64  | Yes      |


Please not that the package is currently being generated based on
Tailscale [static binaries](https://pkgs.tailscale.com/stable/#static)
So if your NAS has any of the supported architectures (x86, x86_64, arm, arm64)
it should theoretically work.

## Installation

Check the [releases](https://github.com/nirev/synology-tailscale/releases)
page for SPKs for your platform. If there is no SPK you have to compile
it yourself using the instructions below.

1.  In the Synology DSM web admin UI, open the Package Center.
2.  Press the *Manual install* button and provide the SPK file.
3.  Follow the wizard until done.
4.  At this point `tailscaled` should be up and running.
5.  SSH into the  machine, and run `sudo tailscale up` so you can authenticate.

## Compiling

This repo is heavily based on [synology-wireguard](https://github.com/runfalk/synology-wireguard).
Likewise, everything is assembled inside docker, as Synology's package building tool `pkgscripts-ng`
clutters the file system quite a bit.

First create the base docker image, which downloads `pkgscripts-ng`:

```bash
git clone https://github.com/nirev/synology-tailscale.git
cd synology-tailscale/
make build-image
```

Now we can build for any platform and DSM version using:

```bash
make build SYNO_PLATFORM=<platform> SYNO_DSM_VERSION=<dsm-ver>
```

You should replace `<platform>` with your NAS's package arch. Using
`this table <https://www.synology.com/en-global/knowledgebase/DSM/tutorial/General/What_kind_of_CPU_does_my_NAS_have>`\_
you can figure out which one to use. Note that the package arch must be
lowercase. `<dsm-ver>` should be replaced with the version of DSM you
are compiling for with just `major.minor` (eg "6.2")

For the DS218+ that I have, the complete command looks like this:

```bash
make build SYNO_PLATFORM=apollolake SYNO_DSM_VERSION=6.2
```

If everything worked you should have a directory called `artifacts` that
contains your SPK files.

## Credits

- [runfalk](https://github.com/runfalk/synology-wireguard) for synology package building skills
- [Tailscale](https://github.com/tailscale) for Tailscale
- https://haugene.github.io/docker-transmission-openvpn/synology-nas/ for the /dev/net/tun thing
