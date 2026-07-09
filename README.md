<div align="center">
  <a href="https://github.com/tosterabgx/marten">
    <img src="./frontend/src/images/dark_branding.png" alt="Marten" width=500>
  </a>
  <br />
  <br />
  <p align="center">
    A self-hosted reverse tunnel, built in Go and inspired by <a href="https://github.com/ekzhang/bore">bore</a>.
  </p>
  <br />
</div>

Expose a local port to the internet without port forwarding, NAT
config, or a static IP - `marten` does the rest.

```
marten tcp 3000
```

## Install

Linux/MacOS:
```
curl -fsSL https://usemarten.tech/install.sh | sh
```

Windows:
```
irm https://usemarten.tech/install.ps1 | iex
```

or build from source:

```
go install github.com/tosterabgx/marten/cmd/marten@latest
```
