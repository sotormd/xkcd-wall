# xkcd-wall

Create wallpapers from [xkcd](https://xkcd.com) comics.

`xkcd-wall` is a simple tool to fetch xkcd comics, recolor them, and generate
wallpapers with a solid background.

`xkcd-wall` is packaged with [Nix](https://nixos.org/download) for
`x86_64-linux` and `aarch64-linux`.

# Usage

Run directly:

```bash
nix run github:sotormd/xkcd-wall -- ./output.png
```

It is also possible to install `xkcd-wall` in a Nix Profile with `nix profile`.

Alternatively, this flake can be added as an input and the following packages
can be consumed:

1. `packages.x86_64-linux.default`
1. `packages.aarch64-linux.default`

For usage without flakes, see [`./default.nix`](./default.nix).

# Options

```bash
xkcd-wall [options] <output>
```

The following command-line arguments are supported:

| option | description                                          | default     |
| ------ | ---------------------------------------------------- | ----------- |
| `-t`   | comic to fetch: `today`, `random`, or a comic number | `today`     |
| `-d`   | output image dimensions                              | `1920x1200` |
| `-b`   | background color in hex format                       | `2e3440`    |
| `-f`   | foreground color in hex format                       | `d8dee9`    |

Examples:

1. Get a random comic, `2000x2000`, white background and black foreground

   ```bash
   xkcd-wall \
     -t random \
     -d 2000x2000 \
     -b ffffff \
     -f 000000 \
     ./output.png
   ```

2. Get comic #1341, `1920x1200`, `2e3440` background and `ebcb8b` foreground

   ```bash
   xkcd-wall \
     -t 1341 \
     -d 1920x1200 \
     -b 2e3440 \
     -f ebcb8b \
     ./output.png
   ```
