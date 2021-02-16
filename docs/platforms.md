# Synology platform names and package architecture

This document exist as a non-exhaustive list of Synology models and their mapping to tailscale package architecture to help users identify which package to use, until the Tailscale package is available officially in Synology's package repository.


### FS Series
| Model  | CPU                        | Synology Platform Name | Tailscale Package Arch |
| ------ | -------------------------- | ---------------------- | ---------------------- |
| FS6400 | Intel Xeon Silver 4110 x 2 | purley                 | amd64                  |
| FS3600 | Intel Xeon D-1567          | broadwellnk            | amd64                  |
| FS3400 | Intel Xeon D-1541          | broadwell              | amd64                  |
| FS3017 | Intel Xeon E5-2620 v3 x 2  | grantley               | amd64                  |
| FS2017 | Intel Xeon D-1541          | broadwell              | amd64                  |
| FS1018 | Intel Pentium D1508        | broadwellnk            | amd64                  |
---
### SA Series
| Model   | CPU                   | Synology Platform Name | Tailscale Package Arch |
| ------- | --------------------- | ---------------------- | ---------------------- |
| SA3600  | Intel Xeon D-1567     | broadwellnk            | amd64                  |
| SA3400  | Intel Xeon D-1541     | broadwellnk            | amd64                  |
| SA3200D | Intel Xeon D-1521 x 2 |                        | amd64                  |
---
### x21 Series
| Model                       | CPU               | Synology Platform Name | Tailscale Package Arch |
| --------------------------- | ----------------- | ---------------------- | ---------------------- |
| RS4021xs+, RS3621xs+        | Intel Xeon D-1541 | broadwellnk            | amd64                  |
| RS3621RPxs                  | Intel Xeon D-1531 | broadwellnk            | amd64                  |
| DS1821+, RS1221+, RS1221RP+ | AMD Ryzen V1500B  | v1000                  | amd64                  |
| DS1621xs+                   | Intel Xeon D-1527 | broadwellnk            | amd64                  |
| DS1621+                     | AMD Ryzen V1500B  | v1000                  | amd64                  |
---
### x20 Series
| Model                   | CPU                 | Synology Platform Name | Tailscale Package Arch |
| ----------------------- | ------------------- | ---------------------- | ---------------------- |
| RS820+, RS820RP+        | Intel Atom C3538    | denverton              | amd64                  |
| DS720+, DS920+, DS1520+ | Intel Celeron J4125 | geminilake             | amd64                  |
| DS620slim               | Intel Celeron J3355 | apollolake             | amd64                  |
| DS220+, DS420+          | Intel Celeron J4025 | geminilake             | amd64                  |
| DS220j, DS420j          | Realtek RTD1296     | rtd1296                | arm64                  |
| DS120j                  | Marvell A3720       | armada37xx             | arm64                  |
---
### x19 Series
| Model            | CPU                  | Synology Platform Name | Tailscale Package Arch |
| ---------------- | -------------------- | ---------------------- | ---------------------- |
| RS1219+          | Intel Atom C2538     | avoton                 | amd64                  |
| RS1619xs+        | Intel Xeon D-1527    | broadwell              | amd64                  |
| RS819            | Realtec RTD1296      | rtd1296                | arm64                  |
| DS1819+, DS2419+ | Intel Atom C3538     | denverton              | amd64                  |
| DS1019+          | Intel Celeron J3455  | apollolake             | amd64                  |
| DS119j           | Marvell A3720        | armada37xx             | arm64                  |
| DS419slim        | Marvell A385 88F6820 | armada38x              | arm                    |
---
### x18 Series
| Model                                  | CPU                        | Synology Platform Name | Tailscale Package Arch |
| -------------------------------------- | -------------------------- | ---------------------- | ---------------------- |
| FS1018, DS3018xs                       | Intel Pentium D1508        | broadwellnk            | amd64                  |
| RS3618xs                               | Intel Xeon D-1521          | broadwell              | amd64                  |
| DS718+, DS918+                         | Intel Celeron J3455        | apollolake             | amd64                  |
| DS218+, DS418play                      | Intel Celeron J3355        | apollolake             | amd64                  |
| RS2818RP+, RS2418+, RS2418RP+, DS1618+ | Intel Atom C3538           | denverton              | amd64                  |
| RS818+, RS818RP+                       | Intel Atom C2538           | avoton                 | amd64                  |
| DS118, DS218, DS218play, DS418         | Realtek RTD1296            | rtd1296                | arm64                  |
| DS418j                                 | Realtek RTD1293            | rtd1296                | arm64                  |
| DS218j                                 | Marvell Armada 385 88F6820 | armada38x              | arm                    |
---
### x17 Series
| Model                 | CPU                          | Synology Platform Name | Tailscale Package Arch |
| --------------------- | ---------------------------- | ---------------------- | ---------------------- |
| FS3017                | Intel Xeon E5-2620 v3 x 2    | grantley               | amd64                  |
| FS2017, RS4017xs+     | Intel Xeon D-1541            | broadwell              | amd64                  |
| RS18017xs+, RS3617xs+ | Intel Xeon D-1531            | broadwell              | amd64                  |
| DS3617xs              | Intel Xeon D-1527            | broadwell              | amd64                  |
| RS3617RPxs            | Intel Xeon D-1521            | broadwell              | amd64                  |
| RS3617xs              | Intel Xeon E3-1230 v2        | bromolow               | amd64                  |
| DS1517+, DS1817+      | Intel Atom C2538             | avoton                 | amd64                  |
| RS217                 | Marvell Armada 385 88F6820   | armada38x              | arm                    |
| DS1517, DS1817        | Annapurna Labs Alpine AL-314 | alpine                 | arm                    |
---
### x16 Series
| Model                                  | CPU                          | Synology Platform Name | Tailscale Package Arch |
| -------------------------------------- | ---------------------------- | ---------------------- | ---------------------- |
| RS18016xs+                             | Intel Xeon E3-1230 v2        | bromolow               | amd64                  |
| RS2416+, RS2416RP+                     | Intel Atom C2538             | avoton                 | amd64                  |
| DS916+                                 | Intel Pentium N3710          | braswell               | amd64                  |
| DS716+II                               | Intel Celeron N3160          | braswell               | amd64                  |
| DS716+                                 | Intel Celeron N3150          | braswell               | amd64                  |
| DS216+II, DS416play                    | Intel Celeron N3060          | braswell               | amd64                  |
| DS216+                                 | Intel Celeron N3050          | braswell               | amd64                  |
| DS416                                  | Annapurna Labs Alpine AL-212 | alpine                 | arm                    |
| DS116, DS216j, DS216, DS416slim, RS816 | Marvell Armada 385 88F6820   | armada38x              | arm                    |
| DS416j                                 | Marvell Armada 388 88F6828   | armada38x              | arm                    |
| DS216se                                | Marvell Armada 370 88F6707   | armada370              | arm                    |
| DS216play                              | STM STiH412                  | monaco                 | arm                    |
---
### x15 Series
| Model                                               | CPU                          | Synology Platform Name | Tailscale Package Arch |
| --------------------------------------------------- | ---------------------------- | ---------------------- | ---------------------- |
| RC18015xs+                                          | Intel Xeon E3-1230 v2        | bromolow               | amd64                  |
| DS3615xs                                            | Intel Core i3-4130           | bromolow               | amd64                  |
| DS415+, DS1515+, DS2415+, DS1815+, RS815+, RS815RP+ | Intel Atom C2538             | avoton                 | amd64                  |
| DS415play                                           | Intel Atom CE5335            | evansport              | 386                    |
| DS2015xs                                            | Annapurna Labs Alpine AL-514 | alpine                 | arm                    |
| DS715, DS1515                                       | Annapurna Labs Alpine AL-314 | alpine                 | arm                    |
| DS215+                                              | Annapurna Labs Alpine AL-212 | alpine                 | arm                    |
| RS815                                               | Marvell Armada XP MV78230    | armadaxp               | arm                    |
| DS115, DS215j                                       | Marvell Armada 375 88F6720   | armada375              | arm                    |
| DS115j                                              | Marvell Armada 370 88F6707   | armada370              | arm                    |
---
### x14 Series
| Model                                   | CPU                        | Synology Platform Name | Tailscale Package Arch |
| --------------------------------------- | -------------------------- | ---------------------- | ---------------------- |
| RS3614xs+                               | Intel Xeon E3-1230 v2      | bromolow               | amd64                  |
| RS3614xs, RS3614RPxs                    | Intel Core i3-4130         | bromolow               | amd64                  |
| RS814+, RS814RP+, RS2414+, RS2414RP+    | Intel Atom D2700           | cedarview              | amd64                  |
| DS214, DS214+, DS414, RS814             | Marvell Armada XP MV78230  | armadaxp               | arm                    |
| DS114, DS214se, DS414slim, RS214, EDS14 | Marvell Armada 370 88F6707 | armada370              | arm                    |
| DS414j                                  | Mindspeed Comcerto C2000   | comcerto2k             | arm                    |
| DS214play                               | Intel Atom CE5335          | evansport              | 386                    |
---
### x13 Series
| Model                             | CPU                        | Synology Platform Name | Tailscale Package Arch |
| --------------------------------- | -------------------------- | ---------------------- | ---------------------- |
| RS3413xs+, RS10613xs+             | Intel Xeon E3-1230 v2      | bromolow               | amd64                  |
| DS713+, DS1513+, DS1813+, DS2413+ | Intel Atom D2700           | cedarview              | amd64                  |
| DS213+, DS413                     | Freescale P1022            | qoriq                  | unsupported            |
| DS213air, DS213, DS413j           | Marvell Kirkwood 88F6282   | 88f6281                | arm                    |
| DS213j                            | Marvell Armada 370 88F6707 | armada370              | arm                    |
---
### x12 Series
| Model                                                          | CPU                      | Synology Platform Name | Tailscale Package Arch |
| -------------------------------------------------------------- | ------------------------ | ---------------------- | ---------------------- |
| DS3612xs, RS3412xs, RS3412RPxs                                 | Intel Core i3-2100       | bromolow               | amd64                  |
| DS412+, DS1512+, DS1812+, RS812+, RS812RP+, RS2212+, RS2212RP+ | Intel Atom D2700         | cedarview              | amd64                  |
| DS712+                                                         | Intel Atom D425          | x86                    | amd64                  |
| DS112, DS112+, DS212, DS212+, RS212, RS812                     | Marvell Kirkwood 88F6282 | 88f6281                | arm                    |
| DS212j                                                         | Marvell Kirkwood 88F6281 | 88f6281                | arm                    |
| DS112j                                                         | Marvell Kirkwood 88F6702 | 88f6281                | arm                    |
---
### x11 Series
| Model                                          | CPU                      | Synology Platform Name | Tailscale Package Arch |
| ---------------------------------------------- | ------------------------ | ---------------------- | ---------------------- |
| DS3611xs, RS3411xs, RS3411RPxs                 | Intel Core i3-2100       | bromolow               | amd64                  |
| DS411+II, DS1511+, DS2411+, RS2211+, RS2211RP+ | Intel Atom D525          | x86                    | amd64                  |
| DS411+                                         | Intel Atom D510          | x86                    | amd64                  |
| DS111, DS211, DS211+, DS411slim, DS411, RS411  | Marvell Kirkwood 88F6282 | 88f6281                | arm                    |
| DS211j, DS411j                                 | Marvell Kirkwood 88F6281 | 88f6281                | arm                    |
---
### x10 Series
| Model                     | CPU                      | Synology Platform Name | Tailscale Package Arch |
| ------------------------- | ------------------------ | ---------------------- | ---------------------- |
| DS1010+, RS810+, RS810RP+ | Intel Atom D510          | x86                    | amd64                  |
| DS710+                    | Intel Atom D410          | x86                    | amd64                  |
| DS110+, DS210+, DS410     | Freescale MPC8533E       | ppc853x                | unsupported            |
| DS110j, DS210j, DS410j    | Marvell Kirkwood 88F6281 | 88f6281                | arm                    |
---
### Routers
| Model    | CPU               | Synology Platform Name | Tailscale Package Arch |
| -------- | ----------------- | ---------------------- | ---------------------- |
| MR2200ac | Qualcomm IPQ4019  | dakota                 | arm                    |
| RT2600ac | Qualcomm IPQ8065  | ipq806x                | arm                    |
| RT1900ac | Broadcom BCM58622 | northstarplus          | arm                    |
---
### Network Video Recorders
| Model   | CPU              | Synology Platform Name | Tailscale Package Arch |
| ------- | ---------------- | ---------------------- | ---------------------- |
| DVA3221 | Intel Atom C3538 | denverton              | amd64                  |
| DVA3219 | Intel Atom C3538 | denverton              | amd64                  |
| NVR1218 | HiSilicon Hi3535 | hi3535                 | arm                    |
| NVR216  | HiSilicon Hi3535 | hi3535                 | arm                    |
---
### Video Stations
| Model   | CPU              | Synology Platform Name | Tailscale Package Arch |
| ------- | ---------------- | ---------------------- | ---------------------- |
| VS960HD | HiSilicon Hi3536 | hi3536                 |                        |
| VS360HD | HiSilicon Hi3535 | hi3535                 | arm                    |
---
