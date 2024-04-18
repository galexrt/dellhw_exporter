## Linux Installation

Either use the container images available or download the binary and run it.

## Windows Installation

The binary supports the proper events and signals for using as a Windows service. Checkout [kardianos/service](https://github.com/kardianos/service) for more information.

Example to add the executable as a service in Windows:

```console
sc.exe create "Dell OMSA Exporter" binPath="C:\Program Files\Dell\dellhw_exporter.exe" start=auto
```

## Kubernetes

A Helm Chart is available at https://github.com/galexrt/dellhw_exporter/tree/main/charts/
