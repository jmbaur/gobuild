{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-21.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }@inputs: {
    overlay = final: prev: {
      gobuild = nixpkgs.legacyPackages.${prev.system}.buildGo117Module {
        pname = "gobuild";
        version = "0.0.1";
        src = builtins.path { path = ./.; };
        CGO_ENABLED = 0;
        vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
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
        buildInputs = with pkgs; [ entr go_1_17 ];
      };
      packages.gobuild = pkgs.gobuild;
      defaultPackage = packages.gobuild;
      apps.gobuild = pkgs.gobuild;
      defaultApp = apps.gobuild;
    });
}
