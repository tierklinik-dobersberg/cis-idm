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


field "object" "integrations" {
    display_name = "Integrations"
    description = "Configure integrations with external services"

    property "object" "discord" {
        display_name = "Discord"
        description = "Link your discord account"

        property "string" "discord" {
            display_name = "Discord ID"
            description = "Enter your discord user ID to link your account"
        }

        property "bool" "notify" {
            display_name = "Notify on Discord"
            description = "Do you want to be notified on discord"
        }
    }

    property "object" "reddit" {
        display_name = "Reddit"
        description = "Link your Reddit account"

        property "string" "discord" {
            display_name = "Reddit User-Handle"
            description = "Enter your Reddit username to link your account"
        }

        property "bool" "notify" {
            display_name = "Notify on Reddit"
            description = "Do you want to be notified on Reddit"
        }
    }
}

field "object" "notificationSettings" {
    visibility = "self"
    writeable = true
    display_name = "Notification Settings"
    description = "Manage your notification settings and preferences"

    property "string" "newsletter" {
        display_name = "Newsletter"
        description = "If and how you want to receive our weekly newsletter"

        value "never" {
            display_name = "Never"
        }

        value "email" {
            display_name = "E-Mail"
        }

        value "sms" {
            display_name = "SMS"
        }

        value "both" {
            display_name = "E-Mail + SMS"
        }
    }

    property "string" "comments" {
        display_name = "Comments & Replies"
        description = "If and how you want to receive notifications about comments and replies"

        value "never" {
            display_name = "Never"
        }

        value "email" {
            display_name = "E-Mail"
        }

        value "sms" {
            display_name = "SMS"
        }

        value "both" {
            display_name = "E-Mail + SMS"
        }
    }
}