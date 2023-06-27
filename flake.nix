{
  description = "Barebones static Website generator from flavored Markdown";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [];
      systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin"];
      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: {
        # configure perSystem's instance of nixpkgs
        _module.args.pkgs = import inputs.nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };

        devShells.default = pkgs.mkShell {
          name = "nyx";
          packages = with pkgs; [
            nil # nix ls
            alejandra # formatter
            git # flakes require git, and so do I
            glow # markdown viewer
            nodePackages.serve # serve static content over http
            go # go language for compiling and running
          ];
        };

        # provide the formatter for nix fmt
        formatter = pkgs.alejandra;

        packages = {
          servejs = pkgs.callPackage ./nix/serve.nix {};
        };
      };
      flake = {
        # The usual flake attributes can be defined here, including system-
        # agnostic ones like nixosModule and system-enumerating ones, although
        # those are more easily expressed in perSystem.
      };
    };
}
