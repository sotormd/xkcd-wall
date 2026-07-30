{
  description = "create wallpapers from xkcd comics";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    inputs:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      forAllSystems =
        apply: inputs.nixpkgs.lib.genAttrs systems (system: apply inputs.nixpkgs.legacyPackages.${system});
    in
    {
      packages = forAllSystems (pkgs: {
        default = import ./default.nix { inherit pkgs; };
      });
    };
}
