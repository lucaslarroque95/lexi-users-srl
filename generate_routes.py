#!/usr/bin/env python3
"""Discovers routes from every endpoint folder's own endpoint.json (method +
path) and writes infra/routes.auto.tfvars.json. Adding an endpoint means
adding a folder with a main.go + endpoint.json — nobody hand-edits a route
table in Terraform.

Run manually while developing, or automatically by
.github/workflows/deploy.yml on every push to main.
"""
import json
import os

HERE = os.path.dirname(os.path.abspath(__file__))


def main() -> None:
    routes = {}
    for entry in sorted(os.listdir(HERE)):
        endpoint_dir = os.path.join(HERE, entry)
        config_path = os.path.join(endpoint_dir, "endpoint.json")
        if not os.path.isfile(config_path):
            continue  # not an endpoint folder (internal/, dist/, infra/, ...)

        with open(config_path) as f:
            config = json.load(f)

        for field in ("method", "path"):
            if field not in config:
                raise ValueError(f"{config_path} is missing required field '{field}'")

        routes[entry] = {"method": config["method"], "path": config["path"]}

    out_path = os.path.join(HERE, "infra", "routes.auto.tfvars.json")
    with open(out_path, "w") as f:
        json.dump({"users_routes": routes}, f, indent=2)
        f.write("\n")

    print(f"{len(routes)} endpoints -> {out_path}")


if __name__ == "__main__":
    main()
