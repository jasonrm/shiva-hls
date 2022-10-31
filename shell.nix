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
    export XDG_DATA_HOME=~/.local/share
    mkdir -p $XDG_DATA_HOME
  '';
}
