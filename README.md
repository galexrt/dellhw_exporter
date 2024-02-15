# dellhw_exporter

![build_release](https://github.com/galexrt/dellhw_exporter/workflows/build_release/badge.svg)

Prometheus exporter for Dell Hardware components using OMSA.

The exporter was originally made by [PrFalken](https://github.com/PrFalken). Due to some issues in the code, I rewrote the whole exporter using the ["node_exporter"](https://github.com/prometheus/node_exporter) pattern and therefore moved it from being a fork out, to a standalone repository.

Omreport parsing functions were borrowed from the [Bosun project](https://github.com/bosun-monitor/bosun/blob/master/cmd/scollector/collectors/dell_hw.go), thank you very much for that, they are the most tedious part of the job.

This exporter wraps the "omreport" command from Dell OMSA. If you can't run omreport on your system, the exporter won't export any metrics.

## Compatibility

### Tested Dell OMSA Compatibility

The dellhw_exporter has been tested with the following OMSA versions:

* `7.4`
* `8.4`
* `9.1`

### Kernel Compatibility

**Please note that only kernel versions that are supported by DELL DSU / OMSA tools are working!**

**State 07.06.2019**: Dell OMSA `DSU_19.05.00` is not compatible with 5.x kernel it seems (e.g., Fedora uses that kernel).

Should you run into issues when using the Container image, please follow the [Troubleshooting - No metrics being exported](#no-metrics-being-exported).

## Collectors

For a list of the available collectors, see [Collectors doc page](docs/collectors.md).

## Configuration

For flags and environment variables, see [Configuration doc page](docs/configuration.md).

## Caching

Optional caching can be enabled to prevent performance issues caused by this exporter, see [Caching doc page](docs/caching.md).

## Running in Container

Container Image available from:

* [Quay.io](https://quay.io/repository/galexrt/dellhw_exporter)
* [GHCR.io](https://github.com/users/galexrt/packages/container/package/dellhw_exporter)

(Docker Hub is not receiving new images begining from release `v1.10.0`)

### Pull the Container Image

Each Git tag is built and published as a release and container image.

#### Quay.io

```console
docker pull quay.io/galexrt/dellhw_exporter
```

#### GHCR.IO

```console
docker pull ghcr.io/galexrt/dellhw_exporter
```

### Run the Container Image

> **NOTE** The `--privileged` flag is required as the OMSA needs to access the host's devices and other components.

```console
docker run -d --name dellhw_exporter --privileged -p 9137:9137 quay.io/galexrt/dellhw_exporter
```

## Running without Docker / Podman

To run without Docker / Podman either download a [release binary](https://github.com/galexrt/dellhw_exporter/releases) or build it (using `make build` command):

```console
./dellhw_exporter
./dellhw_exporter --help
./dellhw_exporter YOUR_FLAGS
```

**The DELL OMSA services must already be running for the exporter to be able to collect metrics!**

E.g., run `/opt/dell/srvadmin/sbin/srvadmin-services.sh start` and / or `systemctl start SERVICE_NAME` (to enable autostart use `systemctl enable SERVICE_NAME`; where `SERVICE_NAME` [are the DELL OMSA service(s) you installed](http://linux.dell.com/repo/hardware/omsa.html)).

## Prometheus

The exporter runs on port `9137` TCP.

Example static Prometheus Job config:

```yaml
[...]
  - job_name: 'dellhw_exporter'
    # Override the global default and scrape targets from this job every 60 seconds.
    scrape_interval: 60s
    static_configs:
      - targets:
        - 'YOUR_SERVER_HERE:9137'
[...]
```

## Monitoring

Checkout the files in the [`contrib/monitoring/`](contrib/monitoring/) directory.

## Installation

See [Installation doc page](docs/installation.md).

## Troubleshooting

See [Troubleshooting doc page](docs/troubleshooting.md).

## Development

Golang version `1.21` is used for testing and building the dellhw_exporter.

`go mod` is used for "vendoring" of the dependencies.

### Creating a new Release

1. Update the version.
    1. [`VERSION` file](VERSION)
    2. Helm chart: `charts/dellhw_exporter/Chart.yaml` `appVersion:` line and bump the Helm chart `version:` by a patch release version.
2. Create an entry in the [`CHANGELOG.md` file](CHANGELOG.md).
    Example of a changelog entry:

    ```
    ## 1.12.0 / 2022-02-02

    * [ENHANCEMENT] Added Pdisk Remaining Rated Write Endurance Metric by @adityaborgaonkar
    * [BUGFIX] ci: fix build routine issues
    ```

    The following "kinds" of entries can be added:

    * `CHANGE`
    * `FEATURE`
    * `ENHANCEMENT`
    * `BUGFIX`
3. Commit the version increase with a commit messages in format: `version: bump to v1.12.0`
4. Create the `git` tag using `git tag v1.12.0`
5. Now push the changes and commit using `git push && git push --tags`
6. In a few minutes the new release should be available for "approval" under the [releases section](https://github.com/galexrt/dellhw_exporter/releases). Edit and save the release on GitHub and the release is complete.
