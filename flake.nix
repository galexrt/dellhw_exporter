# Based upon https://github.com/the-nix-way/dev-templates
{
  description = "Basic flake for Go development";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      goMajorVersion = 1;
      goMinorVersion = 23; # Change this to update the whole stack

      lib = nixpkgs.lib;

      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };
      });
    in
    {
      overlays.default = final: prev: {
        go = final."go_${toString goMajorVersion}_${toString goMinorVersion}";
      };

      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          # Workaround CGO issue https://nixos.wiki/wiki/Go#Using_cgo_on_NixOS
          hardeningDisable = [ "fortify" ];

          packages = with pkgs; [
            # go and tools
            go
            # goimports, godoc, etc.
            gotools
            gofumpt
          ];
        };
      });
    };

}
