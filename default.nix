# nix-build default.nix --arg pkgs 'import <nixpkgs> { }'
# xkcd-wall = import sources.xkcd-wall { inherit pkgs; };

{ pkgs }:

pkgs.callPackage ./package.nix { }
