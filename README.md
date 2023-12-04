# FDF - Wireframe - Fil de Fer

Wireframe visualizer

Two renderer are available: `png` to create an image file and [ebitengine](https://github.com/hajimehoshi)
which in turn can render as a native window or as a WASM website.

## Window mode

```sh
go run go.creack.net/fdf@latest
```

## WASM

### One liner

```sh
env -i HOME=${HOME} PATH=${PATH} go run github.com/hajimehoshi/wasmserve@latest go.creack.net/fdf@latest
```

### Details

Install `wasmer`:

```sh
go install github.com/hajimehoshi/wasmserve@latest
```

Clone this repo:

```sh
git clone https://github.com/creac/fdf
cd fdf
```

Run:

```sh
env -i HOME=${HOME} PATH=${PATH} wasmserve .
```

For development, `wasmer` exposes an endpoint to do live reload.

I recommend [reflex](https://github.com/cespare/reflex). 

Install:

```sh
go install github.com/cespare/reflex@latest
```

Then, with `wasmer` running:

```sh
reflex curl -v http://localhost:8080/_notify
```

## Controls

When running the `ebitengine` renderer, *wasm* or *window* mode, a few keyboard controls are available:

- up/down/left/right/shift left/shift right: Change x/y/z camera angles
- w/a/s/d: Move the image
- 1/2: Change the height
- 3/4: Change the scale

## Examples

### maps/42.fdf

![42.fdf](docs/42.png)

### maps/t1.fdf

![t1.fdf](docs/t1.png)
