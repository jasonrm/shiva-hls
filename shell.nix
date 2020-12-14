{
  pkgs ? import <nixpkgs> { },
  mkShell ? pkgs.mkShell,
}:
let
  callPackage = pkgs.lib.callPackageWith pkgs;

  # this-pkg = callPackage ./default.nix { };
  deps = import ./nix/deps.nix {};
in
mkShell {
  buildInputs = [ ] ++ deps;

  shellHook = ''
    export VIRTUAL_ENV=nix-$(basename $(pwd))
  '';
}