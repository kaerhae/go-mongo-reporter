Name:           reporter
Version:        {VERSION}
Release:        {BUILD_NUMBER}%{?dist}
Summary:        Reporter server written in Go

License:        GPLv3
Source0:        go-mongo-reporter-{VERSION}.tar.gz

%global debug_package %{nil}


%description
Reporter backend service written in Go.



%prep
%setup


%install
rm -rf $RPM_BUILD_ROOT
mkdir -p %{buildroot}/etc/opt/reporter/config
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/etc/systemd/system
install -Dpm 0755 reporter %{buildroot}/usr/local/bin
install -Dpm 644 reporter.service %{buildroot}/etc/systemd/system


%files
%dir /etc/opt/reporter/config
/usr/local/bin/reporter
/etc/systemd/system/reporter.service

%changelog
* Sun Jul 21 2024 vagrant
- 