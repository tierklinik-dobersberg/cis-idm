log_level = "debug"

database_url = "file:/data/idm.db"

enable_dynamic_roles = true

registration = "public"

server {
    secure_cookies = true
    domain = "dobersberg.dev"
    public_listener = ":8080"
    admin_listener = ":8081"
    allowed_origins = [
        "https://dobersberg.dev",
        "https://*.dobersberg.dev"
    ]

    static_files = "http://ui:4200"

    trusted_networks = ["traefik"]

    allowed_redirects = [
        "*.dobersberg.dev"
    ]
}

jwt {
    secret = "some-random-secret"
}

ui {
    site_name = "Example Inc"
    public_url = "https://account.dobersberg.dev"
}

forward_auth {
    default = "deny"

    allow_cors_preflight = "true"
}

policies {
    policy "default" {
        content = <<EOT
        package cisidm.forward_auth

        import rego.v1

        allow if {
            input.subject
        }
        EOT
    }
}

