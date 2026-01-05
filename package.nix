{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "RankPoll";
  version = "1.0.0";
  src = ./.;
  vendorHash = "sha256-mGKxBRU5TPgdmiSx0DHEd0Ys8gsVD/YdBfbDdSVpC3U=";
}
