{ lib, buildGoModule }:

buildGoModule {
  pname = "xkcd-wall";
  version = "0";
  src = ./.;
  subPackages = [ "./cmd/xkcd-wall" ];
  vendorHash = null;

  meta = {
    description = "Create wallpapers from xkcd comics";
    homepage = "https://github.com/sotormd/xkcd-wall";
    license = lib.licenses.gpl3Only;
    mainProgram = "xkcd-wall";
  };
}
