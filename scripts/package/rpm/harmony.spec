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
./harmony --version
exit 0

%pre
getent group harmony >/dev/null || groupadd -r harmony
getent passwd harmony >/dev/null || \
   useradd -r -g harmony -d /home/harmony -m -s /sbin/nologin \
   -c "Posichain validator node account" harmony
mkdir -p /home/harmony/.hmy/blskeys
mkdir -p /home/harmony/.config/rclone
chown -R harmony.harmony /home/harmony
exit 0


%install
install -m 0755 -d ${RPM_BUILD_ROOT}/usr/sbin ${RPM_BUILD_ROOT}/etc/systemd/system ${RPM_BUILD_ROOT}/etc/sysctl.d ${RPM_BUILD_ROOT}/etc/harmony
install -m 0755 -d ${RPM_BUILD_ROOT}/home/harmony/.config/rclone
install -m 0755 harmony ${RPM_BUILD_ROOT}/usr/sbin/
install -m 0755 harmony-setup.sh ${RPM_BUILD_ROOT}/usr/sbin/
install -m 0755 harmony-rclone.sh ${RPM_BUILD_ROOT}/usr/sbin/
install -m 0644 harmony.service ${RPM_BUILD_ROOT}/etc/systemd/system/
install -m 0644 harmony-sysctl.conf ${RPM_BUILD_ROOT}/etc/sysctl.d/99-harmony.conf
install -m 0644 rclone.conf ${RPM_BUILD_ROOT}/etc/harmony/
install -m 0644 harmony.conf ${RPM_BUILD_ROOT}/etc/harmony/
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
/usr/sbin/harmony
/usr/sbin/harmony-setup.sh
/usr/sbin/harmony-rclone.sh
/etc/sysctl.d/99-harmony.conf
/etc/systemd/system/harmony.service
/etc/harmony/harmony.conf
/etc/harmony/rclone.conf
/home/harmony/.config/rclone

%config(noreplace) /etc/harmony/harmony.conf
%config /etc/harmony/rclone.conf
%config /etc/sysctl.d/99-harmony.conf
%config /etc/systemd/system/harmony.service

%doc
%license



%changelog
* Mon Jun 27 2022 Danny <danny@position.exchange> 1.0.0
   - initial version of the posichain node program
