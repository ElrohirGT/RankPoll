{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "RankPoll";
  version = "1.0.0";
  src = ./.;
  vendorHash = null;
}
