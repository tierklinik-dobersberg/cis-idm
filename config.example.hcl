domain = "dobersberg.dev"
jwt_secret = "some-secret"
database_url = "file:/data/idm.db"
trusted_networks = [ "traefik" ]

logo_url = "/files/logo.png"

role "dummy" {
    name = "dummy role"
    description = "Nothing useful"
    permissions = [
        "roster",
        "idm:users:write"
    ]
}

allowed_origins = [
    "https://example.dev",
    "https://*.example.dev",
]

public_url = "https://account.example.dev"

allowed_redirects = [
    "*.example.dev",
]

site_name = "Example"

permissions = [
    "roster:write:create",
    "roster:write:approve",
    "roster:read"
]