# Configures where the SQLite3 database should be stored.
database_url = "file:/data/idm.db"

# Disable self-registration of users. This means that any user account must be
# created by an administrator using the `idmctl` cli utility.
registration = "disabled"

# This configures some defaults for the built-in server.
server {
    # We're running on HTTP for this example
    # so cookies must not be set to "secure"
    secure_cookies = false

    domain = "example.intern"

    allowed_origins = [
        "http://example.intern",
        "http://*.example.intern",
    ]

    trusted_networks = [
        "traefik"
    ]

    allowed_redirects = [
        "example.intern",
        ".example.intern"
    ]
}

jwt {
    secret = "some-secure-random-string"
}

ui {
    site_name = "Example Inc"
    public_url = "http://account.example.intern"
}

forward_auth {
    # Allways allow CORS preflight requests.
    # If this would be set to false you would likely need to
    # account for CORS preflight requests in your policies.
    allow_cors_preflight = true
}

policies {
    debug = false

    policy "default" {
        content = <<EOT
        package cisidm.forward_auth

        import rego.v1

        default allow := false

        allow if {
            # input.subject is only set if the request is authenticated
            input.subject
        }
        EOT
    }
}
