# Decision log

## Base 
For stability we base on debian/trixie. This means a stable suite of tools for the core OS that will receive important security updates over time.

We require a few advanced features from `systemd` that aren't available in trixie yet, so we use the `OBS` profile to pull in systemd from basically "main" of systemd development. [link](https://build.opensuse.org/project/show/system:systemd)