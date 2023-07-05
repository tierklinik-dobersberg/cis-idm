# cisidm: Simple Self-Hosted Identity-Management Server

Welcome to the project page of cisidm, a simple, self-hosted and high-available identity management server.

Please note that cisidm is being actively developed and is not yet complete or ready for production use. Use at your own risk! 

## Features

- Protobuf defined API using [Connect](https://buf.build/blog/connect-a-better-grpc) for interoperability with browsers and gRPC.
- Support for **2FA using TOTP**
- Support for **WebAuthN** and **Passkeys**
- A public listener (which requires authentication)
- A admin/internal listener for un-authenticated use by other micro-services
- Privacy (access to user profile fields) backed into Protobuf (see tierklinik-dobersberg/apis)
- Stateless (uses [rqlite](https://rqlite.io) for storage) so it can be deployed
  multiple times for load-balancing.
- A `/validate` endpoint that can be used for proxy forward authentication as supported in [Traefik](https://doc.traefik.io/traefik/) or [Caddy](https://caddyserver.com/docs/caddyfile/directives/forward_auth).
- A pretty **Self-Service Portal** UI:
  - Update / Manage profile / avatar picture.
  - Change passwords
  - Enroll 2FA
  - Enroll WebAuthN/Passkeys
  - Self registration (may optionally require a registration token) with either Password or WebAuthN
  - Change privacy settings (`in-development`)
  - Manage E-Mail addresses
    - Verification of E-Mail addresses is `in-development`
  - Manage phone numbers
    - Verification of phone numbers is `in-development` (will use Twilio)
  - Manage addresses (delivery/billing/...)

## Quick-Start

To quickly get up and running cisidm for testing purposes you can use the [docker-compose.yml](./docker-compose.yml) file to bring up cisidm a single rqlite node and Traefik (configured for self-signed certificates). Please make sure to update [./config.test.yml](./config.test.yml) before to match your settings. Also, it is recommended to update your `/etc/hosts` so the domain names you use for testing will resolve to 127.0.0.1.

For example, given the following configuration file:

```yaml
audience: example.dev
domain: example.dev
secureCookie: true
jwtSecret: some-secure-string
rqliteURL: http://rqlite:4001/
forwardAuth:
  - url: http(s){0,1}://wiki.example.dev
    required: true
allowedOrigins: 
  - http://example.dev
  - https://example.dev
publicURL: https://account.example.dev
allowedRedirects:
  - wiki.example.dev
```

You should make sure that `example.dev`, `wiki.example.dev` and `account.example.dev` resolve to localhost.

Finally, just launch:

```bash
docker-compose build && docker-compose up
```

### Important Warning

For the time being cisidm depends on tierklinik-dobersberg/apis (for Go) and on '@tkd/apis' (for JS/TypeScript). These are not yet released! The Dockerfile provided in this repository expected a `tkd/apis:latest` image on your machine.

Just clone [tierklinik-dobersberg/apis](https://github.com/tierklinik-dobersberg/apis) and execute
 `docker build -t tkd/apis:latest .` once before running `docker-compose build` from this repo.

This will likely be fixed in the next weeks.

## Documentation

To be done.

## Versioning

Since cisidm is still in early development it has not yet reached a stable API. While we try to avoid breaking changes please expect them to happen at this point!

Once we reach a final v1 the APIs will be frozen and not changed in backwards incompatible ways. Stay tuned ...

## License

For now, this repository is licensed under [MIT License](./LICENSE). While this might be subject to change cis-idm will stay OSS but may start prohibiting unlicensed enterprise use.

Any such changes will be communicated and can be discussed beforehand on the Github Issue Tracker.
