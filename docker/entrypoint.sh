#!/bin/sh

START_DELL_SRVADMIN_SERVICES="${START_DELL_SRVADMIN_SERVICES:-true}"

if [ "${START_DELL_SRVADMIN_SERVICES}" = "true" ]; then
    echo "Starting dell srvadmin services ..."
    /opt/dell/srvadmin/sbin/dsm_sa_datamgrd &
    /opt/dell/srvadmin/sbin/dsm_sa_eventmgrd &
    /opt/dell/srvadmin/sbin/dsm_sa_snmpd &
    /usr/libexec/instsvcdrv-helper start &

    wait
    echo "Started dell srvadmin services."
else
    echo "Skipping start of dell srvadmin services."
fi

exec dellhw_exporter "$@"
