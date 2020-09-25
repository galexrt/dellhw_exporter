## No metrics being exported

If you are not running the Docker container, it is probably that your OMSA / srvadmin services are not running. Start them using the following commands:

```console
/opt/dell/srvadmin/sbin/srvadmin-services.sh status
/opt/dell/srvadmin/sbin/srvadmin-services.sh start
echo "return code: $?"
```
Please note that the return code should be `0`, if not please investigate the logs of srvadmin services.

When running inside the container this most of the time means
Be sure to enter the container and run the following commands to verify if the kernel modules have been loaded:

```console
/usr/libexec/instsvcdrv-helper status
lsmod | grep -iE 'dell|dsu'
```

Should the `lsmod` not contain any module named after `dell_` and / or `dsu_`, be sure to add the following read-only mounts depending on your OS, for the kernel modules directory (`/lib/modules`) and / or the kernel source / headers directory (depends hardly on the OS your are using) to the `dellhw_exporter` Docker container using `-v HOST_PATH:CONTAINER_PATH:ro` flag.
