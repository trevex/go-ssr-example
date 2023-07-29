# Go SSR Example

## Prerequisites

```
make
golang
yarn
nodejs
```
(or use provided nix shell with direnv)

## Usage

```
make deps
make run
```

## TODOs

- [ ] In development mode make sure to force no caching by setting appropriate headers
- [ ] Suffix scripts and styles for production builds
- [ ] Do not minify javascript in development mode
- [ ] Watch mode is prepared but not fully implemented, should be used to avoid long rebuild cycles in development
