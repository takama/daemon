%global go_import_path  github.com/takama/daemon

Name:           daemon
Version:        0.2.2
Release:        1%{?dist}
Summary:        A daemon package for use with Go (golang) services with no dependencies
License:        MIT
URL:            https://%{go_import_path}
Source0:        https://%{go_import_path}/archive/%{version}.tar.gz
BuildRequires:  golang

%description
%{summary}

%package devel
Requires:       golang
Summary:        A daemon package for use with Go (golang) services with no dependencies
Provides:       golang(%{go_import_path}) = %{version}-%{release}

%description devel
%{summary}

%prep

%setup -n daemon-%{version}

%build

%install
install -d %{buildroot}/%{gopath}/src/%{go_import_path}
for i in `ls -1|egrep -iv 'license|readme|\.spec'`; do
cp -ap $i %{buildroot}/%{gopath}/src/%{go_import_path}/
done

%clean
%{__rm} -rf %{buildroot}

%check
#GOPATH=%{buildroot}/%{gopath} go test %{go_import_path}

%files devel
%defattr(-,root,root,-)
%doc README.md LICENSE
%dir %attr(755,root,root) %{gopath}/src/%{go_import_path}
%dir %attr(755,root,root) %{gopath}/src/%{go_import_path}/example
%{gopath}/src/%{go_import_path}/*.go
%{gopath}/src/%{go_import_path}/example/*.go

%changelog
* Mon Oct 20 2014 Igor Dolzhikov - 0.2.2
- fix rpm spec
