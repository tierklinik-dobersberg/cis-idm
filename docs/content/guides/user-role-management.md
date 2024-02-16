---
next:
  text: Additional User Fields
  link: ./extra-user-fields.md
---

# User and Role Management

In `cisidm` identity management is split up into [Users](#users),
[Roles](#roles) and [Permissions](#permissions).

Each user might be assigned multiple roles (or even none) with each role having
a set associated permissions.

Note that cisidm itself does not define any permissions itself since this is up
to the administrator to define. [Permissions](#permissions) are just opaque
strings in cisidm. Instead, the cisidm APIs are already protected by either just
requiring valid authentication or an administrator account with the built-in
`idm_superuser` role. See [Roles](#roles) for more information.

<br>
<hr>

**Contents**

[[toc]]

<hr>

## Users

Users are the main identities used for authentication in cisidm. Users are
identified by a unique ID that is not allowed to change. Note that by default,
`cisidm` permit users to change their username as long as it's not already take.
The unique ID ensures that services that integrate with cisidm have a stable
identifier.

::: tip

When using cisidm just for forward authentication using a supported reverse
proxy like Traefik or Caddy, it's best to disable the username-change feature
using the configuration. This way the username is also expected to be stable and
can thus safely be sent to protected upstream services using the `X-Remote-User`
header.

See [Policies](/guides/policies.md) and the [Configuration File
Reference](./cli-reference.md) for more information on forward authentication.

:::

`cisidm` provides a default set of fields that users can change on their own
using the `tkd.idm.v1.SelfService` API endpoint. Those fields include:

 - First and Lastname
 - A Display name that should/can be used instead of the username
 - Birthday
 - Avatar
 - Multiple addresses: for example billing, delivery, ...
 - Multiple phone numbers
 - Multiple email addresses

Each user can also specify a primary phone and email address that is used for
password-reset codes and any other notification.

In addition, `cisidm` maintains metadata for two-factor authentication (2FA), 
Backup-Codes, WebauthN/Passkeys.

### Creating Users

There are multiple ways to setup user accounts in cisidm. If `registration` is
set to `"public"` in the configuration, any person can create a user account on
your deployment using the built-in web interface.

If `registration` is set to `"token"`, only users with a valid registration
token can register themself. An administrator can either send an invitation mail
to a user or generate a token and distribute manually. Using token based
registration allows the administrator to immediately assign one or more roles as
soon as the user completes the registration.

An administrator may also manually create user accounts and can then send an
account-creation-notice that contains a password reset link so your users can
choose their initial password.

<CodeGroup>
  <CodeGroupItem title="Invitation Mail">

```bash
idmctl users invite [email-addresses ...] --roles [role-name/id ...]

# For example:
idmctl users invite alice@example.com --roles help-desk --roles support
```

  </CodeGroupItem>

  <CodeGroupItem title="Generating Tokens">

```bash
idmctl users generate-registration-token
    --max-usage [number]                    # How often the token can be used
    --ttl [duration]                        # How long the token remains valid 
    --roles [role-name/id ...]              # One or more roles to assign upon
                                            # registration

# For example:
idmctl users generate-registration-token
    --max-usage 3 
    --ttl 24h                                          
    --roles help-desk                                  
    --roles support

# Output:
#   token: some-random-token-to-distribute
```

  </CodeGroupItem>

  <CodeGroupItem title="Manual Account Creation">

```bash
# Manually create a user
idmctl users create
    --name [username]
    --display-name [display-name] 
    --first-name [first-name]
    --last-name [last-name]
    --phone [phone-numbers...]
    --email [email-addresses...]
    --role [roles...]
    --password [plain-text-or-bcrypt-password]

# Send account creation notice to one or more users containing a password-reset
# link.
idmctl users send-account-notice [username]

# Example:
# Note that the first phone number and the first email address will be marked as
# the primary one.
idmctl users create
    --name alice
    --first-name Alice
    --last-name Mustermann
    --phone +4312341234                     
    --phone +49987987987
    --email alice@example.com
    --email alice@help-desk.example.com
    --role help-desk
    --role support

idmctl users send-account-notice alice
```

  </CodeGroupItem>

</CodeGroup>

### Additional User Fields

When integrating cisidm with your own services / applications you might need
additional metadata for users. Refer to the [Additional User
Fields](./extra-user-fields.md) guide for more information.

## Roles

As mentioned at the beginning, each user may be assigned multiple roles. Those
roles may either be used in [Policies](./policies.md) to implement
Role-Based-Access-Control (RBAC) or by services that integrate directly with
cisidm using it's exposed API.

`cisidm` itself does not care about the defined roles but rather allows an
administrator to configure roles based on their requirements.

:::warning Note

There is one special role in `cisidm` called the `idm_superuser`
role. This role actually does imply a set of permissions: Any user with this
role can perform any action on any API endpoint of cisidm and is thus considered
an administrative account.

It's **strongly advised** to only use a `idm_superuser` account for
administrative tasks with multi-factor authentication (TOTP or SMS/E-Mail codes)
enabled and use a separate user account for daily work/authentication.

:::

### Static configuration

Roles may already be specified in the configuration file. In this case, `cisidm`
will prevent modifications to those roles using the API (they are considered
system-roles).

To define a role in the configuration file, you just need to add a `role` block:

```hcl
role "computer-accounts" {
    # The name of the role. This is just for human representation but must be
    # set.
    name = "Computer Accounts"

    # An optional description of the role. This is just for human
    # representation.
    description = "Accounts for shared office computers"

    # A list of permissions that are assigned to this role.
    # See Permissions below.
    permissions = [
        "roster:read"
    ]
}
```

### Dynamic configuration

When no roles are configured in the configuration file or if
`enable_dynamic_roles` is set, then roles may be created dynamically using the
`tkd.idm.v1.RolesService` API.

For example, to create the same role as above using the [idmctl
utility](./cli-reference.md):

```bash
# Create the role, note if --id is not set, a random ID will be generated by
# cisidm.
idmctl roles create
    --id "computer-accounts"
    --name "Computer Accounts"
    --description "Accounts for shared office computers"

# Optionally, assign permissions to the role
# FIXME(ppacher): not yet implemented
idmctl roles add-permissions "computer-accounts" --permission "roster:read"

# Inspect role and it's assigned permissions
idmctl roles "computer-account"
idmctl roles get-permissions "computer-accounts"
```

## Permissions

Last but not least, each role may be assigned multiple permissions. Permissions
are (mostly) opaque text values without any special meaning in `cisidm` but
enable an administrator to implement Permission-Based-Access-Control (PBAC).
Since `cisidm` does not inspect the value of permissions, it's also possible to
implement more sophisticated authorization systems like AWS IAM Statements by
using JSON encoded objects as permission strings and parse/evaluate them in your
[rego policies](./policies.md).

For simple setups, one may enable `permission_trees = true` in the configuration
file. If enabled, permssions strings will be parsed as trees and `cisidm` will
resolve all user permissions based on that tree.

For example, consider the following configuration file:

```hcl
permission_trees = true

permissions = [
    "calendar:events:read"
    "calendar:events:write"
    "calendar:events:write:delete",
    "calendar:events:write:create",
    "calendar:events:write:move",
    "calendar:events:write:update"
]
```

Taken the above permission configuration, `cisidm` will build the following
permission tree:

```
- calendar
    - events
        - read
        - write
            - delete
            - create
            - move
            - update
```
    
When resolving user permissions, cisidm will try to detect inherited child
permissions automatically:

- assigned: `calendar:events:read`  
  resolved: `calendar:events:read`

- assigned: `calendar:events:write`  
  resolved: `calendar:events:write:delete`, `calendar:events:write:create`,
  `calendar:events:write:move`, `calendar:events:write:update`

- assigned: `calendar`  
  resolved: `calendar:events.read`, `calendar:events:write:delete`,
  `calendar:events:write:create`, `calendar:events:write:move`,
  `calendar:events:write:update`

