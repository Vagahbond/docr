{
  stdenvNoCC,
  nodejs,
  extraFlags ? [],
  lib,
  makeWrapper,
  ...
}:
stdenvNoCC.mkDerivation {
  pname = "servejs";
  version = "0.0.1";

  src = ../.;

  buildInputs = [nodejs makeWrapper];

  preInstall = ''
    mkdir -p $out/lib/node_modules/
  '';

  postInstall = ''
    cp -rvf serve.js $out/lib/node_modules
  '';

  buildPhase = ''
    ${nodejs}/bin/node --version
    ${makeWrapper} ${nodejs}/bin/node $out/bin/servejs \
      --add-flags $out/lib/node_modules/servejs \
      --chdir $out/lib/node_modules/servejs \
      ${lib.concatStringsSep " " (map (flag: "--add-flags ${flag}") extraFlags)}
  '';

  meta = with lib; {
    mainProgram = "serve.js";
    description = "Minimalistic NodeJS application to serve static html directories";
    homepage = "https://github.com/notashelf/docr";
    license = licenses.gpl3;
    maintainers = with maintainers; [notashelf];
  };
}
