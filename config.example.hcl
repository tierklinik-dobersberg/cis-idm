# General Configuration Values
# -------------------------------------------------------------------------------------------------------


# Configures the logging level to use. Valid values for this setting are:
#  - debug
#  - info
#  - warn
#  - error
log_level = "info"

# The URL for the SQLite3 database that stores any user related information.
database_url = "file:/data/idm.db"

# Whehter or not roles can be created, modified or deleted via the tkd.idm.v1.RolesService API.
# If unset, the default value for enable_dynamic_roles depends on the presence of one or more
# role blocks (see below).
#
# If at least one role is defined in the configuration file, enable_dynamic_roles defaults to false
# since cisidm expects a static role configuration. 
# If no roles are configures, enable_dynamic_roles defaults to true.
#
# Note that it's possible to disable roles altogheter by not specifying a role block and by explicitly
# setting enable_dynamic_roles to false. Counterwise, it's possible to have some static roles configured
# (which are delete/modify protected) while still allowing roles to be created dynamically.
enable_dynamic_roles = false

# Configures the user registration mode:
#  - disabled (default): users are not allowed to register themself. In this mode the administrator
#                        must create user accounts for all users. User invitations are not possible!
#  - token: users must provide a token for registration. In this mode, it's possible to send invitation
#           mails to your users and also assign them to specific roles upon successful registration.
#  - public: Anyone can create a user account on cisidm. Note that it's still possible to use token registration
#            and user invitations as specified above.
registration = "public"

# The server block configures the built-in HTTP/2 servers.
server {
    # Whether or not cookies issued by the server should enforce 
    # HTTPS. If unset, the default will be set based on the protocol
    # in ui.public_url
    secure_cookies = true

    # Domain configures the cookie domain and the JWT issuer for access and refresh tokens.
    domain = "example.com"

    # The listen address of the HTTP server which requires authentication
    # for API endpoints.
    public_listener = ":8080"

    # The listen address of the admin HTTP server. Requests on this endpoint
    # will always be authenticated with the idm_superuser role assigned.
    #
    # SECURITY: You must make sure that this port is not publically accessible!
    admin_listener = ":8081"

    # Path to the static files to serve the web-interface on the public_listener.
    # Valid values are:
    #    - an empty string (or if the setting is omitted) serves the built-in
    #      web-interface.
    #    - a path to folder which contains web assets to serve.
    #    - a HTTP or HTTPS URL. In this case, cisidm will setup a single host reverse
    #      proxy and forward all UI/web related requests to that server.
    static_files = ""

    # In addition to the static_files setting above, it is also possible to expose additional
    # assets under the /files path. Specify a path to a local folder to expose all folder content
    # on the web.
    #
    # You may include additional assets like a logo or brand image here and set ui.logo_url = "/files/my-logo.png"
    extra_assets = ""

    # Configures additional origins that are allowed to perform cross-origin-resource requests (CORS).
    # Note that ui.public_url and server.domain (http and https) are always added to this list.
    allowed_origins = [
        "https://example.com",
        "https://*.example.com"
    ]

    # A list of CIDR network addresses or hostnames which are considered trusted. For any requests originating
    # from one of the specified networks, cisidm will trust the X-Forwarded-For header to determine the actual
    # client IP.
    trusted_networks = [
        "10.1.1.1/32", # A single host
        "10.1.2.1/24", # A whole sub-net in CIDR notation
        "traefik",     # A hostname, cisidm will resolve the hostname on a regular basis to detect changes to the IP.
                       # In containerized environments with changing container IPs it's best to use this method.
    ]

    # A list of application domains to which cisidm will permit a redirect after a successful login/access token refresh.
    # If this is unset, cisidm will refuse to redirect the user to prevent open-redirect vulnerabilitites.
    allowed_redirects = [
        "example.com", # Allow redirects only to example.com
        ".example.com", # Allow redirects to all subdomains of example.com
    ]
}

# The JWT block configures addtional settings for signing access and refresh tokens.
jwt {
    # The audience for JWT tokens. This defaults to server.domain
    audience = ""

    # The secret used to sign various data and tokens. Rotating this secret will invalidate any access 
    # and refresh tokens.
    #
    # SECURITY: for the time being, cisidm only supports signing tokens in HS512 mode.
    #           Support for public-private keypairs with JWSK (JSON-Web-Signing-Keys) is planned but not
    #           yet implemented.
    secret = "some-secure-random-string"

    # Configures the time-to-live for all access tokens issued by cisidm. This defaults to 1h.
    access_token_ttl = "1h"
    
    # The name of the cookie that will hold the access token. This defaults to cis_idm_access
    access_token_cookie_name = "cis_idm_access"

    # Configures the time-to-live for all refresh tokens issued by cisidm. This defaults to 720h.
    refresh_token_ttl = "720h"

    # The name of the cookie that will hold the refresh token. This defaults to cis_idm_refresh.
    # Note that the refresh cookie is limited to the refresh API endpoint /tkd.idm.v1.AuthService/RefreshToken
    # on server.domain
    refresh_token_cookie_name = "cis_idm_refresh"
}

# The UI block configures settings for all user-facing interface like the web-ui or any
# mail or SMS templates.
ui {
    # The name of your deployment. This will be visible in the built-in user interface and also
    # in compiled mail templates.
    site_name = "Example Inc"

    # The public address at which the cisidm server is reachable (the web-ui).
    public_url = "https://account.example.com"

    # A optional URL that is used on the web-ui and in mail templates when clicking the brand/deployment
    # logo or the name of the deployment.
    # If omitted, this defaults to ui.public_url
    site_name_url = "https://example.com"

    # Configures the resource where to find the the deployment/brand logo. This might be a fully-qualified
    # external URL or a web-resource on the cisidm public server. See server.extra_assets for more information.
    logo_url = "/files/logo.png"

    # Configures the URL template to build the redirect URL for forward authentication.
    # This value should not be set if the built-in web-ui is used.
    login_url = ""

    # Configures the URL template to build the redirect URL for forward authentication if the access token is
    # expired.
    # This value should not be set if the built-in web-ui is used.
    # The built-in web-ui will automatically try to refresh the access token and redirect the user back to
    # the application. If the refresh token is invalid (i.e. expired), the web-ui will render the login screen
    # before redirecting the user back.
    refresh_url = ""

    # Configures the URL template to build the mail verification URL.
    # This value should not be set if the built-in web-ui is used.
    verify_mail_url = ""

    # Configures the URL template to build URL used in user invitations.
    # This value should not be set if the built-in web-ui is used.
    registration_url = ""
}

# The (single) twilio block configures the Twilio integration which allows sending SMS messages
# to your users. This is required for phone-number verification to work.
twilio {
    from = "Example Inc"
    sid = "your-twilio-account-sid"
    token = "your-twilio-account-token"
}

# The (single) mail block configures the SMTP relay server that is used to send mail messages.
# This is required for email-verification, user-invitation and password reset mails to work.
mail {
    # The hostname (FQDN) of your mail server
    host = "smtp.example.com"

    # The port of your mail server
    port = 456

    # The user and password to authenticate
    user = "noreply"
    password = "a-secure-password"

    # The sender name used in e-mails. Make sure the value specified here is actually
    # allowed as a sender on your SMTP server.
    from = "Example Inc <noreply@example.com>"

    # Whether or not SSL/TLS should be used when connecting to the mail server.
    use_tls = true

    # allow_insecure may be set to true to disable certificate validation when use_tls = true.
    #
    # SECURITY: make sure you know what you're doing before setting this value to true.
    allow_insecure = false
}

# The (single) webpush block can be used to configure the Voluntary Application Server Identification (VAPID)
# for end-to-end encrypted WebPush support.
webpush {
    # The email address of the person responsible for this cisidm deployment.
    # This value is sent to web-push gateways upon web-push subscription.
    admin = "admin@example.com"

    # The public VAPID key
    vapid_public_key = "..."

    # The private VAPID key
    vapid_private_key = "..."
}


# The dry_run block may be used to set cisidm into "dry-run" mode. In this mode, any outgoing mail or SMS
# notification will be redirected to addresses specified here.
# This setting is mainly for developers that want to test delivery or new templates.
#
# dry_run {
#     main = "testing@example.com"
#     sms = "+4312312312"
# }

# Custom User Field Configuration
# -------------------------------------------------------------------------------------------------------
#
# cisidm already supports storing private user information like email addresses, phone numbers and addresses (like for delivery/billing).
# Though, in most use-cases where cisidm provides authentication and authorization in micro-service environments
# there will be some need for additional user data. Some use cases might include:
#  - storing user settings (like notification preferences)
#  - adding additional per-user data like company internal phone extensions,
#
# For such cases, cisidm supports definition additional user fields that are stored as JSON blobs in the
# cisidm database. For more information on additional user fields please refer to the documentation.
# 
# Take the following example:

field "string" "internal-phone-extension" {
    # This field may be seen by any authenticated user
    visibility = "public"
    
    # An optional description, this is for documentation purposes only
    description = "The company internal phone extension"

    # This field is populated by an administrator and cannot be changed by the user
    # themself.
    writeable = false

    # A display name for the self-service web-ui. This is currently unused.
    display_name = "Notification settings"
}

field "object" "notification-settings" {
    # Only the user is able to see this field
    visibility = "private"

    # This field may be written by the user
    writeable = true

    property "bool" "sms" {
        description = "Wether the user wants to receive notifications via SMS"
    }

    property "bool" "email" {
        description = "Wether the user wants to receive notifications via EMail"
    }
}

# Role, Permissions and user/role overwrites
# -------------------------------------------------------------------------------------------------------

# A list of permissions that are stored in cisidm. Note that it's completely fine to use the
# permission feature without every specifying some in the configuration. cisidm does treat permissions
# as "unknown" strings since evaluation of privileges is either done using rego policies or by other
# micro-service applications that just request a list of user permissions.
# Though, for permissions specified here, cisidm builds a hierarchical tree so it's easier to
# assign multiple permissions by using a shared prefix. Hierarchical levels of a permission string
# are seperated using colons (":").
#
# The example values below will build the following tree of permissions:
#
#    - roster
#      - write
#         - create
#         - approve
#      - read
#    - calendar
#      - write
#         - create
#         - delete
#         - move
#      - read
#
# With this tree, it's possible to assign the permission string "roster" and "calendar:write" to roles
# which cisidm will resolve to the following permission set:
#
#    roster:write:create
#    roster:write:approve
#    roster:read
#    calendar:write:create
#    calendar:write:delete
#    calendar:write:move
#
# Note: permissions themself do not have any authorizational meaning for cisidm per-se but the can be used
# in rego policies to implement permission/attribute-based access control (ABAC) rather than
# role-based-access-control (RBAC).
permissions = [
    "roster:write:create",
    "roster:write:approve",
    "roster:read",
    "calendar:write:create",
    "calendar:write:delete",
    "calendar:write:move",
    "calendar:read"
]

# Role blocks can be used to configure static roles which cannot be modified or deleted
# via the tkd.idm.v1.RoleService API.
role "computer-accounts" {
    # The name of the role. This is just for human representation but must be set.
    name = "Computer Accounts"

    # An optional description of the role. This is just for human representation.
    description = "Accounts for shared office computers"

    # A list of permissions that are assigned to this role.
    # See permissions above.
    permissions = [
        "roster:read"
    ]
}

# Overwrite blocks allow to overwrite certain settings on a per role or per-user basis.

overwrite "role" "idm_superuser" {
    # This block will set the access_token_ttl and the refresh_token_ttl to a much lower value for all
    # user accounts with the idm_superuser role. Note that roles are matched by ID!
    access_token_ttl = "10m"
    refresh_token_ttl = "2h"
}

overwrite "user" "computer-account-1" {
    # Set access and refresh token TTLs for a user account with ID "computer-account-1".
    # Though, overwriting on per-user basis is less common since one must first figure out
    # the ID of the user.
    access_token_ttl = "1h"
    refresh_token_ttl = "1480h"
}

