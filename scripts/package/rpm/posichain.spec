{{VER=2.0.0}}
{{REL=0}}
# SPEC file overview:
# https://docs.fedoraproject.org/en-US/quick-docs/creating-rpm-packages/#con_rpm-spec-file-overview
# Fedora packaging guidelines:
# https://docs.fedoraproject.org/en-US/packaging-guidelines/

Name:		posichain
Version:	{{ VER }}
Release:	{{ REL }}
Summary:	posichain blockchain validator node program

License:		MIT
URL:			https://posichain.org
Source0:		%{name}-%{version}.tar
BuildArch: 		x86_64
Packager: 		Danny <danny@position.exchange>
Requires(pre): 	shadow-utils
Requires: 		systemd-rpm-macros jq

%description
Posichain is a sharded, fast finality, low fee, PoS public blockchain.
This package contains the validator node program for posichain blockchain.

%global debug_package %{nil}

%prep
%setup -q

%build
exit 0

%check
./posichain --version
exit 0

%pre
getent group posichain >/dev/null || groupadd -r posichain
getent passwd posichain >/dev/null || \
   useradd -r -g posichain -d /home/posichain -m -s /sbin/nologin \
   -c "Posichain validator node account" posichain
mkdir -p /home/posichain/.psc/blskeys
mkdir -p /home/posichain/.config/rclone
chown -R posichain.posichain /home/posichain
exit 0


%install
install -m 0755 -d ${RPM_BUILD_ROOT}/usr/sbin ${RPM_BUILD_ROOT}/etc/systemd/system ${RPM_BUILD_ROOT}/etc/sysctl.d ${RPM_BUILD_ROOT}/etc/posichain
install -m 0755 -d ${RPM_BUILD_ROOT}/home/posichain/.config/rclone
install -m 0755 posichain ${RPM_BUILD_ROOT}/usr/sbin/
install -m 0755 posichain-setup.sh ${RPM_BUILD_ROOT}/usr/sbin/
install -m 0755 posichain-rclone.sh ${RPM_BUILD_ROOT}/usr/sbin/
install -m 0644 posichain.service ${RPM_BUILD_ROOT}/etc/systemd/system/
install -m 0644 posichain-sysctl.conf ${RPM_BUILD_ROOT}/etc/sysctl.d/99-posichain.conf
install -m 0644 rclone.conf ${RPM_BUILD_ROOT}/etc/posichain/
install -m 0644 posichain.conf ${RPM_BUILD_ROOT}/etc/posichain/
exit 0

%post
%systemd_user_post %{name}.service
%sysctl_apply %{name}-sysctl.conf
exit 0

%preun
%systemd_user_preun %{name}.service
exit 0

%postun
%systemd_postun_with_restart %{name}.service
exit 0

%files
/usr/sbin/posichain
/usr/sbin/posichain-setup.sh
/usr/sbin/posichain-rclone.sh
/etc/sysctl.d/99-posichain.conf
/etc/systemd/system/posichain.service
/etc/posichain/posichain.conf
/etc/posichain/rclone.conf
/home/posichain/.config/rclone

%config(noreplace) /etc/posichain/posichain.conf
%config /etc/posichain/rclone.conf
%config /etc/sysctl.d/99-posichain.conf
%config /etc/systemd/system/posichain.service

%doc
%license



%changelog
* Mon Jun 27 2022 Danny <danny@position.exchange> 1.0.0
   - initial version of the posichain node program
