{
  description = "A Nix-flake-based go development environment";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

  outputs = {
    self,
    nixpkgs,
  }: let
    system = "aarch64-darwin";
    pkgs = import nixpkgs {inherit system;};
  in {
    devShells.${system} = {
      default = with pkgs;
        mkShell {
          buildInputs = [
            go
            go-tools
            gopls
          ];
          # Shell hook to symlink discord-ipc as every discord presence solution looks for TMPDIR on macos (Manually created symlink to `projects` folder only works for me)
          shellHook = ''
            ln -s $HOME/Desktop/projects/discord-ipc-0 $TMPDIR
          '';
        };
    };
  };
}
