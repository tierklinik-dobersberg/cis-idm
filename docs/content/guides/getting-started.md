# Getting Started

Welcome to the Getting-Started Guide for `cisidm`. On this page we will walk through a simple docker-compose based setup to get you started.

## Overview

The easiest way to get started using `cisidm` is to deploy it and the database (rqlite) using docker-compose.
In this guide we will setup a docker deployment with the following service:

- [**rqlite**](https://rqlite.io/)  
  rqlite is a database based on sqlite with clustering support using RAFT. It used by cisidm to store all it's
  data like users, roles, webauthn credentials and more.

- [**traefik**](https://doc.traefik.io/traefik/):  
  A flexible and powerful reverse proxy that will handle automatic HTTPS via Let's Encrypt and secure
  our services by enforcing authentication via cisidm.

- [**cisidm**](https://github.com/tierklinik-dobersberg/cis-idm):  
  The identity management server.

- [**echoserver**](https://gcr.io/google_containers/echoserver:1.4):  
  A simple demo application that will be secured using cisidm and traefik.

:::warning Customization Required

Please make sure to update the configurations in this guide to match your domains and e-mail addresses!
You will also need DNS to be setup correctly (i.e. your sub-domains pointing to your server/IP address).

:::

## Setup the Project

Prepare the directory structure for our deployment.

```bash
mkdir -p cisidm-demo/config;
```

The `cisidm-config/config` folder will hold the configuration files for `cisidm` and `traefik`.
Our `docker-compose.yml` file will be placed directly in `cisidm-config`.

## Configuration File

Create the configuration file for `cisidm` in `cisidm-demo/config/idm-config.yml`:

<CodeGroup>
  <CodeGroupItem title="idm-config.yml">

```yaml
# The domain that cisidm is going to protect. If you plan on deploying other services using sub-domains make sure
# to set the this field to the top-level domain. Otherwise, browser will not send the access token cookie for sub-domains.
domain: example.com

# Whether or not the cookie should only be sent for HTTPS. It's best to keep this at true.
secureCookie: true

# The secret key used to sign the JWT access tokens. Choose some secure random string here.
# Note that if you ever change this value any access and refresh tokens will be invalidated and your users
# will need to re-authenticate.
jwtSecret: "some-secure-secret"

# The address for the rqlite database
rqliteURL: http://rqlite:4001/

# The default log level. Possible values are debug, info, warn and error
logLevel: info

# A list of IP networks (in CIDR notation) or hostnames from which cisidm will trust the
# X-Forwarded-For headers.
trustedNetworks:
  - traefik

# Whether or not your users require a registration token. Setting this to false enables public registration.
registrationRequiresToken: true

# The name of your deployment.
siteName: Example Site

# A URL that will be used in email templates and on the self-service UI.
siteNameUrl: https://example.com

# Configuration for the forward authentication support.
forwardAuth:
  # We require authentication for all sub-domains of example.com
  - url: "http(s){0,1}://(.*).example.com"
    required: true

# The public URL under which the login screen and self-service UI can be accessed.
publicURL: https://account.example.com

# Configures how long access tokens issued to your user are valid. You can keep this relatively short since
# users will also get a long-lived refresh token.
accessTokenTTL: "1h"

# cisidm automatically redirects your users back to the protected services after a successful login. To prevent open-redirect
# attacks make sure to restrict the allowed redirects. It's best to just allow the top-level domain and any
# sub-domains of your deployment.
allowedRedirects:
  - example.com
  - .example.com
```

  </CodeGroupItem>
</CodeGroup>

::: tip
Refer to the [Configuration File Reference](../architecture/config-reference.md) for a more detailed explanation of the configuration file.
:::

## Create the docker-compose file

Next, we create the docker-compose file that will contain our service definitions and also tell treafik which sub-domains should be routed to which services and that we want cisidm to enforce authentication.

Here's the complete docker-compose file, we'll break it down and explain everything below:

<CodeGroup>
  <CodeGroupItem title="docker-compose.yml">

```yaml
version: "3"

volumes:
  db:
    driver: local

services:
  rqlite:
    image: rqlite/rqlite:latest
    hostname: "12e94f6300a8"
    restart: unless-stopped
    environment:
      HTTP_ADDR_ADV: localhost
    command:
      - "-on-disk=true"
      - "-node-id=1"
      - "-fk"
    ports:
      - 4001:4001
    volumes:
      - db:/rqlite/file

  # Traefik ###################################################

  traefik:
    image: "traefik:v2.10"
    restart: unless-stopped
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.http.address=:80"
      - "--entrypoints.secure.address=:443"

      - "--entrypoints.http.http.redirections.entrypoint.to=secure"
      - "--entrypoints.http.http.redirections.entrypoint.scheme=https"

      - "--certificatesresolvers.resolver.acme.email=admin@example.com"
      - "--certificatesresolvers.resolver.acme.storage=/secrets/acme.json"
      - "--certificatesresolvers.resolver.acme.httpchallenge.entrypoint=http"

      # disable for production
      # - "--certificatesresolvers.resolver.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory"

    ports:
      - target: 80
        published: 80
        protocol: tcp
        mode: host
      - target: 443
        published: 443
        protocol: tcp
        mode: host

    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock"
      - "./config:/secrets"

  # IDM ############################################################
  cisidm:
    image: ghcr.io/tierklinik-dobersberg/cis-idm:latest
    restart: unless-stopped
    volumes:
      - ./config/idm-config.yml:/etc/config.yml:ro
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.idm.rule=Host(`account.dobersberg.vet`)"
      - "traefik.http.routers.idm.entrypoints=web"
      - "traefik.http.routers.idm.tls=true"
      - "traefik.http.routers.idm.tls.certresolver=resolver"
      - "traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate"
      - "traefik.http.middlewares.auth.forwardauth.authResponseHeaders=X-Remote-User,X-Remote-User-ID,X-Remote-Mail,X-Remote-Mail-Verified,X-Remote-Avatar-URL,X-Remote-Role,X-Remote-User-Display-Name"
    environment:
      CONFIG_FILE: "/etc/config.yml"
    depends_on:
      - rqlite

  echoserver:
    image: gcr.io/google_containers/echoserver:1.4
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.echo.rule=Host(`app.example.com`)"
      - "traefik.http.routers.echo.entrypoints=web"
      - "traefik.http.routers.echo.tls=true"
      - "traefik.http.routers.echo.tls.certresolver=resolver"
      - "traefik.http.routers.echo.middlewares=auth"
      - "traefik.http.services.echo.loadbalancer.server.port=8080"
```

  </CodeGroupItem>
</CodeGroup>

:::tip
If you're already familiar with docker-compose you can skip over the explanation to the next section [here](#start-the-project).
:::

Let's break it down a bit:

```yaml
version: "3"

volumes:
  db:
```

This just tells docker-compose which file-version we're using and that we want a volume named `db` that will hold our rqlite database.

### RQLite Container

```yaml
rqlite:
  image: rqlite/rqlite:latest
  hostname: "rqlite-node-1"
  restart: unless-stopped
  environment:
    HTTP_ADDR_ADV: localhost
  command:
    - "-on-disk=true"
    - "-node-id=1"
    - "-fk"
  ports:
    - 4001:4001
  volumes:
    - db:/rqlite/file
```

This defines the rqlite database service. Note that rqlite needs a stable hostname that's why we force the hostname to `rqlite-node-1`:

We also pass a few parameters to the rqlite server:

- `-on-disk=true`: Tells rqlite to store the sqlite database on the disk. By default, rqlite just keeps the database in memory as it's expected to run in cluster mode. For this simple setup, we just run one rqlite node so we want that data to be persisted to disk.
- `-node-id=1`: Each rqlite node in the cluster needs a stable ID. Even though we don't use clustering support here we still need to set a node ID.
- `-fk`: Enables foreign key support in the underlying sqlite database. This is required as the database schema from `cisidm` relies on cascading deletes.

### Traefik Container

```yaml
traefik:
  image: "traefik:v2.10"
  restart: unless-stopped
  command:
    - "--api.insecure=true"
    - "--providers.docker=true"
    - "--providers.docker.exposedbydefault=false"
    - "--entrypoints.http.address=:80"
    - "--entrypoints.secure.address=:443"

    - "--entrypoints.http.http.redirections.entrypoint.to=secure"
    - "--entrypoints.http.http.redirections.entrypoint.scheme=https"

    - "--certificatesresolvers.resolver.acme.email=admin@example.com"
    - "--certificatesresolvers.resolver.acme.storage=/secrets/acme.json"
    - "--certificatesresolvers.resolver.acme.httpchallenge.entrypoint=http"

    ports:
      - target: 8080
        published: 8080
        protocol: tcp
        mode: host
      - target: 80
        published: 80
        protocol: tcp
        mode: host
      - target: 443
        published: 443
        protocol: tcp
        mode: host

    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock"
      - "./config:/secrets"
```

The above service snippet configures Traefik as our reverse proxy. Let's go through the arguments passed to Traefik:

- `--api.insecure=true`  
  This enables the management dashboard of Traefik. This is useful for debugging service configurations. The management dashboard will be accessible on port 8080:

- `--providers.docker=true`,  
  `--providers.docker.exposedbydefault=false`  
  This enables the docker provider so we can configure Traefik using container labels. This will be explained in more detail below.

- `--entrypoints.http.address=:80`,  
  `--entrypoints.secure.address=:443`  
  Tell traefik that we want to have two entrypoints, one on port 80 (named `http`) and one on port 443 (named `secure`):

- `--entrypoints.http.http.redirections.entrypoint.to=secure`,  
  `--entrypoints.http.http.redirections.entrypoint.scheme=https`  
  Configures a redirection from insecure, plain-text HTTP to TLS encrypted HTTPS. Whenever a user tries to access
  a service using HTTP, Traefik will automatically redirect the browser to the secure endpoint.

- `--certificatesresolvers.resolver.acme.email=admin@example.com`,  
  `--certificatesresolvers.resolver.acme.storage=/secrets/acme.json`,  
  `--certificatesresolvers.resolver.acme.httpchallenge.entrypoint=http`  
  Configure a certificate resolver that uses ACME to request certificates from Let's Encrypt.

### Identity Management Server

```yaml
cisidm:
  image: ghcr.io/tierklinik-dobersberg/cis-idm:latest
  restart: unless-stopped
  volumes:
    - ./config/idm-config.yml:/etc/config.yml:ro
  labels:
    - "traefik.enable=true"
    - "traefik.http.routers.idm.rule=Host(`account.example.com`)"
    - "traefik.http.routers.idm.entrypoints=web"
    - "traefik.http.routers.idm.tls=true"
    - "traefik.http.routers.idm.tls.certresolver=resolver"
    - "traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate"
    - "traefik.http.middlewares.auth.forwardauth.authResponseHeaders=X-Remote-User,X-Remote-User-ID,X-Remote-Mail,X-Remote-Mail-Verified,X-Remote-Avatar-URL,X-Remote-Role,X-Remote-User-Display-Name"
  environment:
    CONFIG_FILE: "/etc/config.yml"
  depends_on:
    - rqlite
```

This configures the `cisidm` docker container, mounts the configuration file to `/etc/config.yml` and tells traefik that is should be reachable at `account.example.com`.

It also configures a new HTTP middleware `auth` that uses the forward-auth feature of traefik:

- `traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate`  
  This configures the Forward-Auth middleware (named `auth` here) to forward any HTTP request to cisidm to determine if the user is actually allowed to access the requested resource. If cisidm replies with a HTTP success status code (2xx) than traefik will forward the original request to the actual service container. If cisidm replies with an error code, traefik will immediately return the response from `cisidm` to the user. This is used to redirect the user to the login page in case the request is unauthenticated.

- `traefik.http.middlewares.auth.forwardauth.authResponseHeaders=....`  
  When `cisidm` successfully authenticated a request, it will return a set of headers that contain information about the logged in user. With this setting, we tell traefik to forward those headers to the actual service container. This enables service containers to know which user performs the access without the need to parse and validate the JWT token issued by `cisidm` for every successful authentication.

## First Start

Finally it's time to start our services by calling the following command from the project directory (the one that contains the docker-compose file)

```
docker-compose up -d
```

:::tip Admin User
Whenever `cisidm` starts, it checks if a user with super-user privileges (member of the `iam_superuser` role) exists. If not, a new registration token will be created and logged to stdout.

To create your initial admin user, copy that token from the log output (`docker-compose logs cisidm | grep "superuser account"`) and visit `https://account.example.com/registration?token=YOUR_TOKEN` and replace `YOUR_TOKEN` with the token from the log output.

It's also possible to use the cli tool for the registration: `idmctl register-user --registration-token YOUR-TOKEN your-username`.
:::

<br />

---

<br />

:::tip Congratulations
You just finished setting up cisidm with a HTTPS enabled reverse proxy that will now protect your services using Proxy/Forward Authentication.

Now it's time to check the [User and Role Administration Guide](./user-role-management.md) or the [Command Line Reference](./cli-reference.md).
:::
