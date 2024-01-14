---
next:
  text: User and Role Management
  link: ./user-role-management.md
---

# Getting Started

Welcome to the Getting-Started Guide for `cisidm`. On this page we will walk
through a simple docker-compose based setup to get you started.

## Overview

The easiest way to get started using `cisidm` is to deploy it using
docker-compose. In this guide we will setup a docker deployment with the
following service:

- [**traefik**](https://doc.traefik.io/traefik/):  
  A flexible and powerful reverse proxy that will handle automatic HTTPS via
  Let's Encrypt and secure our services by enforcing authentication via cisidm.

- [**cisidm**](https://github.com/tierklinik-dobersberg/cis-idm):  
  The identity management server.

- [**echoserver**](https://gcr.io/google_containers/echoserver:1.4):  
  A simple demo application that will be secured using cisidm and traefik.

:::warning Customization Required

Please make sure to update the configurations in this guide to match your
domains and e-mail addresses! You will also need DNS to be setup correctly (i.e.
your sub-domains pointing to your server/IP address).

:::

## Setup the Project

Prepare the directory structure for our deployment.

```bash
mkdir -p cisidm-demo/config
```

The `cisidm-config/config` folder will hold the configuration files for `cisidm`
and `traefik`. Our `docker-compose.yml` file will be placed directly in
`cisidm-config`.

## Configuration File

Create the configuration file for `cisidm` in `cisidm-demo/config/idm-config.hcl`:

<CodeGroup>
  <CodeGroupItem title="idm-config.yml">

```hcl
# Configures where the SQLite3 database should be stored.
database_url = "file:/data/idm.db"

# Disable self-registration of users. This means that any user account must be
# created by an administrator using the `idmctl` cli utility.
registration = "disabled"

# This configures some defaults for the built-in server.
server {
    secure_cookies = true
    domain = "example.com"
    allowed_origins = [
        "https://example.com",
        "https://*.example.com",
    ]
    trusted_networks = [
        "traefik"
    ]
    allowed_redirects = [
        "example.com",
        ".example.com"
    ]
}

jwt {
    secret = "some-secure-random-string"
}

ui {
    site_name = "Example Inc"
    public_url = "https://account.example.com"
}

policies {
    debug = false

    policy "superuser" {
        content = <<EOT
        package cisidm.forward_auth

        import future.keywords.in

        allow {
            user_is_superuser
        }

        user_is_superuser {
            input.subject
            input.subject.roles

            some role in input.subject.roles
            role.ID = "idm_superuser"
        }
        EOT
    }
}

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

      # disable for production, this enables the use of Let's Encrypt staging servers
      # so any misconfiguration will not get you rate-limited.
      - "--certificatesresolvers.resolver.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory"

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
      - ./config/idm-config.yml:/etc/config.hcl:ro
      - db:/data
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.idm.rule=Host(`account.dobersberg.vet`)"
      - "traefik.http.routers.idm.entrypoints=web"
      - "traefik.http.routers.idm.tls=true"
      - "traefik.http.routers.idm.tls.certresolver=resolver"
      - "traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate"
      - "traefik.http.middlewares.auth.forwardauth.authResponseHeaders=X-Remote-User,X-Remote-User-ID,X-Remote-Mail,X-Remote-Mail-Verified,X-Remote-Avatar-URL,X-Remote-Role,X-Remote-User-Display-Name"
    environment:
      CONFIG_FILE: "/etc/config.hcl"

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

This just tells docker-compose which file-version we're using and that we want a volume named `db` that will hold our sqlite3 database.


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
```

This configures the `cisidm` docker container, mounts the configuration file to `/etc/config.hcl` and tells traefik that is should be reachable at `account.example.com`.

It also configures a new HTTP middleware `auth` that uses the forward-auth feature of traefik:

- `traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate`  
  This configures the Forward-Auth middleware (named `auth` here) to forward any HTTP request to cisidm to determine if the user is actually allowed to access the requested resource. If cisidm replies with a HTTP success status code (2xx) than traefik will forward the original request to the actual service container. If cisidm replies with an error code, traefik will immediately return the response from `cisidm` to the user. This is used to redirect the user to the login page in case the request is unauthenticated.

- `traefik.http.middlewares.auth.forwardauth.authResponseHeaders=....`  
  When `cisidm` successfully authenticated a request, it will return a set of headers that contain information about the logged in user. With this setting, we tell traefik to forward those headers to the actual service container. This enables service containers to know which user performs the access without the need to parse and validate the JWT token issued by `cisidm` for every successful authentication.

### Echoserver



## First Start

Finally it's time to start our services by calling the following command from the project directory (the one that contains the docker-compose file)

```
docker-compose up -d
```

:::tip Admin User

The first user that registers it self on the web-interface will be granted admin privileges (i.e. the `idm_superuser` role is assigned). Note that cisidm permits at least one registration even if `registration` mode is set to `token` or `disabled`.

Open `https://account.example.com/register` in your browser to create your initial admin user.
It's also possible to use the `idmctl` cli utility:

```bash
idmctl register my-username --password my-password
```

:::

<br />

---

<br />

:::tip Congratulations
You just finished setting up cisidm with a HTTPS enabled reverse proxy that will now protect your services using Proxy/Forward Authentication.

Now it's time to check the [User and Role Administration Guide](./user-role-management.md) or the [Command Line Reference](./cli-reference.md).
:::
