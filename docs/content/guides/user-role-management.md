# User and Role Management

`cisidm` implements the concept of users and roles where each user can be assigned to any number of roles. Note that roles per-se do not permit any special permissions. For simple authentication purposes, an administrator can use role definitions to configure forward-auth permissions. For a micro-service environment, developers may grant permissions based on user roles but `cisidm` itself does not implement RBAC (Role Based Access Control). For the time being, roles are a flat list instead of being hierarchical so it's not possible that a single role assignment automatically includes other roles. However, this might change in a future release.

:::warning Note
There is one special role in `cisidm` called the `idm_superuser` role. This role actually does imply a set of permissions: Any user with this role can perform any action on any API endpoint of cisidm and is thus considered an administrative account.

It's **strongly advised** to only use a `idm_superuser` account for administrative tasks and multi factor authentication (TOTP or SMS/E-Mail codes) and use a separate user account for daily work/authentication.
:::
