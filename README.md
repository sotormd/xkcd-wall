# xkcd-wall

Create wallpapers from [xkcd](https://xkcd.com) comics.

`xkcd-wall` is a simple tool to fetch xkcd comics, recolor them, and generate
wallpapers with a solid background.

`xkcd-wall` is packaged with [Nix](https://nixos.org/download) for
`x86_64-linux` and `aarch64-linux`.

# Usage

```bash
nix run github:sotormd/xkcd-wall -- ./output.png
```

# Options

```bash
nix run github:sotormd/xkcd-wall -- [options] <output>
```

The following command-line arguments are supported:

| option | description                                          | default     |
| ------ | ---------------------------------------------------- | ----------- |
| `-t`   | comic to fetch: `today`, `random`, or a comic number | `today`     |
| `-d`   | output image dimensions                              | `1920x1080` |
| `-b`   | background color in hex format                       | `2e3440`    |
| `-f`   | foreground color in hex format                       | `d8dee9`    |

Examples:

1. Get a random comic, `2000x2000`, white background and black foreground

```bash
nix run github:sotormd/xkcd-wall -- \
  -t random \
  -d 2000x2000 \
  -b ffffff \
  -f 000000 \
  ./output.png
```

2. Get comic #1341, `1920x1200`, `2e3440` background and `ebcb8b` background

```bash
nix run github:sotormd/xkcd-wall -- \
  -t 1341 \
  -d 1920x1200 \
  -b 2e3440 \
  -f ebcb8b \
  ./output.png
```
