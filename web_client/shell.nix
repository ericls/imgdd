# shell.nix
{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.nodejs_20
    pkgs.pnpm_10
  ];

  shellHook = ''
    echo "Node.js version $(node -v) is available"
  '';
}
