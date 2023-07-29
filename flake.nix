{
  description = "go-ssr-example";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    let
      overlays = [
        (final: prev: {
          # currently none
        })
      ];
    in
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system overlays; };
      in
      rec {
        devShell = pkgs.mkShell rec {
          name = "go-ssr-example";

          buildInputs = with pkgs; [
            go
            golangci-lint
            delve
            yarn
            nodejs
            gnumake
            tmux
          ];
        };
      }
    );
}
