{
  lib,
  buildGoModule,
  fetchFromGitHub,
}:
buildGoModule rec {
  pname = "docr";
  version = "0.0.2";

  src = fetchFromGitHub {
    owner = "notashelf";
    repo = "docr";
    rev = "v${version}";
    hash = "sha256-zjRGv95Z2l14J4UMhW3ohlDPrisuy1DjUkHMrmW79l8=";
  };

  vendorHash = "sha256-qcOBYuXoG5rFi3lF3SK5d/xOcKSH86++ly2BrNM0HAE=";

  ldflags = ["-s" "-w"];

  meta = with lib; {
    description = "Barebones static site generator in Go";
    homepage = "https://github.com/notashelf/docr";
    license = with licenses; [gpl3];
    maintainers = with maintainers; [notashelf];
  };
}
