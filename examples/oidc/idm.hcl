database_url = "file:/var/idm/idm.db"

registration = "disabled"

server {
    # We're not running under HTTPS in this example
    # so cookies must not have the "secure" attribute
    # set
    secure_cookies = false

    domain = "example.intern"
    public_listener = ":8080"
    admin_listener = ":8081"
    allowed_origins = [
        "http://example.intern",
        "http://*.example.intern"
    ]
    allowed_redirects = [
        ".example.intern",
        "example.intern"
    ]
}

jwt {
    secret = "some-random-secret"
}

ui {
    site_name = "Example Inc"
    public_url = "http://account.example.intern"
}

forward_auth {
    allow_cors_preflight = true
}

policies {
    # A default policy to allow all requests to upstream services
    # the there's an authenticated user (i.e. input.subject is available)
    policy "default" {
        content = <<EOT
        package cisidm.forward_auth

        import rego.v1

        default allow := false

        allow if {
            input.subject
        }
        EOT
    }
}