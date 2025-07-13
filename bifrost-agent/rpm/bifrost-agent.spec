Name:           bifrost-agent
Version:        1.0
Release:        1%{?dist}
Summary:        Bifrost Agent Python collector

License:        MIT
URL:            https://github.com/fmendonca/bifrostv2
Source0:        bifrost-agent-1.0.tar.gz

BuildArch:      noarch
Requires:       python3, libvirt

%description
Python agent that collects VM data from Libvirt and sends it to a remote Bifrost API.

%prep

%build

%install
mkdir -p %{buildroot}/opt/bifrost-agent
mkdir -p %{buildroot}/var/log/bifrost
mkdir -p %{buildroot}/etc/systemd/system
mkdir -p %{buildroot}/etc/logrotate.d

install -m 0755 bifrost-agent.py %{buildroot}/opt/bifrost-agent/bifrost-agent.py
install -m 0644 bifrost-agent.service %{buildroot}/etc/systemd/system/bifrost-agent.service
install -m 0644 bifrost-agent.logrotate %{buildroot}/etc/logrotate.d/bifrost-agent

%post
/usr/bin/systemctl daemon-reload
/usr/bin/systemctl enable --now bifrost-agent

%preun
if [ $1 -eq 0 ]; then
  /usr/bin/systemctl disable --now bifrost-agent
fi

%files
/opt/bifrost-agent/bifrost-agent.py
/etc/systemd/system/bifrost-agent.service
/etc/logrotate.d/bifrost-agent
/var/log/bifrost

%changelog
* Tue Jul 16 2025 You filipecm@gmail.com - 1.0-1
- Initial RPM package
