# E-Mail Setup

To enalbe email verification and support for password-reset links, account
invitations etc it's require to configure an outgoing E-Mail (SMTP) server in
your `cisidm` configuration:

```hcl

# The (single) mail block configures the SMTP relay server that is used to send
# mail messages. This is required for email-verification, user-invitation and
# password reset mails to work.
mail {
    # The hostname (FQDN) of your mail server
    host = "smtp.example.com"

    # The port of your mail server
    port = 456

    # The user and password to authenticate
    user = "noreply"
    password = "a-secure-password"

    # The sender name used in e-mails. Make sure the value specified here is
    # actually allowed as a sender on your SMTP server.
    from = "Example Inc <noreply@example.com>"

    # Whether or not SSL/TLS should be used when connecting to the mail server.
    use_tls = true

    # allow_insecure may be set to true to disable certificate validation when
    # use_tls = true.
    #
    # SECURITY: make sure you know what you're doing before setting this value
    # to true.
    allow_insecure = false
}

```
