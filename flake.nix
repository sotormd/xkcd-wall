{
  description = "create wallpapers from xkcd comics";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    wallpapers.url = "github:sotormd/wallpapers";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in
    {
      packages.x86_64-linux.default = pkgs.buildGoModule {
        pname = "xkcd-wall";
        version = "0.1.0";
        src = ./.;
        subPackages = [ "./cmd/xkcd-wall" ];
        vendorHash = null;
      };

      apps.x86_64-linux.default = {
        type = "app";
        program = "${self.packages.x86_64-linux.default}/bin/xkcd-wall";
      };

      nixosModules.xkcd =
        {
          config,
          inputs,
          lib,
          ...
        }:
        let
          cfg = config.xkcd;
        in
        {
          options = {
            xkcd = {
              enable = lib.mkOption {
                type = lib.types.bool;
                default = false;
                description = "Install the xkcd-wall timer, service and package";
                example = true;
              };

              package = lib.mkOption {
                type = lib.types.package;
                default = self.packages.x86_64-linux.default;
                description = "xkcd-wall package";
              };

              type = lib.mkOption {
                type = lib.types.enum [
                  "today"
                  "random"
                ];
                default = "today";
                description = "Which comic to fetch";
                example = "random";
              };

              background-colors = lib.mkOption {
                type = lib.types.listOf lib.types.str;
                default = [ "#2e3440" ];
                description = "Background colors to choose from";
                example = [
                  "#2e3440"
                  "#3b4252"
                ];
              };

              foreground-colors = lib.mkOption {
                type = lib.types.listOf lib.types.str;
                default = [ "#81a1c1" ];
                description = "Foreground colors to choose from";
                example = [
                  "#bf616a"
                  "#d08770"
                  "#ebcb8b"
                ];
              };

              dimensions = lib.mkOption {
                type = lib.types.str;
                default = "1920x1080";
                description = "Dimensions of output image";
                example = "1920x1200";
              };

              target = lib.mkOption {
                type = lib.types.str;
                default = "/tmp/xkcd.png";
                description = "Path to save the final image";
                example = "/home/user/Pictures/xkcd.png";
              };

              cache = lib.mkOption {
                type = lib.types.str;
                default = "/tmp/xkcd-wall-cache";
                description = "Path to use as cache directory";
                example = "/home/user/.cache/xkcd-wall";
              };

              interval = lib.mkOption {
                type = lib.types.str;
                default = "";
                description = "How often to run the systemd timer";
                example = "1h";
              };

              fallback = lib.mkOption {
                type = lib.types.path;
                default = inputs.wallpapers.lib.wallpapers.nord.nixos;
                description = "Fallback image";
              };
            };
          };
          config = lib.mkIf cfg.enable (
            let
              configuration = pkgs.writeText "config.json" (
                builtins.toJSON {
                  inherit (cfg)
                    background-colors
                    foreground-colors
                    dimensions
                    cache
                    ;
                }
              );

              xkcd-wall-package = pkgs.writeShellScriptBin "xkcd-wall" ''
                TARGET="${cfg.target}"
                TMP="/tmp/xkcd-tmp.png"

                ${cfg.package}/bin/xkcd-wall -c ${configuration} -t ${cfg.type} $TMP || \
                ln -sf ${cfg.fallback} $TMP

                rm $TARGET
                cp $TMP $TARGET
                unlink $TMP
              '';
            in
            {
              systemd.services.xkcd-wall = {
                description = "Fetch XKCD comic and generate wallpaper";
                wantedBy = [ "multi-user.target" ];
                serviceConfig = {
                  Type = "oneshot";
                  ExecStart = "${xkcd-wall-package}/bin/xkcd-wall";
                  Restart = "on-failure";
                };
              };

              systemd.timers.xkcd-wall = {
                description = "Run xkcd-wall periodically";
                wantedBy = [ "timers.target" ];
                timerConfig = lib.mkIf (cfg.interval != "") {
                  OnUnitActiveSec = cfg.interval;
                  Persistent = true;
                };
              };
            }
          );
        };

    };
}
