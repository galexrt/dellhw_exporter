FROM centos:7

ARG BUILD_DATE="N/A"
ARG REVISION="N/A"

ARG DELLHW_EXPORTER_VERSION="N/A"

LABEL org.opencontainers.image.authors="Alexander Trost <galexrt@googlemail.com>" \
    org.opencontainers.image.created="${BUILD_DATE}" \
    org.opencontainers.image.title="galexrt/dellhw_exporter" \
    org.opencontainers.image.description="Prometheus exporter for Dell Hardware components using OMSA." \
    org.opencontainers.image.documentation="https://github.com/galexrt/dellhw_exporter/blob/main/README.md" \
    org.opencontainers.image.url="https://github.com/galexrt/dellhw_exporter" \
    org.opencontainers.image.source="https://github.com/galexrt/dellhw_exporter" \
    org.opencontainers.image.revision="${REVISION}" \
    org.opencontainers.image.vendor="galexrt" \
    org.opencontainers.image.version="${DELLHW_EXPORTER_VERSION}"

# Environment variables
ENV PATH="$PATH:/opt/dell/srvadmin/bin:/opt/dell/srvadmin/sbin" \
    SYSTEMCTL_SKIP_REDIRECT="1" \
    START_DELL_SRVADMIN_SERVICES="true"

# Do overall update and install missing packages needed for OpenManage
RUN yum -y update && \
    yum -y install sysvinit-tools wget perl passwd gcc which tar libstdc++.so.6 compat-libstdc++-33.i686 glibc.i686 make && \
    wget -q -O - "https://linux.dell.com/repo/hardware/dsu/bootstrap.cgi" | bash && \
    rpm --import "https://linux.dell.com/repo/pgp_pubkeys/0x1285491434D8786F.asc" && \
    yum -y install srvadmin-base srvadmin-storageservices && \
    yum clean all

EXPOSE 9137/tcp

ADD container/entrypoint.sh /bin/entrypoint

RUN chmod +x /bin/entrypoint

ADD .build/linux-amd64/dellhw_exporter /bin/dellhw_exporter

ENTRYPOINT ["/bin/entrypoint"]
