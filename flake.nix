{
  description = "A basic flake with a shell";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.flake-compat = {
    url = "github:edolstra/flake-compat";
    flake = false;
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in rec {
      packages = rec {
        shiva-hls = pkgs.callPackage ./nix/package.nix {};
      };

      apps.default = self.apps.${system}.shiva-hls;
      apps.shiva-hls = {
        type = "app";
        program = "${packages.shiva-hls}/bin/shiva-hls";
      };

      devShells.default = pkgs.mkShell {
        nativeBuildInputs = [
          pkgs.bashInteractive
          pkgs.yt-dlp
        ];
        buildInputs = [
          pkgs.go
        ];
      };
    })
    // {
      overlays.default = final: prev: {
        shiva-hls = final.callPackage ./nix/package.nix {};
      };
      nixosModules = {
        default = import ./nix/module.nix;
      };
    };
}
