{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }@inputs: {
    overlay = final: prev: {
      gobuild = nixpkgs.legacyPackages.${prev.system}.buildGo118Module {
        pname = "gobuild";
        version = "0.0.1";
        src = builtins.path { path = ./.; };
        CGO_ENABLED = 0;
        vendorSha256 = "sha256-jfwJkPKpfKOkXafuKkB2rH04F4uRiR6u2bW+MEU7uSM=";
      };
    };
  } // flake-utils.lib.eachDefaultSystem (system:
    let pkgs = import nixpkgs {
      inherit system;
      overlays = [ self.overlay ];
    };
    in
    rec {
      devShell = pkgs.mkShell {
        buildInputs = with pkgs; [ entr go_1_18 ];
      };
      packages.gobuild = pkgs.gobuild;
      defaultPackage = packages.gobuild;
      apps.gobuild = pkgs.gobuild;
      defaultApp = apps.gobuild;
    });
}
