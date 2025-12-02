# shell.nix
{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.nodejs_20
    pkgs.corepack_20
  ];

  shellHook = ''
    echo "Node.js version $(node -v) is available. Use corepack to install pnpm"
  '';
}
