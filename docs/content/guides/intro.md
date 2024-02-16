# Welcome to `cisidm`

Welcome to `cisidm`, a simple Identity Management Service with Self-Service capabilities for custom application development or to secure self-hosted applications with support for Single-Sign-On (SSO) using [Proxy/Forward Authentication](./policies.md) and/or [OpenID Connect (using DexIdP)](./setup-oidc.md).


::: tip Features

If you just want to check out the list of features provided by `cisidm`, head over to the [Features](#features) section

:::

<br>
<hr>

**Contents**

[[toc]]

<hr>

## Use Cases

`cisidm` has been developed with the following two use-cases in mind:

### SSO Authentication / Proxy Auth

`cisidm` can be used to provide proxy authentication to your services just like
[Authelia](https://www.authelia.com/) or [Authentik](https://goauthentik.io/)
when used with a supported reverse proxy like
[Traefik](https://doc.traefik.io/traefik/).

In this mode your reverse proxy will be configured to authenticate requests to
your upstream services by quering the `/validate` endpoint of `cisidm`. If the
request is allowed, `cisidm` will add additional headers to the request so your
upstream services can authenticate users by just looking at the request headers.

Authorization is implemented using [Rego
Policies](https://www.openpolicyagent.org/docs/latest/policy-language/) from the
Open-Policy-Agent project. Checkout our [Policies Guide](./policies.md) for more
examples.


### Custom Application / Service Development

`cisidm` can also be used as a standalone authentication and user-management
micro-service so you don't need to roll your own.

For this, `cisidm` provides an extensive API based on the awesome `Connect-RPC`
system which brings compatability with Browser HTTP and GRPC. For extended
use-cases, `cisidm` can also be configured to store and manage custom user
metadata that may even be visible and writeable by your users. Check the
[Additional User Fields Guide](./extra-user-fields.md) for more information and
examples.

## Features

The following is a likely incomplete list of features currently implemented by
`cisidm`:

- **Authentication**
  - Password-based login
  - Second-Factor (2FA) using
    - Time-Based-One-Time Password (TOTP)
    - Backup Recovery Codes
    - SMS `work-in-progress`
  - Password-less authentication:
    - WebauthN / Passkeys
    - E-Mail magic links (`work-in-progress`)
    - One-Time passwords using SMS (`work-in-progress`)
- **Self-Service-Portal (Web-UI)**
  - Manage profile:
    - Avatar
    - Username
    - Display Name
    - First and Givenname
    - Birthday
  - User Addresses for Work-place, Delivery, ...
  - E-Mail Addresses
    - With mail-verification
  - Phone Numbers
    - With SMS verification
  - Custom user fields: see [here](./extra-user-fields.md)
  - Per-User API Tokens
  - Password reset mails
- User Invitations by E-Mail (including a "password-reset" link)
- User Roles and Permissions
- **Support for proxy authentication (aka Forward-Auth)**
  - Uses Open-Policy-Agent policies (rego) for authorization
  - Supports RBAC, ABAC, PBAC, AWS-IAM style policies ...
- **API / Integration Support**
  - Defined using
    [`protobuf`](https://github.com/tierklinik-dobersberg/apis/tree/main/proto/tkd/idm/v1)
    with clients already available for `Go` and `Typescript`/`Javascript`
  - Every `cisidm` API can be consumed using `Connect-RPC`, `Browser HTTP` and `GRPC`
  - Users-API (manage cisidm users)
  - Role-Management API
  - SelfService API
  - Authentication API
  - Notification API with support for SMS, E-Mail (with HTML E-Mail templates)
    and WebPush (using VAPID)


## Comparisons

::: tip To be done
:::

### Authelia

### Kanidm

### DexIdP
