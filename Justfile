# Default recipe
default: build

build: snowctl
  mkosi build

force: snowctl
  mkosi build -f

build-compress: snowctl
  mkosi build --compress-output=yes

bump:
  mkosi bump

clean:
  mkosi clean -ff


launch:
  ./scripts/launch-incus.sh

snowctl:
  cd snowctl && go build -o snowctl .

main:
  git checkout main
  git pull --recurse-submodules

dev branch:
  git checkout -b {{branch}}
