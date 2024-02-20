## 1.13.9 / 2024-02-16

* [BUGFIX] log the command that failed to execute

## 1.13.8 / 2024-02-16

* [BUGFIX] fix vdisk for (older?) omreport outputs

## 1.13.7 / 2024-02-15

* [CHORE] updated minimum go version to 1.21

## 1.13.6 / 2024-02-15

* [BUGFIX] add vdisk read and write policy to vdisk collector to address final parts of [#93](https://github.com/galexrt/dellhw_exporter/issues/93)

## 1.13.5 / 2024-02-15

* [BUGFIX] add "Non-Raid" state to pdisk collector to address parts of [#93](https://github.com/galexrt/dellhw_exporter/issues/93)
* [BUGFIX] add logging to vdisk collector

## 1.13.4 / 2024-02-13

* [BUGFIX] add logging to pdisk collector

## 1.13.3 / 2024-02-13

* [BUGFIX] improve log lines to better be able to pin point the recent parsing issues

## 1.13.2 / 2024-02-06

* [BUGFIX] [Consider 'Not Applicable' as healthy for Nic status #95](https://github.com/galexrt/dellhw_exporter/pull/95)
    * Thanks to [@B0go](https://github.com/B0go) for fixing this issue!

## 1.13.1 / 2023-12-07

* [BUGFIX] Fix container image build issue caused by wget, use curl now

## 1.13.0 / 2023-12-07

* [ENHANCEMENT] [Allow for user to specify a list of interfaces to monitor #89](https://github.com/galexrt/dellhw_exporter/pull/89)
* [ENHANCEMENT] Added Storage Pdisk Hardware Encryption status
* [SECURITY] Updated dependencies to latest version

## 1.12.2 / 2022-05-31

* [SECURITY] update gopkg.in/yaml.v3 to v3.0.1 (CVE-2022-28948)

## 1.12.1 / 2022-05-04

* [ENHANCEMENT] update deps to latest version
* [ENHANCEMENT] updated minimum go version to 1.18

## 1.12.0 / 2022-02-02

* [ENHANCEMENT] Added Pdisk Remaining Rated Write Endurance Metric by @adityaborgaonkar

## 1.11.1 / 2021-10-12

* [ENHANCEMENT] update go version to 1.16

## 1.11.0 / 2021-09-12

* [ENHANCEMENT] add vdisk raid level metric
  * This adds `dell_hw_storage_vdisk_raidlevel` metric, which holds the RAID
    level of the VDISK.
    Additionally the controller ID label was added to some metrics missing
    it. Resolves #8

## 1.10.0 / 2021-08-30

* [ENHANCEMENT] add pdisk "predicted failure" metric

## 1.9.0 / 2021-08-29

* [ENHANCEMENT] update go version and deps

## 1.8.0 / 2020-10-07

* [ENHANCEMENT] Windows Service Support
  * Thanks to [@kyle-williams-1](https://github.com/kyle-williams-1) for adding this feature!
* [ENHANCEMENT] Kubernetes Helm chart

## 1.7.0 / 2020-09-29

* [ENHANCEMENT] Metric results can be cached to improve performance.
  * Thanks to [@Phil1602](https://github.com/Phil1602) for adding this as a feature!
* [ENHANCEMENT] The default value of the `--collectors-omreport` flag is now dependent on the OS for Linux and Windows.
  * Thanks to [@kyle-williams-1](https://github.com/kyle-williams-1) for adding this as a feature!
* [ENHANCEMENT] Enabled `windows/amd64` release binary builds.
* [ENHANCEMENT] Golang 1.15 is used by default for CI and build.
* [ENHANCEMENT] Updated LICENSE file and go code file headers.
* [ENHANCEMENT] Created documentation page using [mkdocs](https://www.mkdocs.org/), available at [dellhw-exporter.galexrt.moe](https://dellhw-exporter.galexrt.moe/).

## 1.6.0 / 2020-06-09

* [ENHANCEMENT] Add support for firmware versions #43 (PR #44).
  * Thanks to [@sfudeus](https://github.com/sfudeus) for implementing this!
* [ENHANCEMENT] docker: added expose for 9137/tcp exporter port

## 1.5.19 / 2020-06-07

* [BUGFIX] ci: debug using tmate action

## 1.5.18 / 2020-06-07

* [BUGFIX] ci: debug using tmate action

## 1.5.17 / 2020-06-07

* [BUGFIX] ci: fix build routine issues #42

## 1.5.16 / 2020-06-07

* [ENHANCEMENT] ci: no need to specify docker build dir

## 1.5.15 / 2020-06-07

* [ENHANCEMENT] docker: fix copy path for binary

## 1.5.14 / 2020-06-07

* [ENHANCEMENT] ci: use correct env vars for image name

## 1.5.13 / 2020-06-07

* [ENHANCEMENT] ci: use correct env vars for image name

## 1.5.11 / 2020-06-07

* [ENHANCEMENT] ci: use github actions

## 1.5.9 / 2020-06-07

* [ENHANCEMENT] ci: use github actions

## 1.4.2 / 2020-02-24

* [ENHANCEMENT] ci: fixed CI release upload
