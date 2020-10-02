## Linux Installation
<!-- TODO -->

## Windows Installation

The binary supports the proper events and signals for using as a Windows service. Checkout [kardianos/service](https://github.com/kardianos/service) for more information.

Example to add the executable as a service in Windows:

```console
sc.exe create "Dell OMSA Exporter" binPath="C:\Program Files\Dell\dellhw_exporter.exe" start=auto
```