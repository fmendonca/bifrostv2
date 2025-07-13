Name:           bifrost-agent
Version:        1.0
Release:        1%{?dist}
Summary:        Bifrost Agent Python collector

License:        MIT
URL:            https://github.com/fmendonca/bifrostv2
Source0:        %{name}-%{version}.tar.gz

BuildArch:      noarch
Requires:       python3, libvirt

%description
Python agent that collects VM data from Libvirt and sends it to a remote Bifrost API.

%prep
%autosetup

%build
# Nothing to build

%install
mkdir -p %{buildroot}/opt/bifrost-agent
mkdir -p %{buildroot}/var/log/bifrost
mkdir -p %{buildroot}/etc/systemd/system
mkdir -p %{buildroot}/etc/logrotate.d

install -m 0755 bifrost-agent.py %{buildroot}/opt/bifrost-agent/bifrost-agent.py
install -m 0644 bifrost-agent.service %{buildroot}/etc/systemd/system/bifrost-agent.service
install -m 0644 bifrost-agent.logrotate %{buildroot}/etc/logrotate.d/bifrost-agent

%post
%systemd_post bifrost-agent.service

%preun
%systemd_preun bifrost-agent.service

%postun
%systemd_postun_with_restart bifrost-agent.service

%files
%license
/opt/bifrost-agent/bifrost-agent.py
/etc/systemd/system/bifrost-agent.service
/etc/logrotate.d/bifrost-agent
%dir /var/log/bifrost

%changelog
* Tue Jul 16 2024 Filipe Mendonça <filipecm@gmail.com> - 1.0-1
- Initial RPM package