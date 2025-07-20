{

  description = "Go project template";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flakelight.url = "github:nix-community/flakelight";
  };

  outputs = { flakelight, nixpkgs, ... }:
    flakelight ./. {

      inputs.nixpkgs = nixpkgs;

      devShell.packages = pkgs:
        with pkgs; [

          coreutils

          lefthook

          go
          golangci-lint
          gofumpt

        ];
    };

}
