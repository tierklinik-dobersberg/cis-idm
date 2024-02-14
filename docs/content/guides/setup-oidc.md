# OpenID Connect Setup

While `cisidm` does not implement OIDC itself it is very easy to add support for Open-ID-Connect using the [DexIdp](https://dexidp.io/) project.
It integrates with `cisidm` using the [**AuthProxy**](https://dexidp.io/docs/connectors/authproxy/) provider.

In this mode, OIDC clients ask DexIdp for authentication/authorization which in turn takes authenticated user information from
HTTP request headers that have been added by `cisidm` during the forward-authentication.

For self-hosted environments it is also possible to configure DexIdp to skip the consent screen and directly redirect the user to the 
requested application providing a seamless Single-Sign-On experience.

We have tested OIDC support based on DexIdp with [Rallly](https://rallly.co), [WikiJS](https://js.wiki/), [Outline](https://www.getoutline.com/) and
[Nextcloud](https://nextcloud.com). It's expected that every OIDC client will work with this setup.
