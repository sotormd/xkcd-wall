{
  description = "create wallpapers from xkcd comics";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
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
    };
}
