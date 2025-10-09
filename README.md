# SnowLinux

SnowLinux is a fully customizable immutable distribution implementing the
concepts described in
[Fitting Everything Together](https://0pointer.net/blog/fitting-everything-together.html).

Note that SnowLinux is still in development, and we don't provide any backwards
compatibility guarantees at all.

SnowLinux is inspired by and derived from [ParticleOS](https://github.com/systemd/particleos)

The SnowLinux image is built using [mkosi](https://github.com/systemd/mkosi).
You will need to install the current main branch of mkosi to build current
SnowLinux images.

## Signing keys

SnowLinux images are signed for Secure Boot with the user's keys. To generate a new key,
run `mkosi genkey`. The key must be stored safely, it will be required to sign updates.

The key can be stored in a smartcard. Then you have to set the key in `mkosi.local.conf`:

```
[Validation]
SecureBootKey=pkcs11:object=Private key 1;type=private
SecureBootKeySource=provider:pkcs11
SignExpectedPcrKey=pkcs11:object=Private key 1;type=private
SignExpectedPcrKeySource=provider:pkcs11
VerityKey=pkcs11:object=Private key 1;type=private
VerityKeySource=provider:pkcs11
```

## Installation

Before installing SnowLinux, make sure that Secure Boot is in setup mode on the
target system. The Secure Boot mode can be configured in the UEFI firmware
interface of the target system. If there's an existing Linux installation on the
target system already, run `systemctl reboot --firmware-setup` to reboot into
the UEFI firmware interface. At the same time, make sure the UEFI firmware
interface is password protected so an attacker cannot just disable Secure Boot
again.

To install SnowLinux with a USB drive, first build the image on an existing
Linux system as described above. Then, burn it to the USB drive with
`mkosi burn /dev/<usb>`. Once burned to the USB drive, plug the USB drive into
the system onto which you'd like to install SnowLinux and boot into the USB
drive via the firmware. Then, boot into the "Installer" UKI profile. When you
end up in the root shell, run
`snowctl install`
to install SnowLinux to the system's drive. Finally, reboot into the target
drive (not the USB) and the regular profile (not the installer one) to complete
the installation.

## LUKS recovery key

systemd doesn't support adding a recovery key to a partition enrolled with a token
only (tpm/fido2). It is possible to use cryptenroll to add a recovery password
to the root partition: `cryptsetup luksAddKey --token-type systemd-tpm2 /dev/<id>`

## Firmwares

Only firmwares that are dependencies of a kernel module are included, but some
modules don't declare their dependencies properly. Dependencies of a module can be
found with `modinfo`. If you experience missing firmwares, you should report
this to the module maintainer. `FirmwareInclude=` can be added in `mkosi.local.conf`
to include the firmware regardless of whether a module depends on it.

## Software and Tool Installation

- snap preinstalled
- flatpak preinstalled
- docker preinstalled
- brew (user installable)
- systemd-sysext extensions (not started)
