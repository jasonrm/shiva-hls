{
  pkgs ? import <nixpkgs> {}
}:
with pkgs; [
    gcc
    glfw
    go
    imagemagick
    jq
    libpcap
    python3
    pkg-config
    xlibs.libXext.dev
    xorg.libX11
    xorg.libXcursor
    xorg.libXi
    xorg.libXinerama
    xorg.libXrandr
    xorg.libXxf86vm
    xorg.xinput
]