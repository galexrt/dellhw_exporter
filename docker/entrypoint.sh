#!/bin/sh

echo "Starting dell srvadmin services ..."
/opt/dell/srvadmin/sbin/dsm_sa_datamgrd &
/opt/dell/srvadmin/sbin/dsm_sa_eventmgrd &
/opt/dell/srvadmin/sbin/dsm_sa_snmpd &
/usr/libexec/instsvcdrv-helper start &

wait
echo "Started dell srvadmin services."

exec dellhw_exporter "$@"
