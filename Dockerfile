FROM ubuntu:18.04

ENV IS_IN_CONTAINER 1

WORKDIR /

RUN apt-get update \
  && apt-get -qy install git python3 wget ca-certificates

RUN git clone https://github.com/nirev/synology-pkgscripts-ng /pkgscripts-ng

COPY build /source/WireGuard

ENTRYPOINT exec /source/WireGuard/build.sh
