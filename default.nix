{
  pkgs ? import <nixpkgs> {}
, lib ? pkgs.lib
, makeScope ? lib.makeScope
, newScope ? pkgs.newScope
}:
let
  deps = import ./nix/deps.nix {};
in
pkgs.buildEnv {
  name ="shiva-hls";
  paths = with pkgs; [] ++ deps;
}