# Default recipe
default: build

build:
  mkdir -p mkosi.extra.common/boot/EFI/
  openssl x509 -in mkosi.crt -out mkosi.extra.common/boot/EFI/mkosi.der -outform DER

  mkosi build
  mkosi --profile=obs,live --image-id=SnowLive build

force:
  mkosi build -ff
  mkosi --profile=obs,live --image-id=SnowLive build

build-compress: snowctl
  mkosi build --compress-output=yes
  mkosi --compress-output=yes --profile=obs,live --image-id=SnowLive build

bump:
  mkosi bump

clean:
  mkosi clean -ff


launch:
  ./scripts/launch-incus.sh


main:
  git checkout main
  git pull --recurse-submodules

dev branch:
  git checkout -b {{branch}}
