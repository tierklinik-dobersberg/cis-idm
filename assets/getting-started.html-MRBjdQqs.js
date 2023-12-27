import{_ as u,r as i,o as p,c as d,a as e,b as s,d as n,w as a,e as r}from"./app-IXazZNe3.js";const m={},k=e("h1",{id:"getting-started",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#getting-started","aria-hidden":"true"},"#"),n(" Getting Started")],-1),v=e("p",null,[n("Welcome to the Getting-Started Guide for "),e("code",null,"cisidm"),n(". On this page we will walk through a simple docker-compose based setup to get you started.")],-1),h=e("h2",{id:"overview",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#overview","aria-hidden":"true"},"#"),n(" Overview")],-1),b=e("p",null,[n("The easiest way to get started using "),e("code",null,"cisidm"),n(" is to deploy it and the database (rqlite) using docker-compose. In this guide we will setup a docker deployment with the following service:")],-1),g={href:"https://rqlite.io/",target:"_blank",rel:"noopener noreferrer"},f=e("strong",null,"rqlite",-1),y=e("br",null,null,-1),w={href:"https://doc.traefik.io/traefik/",target:"_blank",rel:"noopener noreferrer"},_=e("strong",null,"traefik",-1),q=e("br",null,null,-1),T={href:"https://github.com/tierklinik-dobersberg/cis-idm",target:"_blank",rel:"noopener noreferrer"},x=e("strong",null,"cisidm",-1),R=e("br",null,null,-1),I={href:"https://gcr.io/google_containers/echoserver:1.4",target:"_blank",rel:"noopener noreferrer"},C=e("strong",null,"echoserver",-1),L=e("br",null,null,-1),P=r(`<div class="custom-container warning"><p class="custom-container-title">Customization Required</p><p>Please make sure to update the configurations in this guide to match your domains and e-mail addresses! You will also need DNS to be setup correctly (i.e. your sub-domains pointing to your server/IP address).</p></div><h2 id="setup-the-project" tabindex="-1"><a class="header-anchor" href="#setup-the-project" aria-hidden="true">#</a> Setup the Project</h2><p>Prepare the directory structure for our deployment.</p><div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code><span class="token function">mkdir</span> <span class="token parameter variable">-p</span> cisidm-demo/config<span class="token punctuation">;</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><p>The <code>cisidm-config/config</code> folder will hold the configuration files for <code>cisidm</code> and <code>traefik</code>. Our <code>docker-compose.yml</code> file will be placed directly in <code>cisidm-config</code>.</p><h2 id="configuration-file" tabindex="-1"><a class="header-anchor" href="#configuration-file" aria-hidden="true">#</a> Configuration File</h2><p>Create the configuration file for <code>cisidm</code> in <code>cisidm-demo/config/idm-config.yml</code>:</p>`,7),U=e("div",{class:"language-yaml line-numbers-mode","data-ext":"yml"},[e("pre",{class:"language-yaml"},[e("code",null,[e("span",{class:"token comment"},"# The domain that cisidm is going to protect. If you plan on deploying other services using sub-domains make sure"),n(`
`),e("span",{class:"token comment"},"# to set the this field to the top-level domain. Otherwise, browser will not send the access token cookie for sub-domains."),n(`
`),e("span",{class:"token key atrule"},"domain"),e("span",{class:"token punctuation"},":"),n(` example.com

`),e("span",{class:"token comment"},"# Whether or not the cookie should only be sent for HTTPS. It's best to keep this at true."),n(`
`),e("span",{class:"token key atrule"},"secureCookie"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token boolean important"},"true"),n(`

`),e("span",{class:"token comment"},"# The secret key used to sign the JWT access tokens. Choose some secure random string here."),n(`
`),e("span",{class:"token comment"},"# Note that if you ever change this value any access and refresh tokens will be invalidated and your users"),n(`
`),e("span",{class:"token comment"},"# will need to re-authenticate."),n(`
`),e("span",{class:"token key atrule"},"jwtSecret"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"some-secure-secret"'),n(`

`),e("span",{class:"token comment"},"# The address for the rqlite database"),n(`
`),e("span",{class:"token key atrule"},"rqliteURL"),e("span",{class:"token punctuation"},":"),n(" http"),e("span",{class:"token punctuation"},":"),n("//rqlite"),e("span",{class:"token punctuation"},":"),n(`4001/

`),e("span",{class:"token comment"},"# The default log level. Possible values are debug, info, warn and error"),n(`
`),e("span",{class:"token key atrule"},"logLevel"),e("span",{class:"token punctuation"},":"),n(` info

`),e("span",{class:"token comment"},"# A list of IP networks (in CIDR notation) or hostnames from which cisidm will trust the"),n(`
`),e("span",{class:"token comment"},"# X-Forwarded-For headers."),n(`
`),e("span",{class:"token key atrule"},"trustedNetworks"),e("span",{class:"token punctuation"},":"),n(`
  `),e("span",{class:"token punctuation"},"-"),n(` traefik

`),e("span",{class:"token comment"},"# Whether or not your users require a registration token. Setting this to false enables public registration."),n(`
`),e("span",{class:"token key atrule"},"registrationRequiresToken"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token boolean important"},"true"),n(`

`),e("span",{class:"token comment"},"# The name of your deployment."),n(`
`),e("span",{class:"token key atrule"},"siteName"),e("span",{class:"token punctuation"},":"),n(` Example Site

`),e("span",{class:"token comment"},"# A URL that will be used in email templates and on the self-service UI."),n(`
`),e("span",{class:"token key atrule"},"siteNameUrl"),e("span",{class:"token punctuation"},":"),n(" https"),e("span",{class:"token punctuation"},":"),n(`//example.com

`),e("span",{class:"token comment"},"# Configuration for the forward authentication support."),n(`
`),e("span",{class:"token key atrule"},"forwardAuth"),e("span",{class:"token punctuation"},":"),n(`
  `),e("span",{class:"token comment"},"# We require authentication for all sub-domains of example.com"),n(`
  `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token key atrule"},"url"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"http(s){0,1}://(.*).example.com"'),n(`
    `),e("span",{class:"token key atrule"},"required"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token boolean important"},"true"),n(`

`),e("span",{class:"token comment"},"# The public URL under which the login screen and self-service UI can be accessed."),n(`
`),e("span",{class:"token key atrule"},"publicURL"),e("span",{class:"token punctuation"},":"),n(" https"),e("span",{class:"token punctuation"},":"),n(`//account.example.com

`),e("span",{class:"token comment"},"# Configures how long access tokens issued to your user are valid. You can keep this relatively short since"),n(`
`),e("span",{class:"token comment"},"# users will also get a long-lived refresh token."),n(`
`),e("span",{class:"token key atrule"},"accessTokenTTL"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"1h"'),n(`

`),e("span",{class:"token comment"},"# cisidm automatically redirects your users back to the protected services after a successful login. To prevent open-redirect"),n(`
`),e("span",{class:"token comment"},"# attacks make sure to restrict the allowed redirects. It's best to just allow the top-level domain and any"),n(`
`),e("span",{class:"token comment"},"# sub-domains of your deployment."),n(`
`),e("span",{class:"token key atrule"},"allowedRedirects"),e("span",{class:"token punctuation"},":"),n(`
  `),e("span",{class:"token punctuation"},"-"),n(` example.com
  `),e("span",{class:"token punctuation"},"-"),n(` .example.com
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),N={class:"custom-container tip"},H=e("p",{class:"custom-container-title"},"TIP",-1),A=e("h2",{id:"create-the-docker-compose-file",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#create-the-docker-compose-file","aria-hidden":"true"},"#"),n(" Create the docker-compose file")],-1),S=e("p",null,"Next, we create the docker-compose file that will contain our service definitions and also tell treafik which sub-domains should be routed to which services and that we want cisidm to enforce authentication.",-1),D=e("p",null,"Here's the complete docker-compose file, we'll break it down and explain everything below:",-1),E=e("div",{class:"language-yaml line-numbers-mode","data-ext":"yml"},[e("pre",{class:"language-yaml"},[e("code",null,[e("span",{class:"token key atrule"},"version"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"3"'),n(`

`),e("span",{class:"token key atrule"},"volumes"),e("span",{class:"token punctuation"},":"),n(`
  `),e("span",{class:"token key atrule"},"db"),e("span",{class:"token punctuation"},":"),n(`
    `),e("span",{class:"token key atrule"},"driver"),e("span",{class:"token punctuation"},":"),n(` local

`),e("span",{class:"token key atrule"},"services"),e("span",{class:"token punctuation"},":"),n(`
  `),e("span",{class:"token key atrule"},"rqlite"),e("span",{class:"token punctuation"},":"),n(`
    `),e("span",{class:"token key atrule"},"image"),e("span",{class:"token punctuation"},":"),n(" rqlite/rqlite"),e("span",{class:"token punctuation"},":"),n(`latest
    `),e("span",{class:"token key atrule"},"hostname"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"12e94f6300a8"'),n(`
    `),e("span",{class:"token key atrule"},"restart"),e("span",{class:"token punctuation"},":"),n(" unless"),e("span",{class:"token punctuation"},"-"),n(`stopped
    `),e("span",{class:"token key atrule"},"environment"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token key atrule"},"HTTP_ADDR_ADV"),e("span",{class:"token punctuation"},":"),n(` localhost
    `),e("span",{class:"token key atrule"},"command"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"-on-disk=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"-node-id=1"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"-fk"'),n(`
    `),e("span",{class:"token key atrule"},"ports"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(" 4001"),e("span",{class:"token punctuation"},":"),e("span",{class:"token number"},"4001"),n(`
    `),e("span",{class:"token key atrule"},"volumes"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(" db"),e("span",{class:"token punctuation"},":"),n(`/rqlite/file

  `),e("span",{class:"token comment"},"# Traefik ###################################################"),n(`

  `),e("span",{class:"token key atrule"},"traefik"),e("span",{class:"token punctuation"},":"),n(`
    `),e("span",{class:"token key atrule"},"image"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"traefik:v2.10"'),n(`
    `),e("span",{class:"token key atrule"},"restart"),e("span",{class:"token punctuation"},":"),n(" unless"),e("span",{class:"token punctuation"},"-"),n(`stopped
    `),e("span",{class:"token key atrule"},"command"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--api.insecure=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--providers.docker=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--providers.docker.exposedbydefault=false"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--entrypoints.http.address=:80"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--entrypoints.secure.address=:443"'),n(`

      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--entrypoints.http.http.redirections.entrypoint.to=secure"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--entrypoints.http.http.redirections.entrypoint.scheme=https"'),n(`

      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--certificatesresolvers.resolver.acme.email=admin@example.com"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--certificatesresolvers.resolver.acme.storage=/secrets/acme.json"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--certificatesresolvers.resolver.acme.httpchallenge.entrypoint=http"'),n(`

      `),e("span",{class:"token comment"},"# disable for production"),n(`
      `),e("span",{class:"token comment"},'# - "--certificatesresolvers.resolver.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory"'),n(`

    `),e("span",{class:"token key atrule"},"ports"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token key atrule"},"target"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token number"},"80"),n(`
        `),e("span",{class:"token key atrule"},"published"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token number"},"80"),n(`
        `),e("span",{class:"token key atrule"},"protocol"),e("span",{class:"token punctuation"},":"),n(` tcp
        `),e("span",{class:"token key atrule"},"mode"),e("span",{class:"token punctuation"},":"),n(` host
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token key atrule"},"target"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token number"},"443"),n(`
        `),e("span",{class:"token key atrule"},"published"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token number"},"443"),n(`
        `),e("span",{class:"token key atrule"},"protocol"),e("span",{class:"token punctuation"},":"),n(` tcp
        `),e("span",{class:"token key atrule"},"mode"),e("span",{class:"token punctuation"},":"),n(` host

    `),e("span",{class:"token key atrule"},"volumes"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"/var/run/docker.sock:/var/run/docker.sock"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"./config:/secrets"'),n(`

  `),e("span",{class:"token comment"},"# IDM ############################################################"),n(`
  `),e("span",{class:"token key atrule"},"cisidm"),e("span",{class:"token punctuation"},":"),n(`
    `),e("span",{class:"token key atrule"},"image"),e("span",{class:"token punctuation"},":"),n(" ghcr.io/tierklinik"),e("span",{class:"token punctuation"},"-"),n("dobersberg/cis"),e("span",{class:"token punctuation"},"-"),n("idm"),e("span",{class:"token punctuation"},":"),n(`latest
    `),e("span",{class:"token key atrule"},"restart"),e("span",{class:"token punctuation"},":"),n(" unless"),e("span",{class:"token punctuation"},"-"),n(`stopped
    `),e("span",{class:"token key atrule"},"volumes"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(" ./config/idm"),e("span",{class:"token punctuation"},"-"),n("config.yml"),e("span",{class:"token punctuation"},":"),n("/etc/config.yml"),e("span",{class:"token punctuation"},":"),n(`ro
    `),e("span",{class:"token key atrule"},"labels"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.enable=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.idm.rule=Host(`account.dobersberg.vet`)"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.idm.entrypoints=web"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.idm.tls=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.idm.tls.certresolver=resolver"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.middlewares.auth.forwardauth.authResponseHeaders=X-Remote-User,X-Remote-User-ID,X-Remote-Mail,X-Remote-Mail-Verified,X-Remote-Avatar-URL,X-Remote-Role,X-Remote-User-Display-Name"'),n(`
    `),e("span",{class:"token key atrule"},"environment"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token key atrule"},"CONFIG_FILE"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"/etc/config.yml"'),n(`
    `),e("span",{class:"token key atrule"},"depends_on"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(` rqlite

  `),e("span",{class:"token key atrule"},"echoserver"),e("span",{class:"token punctuation"},":"),n(`
    `),e("span",{class:"token key atrule"},"image"),e("span",{class:"token punctuation"},":"),n(" gcr.io/google_containers/echoserver"),e("span",{class:"token punctuation"},":"),e("span",{class:"token number"},"1.4"),n(`
    `),e("span",{class:"token key atrule"},"labels"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.enable=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.echo.rule=Host(`app.example.com`)"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.echo.entrypoints=web"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.echo.tls=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.echo.tls.certresolver=resolver"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.echo.middlewares=auth"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.services.echo.loadbalancer.server.port=8080"'),n(`
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),X=r(`<div class="custom-container tip"><p class="custom-container-title">TIP</p><p>If you&#39;re already familiar with docker-compose you can skip over the explanation to the next section <a href="#start-the-project">here</a>.</p></div><p>Let&#39;s break it down a bit:</p><div class="language-yaml line-numbers-mode" data-ext="yml"><pre class="language-yaml"><code><span class="token key atrule">version</span><span class="token punctuation">:</span> <span class="token string">&quot;3&quot;</span>

<span class="token key atrule">volumes</span><span class="token punctuation">:</span>
  <span class="token key atrule">db</span><span class="token punctuation">:</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>This just tells docker-compose which file-version we&#39;re using and that we want a volume named <code>db</code> that will hold our rqlite database.</p><h3 id="rqlite-container" tabindex="-1"><a class="header-anchor" href="#rqlite-container" aria-hidden="true">#</a> RQLite Container</h3><div class="language-yaml line-numbers-mode" data-ext="yml"><pre class="language-yaml"><code><span class="token key atrule">rqlite</span><span class="token punctuation">:</span>
  <span class="token key atrule">image</span><span class="token punctuation">:</span> rqlite/rqlite<span class="token punctuation">:</span>latest
  <span class="token key atrule">hostname</span><span class="token punctuation">:</span> <span class="token string">&quot;rqlite-node-1&quot;</span>
  <span class="token key atrule">restart</span><span class="token punctuation">:</span> unless<span class="token punctuation">-</span>stopped
  <span class="token key atrule">environment</span><span class="token punctuation">:</span>
    <span class="token key atrule">HTTP_ADDR_ADV</span><span class="token punctuation">:</span> localhost
  <span class="token key atrule">command</span><span class="token punctuation">:</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;-on-disk=true&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;-node-id=1&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;-fk&quot;</span>
  <span class="token key atrule">ports</span><span class="token punctuation">:</span>
    <span class="token punctuation">-</span> 4001<span class="token punctuation">:</span><span class="token number">4001</span>
  <span class="token key atrule">volumes</span><span class="token punctuation">:</span>
    <span class="token punctuation">-</span> db<span class="token punctuation">:</span>/rqlite/file
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>This defines the rqlite database service. Note that rqlite needs a stable hostname that&#39;s why we force the hostname to <code>rqlite-node-1</code>:</p><p>We also pass a few parameters to the rqlite server:</p><ul><li><code>-on-disk=true</code>: Tells rqlite to store the sqlite database on the disk. By default, rqlite just keeps the database in memory as it&#39;s expected to run in cluster mode. For this simple setup, we just run one rqlite node so we want that data to be persisted to disk.</li><li><code>-node-id=1</code>: Each rqlite node in the cluster needs a stable ID. Even though we don&#39;t use clustering support here we still need to set a node ID.</li><li><code>-fk</code>: Enables foreign key support in the underlying sqlite database. This is required as the database schema from <code>cisidm</code> relies on cascading deletes.</li></ul><h3 id="traefik-container" tabindex="-1"><a class="header-anchor" href="#traefik-container" aria-hidden="true">#</a> Traefik Container</h3><div class="language-yaml line-numbers-mode" data-ext="yml"><pre class="language-yaml"><code><span class="token key atrule">traefik</span><span class="token punctuation">:</span>
  <span class="token key atrule">image</span><span class="token punctuation">:</span> <span class="token string">&quot;traefik:v2.10&quot;</span>
  <span class="token key atrule">restart</span><span class="token punctuation">:</span> unless<span class="token punctuation">-</span>stopped
  <span class="token key atrule">command</span><span class="token punctuation">:</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;--api.insecure=true&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;--providers.docker=true&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;--providers.docker.exposedbydefault=false&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;--entrypoints.http.address=:80&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;--entrypoints.secure.address=:443&quot;</span>

    <span class="token punctuation">-</span> <span class="token string">&quot;--entrypoints.http.http.redirections.entrypoint.to=secure&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;--entrypoints.http.http.redirections.entrypoint.scheme=https&quot;</span>

    <span class="token punctuation">-</span> <span class="token string">&quot;--certificatesresolvers.resolver.acme.email=admin@example.com&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;--certificatesresolvers.resolver.acme.storage=/secrets/acme.json&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;--certificatesresolvers.resolver.acme.httpchallenge.entrypoint=http&quot;</span>

    <span class="token key atrule">ports</span><span class="token punctuation">:</span>
      <span class="token punctuation">-</span> <span class="token key atrule">target</span><span class="token punctuation">:</span> <span class="token number">8080</span>
        <span class="token key atrule">published</span><span class="token punctuation">:</span> <span class="token number">8080</span>
        <span class="token key atrule">protocol</span><span class="token punctuation">:</span> tcp
        <span class="token key atrule">mode</span><span class="token punctuation">:</span> host
      <span class="token punctuation">-</span> <span class="token key atrule">target</span><span class="token punctuation">:</span> <span class="token number">80</span>
        <span class="token key atrule">published</span><span class="token punctuation">:</span> <span class="token number">80</span>
        <span class="token key atrule">protocol</span><span class="token punctuation">:</span> tcp
        <span class="token key atrule">mode</span><span class="token punctuation">:</span> host
      <span class="token punctuation">-</span> <span class="token key atrule">target</span><span class="token punctuation">:</span> <span class="token number">443</span>
        <span class="token key atrule">published</span><span class="token punctuation">:</span> <span class="token number">443</span>
        <span class="token key atrule">protocol</span><span class="token punctuation">:</span> tcp
        <span class="token key atrule">mode</span><span class="token punctuation">:</span> host

    <span class="token key atrule">volumes</span><span class="token punctuation">:</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;/var/run/docker.sock:/var/run/docker.sock&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;./config:/secrets&quot;</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>The above service snippet configures Traefik as our reverse proxy. Let&#39;s go through the arguments passed to Traefik:</p><ul><li><p><code>--api.insecure=true</code><br> This enables the management dashboard of Traefik. This is useful for debugging service configurations. The management dashboard will be accessible on port 8080:</p></li><li><p><code>--providers.docker=true</code>,<br><code>--providers.docker.exposedbydefault=false</code><br> This enables the docker provider so we can configure Traefik using container labels. This will be explained in more detail below.</p></li><li><p><code>--entrypoints.http.address=:80</code>,<br><code>--entrypoints.secure.address=:443</code><br> Tell traefik that we want to have two entrypoints, one on port 80 (named <code>http</code>) and one on port 443 (named <code>secure</code>):</p></li><li><p><code>--entrypoints.http.http.redirections.entrypoint.to=secure</code>,<br><code>--entrypoints.http.http.redirections.entrypoint.scheme=https</code><br> Configures a redirection from insecure, plain-text HTTP to TLS encrypted HTTPS. Whenever a user tries to access a service using HTTP, Traefik will automatically redirect the browser to the secure endpoint.</p></li><li><p><code>--certificatesresolvers.resolver.acme.email=admin@example.com</code>,<br><code>--certificatesresolvers.resolver.acme.storage=/secrets/acme.json</code>,<br><code>--certificatesresolvers.resolver.acme.httpchallenge.entrypoint=http</code><br> Configure a certificate resolver that uses ACME to request certificates from Let&#39;s Encrypt.</p></li></ul><h3 id="identity-management-server" tabindex="-1"><a class="header-anchor" href="#identity-management-server" aria-hidden="true">#</a> Identity Management Server</h3><div class="language-yaml line-numbers-mode" data-ext="yml"><pre class="language-yaml"><code><span class="token key atrule">cisidm</span><span class="token punctuation">:</span>
  <span class="token key atrule">image</span><span class="token punctuation">:</span> ghcr.io/tierklinik<span class="token punctuation">-</span>dobersberg/cis<span class="token punctuation">-</span>idm<span class="token punctuation">:</span>latest
  <span class="token key atrule">restart</span><span class="token punctuation">:</span> unless<span class="token punctuation">-</span>stopped
  <span class="token key atrule">volumes</span><span class="token punctuation">:</span>
    <span class="token punctuation">-</span> ./config/idm<span class="token punctuation">-</span>config.yml<span class="token punctuation">:</span>/etc/config.yml<span class="token punctuation">:</span>ro
  <span class="token key atrule">labels</span><span class="token punctuation">:</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;traefik.enable=true&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.routers.idm.rule=Host(\`account.example.com\`)&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.routers.idm.entrypoints=web&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.routers.idm.tls=true&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.routers.idm.tls.certresolver=resolver&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate&quot;</span>
    <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.middlewares.auth.forwardauth.authResponseHeaders=X-Remote-User,X-Remote-User-ID,X-Remote-Mail,X-Remote-Mail-Verified,X-Remote-Avatar-URL,X-Remote-Role,X-Remote-User-Display-Name&quot;</span>
  <span class="token key atrule">environment</span><span class="token punctuation">:</span>
    <span class="token key atrule">CONFIG_FILE</span><span class="token punctuation">:</span> <span class="token string">&quot;/etc/config.yml&quot;</span>
  <span class="token key atrule">depends_on</span><span class="token punctuation">:</span>
    <span class="token punctuation">-</span> rqlite
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>This configures the <code>cisidm</code> docker container, mounts the configuration file to <code>/etc/config.yml</code> and tells traefik that is should be reachable at <code>account.example.com</code>.</p><p>It also configures a new HTTP middleware <code>auth</code> that uses the forward-auth feature of traefik:</p><ul><li><p><code>traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate</code><br> This configures the Forward-Auth middleware (named <code>auth</code> here) to forward any HTTP request to cisidm to determine if the user is actually allowed to access the requested resource. If cisidm replies with a HTTP success status code (2xx) than traefik will forward the original request to the actual service container. If cisidm replies with an error code, traefik will immediately return the response from <code>cisidm</code> to the user. This is used to redirect the user to the login page in case the request is unauthenticated.</p></li><li><p><code>traefik.http.middlewares.auth.forwardauth.authResponseHeaders=....</code><br> When <code>cisidm</code> successfully authenticated a request, it will return a set of headers that contain information about the logged in user. With this setting, we tell traefik to forward those headers to the actual service container. This enables service containers to know which user performs the access without the need to parse and validate the JWT token issued by <code>cisidm</code> for every successful authentication.</p></li></ul><h2 id="first-start" tabindex="-1"><a class="header-anchor" href="#first-start" aria-hidden="true">#</a> First Start</h2><p>Finally it&#39;s time to start our services by calling the following command from the project directory (the one that contains the docker-compose file)</p><div class="language-text line-numbers-mode" data-ext="text"><pre class="language-text"><code>docker-compose up -d
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><div class="custom-container tip"><p class="custom-container-title">Admin User</p><p>Whenever <code>cisidm</code> starts, it checks if a user with super-user privileges (member of the <code>iam_superuser</code> role) exists. If not, a new registration token will be created and logged to stdout.</p><p>To create your initial admin user, copy that token from the log output (<code>docker-compose logs cisidm | grep &quot;superuser account&quot;</code>) and visit <code>https://account.example.com/registration?token=YOUR_TOKEN</code> and replace <code>YOUR_TOKEN</code> with the token from the log output.</p><p>It&#39;s also possible to use the cli tool for the registration: <code>idmctl register-user --registration-token YOUR-TOKEN your-username</code>.</p></div><br><hr><br>`,25),j={class:"custom-container tip"},F=e("p",{class:"custom-container-title"},"Congratulations",-1),O=e("p",null,"You just finished setting up cisidm with a HTTPS enabled reverse proxy that will now protect your services using Proxy/Forward Authentication.",-1);function W(G,V){const t=i("ExternalLinkIcon"),l=i("CodeGroupItem"),c=i("CodeGroup"),o=i("RouterLink");return p(),d("div",null,[k,v,h,b,e("ul",null,[e("li",null,[e("p",null,[e("a",g,[f,s(t)]),y,n(" rqlite is a database based on sqlite with clustering support using RAFT. It used by cisidm to store all it's data like users, roles, webauthn credentials and more.")])]),e("li",null,[e("p",null,[e("a",w,[_,s(t)]),n(":"),q,n(" A flexible and powerful reverse proxy that will handle automatic HTTPS via Let's Encrypt and secure our services by enforcing authentication via cisidm.")])]),e("li",null,[e("p",null,[e("a",T,[x,s(t)]),n(":"),R,n(" The identity management server.")])]),e("li",null,[e("p",null,[e("a",I,[C,s(t)]),n(":"),L,n(" A simple demo application that will be secured using cisidm and traefik.")])])]),P,s(c,null,{default:a(()=>[s(l,{title:"idm-config.yml"},{default:a(()=>[U]),_:1})]),_:1}),e("div",N,[H,e("p",null,[n("Refer to the "),s(o,{to:"/architecture/config-reference.html"},{default:a(()=>[n("Configuration File Reference")]),_:1}),n(" for a more detailed explanation of the configuration file.")])]),A,S,D,s(c,null,{default:a(()=>[s(l,{title:"docker-compose.yml"},{default:a(()=>[E]),_:1})]),_:1}),X,e("div",j,[F,O,e("p",null,[n("Now it's time to check the "),s(o,{to:"/guides/user-role-management.html"},{default:a(()=>[n("User and Role Administration Guide")]),_:1}),n(" or the "),s(o,{to:"/guides/cli-reference.html"},{default:a(()=>[n("Command Line Reference")]),_:1}),n(".")])])])}const Y=u(m,[["render",W],["__file","getting-started.html.vue"]]);export{Y as default};
