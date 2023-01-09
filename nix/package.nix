{
  lib,
  buildGoModule,
  yt-dlp,
}:
buildGoModule rec {
  pname = "shiva-hls";
  version = "1.0.0";

  src = ../.;

  vendorSha256 = "sha256-owmyrctnL3p0uKjRhuOn0b/Wd59ZTUNrYwCufyLpMsc=";
  # vendorSha256 = lib.fakeHash;

  buildInputs = [yt-dlp];

  meta = with lib; {
    description = "Download Twitch HLS streams";
    license = licenses.mit;
    homepage = "https://github.com/jasonrm/shiva-hls";
    maintainer = ["jason@mcneil.dev"];
  };
}
