{
  config,
  pkgs,
  lib,
  ...
}:
with lib; let
  cfg = config.services.shiva-hls;

  shiva-hls-bulk = pkgs.writeShellScriptBin "shiva-hls-bulk" ''
    for CHANNEL in $@; do
      echo "Checking $CHANNEL"
      ${pkgs.shiva-hls}/bin/shiva-hls --out ${cfg.downloadDirectory} ${cfg.additionalArgs} $CHANNEL
    done
  '';
in {
  options = with lib.types; {
    services.shiva-hls = {
      enable = mkEnableOption "Enable Twitch HLS stream downloader";
      channels = mkOption {
        type = listOf str;
        default = [];
      };
      additionalArgs = mkOption {
        type = str;
        default = "";
      };
      twitchClientId = mkOption {
        type = str;
      };
      twitchClientSecret = mkOption {
        type = str;
      };
      downloadDirectory = mkOption {
        type = str;
      };
      frequencySec = mkOption {
        type = str;
      };
    };
  };

  config = {
    systemd.services.shiva-hls = mkIf cfg.enable {
      description = "Twitch HLS stream downloader";
      wants = ["network-online.target"];
      after = ["network.target" "network-online.target"];
      path = [pkgs.yt-dlp];
      # TODO: Use EnvironmentFile instead
      environment = {
        TWITCH_CLIENT_ID = cfg.twitchClientId;
        TWITCH_CLIENT_SECRET = cfg.twitchClientSecret;
      };
      serviceConfig = {
        ExecStart = "${shiva-hls-bulk}/bin/shiva-hls-bulk ${concatStringsSep " " cfg.channels}";
      };
    };
    systemd.timers.shiva-hls = {
      description = "Twitch Downloader Timer";
      wantedBy = ["timers.target"];
      timerConfig = {
        OnBootSec = "${cfg.frequencySec}";
        OnUnitActiveSec = "${cfg.frequencySec}";
      };
    };
  };
}
