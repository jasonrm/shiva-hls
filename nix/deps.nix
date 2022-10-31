{
  pkgs ? import <nixpkgs> {}
}:
with pkgs; [
    go
    youtube-dl
]
