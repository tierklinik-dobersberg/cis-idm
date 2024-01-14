---
next:
  text: Policies
  link: ./policies.md
---

# Additional User Fields

Although cisidm already provides common user information (see [Users](#users)),
there may be use-cases where additional user metadata needs to be stored. In
order to avoid requiring a dedicated service to store this metadata, `cisidm`
provides support for custom user fields. Those fields are stored in the database
as a JSON blob.

Those fields may either be used by services that directly integrate with
`cisidm` using the API or may be used in [Policies](./policies.md) to implement
Attribute-Based-Access-Control (ABAC).

Additional user fields need to be configured in the configuration file. Below is
a short example on how to configure a custom field. Refer to the [Configuration
File Reference](./config-reference.md) for more information.

```hcl
field "string" "internal-phone-extension" {
    # This field may be seen by any authenticated user
    visibility = "public"
    
    # An optional description, this is for documentation purposes only
    description = "The company internal phone extension"

    # This field is populated by an administrator and cannot be changed by the
    # user themself.
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
```

Note that even if `cisidm` allows users to update their additional fields (as
long as it's set to `writable`), the built-in web-UI does not yet support
displaying or manipulating those fields.

As an administrator, you may set those fields using the idmctl commandline utility:

```bash
idmctl users set-extra [user] [path] [value]

# Examples:
idmctl users set-extra alice "notification-settings.sms" true
idmctl users set-extra alice "internal-phone-extenstion" '"34"'
```

:::tip Note

Values should be encoded in their JSON representation. For example, the
following should also work:

```bash
idmctl users set-extra alice                    \
    "notification-settings"                     \
    '{"sms": true, "email": false}'
```

:::

