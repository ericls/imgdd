{
  description = "imgdd dev";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_25
            air

            nodejs_24
            corepack_24

            git
            lefthook
          ];

          # Use the Go from this shell; fail instead of silently downloading a newer
          # toolchain (and bumping go.mod) when a dependency requires one.
          GOTOOLCHAIN = "local";

          shellHook = ''
            lefthook install
            echo "Go version:   $(go version)"
            echo "Node version: $(node -v)"
            echo ""
            echo "Notes: "
            echo "For backend: Expects docker and docker-compose for local dev services and tests"
            echo "For web_client: Use corepack to install pnpm"
          '';
        };
      }
    );
}
