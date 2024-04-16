import{_ as p,a as d}from"./03_login-kWPkIgdr.js";import{_ as m,r as l,o as k,c as b,a as n,d as e,b as s,w as a,e as r}from"./app-wNEdNvye.js";const v="/cis-idm/assets/02_rallly-WuxM1DZp.png",h={},y=n("h1",{id:"openid-connect-setup",tabindex:"-1"},[n("a",{class:"header-anchor",href:"#openid-connect-setup","aria-hidden":"true"},"#"),e(" OpenID Connect Setup")],-1),g=n("code",null,"cisidm",-1),f={href:"https://dexidp.io/",target:"_blank",rel:"noopener noreferrer"},_=n("code",null,"cisidm",-1),x={href:"https://dexidp.io/docs/connectors/authproxy/",target:"_blank",rel:"noopener noreferrer"},w=n("strong",null,"AuthProxy",-1),I=n("p",null,[e("In this mode, OIDC clients ask DexIdp for authentication/authorization which in turn takes authenticated user information from HTTP request headers that have been added by "),n("code",null,"cisidm"),e(" during the forward-authentication.")],-1),D=n("p",null,"For self-hosted environments it is also possible to configure DexIdp to skip the consent screen and directly redirect the user to the requested application providing a seamless Single-Sign-On experience.",-1),R={href:"https://rallly.co",target:"_blank",rel:"noopener noreferrer"},S={href:"https://js.wiki/",target:"_blank",rel:"noopener noreferrer"},C={href:"https://www.getoutline.com/",target:"_blank",rel:"noopener noreferrer"},O={href:"https://nextcloud.com",target:"_blank",rel:"noopener noreferrer"},T=n("h2",{id:"contents",tabindex:"-1"},[n("a",{class:"header-anchor",href:"#contents","aria-hidden":"true"},"#"),e(" Contents")],-1),E={class:"table-of-contents"},P=n("h2",{id:"example-setup",tabindex:"-1"},[n("a",{class:"header-anchor",href:"#example-setup","aria-hidden":"true"},"#"),e(" Example Setup")],-1),L=n("p",null,[e("This section presents an example setup of using "),n("code",null,"cisidm"),e(" together with DexIdp to enable OIDC support. We will also deploy Rallly and configure it to use Single-Sign-On usig OIDC.")],-1),A={href:"https://github.com/tierklinik-dobersberg/cis-idm/tree/main/examples/oidc",target:"_blank",rel:"noopener noreferrer"},q=n("div",{class:"custom-container warning"},[n("p",{class:"custom-container-title"},"Network Mode"),n("p",null,[e("Do to the nature of how OAuth2 and OIDC work (with a couple of redirects), the following example requires that all services in the docker-compose.yml are configured to run in "),n("code",null,"network_mode: host"),e(".")]),n("p",null,"If you have public DNS entries it would work with normal docker networking as well!")],-1),N=n("p",null,"Let's create the different configuration files we need:",-1),U=n("div",{class:"language-yaml line-numbers-mode","data-ext":"yml"},[n("pre",{class:"language-yaml"},[n("code",null,[n("span",{class:"token key atrule"},"version"),n("span",{class:"token punctuation"},":"),e(),n("span",{class:"token string"},'"3"'),e(`

`),n("span",{class:"token comment"},"# Volume definitions for our databases"),e(`
`),n("span",{class:"token key atrule"},"volumes"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token key atrule"},"dexdb"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token key atrule"},"idmdb"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token key atrule"},"postgresdb"),n("span",{class:"token punctuation"},":"),e(`

`),n("span",{class:"token key atrule"},"services"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token comment"},"# We use traefik as our reverse proxy. It will get configured by docker container"),e(`
  `),n("span",{class:"token comment"},"# labels."),e(`
  `),n("span",{class:"token key atrule"},"traefik"),n("span",{class:"token punctuation"},":"),e(`
    `),n("span",{class:"token key atrule"},"image"),n("span",{class:"token punctuation"},":"),e(),n("span",{class:"token string"},'"traefik:v2.10"'),e(`
    `),n("span",{class:"token key atrule"},"restart"),n("span",{class:"token punctuation"},":"),e(" unless"),n("span",{class:"token punctuation"},"-"),e(`stopped
    `),n("span",{class:"token key atrule"},"network_mode"),n("span",{class:"token punctuation"},":"),e(` host
    `),n("span",{class:"token key atrule"},"command"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"--providers.docker=true"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"--providers.docker.exposedbydefault=false"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"--entrypoints.http.address=:80"'),e(`

    `),n("span",{class:"token key atrule"},"ports"),n("span",{class:"token punctuation"},":"),e(`
    `),n("span",{class:"token comment"},"# HTTP port, this example does not support HTTPS since it would make testing"),e(`
    `),n("span",{class:"token comment"},"# much harder."),e(`
    `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token key atrule"},"target"),n("span",{class:"token punctuation"},":"),e(),n("span",{class:"token number"},"80"),e(`
      `),n("span",{class:"token key atrule"},"published"),n("span",{class:"token punctuation"},":"),e(),n("span",{class:"token number"},"80"),e(`
      `),n("span",{class:"token key atrule"},"protocol"),n("span",{class:"token punctuation"},":"),e(` tcp
      `),n("span",{class:"token key atrule"},"mode"),n("span",{class:"token punctuation"},":"),e(` host

    `),n("span",{class:"token comment"},"# Traefik needs access to the docker daemon so it can watch for containers"),e(`
    `),n("span",{class:"token comment"},"# and auto-configure itself based on container labels."),e(`
    `),n("span",{class:"token key atrule"},"volumes"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"/var/run/docker.sock:/var/run/docker.sock"'),e(`

  `),n("span",{class:"token comment"},"# Rallly needs a postgres database so we provide one here"),e(`
  `),n("span",{class:"token key atrule"},"postgres"),n("span",{class:"token punctuation"},":"),e(`
    `),n("span",{class:"token key atrule"},"image"),n("span",{class:"token punctuation"},":"),e(" postgres"),n("span",{class:"token punctuation"},":"),e(`latest
    `),n("span",{class:"token key atrule"},"restart"),n("span",{class:"token punctuation"},":"),e(` always
    `),n("span",{class:"token key atrule"},"network_mode"),n("span",{class:"token punctuation"},":"),e(` host
    `),n("span",{class:"token key atrule"},"environment"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token key atrule"},"POSTGRES_USER"),n("span",{class:"token punctuation"},":"),e(` postgres
      `),n("span",{class:"token key atrule"},"POSTGRES_PASSWORD"),n("span",{class:"token punctuation"},":"),e(` password
      `),n("span",{class:"token key atrule"},"POSTGRES_DB"),n("span",{class:"token punctuation"},":"),e(` rallly
    `),n("span",{class:"token key atrule"},"volumes"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(" postgresdb"),n("span",{class:"token punctuation"},":"),e(`/var/lib/postgresql/data

  `),n("span",{class:"token comment"},"# Service definition for dex"),e(`
  `),n("span",{class:"token key atrule"},"dex"),n("span",{class:"token punctuation"},":"),e(`
    `),n("span",{class:"token key atrule"},"image"),n("span",{class:"token punctuation"},":"),e(" dexidp/dex"),n("span",{class:"token punctuation"},":"),e(`latest
    `),n("span",{class:"token key atrule"},"network_mode"),n("span",{class:"token punctuation"},":"),e(` host
    `),n("span",{class:"token key atrule"},"depends_on"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(` cisidm
    `),n("span",{class:"token key atrule"},"volumes"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(" ./dex.yml"),n("span",{class:"token punctuation"},":"),e(`/etc/dex/config.docker.yaml
      `),n("span",{class:"token punctuation"},"-"),e(" dexdb"),n("span",{class:"token punctuation"},":"),e(`/var/dex
    `),n("span",{class:"token key atrule"},"labels"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.enable=true"'),e(`

      `),n("span",{class:"token comment"},"# Access to dex itself must not be protected by cisidm as connected"),e(`
      `),n("span",{class:"token comment"},"# OIDC clients need to be able to access them."),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.routers.dex.rule=Host(`oidc.example.intern`)"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.routers.dex.entrypoints=http"'),e(`
      `),n("span",{class:"token comment"},'# The dex-container does not have an "EXPOSE" port configuration so we'),e(`
      `),n("span",{class:"token comment"},"# manually need to tell traefik which port dex is using."),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.services.dex.loadbalancer.server.port=5556"'),e(`

      `),n("span",{class:"token comment"},"# But we need to enable the forward-auth middleware (configured by the cisidm"),e(`
      `),n("span",{class:"token comment"},"# container) for the dex callback URI."),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.routers.dexcallback.rule=Host(`oidc.example.intern`) && PathPrefix(`/callback/`)"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.routers.dexcallback.middlewares=auth"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.routers.dexcallback.entrypoints=http"'),e(`

  `),n("span",{class:"token comment"},"# Configuration of rallly as an example OIDC client."),e(`
  `),n("span",{class:"token key atrule"},"rallly"),n("span",{class:"token punctuation"},":"),e(`
    `),n("span",{class:"token key atrule"},"image"),n("span",{class:"token punctuation"},":"),e(" lukevella/rallly"),n("span",{class:"token punctuation"},":"),e(`latest
    `),n("span",{class:"token key atrule"},"network_mode"),n("span",{class:"token punctuation"},":"),e(` host
    `),n("span",{class:"token key atrule"},"restart"),n("span",{class:"token punctuation"},":"),e(` always
    `),n("span",{class:"token key atrule"},"hostname"),n("span",{class:"token punctuation"},":"),e(` rallly.example.intern
    `),n("span",{class:"token key atrule"},"ports"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(" 3000"),n("span",{class:"token punctuation"},":"),n("span",{class:"token number"},"3000"),e(`
    `),n("span",{class:"token key atrule"},"depends_on"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(` postgres
    `),n("span",{class:"token key atrule"},"environment"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(" DATABASE_URL=postgres"),n("span",{class:"token punctuation"},":"),e("//postgres"),n("span",{class:"token punctuation"},":"),e("password@127.0.0.1"),n("span",{class:"token punctuation"},":"),e(`5432/rallly
      `),n("span",{class:"token punctuation"},"-"),e(" SECRET_PASSWORD=a"),n("span",{class:"token punctuation"},"-"),e("random"),n("span",{class:"token punctuation"},"-"),e("string"),n("span",{class:"token punctuation"},"-"),e("with"),n("span",{class:"token punctuation"},"-"),e("32"),n("span",{class:"token punctuation"},"-"),e("chars"),n("span",{class:"token punctuation"},"-"),n("span",{class:"token number"},"-1"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(" NEXT_PUBLIC_BASE_URL=http"),n("span",{class:"token punctuation"},":"),e(`//rallly.example.intern
      `),n("span",{class:"token punctuation"},"-"),e(` NOREPLY_MAIL=noreply@example.intern
      `),n("span",{class:"token punctuation"},"-"),e(` SUPPORT_MAIL=noreply@example.intern
      `),n("span",{class:"token punctuation"},"-"),e(` SMTP_HOST=smtp.example.intern
      `),n("span",{class:"token punctuation"},"-"),e(` SMTP_PORT=465
      `),n("span",{class:"token punctuation"},"-"),e(` SMTP_SECURE=true
      `),n("span",{class:"token punctuation"},"-"),e(` SMTP_USER=noreply@example.intern
      `),n("span",{class:"token punctuation"},"-"),e(` SMTP_PWD=password
      `),n("span",{class:"token punctuation"},"-"),e(` OIDC_NAME=Example Inc

      `),n("span",{class:"token comment"},"# Dex is running on the same domain as idm but under the /dex"),e(`
      `),n("span",{class:"token comment"},"# path"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(" OIDC_DISCOVERY_URL=http"),n("span",{class:"token punctuation"},":"),e("//oidc.example.intern/.well"),n("span",{class:"token punctuation"},"-"),e("known/openid"),n("span",{class:"token punctuation"},"-"),e(`configuration

      `),n("span",{class:"token comment"},"# The following values must match those from dex.yml staticClients"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(` OIDC_CLIENT_ID=rallly
      `),n("span",{class:"token punctuation"},"-"),e(" OIDC_CLIENT_SECRET=some"),n("span",{class:"token punctuation"},"-"),e("secure"),n("span",{class:"token punctuation"},"-"),e("random"),n("span",{class:"token punctuation"},"-"),e(`string
    `),n("span",{class:"token key atrule"},"labels"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token comment"},"# There's no need to protect access to rallly with the forward-auth"),e(`
      `),n("span",{class:"token comment"},"# middleware since authentication is already handeled via OIDC"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.enable=true"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.routers.rallly.rule=Host(`rallly.example.intern`)"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.routers.rallly.entrypoints=http"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.services.rallly.loadbalancer.server.port=3000"'),e(`

  `),n("span",{class:"token comment"},"# Our cisidm container"),e(`
  `),n("span",{class:"token key atrule"},"cisidm"),n("span",{class:"token punctuation"},":"),e(`
    `),n("span",{class:"token key atrule"},"build"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token key atrule"},"context"),n("span",{class:"token punctuation"},":"),e(` ../../
    `),n("span",{class:"token key atrule"},"restart"),n("span",{class:"token punctuation"},":"),e(" unless"),n("span",{class:"token punctuation"},"-"),e(`stopped
    `),n("span",{class:"token key atrule"},"network_mode"),n("span",{class:"token punctuation"},":"),e(` host
    `),n("span",{class:"token key atrule"},"volumes"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(" ./idm.hcl"),n("span",{class:"token punctuation"},":"),e("/etc/config.hcl"),n("span",{class:"token punctuation"},":"),e(`ro
      `),n("span",{class:"token punctuation"},"-"),e(" idmdb"),n("span",{class:"token punctuation"},":"),e(`/var/idm
    `),n("span",{class:"token key atrule"},"labels"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.enable=true"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.routers.idm.rule=Host(`account.example.intern`)"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.routers.idm.entrypoints=http"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.services.cisidm.loadbalancer.server.port=8080"'),e(`

      `),n("span",{class:"token comment"},"# Configure the forward auth middleware for traefik and also set a list"),e(`
      `),n("span",{class:"token comment"},"# of headers that cisidm will add to upstream requests."),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.middlewares.auth.forwardauth.address=http://127.0.0.1:8080/validate"'),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},'"traefik.http.middlewares.auth.forwardauth.authResponseHeaders=X-Remote-User,X-Remote-User-ID,X-Remote-Mail,X-Remote-Mail-Verified,X-Remote-Avatar-URL,X-Remote-Role,X-Remote-User-Display-Name"'),e(`
    `),n("span",{class:"token key atrule"},"environment"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token key atrule"},"CONFIG_FILE"),n("span",{class:"token punctuation"},":"),e(),n("span",{class:"token string"},'"/etc/config.hcl"'),e(`
`)])]),n("div",{class:"line-numbers","aria-hidden":"true"},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),H=n("div",{class:"language-hcl line-numbers-mode","data-ext":"hcl"},[n("pre",{class:"language-hcl"},[n("code",null,[n("span",{class:"token property"},"database_url"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token string"},'"file:/var/idm/idm.db"'),e(`

`),n("span",{class:"token property"},"registration"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token string"},'"disabled"'),e(`

`),n("span",{class:"token keyword"},"server"),e(),n("span",{class:"token punctuation"},"{"),e(`
    `),n("span",{class:"token comment"},"# We're not running under HTTPS in this example"),e(`
    `),n("span",{class:"token comment"},'# so cookies must not have the "secure" attribute'),e(`
    `),n("span",{class:"token comment"},"# set"),e(`
    `),n("span",{class:"token property"},"secure_cookies"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token boolean"},"false"),e(`

    `),n("span",{class:"token property"},"domain"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token string"},'"example.intern"'),e(`
    `),n("span",{class:"token property"},"public_listener"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token string"},'":8080"'),e(`
    `),n("span",{class:"token property"},"admin_listener"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token string"},'":8081"'),e(`
    `),n("span",{class:"token property"},"allowed_origins"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token punctuation"},"["),e(`
        `),n("span",{class:"token string"},'"http://example.intern"'),e(`,
        `),n("span",{class:"token string"},'"http://*.example.intern"'),e(`
    `),n("span",{class:"token punctuation"},"]"),e(`
    `),n("span",{class:"token property"},"allowed_redirects"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token punctuation"},"["),e(`
        `),n("span",{class:"token string"},'".example.intern"'),e(`,
        `),n("span",{class:"token string"},'"example.intern"'),e(`
    `),n("span",{class:"token punctuation"},"]"),e(`
`),n("span",{class:"token punctuation"},"}"),e(`

`),n("span",{class:"token keyword"},"jwt"),e(),n("span",{class:"token punctuation"},"{"),e(`
    `),n("span",{class:"token property"},"secret"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token string"},'"some-random-secret"'),e(`
`),n("span",{class:"token punctuation"},"}"),e(`

`),n("span",{class:"token keyword"},"ui"),e(),n("span",{class:"token punctuation"},"{"),e(`
    `),n("span",{class:"token property"},"site_name"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token string"},'"Example Inc"'),e(`
    `),n("span",{class:"token property"},"public_url"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token string"},'"http://account.example.intern"'),e(`
`),n("span",{class:"token punctuation"},"}"),e(`

`),n("span",{class:"token keyword"},"forward_auth"),e(),n("span",{class:"token punctuation"},"{"),e(`
    `),n("span",{class:"token property"},"allow_cors_preflight"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token boolean"},"true"),e(`
`),n("span",{class:"token punctuation"},"}"),e(`

`),n("span",{class:"token keyword"},"policies"),e(),n("span",{class:"token punctuation"},"{"),e(`
    `),n("span",{class:"token comment"},"# A default policy to allow all requests to upstream services"),e(`
    `),n("span",{class:"token comment"},"# the there's an authenticated user (i.e. input.subject is available)"),e(`
    policy `),n("span",{class:"token string"},'"default"'),e(),n("span",{class:"token punctuation"},"{"),e(`
        `),n("span",{class:"token property"},"content"),e(),n("span",{class:"token punctuation"},"="),e(),n("span",{class:"token heredoc string"},`<<EOT
        package cisidm.forward_auth

        import rego.v1

        default allow := false

        allow if {
            input.subject
        }
        EOT`),e(`
    `),n("span",{class:"token punctuation"},"}"),e(`
`),n("span",{class:"token punctuation"},"}"),e(`
`)])]),n("div",{class:"line-numbers","aria-hidden":"true"},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),M=n("div",{class:"language-yaml line-numbers-mode","data-ext":"yml"},[n("pre",{class:"language-yaml"},[n("code",null,[n("span",{class:"token comment"},"# The issuer for OIDC access and refresh tokens"),e(`
`),n("span",{class:"token key atrule"},"issuer"),n("span",{class:"token punctuation"},":"),e(" http"),n("span",{class:"token punctuation"},":"),e(`//oidc.example.intern

`),n("span",{class:"token comment"},"# Where dexIdp should store it's data, we use sqlite3 for simplicity here"),e(`
`),n("span",{class:"token key atrule"},"storage"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token key atrule"},"type"),n("span",{class:"token punctuation"},":"),e(` sqlite3
  `),n("span",{class:"token key atrule"},"config"),n("span",{class:"token punctuation"},":"),e(`
    `),n("span",{class:"token key atrule"},"file"),n("span",{class:"token punctuation"},":"),e(` /var/dex/dex.db

`),n("span",{class:"token comment"},"# Where dexidp should listen for incoming requests"),e(`
`),n("span",{class:"token key atrule"},"web"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token key atrule"},"http"),n("span",{class:"token punctuation"},":"),e(),n("span",{class:"token string"},'"0.0.0.0:5556"'),e(`

`),n("span",{class:"token comment"},"# Token expiries"),e(`
`),n("span",{class:"token key atrule"},"expiry"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token key atrule"},"deviceRequests"),n("span",{class:"token punctuation"},":"),e(` 5m
  `),n("span",{class:"token key atrule"},"signingKeys"),n("span",{class:"token punctuation"},":"),e(` 6h
  `),n("span",{class:"token key atrule"},"idTokens"),n("span",{class:"token punctuation"},":"),e(` 24h
  `),n("span",{class:"token key atrule"},"authRequests"),n("span",{class:"token punctuation"},":"),e(` 24h

`),n("span",{class:"token comment"},"# Configuration for oauth2"),e(`
`),n("span",{class:"token key atrule"},"oauth2"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token comment"},"# A list of response-types dexidp should support,"),e(`
  `),n("span",{class:"token comment"},"# it's best to keep this list as it is."),e(`
  `),n("span",{class:"token key atrule"},"responseTypes"),n("span",{class:"token punctuation"},":"),e(`
    `),n("span",{class:"token punctuation"},"-"),e(` code
    `),n("span",{class:"token punctuation"},"-"),e(` token
    `),n("span",{class:"token punctuation"},"-"),e(` id_token

  `),n("span",{class:"token comment"},'# For self-hosted environments you will likely always "approve" OIDC clients'),e(`
  `),n("span",{class:"token comment"},"# to use your login data so we can instruct dex to assume approval and skip "),e(`
  `),n("span",{class:"token comment"},"# the consent screen altogether"),e(`
  `),n("span",{class:"token key atrule"},"skipApprovalScreen"),n("span",{class:"token punctuation"},":"),e(),n("span",{class:"token boolean important"},"true"),e(`
  `),n("span",{class:"token key atrule"},"alwaysShowLoginScreen"),n("span",{class:"token punctuation"},":"),e(),n("span",{class:"token boolean important"},"false"),e(`

`),n("span",{class:"token comment"},"# A list of upstream identity connectors for dex. We use the authproxy"),e(`
`),n("span",{class:"token comment"},"# connector for cisidm. That is, any request to the dex OIDC callback"),e(`
`),n("span",{class:"token comment"},"# will be authenticated by cisidm and dex will use the user-information added"),e(`
`),n("span",{class:"token comment"},"# to the request headers."),e(`
`),n("span",{class:"token key atrule"},"connectors"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token key atrule"},"type"),n("span",{class:"token punctuation"},":"),e(` authproxy
    `),n("span",{class:"token key atrule"},"id"),n("span",{class:"token punctuation"},":"),e(` cisidm
    `),n("span",{class:"token key atrule"},"name"),n("span",{class:"token punctuation"},":"),e(" CIS"),n("span",{class:"token punctuation"},"-"),e(`IDM
    `),n("span",{class:"token key atrule"},"config"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token key atrule"},"emailHeader"),n("span",{class:"token punctuation"},":"),e(" X"),n("span",{class:"token punctuation"},"-"),e("Remote"),n("span",{class:"token punctuation"},"-"),e(`Mail
      `),n("span",{class:"token key atrule"},"groupHeader"),n("span",{class:"token punctuation"},":"),e(" X"),n("span",{class:"token punctuation"},"-"),e("Remote"),n("span",{class:"token punctuation"},"-"),e(`Role
      `),n("span",{class:"token key atrule"},"userHeader"),n("span",{class:"token punctuation"},":"),e(" X"),n("span",{class:"token punctuation"},"-"),e("Remote"),n("span",{class:"token punctuation"},"-"),e(`User
      `),n("span",{class:"token key atrule"},"userIDHeader"),n("span",{class:"token punctuation"},":"),e(" X"),n("span",{class:"token punctuation"},"-"),e("Remote"),n("span",{class:"token punctuation"},"-"),e("User"),n("span",{class:"token punctuation"},"-"),e(`ID

`),n("span",{class:"token comment"},"# A list of static OIDC clients. We just add an example entry for rallly here."),e(`
`),n("span",{class:"token key atrule"},"staticClients"),n("span",{class:"token punctuation"},":"),e(`
  `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token key atrule"},"id"),n("span",{class:"token punctuation"},":"),e(` rallly
    `),n("span",{class:"token key atrule"},"redirectURIs"),n("span",{class:"token punctuation"},":"),e(`
      `),n("span",{class:"token punctuation"},"-"),e(),n("span",{class:"token string"},"'http://rallly.example.intern/api/auth/callback/oidc'"),e(`
    `),n("span",{class:"token key atrule"},"name"),n("span",{class:"token punctuation"},":"),e(` Rallly
    `),n("span",{class:"token key atrule"},"secret"),n("span",{class:"token punctuation"},":"),e(" some"),n("span",{class:"token punctuation"},"-"),e("secure"),n("span",{class:"token punctuation"},"-"),e("random"),n("span",{class:"token punctuation"},"-"),e(`string
`)])]),n("div",{class:"line-numbers","aria-hidden":"true"},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),W=r(`<p>In this example we&#39;re using the domains <code>account.example.intern</code>, <code>oidc.example.intern</code> and <code>rallly.example.intern</code>. For testing purposes you should add them to your <code>/etc/hosts</code> file. Add the following entry:</p><div class="language-text line-numbers-mode" data-ext="text"><pre class="language-text"><code>127.0.0.1   oidc.example.intern account.example.intern rallly.example.intern example.intern
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><h2 id="starting-the-example" tabindex="-1"><a class="header-anchor" href="#starting-the-example" aria-hidden="true">#</a> Starting the example</h2><p>First, clone the cis-idm repo to your local machine and enter the example directory:</p><div class="language-text line-numbers-mode" data-ext="text"><pre class="language-text"><code>git clone https://github.com/tierklinik-dobersberg/cis-idm

cd cis-idm/examples/oidc
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>Before we can run <code>docker-compose up</code> to bring up the example we first need to build the <code>cisidm</code> container (we could also use the official container image from the ghcr.io registry):</p><div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code><span class="token function">docker-compose</span> build
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><p>Once the container build is done we can start up the example:</p><div class="language-text line-numbers-mode" data-ext="text"><pre class="language-text"><code>docker-compose up
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><p>After a couple of minutes (depending on the required download size for the docker containers), everything should be up and running.</p><p>Now, we need to create a new cisidm admin user. Even though registration is set to <code>disabled</code> in <code>idm.hcl</code>, cisidm will at least permit the very first user to register and will assign administrative priviledges.</p>`,11),X={href:"http://account.example.intern/registration",target:"_blank",rel:"noopener noreferrer"},B=n("p",null,[n("img",{src:p,alt:"Registration"})],-1),F=n("p",null,"Once you completed the registration you will be redirected to your profile page in cisidm.",-1),G=n("div",{class:"custom-container tip"},[n("p",{class:"custom-container-title"},"Logout"),n("p",null,'For testing and demonstration purposes you should now click on "Logout".')],-1),j={href:"http://rallly.example.intern",target:"_blank",rel:"noopener noreferrer"},V=n("p",null,[n("img",{src:v,alt:"Rallly"})],-1),Y=n("p",null,[e(`Finnaly, you just need to click on "Login with Example Inc". Rallly will now redirect to DexIDP for authentication. Since you're not logged in `),n("code",null,"cisidm"),e(" will intercept the redirect to the DexIDP "),n("code",null,"/callback"),e(" and redirect you to the login page.")],-1),z=n("p",null,'Now enter your credentials and press "Login".',-1),J=n("p",null,[n("img",{src:d,alt:"Login"})],-1),K=r(`<p>Once you finish the login process, cisidm will redirect you to the DexIDP callback endpoint but this time, the request is authenticated and DexIDP will see the remote user headers.</p><p>Since the request is authenticated, DexIDP will issue a new access/refresh token and redirect you back to Rallly.</p><p>Sounds complicated? <strong>It is!</strong></p><p>But as a user you don&#39;t really see any of those redirect. You just click on &quot;Login with Example Inc&quot;, enter your cisidm credentials and immediately are logged into Rallly.</p><h3 id="cleanup" tabindex="-1"><a class="header-anchor" href="#cleanup" aria-hidden="true">#</a> Cleanup</h3><p>Finally, hit <code>CTRL-C</code> in the terminal that is running <code>docker-compose up</code>. Once all containers are stopped, execute the following to clean-up everything (including the databases!):</p><div class="language-bash" data-ext="sh"><pre class="language-bash"><code><span class="token function">docker-compose</span> down <span class="token parameter variable">-v</span>
</code></pre></div>`,7);function Z(Q,$){const t=l("ExternalLinkIcon"),i=l("router-link"),o=l("CodeGroupItem"),u=l("CodeGroup"),c=l("center");return k(),b("div",null,[y,n("p",null,[e("While "),g,e(" does not implement OIDC itself it is very easy to add support for Open-ID-Connect using the "),n("a",f,[e("DexIdp"),s(t)]),e(" project. It integrates with "),_,e(" using the "),n("a",x,[w,s(t)]),e(" provider.")]),I,D,n("p",null,[e("We have tested OIDC support based on DexIdp with "),n("a",R,[e("Rallly"),s(t)]),e(", "),n("a",S,[e("WikiJS"),s(t)]),e(", "),n("a",C,[e("Outline"),s(t)]),e(" and "),n("a",O,[e("Nextcloud"),s(t)]),e(". It's expected that every OIDC client will work with this setup.")]),T,n("nav",E,[n("ul",null,[n("li",null,[s(i,{to:"#contents"},{default:a(()=>[e("Contents")]),_:1})]),n("li",null,[s(i,{to:"#example-setup"},{default:a(()=>[e("Example Setup")]),_:1})]),n("li",null,[s(i,{to:"#starting-the-example"},{default:a(()=>[e("Starting the example")]),_:1}),n("ul",null,[n("li",null,[s(i,{to:"#cleanup"},{default:a(()=>[e("Cleanup")]),_:1})])])])])]),P,L,n("p",null,[n("strong",null,[e("The whole example can be found in the "),n("a",A,[e("GitHub repository"),s(t)]),e(".")])]),q,N,s(u,null,{default:a(()=>[s(o,{title:"docker-compose.yml"},{default:a(()=>[U]),_:1}),s(o,{title:"idm.hcl"},{default:a(()=>[H]),_:1}),s(o,{title:"dex.yml"},{default:a(()=>[M]),_:1})]),_:1}),W,n("p",null,[e("Open the following URL in your web-browser: "),n("a",X,[e("http://account.example.intern/registration"),s(t)]),e(".")]),s(c,null,{default:a(()=>[B]),_:1}),F,G,n("p",null,[e("Next, we can open up Rallly: "),n("a",j,[e("http://rallly.example.intern"),s(t)]),e(":")]),s(c,null,{default:a(()=>[V]),_:1}),Y,z,s(c,null,{default:a(()=>[J]),_:1}),K])}const sn=m(h,[["render",Z],["__file","setup-oidc.html.vue"]]);export{sn as default};
