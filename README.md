# xkcd-wall

Get [xkcd](https://xkcd.com) on YOUR wallpaper today.

`xkcd-wall` is a simple tool to fetch xkcd comics, recolor them, and generate
wallpapers with a solid background.

It supports fetching the latest comic, a specific comic by number, or a random
comic.

The tool is configured via JSON and packaged with
[Nix](https://nixos.org/download).

# Usage

```bash
nix run github:sotormd/xkcd-wall -- -t <today|random|<number>> -c /path/to/config.json
```

<details>

<summary>Click to expand: examples</summary>

1. Get today's comic

   ```bash
   nix run github:sotormd/xkcd-wall -- -t today
   ```

2. Get a random comic

   ```bash
   nix run github:sotormd/xkcd-wall -- -t random
   ```

3. Get a specific comic

   ```bash
   nix run github:sotormd/xkcd-wall -- -t 1341
   ```

</details>

# Configuration

The tool looks for configuration in `$HOME/.config/xkcd-wall/config.json`.

If this does not exist, a default configuration file is created.

A configuration file path can also be passed using the `-c` flag.

<details>

<summary>Click to expand: default configuration values</summary>

```json
{
  "background-colors": ["#2e3440"],
  "foreground-colors": ["#81a1c1"],
  "dimensions": "1920x1080",
  "target": "/tmp/xkcd.png",
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
| `target`            | path to save the final image                      | `"/home/username/.local/share/backgrounds/xkcd.png"`      |
| `cache`             | cache directory                                   | `"/home/username/.cache"`                                 |
