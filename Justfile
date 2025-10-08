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
  ./scripts/launch.sh

snowctl:
  cd snowctl && go build -o snowctl .
  mkdir -p mkosi.extra/usr/local/bin/
  cp snowctl/snowctl mkosi.extra/usr/local/bin/
