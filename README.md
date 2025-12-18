# xkcd-wall

Get [xkcd](https://xkcd.com) on YOUR wallpaper today.

`xkcd-wall` is a simple tool to fetch xkcd comics, recolor them, and generate
wallpapers with a solid background.

It supports fetching the latest comic, a specific comic by number, or a random
comic.

The tool is configured via JSON and packaged with
[Nix](https://nixos.org/download).

A NixOS module is also included, which sets up a systemd timer.

```bash
github:sotormd/xkcd-wall
├───apps
│   └───x86_64-linux
│       └───default: app
├───nixosModules
│   └───xkcd: NixOS module
└───packages
    └───x86_64-linux
        └───default: package 'xkcd-wall-0.1.0'
```

# Usage

1. Get today's comic

   ```bash
   nix run github:sotormd/xkcd-wall -- -t today ./output.png
   ```

2. Get a random comic

   ```bash
   nix run github:sotormd/xkcd-wall -- -t random ./output.png
   ```

3. Get a specific comic

   ```bash
   nix run github:sotormd/xkcd-wall -- -t 1341 ./output.png
   ```

# Configuration

The tool looks for configuration in `$HOME/.config/xkcd-wall/config.json`.

If this does not exist, a default configuration file is created.

A configuration file path can also be passed using the `-c` flag.

```bash
nix run github:sotormd/xkcd-wall -- -t today -c /tmp/config.json ./output.png
```

<details>

<summary>Click to expand: default configuration values</summary>

```json
{
  "background-colors": ["#2e3440"],
  "foreground-colors": ["#d8dee9"],
  "dimensions": "1920x1080",
  "cache": "/tmp/xkcd-wall-cache"
}
```

</details>

The configuration values are explained here:

| value               | explanation                                       | example                                                   |
| ------------------- | ------------------------------------------------- | --------------------------------------------------------- |
| `background-colors` | list of background colors to randomly choose from | `["#2e3440", "#3b4252"]`                                  |
| `foreground-colors` | list of foreground colors to randomly choose from | `["#bf616a", "#d08770", "#ebcb8b", "#a3be8c", "#b48ead"]` |
| `dimensions`        | output image dimensions                           | `"1920x1200"`                                             |
| `cache`             | cache directory                                   | `"/home/username/.cache"`                                 |
