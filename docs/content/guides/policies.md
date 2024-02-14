---
next:
  text: CLI Reference
  link: ./cli-reference.md
---

# Policies

Authorization in `cisidm` is implemented using [Rego
Policies](https://www.openpolicyagent.org/docs/latest/policy-language/) from the
[Open-Policy-Agent](https://www.openpolicyagent.org/) project.

:::tip Note

For the time being, only policies for forward-authentication can be specified.
In the future, `cisidm` will support dynamic policy creation using the API and
also provide endpoints to query policy decision for services that directly 
integrate with `cisidm`.

:::

## Configuration

Policies can either be loaded from a directory or may be specified inline in the
configuration file:

```hcl

policies {
    # Load all .rego files from the following directories:
    directories = [
        "./policies"
    ]

    # Inline specification of a policy name "superuser"
    policy "superuser" {
        content = <<EOT
        package cisidm.forward_auth

        import rego.v1

        allow if {
            # input.subject is only set when the request is authenticated
            input.subject
            input.subject.roles

            # The user must have the idm_superuser role to be granted access
            # regardless of the requested resource.
            some role in input.subject.roles
            role.ID = "idm_superuser"
        }

        EOT
    }
}

```

## Policy Structure

Each policy starts with a `package` declaration. For forward-authentication policies,
`cisidm` queries `data.cisidm.forward_auth` by default. Though, if you want to
use a different policy package name, you can set `forward_auth_query` in your
policy configuration block.

Here's a simple example that permits requests to the `/protected` resource on 
`app.example.com` only if the request is authenticated and the user has the 
`app-user` role assigned. If the user has the `idm_superuser` role assigned, all
requests regardless of the requested host/resource path are allowed.

You can test the following rule on the OPA Rego playground: [https://play.openpolicyagent.org/p/eEPkMSqjrh](https://play.openpolicyagent.org/p/eEPkMSqjrh)

```rego
package cisidm.forward_auth

import rego.v1

# Allow the request if all conditions inside the allow rule evaluate to true
allow if {
    input.path = "/protected"
    input.host = "app.example.com"

    app_user_role
}

allow if {
    user_is_superuser
}

user_is_superuser if {
    some role in user_role_ids
    id = "idm_superuser"
}

# A rule that checks if the user has at least one role with ID "app-user"
# assigned
app_user_role if {
    input.subject
    input.subject.roles

    some role in user_role_ids
    role = "app-user"
}

user_role_ids contains id if {
    some role in input.subject.roles
    id := role.ID
}

```

Whenever the forward auth endpoint (`/validate`) is queried by the reverse
proxy, `cisidm` will evaluate the `forward_auth_query` (wich defaults to `data.cisidm.forward_auth` as stated above) to get a decision if the request should be allowed or not.

For evaluation of policies, `cisidm` constructs an input document which can be
accessed using the `input` variable. See [Input Document](#input-document) at the
end of the section for a complete document reference.


As another example, consider the following policy which allows access to `helpdesk.example.com`
only if the `job` attribute of the authenticated user is set to `"Support"`. This 
is a simple implementation of Attribute-Based-Access-Control

```rego
package cisidm.forward_auth

import rego.v1

allow if {
    input.host = "helpdesk.example.com" # Evaluate this rule only for requests
                                        # to helpdesk.example.com

    input.subject                       # request must be authenticated
    input.subject.fields                # Ensure the user has custom fields
                                        # populated

    # Ensure the user's job is Support
    input.subject.fields["job"] = "Support"
}
```

:::tip Examples

Refer to our [Policy Examples](./policy-examples/) for more sophisticated examples
includes RBAC, ABAC and even some complex permission examples.

:::

## Policy Results

After policy evaluation, `cisidm` checks for the following properties:

- `allow`: Whether or not the request should be allowed
- `status_code`: When the request is denied and `status_code` is set to a non-zero
                 value, `cisidm` immediately replies with the specified code and
                 `response_body` instead of trying to figure out an appropriate return
                 code. 
- `headers`: A map of HTTP headers (`map[string][]string`) that should be added to
             the response. If the request is allowed, those headers can be forwarded
             to the upstream service.  
             If the request is denied, those headers are forwarded to the user-agent
             the performed the initial request.
             Note that forwarding of custom headers might require configuration on
             your reverse proxy.


Below is an example that denies access to `/private`  even if the user is authenticated, still allowing the `idm_superuser` to access:

Playground: [https://play.openpolicyagent.org/p/ituLyXKA3K](https://play.openpolicyagent.org/p/ituLyXKA3K)

```rego
package cisidm.forward_auth

import rego.v1

# Deny access if the resource is protected
allow := false if {
    input.path = "/private"
}

# Always allow access to superusers
allow if {
    is_superuser
}

# The HTTP response headers if the resource is protected
headers := {
    "Content-Type": ["application/json"]
} if is_protected_resource

# The response body if the resource is protected
response_body := json.encode({
    "error": "sorry, you're not allowed to perform this operation"
}) if is_protected_resource

# Return status code 403 if the resource is protected
status_code := 403 if is_protected_resource

# Helper-Rule that evaluates to true if the request is not allowed and
# path matches /private
is_protected_resource if {
    input.path = "/private"
    not allow
}

# HelperRule that evaluates to true if the request is authenticated and the user
# is part of the idm_superuser role
is_superuser if {
    input.subject
    input.subject.roles

    some role in input.subject.roles
    role.ID = "idm_superuser"
}
```

## Input Document

The following input document is constructed by `cisidm` when evaluating forward-
authentication policies and can be accessed in policies using the `input` variable:

```hcl
input = {
    # The following data is directly copied from the user request that was
    # forwarded by the reverse proxy:

    # The requested resource path.
    # example: /protected/index.html
    path = "<requested-resource-path>"

    # The host name of the request.
    # example: app.example.com
    host = "<requested-hostname>"

    # The request method.
    # example: GET, PUT, ...
    method = "<request-method>"

    # All request headers as a map
    headers = {
        "Content-Type" = ["application/json"]
    }

    # Any query parameters of the request
    query = {}
    
    # The IP address of the client. Note that the reverse proxy must set the
    # X-Forwarded-For header and the IP address of the reverse proxy must be
    # in the trusted_networks setting. Otherwise the IP address of the reverse
    # proxy will be set for this field.
    client_ip = "<client-ip-address>"

    # When the request contains a valid access token than cisidm will also resolve
    # the requesting user and populate the subject field as follows:
    subject = {
        # The unique ID of the user
        id = "<unique-user-id>"

        # The username of the authenticated user.
        # SECURITY: only use this field if the username-change feature is disabled!
        username = "<username>"

        # A list of assigned roles for the authenticated user.
        # If the access token is a user API token, than only token roles are
        # reported.
        roles = [
            {
                ID = "<role-id>",
                Name = "<role-name>",
                Description = "<role-description>"
            }
        ]

        # A list of resolved permissions from all assigned roles.
        # If the access token is a user API token, than all permissions from the
        # roles assigned to the token are set.
        permissions = [
            "<list-of-all-role-permissions>"
        ]

        # Any custom user fields. Usable for attribute-based-access-control (ABAC)
        fields = {
            # Custom user fields
        }

        # The user's primary email address, if any
        email = "<primary-user-mail>"

        # The user's configured display name, if any
        display_name = "<user-display-name>"

        # The access token kind. This may be one of the following values:
        #  - password: Token was obtained using password-authentication only
        #  - mfa: Token was obtained by using two or multi-factor authentication
        #  - webauthn: Token was obtained using Webauthn or Passkey
        #  - api: A user generate API token.
        token_kind = "<token-kind>"
    }
}
```

The definition of the input object passed to forward_auth queries can be found
[here](https://github.com/tierklinik-dobersberg/cis-idm/blob/main/internal/services/auth/types.go)
