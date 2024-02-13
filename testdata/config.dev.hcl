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

field "object" "notification" {
    visibility = "self"
    writeable = true
    display_name = "Benachrichtigungen"

    property "object" "roster" {
        display_name = "Dienstplan"

        property "bool" "sms" {
            display_name = "SMS"
        }

        property "bool" "mail" {
            display_name = "E-Mail"
        }
    }

    property "string" "offtime" {
        display_name = "Urlaubsanträge"

        value "sms" {
            display_name = "SMS"
        }

        value "mail" {
            display_name = "E-Mail"
        }

        value "both" {
            display_name = "SMS + E-Mail"
        }
    }
}

field "string" "string" {
    visibility = "self"
    writeable = true
    description = "A custom description"
    display_name = "String Value"
}

field "date" "date" {
    visibility = "self"
    writeable = true
    description = "A custom description"
    display_name = "Date Value"
}

field "number" "number" {
    visibility = "self"
    writeable = true
    description = "A custom description"
    display_name = "Number Value"
}

field "bool" "bool" {
    visibility = "self"
    writeable = true
    description = "A custom description"
    display_name = "Boolean Value"
}

field "object" "object" {
    visibility = "self"
    writeable = true
    description = "A custom description"
    display_name = "Object Value"

    property "string" "string" {
        description = "A string property"
        display_name = "String Property"
    }

    property "bool" "bool" {
        visibility = "self"
        writeable = true
        description = "A custom description"
        display_name = "Boolean Value"
    }

    property "object" "object" {
        visibility = "self"
        writeable = true
        description = "A custom description"
        display_name = "Sub Object Value"

        property "string" "string" {
            description = "A string property"
            display_name = "String Property"
        }

        property "bool" "bool" {
            writeable = true
            description = "A custom description"
            display_name = "Boolean Value"
        }
    
        property "object" "object" {
            writeable = false
            description = "A custom description"
            display_name = "Sub-Sub Object Value"

            property "string" "string" {
                description = "A string property"
                display_name = "String Property"
            }

            property "bool" "bool" {
                writeable = false
                description = "A custom description"
                display_name = "Boolean Value"
            }
        }
    }
}

field "list" "list" {
    visibility = "self"
    writeable = true
    description = "A custom description"
    display_name = "List Value"

    element_type "string" "" {}
}
