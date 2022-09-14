Name:           dellhw_exporter
Version:        1.12.2
Release:        1%{?dist}
Summary:        Prometheus exporter for Dell Hardware components using OMSA.

License:        Apache-2.0
URL:            https://github.com/galexrt/dellhw_exporter
Source0:        dellhw_exporter-1.12.2.tar.gz

Prefix: /usr
BuildRequires: golang-bin
#Requires:       

%define  debug_package %{nil}


%description

Prometheus exporter for Dell Hardware components using OMSA.


%prep
%setup -q


%build
make tree


%install
rm -rf $RPM_BUILD_ROOT
make install DESTDIR=$RPM_BUILD_ROOT

%files
%doc
%config(noreplace) %{_sysconfdir}/sysconfig/dellhw_exporter
%{prefix}/lib/systemd/system/dellhw_exporter.service
%{prefix}/sbin/dellhw_exporter

%changelog
