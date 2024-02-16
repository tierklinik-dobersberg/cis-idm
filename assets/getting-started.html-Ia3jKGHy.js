import{_ as d,a as m}from"./03_login-kWPkIgdr.js";import{_ as h,r as o,o as k,c as v,a as e,d as n,b as s,w as a,e as c}from"./app-wNEdNvye.js";const b={},g=e("h1",{id:"getting-started",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#getting-started","aria-hidden":"true"},"#"),n(" Getting Started")],-1),f=e("p",null,[n("Welcome to the Getting-Started Guide for "),e("code",null,"cisidm"),n(". On this page we will walk through a simple docker-compose based setup to get you started.")],-1),y={class:"custom-container tip"},_=e("p",{class:"custom-container-title"},"Introduction",-1),w=e("code",null,"cisidm",-1),x=e("br",null,null,-1),q=e("hr",null,null,-1),T=e("p",null,[e("strong",null,"Contents")],-1),R={class:"table-of-contents"},I=e("hr",null,null,-1),C=e("h2",{id:"overview",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#overview","aria-hidden":"true"},"#"),n(" Overview")],-1),E=e("p",null,[n("The easiest way to get started using "),e("code",null,"cisidm"),n(" is to deploy it using docker-compose. In this guide we will setup a docker deployment with the following service:")],-1),S={href:"https://doc.traefik.io/traefik/",target:"_blank",rel:"noopener noreferrer"},O=e("strong",null,"traefik",-1),A=e("br",null,null,-1),L={href:"https://github.com/tierklinik-dobersberg/cis-idm",target:"_blank",rel:"noopener noreferrer"},D=e("strong",null,"cisidm",-1),j=e("br",null,null,-1),H={href:"https://gcr.io/google_containers/echoserver:1.4",target:"_blank",rel:"noopener noreferrer"},X=e("strong",null,"echoserver",-1),G=e("br",null,null,-1),U=c(`<div class="custom-container warning"><p class="custom-container-title">Domain Names</p><p>This example uses the following domain-names for testing purposes. You will likely need to add those domains to your <code>/etc/hosts</code> file!</p><p>Add the following content to <code>/etc/hosts</code>:</p><div class="language-plain line-numbers-mode" data-ext="plain"><pre class="language-plain"><code>127.0.0.1 example.intern account.example.intern app.example.intern
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div></div><h2 id="setup-the-project" tabindex="-1"><a class="header-anchor" href="#setup-the-project" aria-hidden="true">#</a> Setup the Project</h2>`,2),F={href:"https://github.com/tierklinik-dobersberg/cis-idm/tree/main/examples/getting-started",target:"_blank",rel:"noopener noreferrer"},N=c(`<p>Clone the <code>cisidm</code> repository and enter the project directory.</p><div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code><span class="token function">git</span> clone https://github.com/tierklinik-dobersberg/cis-idm

<span class="token builtin class-name">cd</span> cis-idm/examples/getting-started
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>All configuration files in this example are ready to go, for reference, we will still show them below and explain the important pieces.</p><p>Before you start, you should make sure that your <code>/etc/hosts</code> file contains the required &quot;DNS&quot; entries for this example:</p><div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code><span class="token function">sudo</span> <span class="token function">bash</span> <span class="token parameter variable">-c</span> <span class="token string">&quot;cat ./hosts &gt;&gt; /etc/hosts&quot;</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><h2 id="cisidm-configuraton-file" tabindex="-1"><a class="header-anchor" href="#cisidm-configuraton-file" aria-hidden="true">#</a> cisidm Configuraton File</h2><p>The configuration file for <code>cisidm</code> is located in the <code>idm.hcl</code> file and has the following content:</p>`,7),P=e("div",{class:"language-hcl line-numbers-mode","data-ext":"hcl"},[e("pre",{class:"language-hcl"},[e("code",null,[e("span",{class:"token comment"},"# Configures where the SQLite3 database should be stored."),n(`
`),e("span",{class:"token property"},"database_url"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token string"},'"file:/data/idm.db"'),n(`

`),e("span",{class:"token comment"},"# Disable self-registration of users. This means that any user account must be"),n(`
`),e("span",{class:"token comment"},"# created by an administrator using the `idmctl` cli utility."),n(`
`),e("span",{class:"token property"},"registration"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token string"},'"disabled"'),n(`

`),e("span",{class:"token comment"},"# This configures some defaults for the built-in server."),n(`
`),e("span",{class:"token keyword"},"server"),n(),e("span",{class:"token punctuation"},"{"),n(`
    `),e("span",{class:"token comment"},"# We're running on HTTP for this example"),n(`
    `),e("span",{class:"token comment"},'# so cookies must not be set to "secure"'),n(`
    `),e("span",{class:"token property"},"secure_cookies"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token boolean"},"false"),n(`

    `),e("span",{class:"token property"},"domain"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token string"},'"example.intern"'),n(`

    `),e("span",{class:"token property"},"allowed_origins"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token punctuation"},"["),n(`
        `),e("span",{class:"token string"},'"http://example.intern"'),n(`,
        `),e("span",{class:"token string"},'"http://*.example.intern"'),n(`,
    `),e("span",{class:"token punctuation"},"]"),n(`

    `),e("span",{class:"token property"},"trusted_networks"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token punctuation"},"["),n(`
        `),e("span",{class:"token string"},'"traefik"'),n(`
    `),e("span",{class:"token punctuation"},"]"),n(`

    `),e("span",{class:"token property"},"allowed_redirects"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token punctuation"},"["),n(`
        `),e("span",{class:"token string"},'"example.intern"'),n(`,
        `),e("span",{class:"token string"},'".example.intern"'),n(`
    `),e("span",{class:"token punctuation"},"]"),n(`
`),e("span",{class:"token punctuation"},"}"),n(`

`),e("span",{class:"token keyword"},"jwt"),n(),e("span",{class:"token punctuation"},"{"),n(`
    `),e("span",{class:"token property"},"secret"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token string"},'"some-secure-random-string"'),n(`
`),e("span",{class:"token punctuation"},"}"),n(`

`),e("span",{class:"token keyword"},"ui"),n(),e("span",{class:"token punctuation"},"{"),n(`
    `),e("span",{class:"token property"},"site_name"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token string"},'"Example Inc"'),n(`
    `),e("span",{class:"token property"},"public_url"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token string"},'"http://account.example.intern"'),n(`
`),e("span",{class:"token punctuation"},"}"),n(`

`),e("span",{class:"token keyword"},"forward_auth"),n(),e("span",{class:"token punctuation"},"{"),n(`
    `),e("span",{class:"token comment"},"# Allways allow CORS preflight requests."),n(`
    `),e("span",{class:"token comment"},"# If this would be set to false you would likely need to"),n(`
    `),e("span",{class:"token comment"},"# account for CORS preflight requests in your policies."),n(`
    `),e("span",{class:"token property"},"allow_cors_preflight"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token boolean"},"true"),n(`
`),e("span",{class:"token punctuation"},"}"),n(`

`),e("span",{class:"token keyword"},"policies"),n(),e("span",{class:"token punctuation"},"{"),n(`
    `),e("span",{class:"token property"},"debug"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token boolean"},"false"),n(`

    policy `),e("span",{class:"token string"},'"default"'),n(),e("span",{class:"token punctuation"},"{"),n(`
        `),e("span",{class:"token property"},"content"),n(),e("span",{class:"token punctuation"},"="),n(),e("span",{class:"token heredoc string"},`<<EOT
        package cisidm.forward_auth

        import rego.v1

        default allow := false

        allow if {
            # input.subject is only set if the request is authenticated
            input.subject
        }
        EOT`),n(`
    `),e("span",{class:"token punctuation"},"}"),n(`
`),e("span",{class:"token punctuation"},"}"),n(`
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),V={class:"custom-container tip"},W=e("p",{class:"custom-container-title"},"TIP",-1),B=e("h2",{id:"the-docker-compose-file",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#the-docker-compose-file","aria-hidden":"true"},"#"),n(" The docker-compose file")],-1),M=e("p",null,[n("We also need a "),e("code",null,"docker-compose.yml"),n(" file that contains our service and container definitions. It has the following content:")],-1),Y=e("div",{class:"language-yaml line-numbers-mode","data-ext":"yml"},[e("pre",{class:"language-yaml"},[e("code",null,[e("span",{class:"token key atrule"},"version"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"3"'),n(`

`),e("span",{class:"token key atrule"},"volumes"),e("span",{class:"token punctuation"},":"),n(`
  `),e("span",{class:"token key atrule"},"db"),e("span",{class:"token punctuation"},":"),n(`

`),e("span",{class:"token key atrule"},"services"),e("span",{class:"token punctuation"},":"),n(`
  `),e("span",{class:"token comment"},"# Traefik is our reverse proxy that will be configured to authenticate requests"),n(`
  `),e("span",{class:"token comment"},"# to upstream services using the forwardauth middleware."),n(`
  `),e("span",{class:"token comment"},"# See the cisidm container definition for more details"),n(`
  `),e("span",{class:"token key atrule"},"traefik"),e("span",{class:"token punctuation"},":"),n(`
    `),e("span",{class:"token key atrule"},"image"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"traefik:v2.10"'),n(`
    `),e("span",{class:"token key atrule"},"restart"),e("span",{class:"token punctuation"},":"),n(" unless"),e("span",{class:"token punctuation"},"-"),n(`stopped
    `),e("span",{class:"token key atrule"},"command"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--providers.docker=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--providers.docker.exposedbydefault=false"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"--entrypoints.http.address=:80"'),n(`

    `),e("span",{class:"token key atrule"},"ports"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token key atrule"},"target"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token number"},"80"),n(`
        `),e("span",{class:"token key atrule"},"published"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token number"},"80"),n(`
        `),e("span",{class:"token key atrule"},"protocol"),e("span",{class:"token punctuation"},":"),n(` tcp
        `),e("span",{class:"token key atrule"},"mode"),e("span",{class:"token punctuation"},":"),n(` host

    `),e("span",{class:"token key atrule"},"volumes"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"/var/run/docker.sock:/var/run/docker.sock"'),n(`

  `),e("span",{class:"token comment"},"# Our cisidm container"),n(`
  `),e("span",{class:"token key atrule"},"cisidm"),e("span",{class:"token punctuation"},":"),n(`
    `),e("span",{class:"token key atrule"},"image"),e("span",{class:"token punctuation"},":"),n(" ghcr.io/tierklinik"),e("span",{class:"token punctuation"},"-"),n("dobersberg/cis"),e("span",{class:"token punctuation"},"-"),n("idm"),e("span",{class:"token punctuation"},":"),n(`latest
    `),e("span",{class:"token key atrule"},"build"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token key atrule"},"context"),e("span",{class:"token punctuation"},":"),n(` ../../
    `),e("span",{class:"token key atrule"},"restart"),e("span",{class:"token punctuation"},":"),n(" unless"),e("span",{class:"token punctuation"},"-"),n(`stopped
    `),e("span",{class:"token key atrule"},"volumes"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(" ./idm.hcl"),e("span",{class:"token punctuation"},":"),n(`/etc/idm.hcl
      `),e("span",{class:"token punctuation"},"-"),n(" db"),e("span",{class:"token punctuation"},":"),n(`/data
    `),e("span",{class:"token key atrule"},"labels"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.enable=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.idm.rule=Host(`account.example.intern`)"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.idm.entrypoints=http"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.middlewares.auth.forwardauth.authResponseHeaders=X-Remote-User,X-Remote-User-ID,X-Remote-Mail,X-Remote-Mail-Verified,X-Remote-Avatar-URL,X-Remote-Role,X-Remote-User-Display-Name"'),n(`

    `),e("span",{class:"token key atrule"},"environment"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token key atrule"},"CONFIG_FILE"),e("span",{class:"token punctuation"},":"),n(),e("span",{class:"token string"},'"/etc/idm.hcl"'),n(`

  `),e("span",{class:"token key atrule"},"echoserver"),e("span",{class:"token punctuation"},":"),n(`
    `),e("span",{class:"token key atrule"},"image"),e("span",{class:"token punctuation"},":"),n(" gcr.io/google_containers/echoserver"),e("span",{class:"token punctuation"},":"),e("span",{class:"token number"},"1.4"),n(`
    `),e("span",{class:"token key atrule"},"labels"),e("span",{class:"token punctuation"},":"),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.enable=true"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.echo.rule=Host(`app.example.intern`)"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.echo.entrypoints=http"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.routers.echo.middlewares=auth"'),n(`
      `),e("span",{class:"token punctuation"},"-"),n(),e("span",{class:"token string"},'"traefik.http.services.echo.loadbalancer.server.port=8080"'),n(`
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),z=c(`<p>We&#39;ll explain the important bits below:</p><div class="custom-container tip"><p class="custom-container-title">TIP</p><p>If you&#39;re already familiar with docker-compose you can skip over the explanation to the next section <a href="#start-the-project">here</a>.</p></div><p>Let&#39;s break it down a bit:</p><div class="language-yaml line-numbers-mode" data-ext="yml"><pre class="language-yaml"><code><span class="token key atrule">version</span><span class="token punctuation">:</span> <span class="token string">&quot;3&quot;</span>

<span class="token key atrule">volumes</span><span class="token punctuation">:</span>
  <span class="token key atrule">db</span><span class="token punctuation">:</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>This just tells docker-compose which file-version we&#39;re using and that we want a volume named <code>db</code> that will hold our sqlite3 database.</p><h3 id="traefik-container" tabindex="-1"><a class="header-anchor" href="#traefik-container" aria-hidden="true">#</a> Traefik Container</h3><div class="language-yaml" data-ext="yml"><pre class="language-yaml"><code>  <span class="token comment"># Traefik is our reverse proxy that will be configured to authenticate requests</span>
  <span class="token comment"># to upstream services using the forwardauth middleware.</span>
  <span class="token comment"># See the cisidm container definition for more details</span>
  <span class="token key atrule">traefik</span><span class="token punctuation">:</span>
    <span class="token key atrule">image</span><span class="token punctuation">:</span> <span class="token string">&quot;traefik:v2.10&quot;</span>
    <span class="token key atrule">restart</span><span class="token punctuation">:</span> unless<span class="token punctuation">-</span>stopped
    <span class="token key atrule">command</span><span class="token punctuation">:</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;--providers.docker=true&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;--providers.docker.exposedbydefault=false&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;--entrypoints.http.address=:80&quot;</span>

    <span class="token key atrule">ports</span><span class="token punctuation">:</span>
      <span class="token punctuation">-</span> <span class="token key atrule">target</span><span class="token punctuation">:</span> <span class="token number">80</span>
        <span class="token key atrule">published</span><span class="token punctuation">:</span> <span class="token number">80</span>
        <span class="token key atrule">protocol</span><span class="token punctuation">:</span> tcp
        <span class="token key atrule">mode</span><span class="token punctuation">:</span> host

    <span class="token key atrule">volumes</span><span class="token punctuation">:</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;/var/run/docker.sock:/var/run/docker.sock&quot;</span>
</code></pre></div><p>The above service snippet configures Traefik as our reverse proxy. Let&#39;s go through the arguments passed to Traefik:</p><ul><li><p><code>--api.insecure=true</code><br> This enables the management dashboard of Traefik. This is useful for debugging service configurations. The management dashboard will be accessible on port 8080:</p></li><li><p><code>--providers.docker=true</code>,<br><code>--providers.docker.exposedbydefault=false</code><br> This enables the docker provider so we can configure Traefik using container labels. This will be explained in more detail below.</p></li><li><p><code>--entrypoints.http.address=:80</code>,<br> Tell traefik that we want to have one entrypoint on port 80 named <code>http</code></p></li></ul><h3 id="cisidm-container" tabindex="-1"><a class="header-anchor" href="#cisidm-container" aria-hidden="true">#</a> cisidm Container</h3><div class="language-yaml" data-ext="yml"><pre class="language-yaml"><code>  <span class="token comment"># Our cisidm container</span>
  <span class="token key atrule">cisidm</span><span class="token punctuation">:</span>
    <span class="token key atrule">image</span><span class="token punctuation">:</span> ghcr.io/tierklinik<span class="token punctuation">-</span>dobersberg/cis<span class="token punctuation">-</span>idm<span class="token punctuation">:</span>latest
    <span class="token key atrule">build</span><span class="token punctuation">:</span>
      <span class="token key atrule">context</span><span class="token punctuation">:</span> ../../
    <span class="token key atrule">restart</span><span class="token punctuation">:</span> unless<span class="token punctuation">-</span>stopped
    <span class="token key atrule">volumes</span><span class="token punctuation">:</span>
      <span class="token punctuation">-</span> ./idm.hcl<span class="token punctuation">:</span>/etc/idm.hcl
      <span class="token punctuation">-</span> db<span class="token punctuation">:</span>/data
    <span class="token key atrule">labels</span><span class="token punctuation">:</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.enable=true&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.routers.idm.rule=Host(\`account.example.intern\`)&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.routers.idm.entrypoints=http&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.middlewares.auth.forwardauth.authResponseHeaders=X-Remote-User,X-Remote-User-ID,X-Remote-Mail,X-Remote-Mail-Verified,X-Remote-Avatar-URL,X-Remote-Role,X-Remote-User-Display-Name&quot;</span>

    <span class="token key atrule">environment</span><span class="token punctuation">:</span>
      <span class="token key atrule">CONFIG_FILE</span><span class="token punctuation">:</span> <span class="token string">&quot;/etc/idm.hcl&quot;</span>
</code></pre></div><p>This configures the <code>cisidm</code> docker container, mounts the configuration file to <code>/etc/config.hcl</code> and tells traefik that is should be reachable at <code>account.example.intern</code>.</p><p>It also configures a new HTTP middleware <code>auth</code> that uses the forward-auth feature of traefik:</p><ul><li><p><code>traefik.http.middlewares.auth.forwardauth.address=http://cisidm:8080/validate</code><br> This configures the Forward-Auth middleware (named <code>auth</code> here) to forward any HTTP request to cisidm to determine if the user is actually allowed to access the requested resource. If cisidm replies with a HTTP success status code (2xx) than traefik will forward the original request to the actual service container. If cisidm replies with an error code, traefik will immediately return the response from <code>cisidm</code> to the user. This is used to redirect the user to the login page in case the request is unauthenticated.</p></li><li><p><code>traefik.http.middlewares.auth.forwardauth.authResponseHeaders=....</code><br> When <code>cisidm</code> successfully authenticated a request, it will return a set of headers that contain information about the logged in user. With this setting, we tell traefik to forward those headers to the actual service container. This enables service containers to know which user performs the access without the need to parse and validate the JWT token issued by <code>cisidm</code> for every successful authentication.</p></li></ul><h3 id="echoserver" tabindex="-1"><a class="header-anchor" href="#echoserver" aria-hidden="true">#</a> Echoserver</h3><p>Finnally, our <code>docker-compose.yml</code> file contains a simple echo server that we want to be protected by cisidm:</p><div class="language-yaml" data-ext="yml"><pre class="language-yaml"><code>  <span class="token key atrule">echoserver</span><span class="token punctuation">:</span>
    <span class="token key atrule">image</span><span class="token punctuation">:</span> gcr.io/google_containers/echoserver<span class="token punctuation">:</span><span class="token number">1.4</span>
    <span class="token key atrule">labels</span><span class="token punctuation">:</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.enable=true&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.routers.echo.rule=Host(\`app.example.intern\`)&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.routers.echo.entrypoints=http&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.routers.echo.middlewares=auth&quot;</span>
      <span class="token punctuation">-</span> <span class="token string">&quot;traefik.http.services.echo.loadbalancer.server.port=8080&quot;</span>
</code></pre></div><h2 id="start-the-example" tabindex="-1"><a class="header-anchor" href="#start-the-example" aria-hidden="true">#</a> Start the Example</h2><p>Finally it&#39;s time to start our services by calling the following command from the example directory (<code>examples/getting-started</code>):</p><div class="language-bash" data-ext="sh"><pre class="language-bash"><code><span class="token function">docker-compose</span> up <span class="token parameter variable">-d</span>
</code></pre></div><p>After a couple of minutes (depending on how many docker containers need to be downloaded) you example project should be up and running!</p><p>Since this is the first-start of <code>cisidm</code>, there&#39;s no administrative user yet. Even though registration is <code>disabled</code> in the cisidm configuration it is still possible to register once. This user will automatically be assigned the <code>idm_superuser</code> role.</p>`,22),J={href:"http://account.example.intern/registration",target:"_blank",rel:"noopener noreferrer"},Q=e("p",null,[e("img",{src:d,alt:"Registration"})],-1),K=c(`<p>After a successful registration you will be redirected to your cisidm profile page where you can complete your profile information, upload a user avatar and more.</p><div class="custom-container tip"><p class="custom-container-title">Registration using the cli</p><p>It&#39;s also possible to use the <code>idmctl</code> cli utility to register your account:</p><div class="language-bash" data-ext="sh"><pre class="language-bash"><code>idmctl register my-username <span class="token parameter variable">--password</span> my-password
</code></pre></div></div><p><strong>For demonstration purposes it&#39;s best to click the &quot;Logout&quot; button now so you can see the full authentication flow when accessing the demo application.</strong></p>`,3),Z={href:"http://app.example.intern",target:"_blank",rel:"noopener noreferrer"},$=e("p",null,'If you clicked "logout" on the profile view before you will now be redirected to the login screen of cisidm:',-1),ee=e("p",null,[e("img",{src:m,alt:"Login Screen"})],-1),ne=c(`<p>Once you complete the login flow, cisidm will redirect to back to the protected application.</p><p>The <code>echoserver</code> application just dumps the HTTP request that it received (your response might look a bit different):</p><div class="language-plain line-numbers-mode" data-ext="plain"><pre class="language-plain"><code>CLIENT VALUES:
client_address=172.26.0.3
command=GET
real path=/
query=nil
request_version=1.1
request_uri=http://app.example.intern:8080/

SERVER VALUES:
server_version=nginx: 1.10.0 - lua: 10001

HEADERS RECEIVED:
accept=text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
accept-encoding=gzip, deflate
accept-language=de,en-US;q=0.9,en;q=0.8
cookie=cis_idm_access=&lt;REDACTED&gt;
host=app.example.intern
referer=http://account.example.intern/
upgrade-insecure-requests=1
x-forwarded-for=172.26.0.1
x-forwarded-host=app.example.intern
x-forwarded-port=80
x-forwarded-proto=http
x-forwarded-server=db8c46e822b1
x-real-ip=172.26.0.1
x-remote-avatar-url=http://account.example.intern/avatar/23fa762e-0267-41e5-aa7c-7bc28a462e8a
x-remote-mail=admin@example.intern
x-remote-role=idm_superuser
x-remote-user=admin
x-remote-user-id=23fa762e-0267-41e5-aa7c-7bc28a462e8a
BODY:
-no body in request-
</code></pre><div class="highlight-lines"><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><div class="highlight-line"> </div><div class="highlight-line"> </div><div class="highlight-line"> </div><div class="highlight-line"> </div><div class="highlight-line"> </div><br><br></div><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>As you can seen between lines <strong>26 and 30</strong>, <code>cisidm</code> authenticated the request and added remote-headers so your applications and services can just rely on the presence of those headers to know which user is logged-in.</p><br>`,5),se={class:"custom-container warning"},ae=e("p",{class:"custom-container-title"},"Open-ID-Connect (OIDC)",-1),te=e("p",null,"Unfortunately there are lots of services that don't support user authentication using remote headers. Though, if this is an issue you might still be able to provide a Single-Sign-On experience to your users by using Open-ID-Connect (OIDC).",-1),ie=e("strong",null,"Checkout our OIDC Guide for an example!",-1),oe=e("hr",null,null,-1),le=e("br",null,null,-1),ce={class:"custom-container tip"},re=e("p",{class:"custom-container-title"},"Congratulations",-1),ue=e("p",null,"You just finished setting up cisidm with a reverse proxy that will now protect your services using Proxy/Forward Authentication.",-1);function pe(de,me){const l=o("RouterLink"),t=o("router-link"),i=o("ExternalLinkIcon"),r=o("CodeGroupItem"),u=o("CodeGroup"),p=o("center");return k(),v("div",null,[g,f,e("div",y,[_,e("p",null,[n("If you just want to find out what "),w,n(" is and what it was designed for head over to the "),s(l,{to:"/guides/intro.html"},{default:a(()=>[n("Introduction")]),_:1})])]),x,q,T,e("nav",R,[e("ul",null,[e("li",null,[s(t,{to:"#overview"},{default:a(()=>[n("Overview")]),_:1})]),e("li",null,[s(t,{to:"#setup-the-project"},{default:a(()=>[n("Setup the Project")]),_:1})]),e("li",null,[s(t,{to:"#cisidm-configuraton-file"},{default:a(()=>[n("cisidm Configuraton File")]),_:1})]),e("li",null,[s(t,{to:"#the-docker-compose-file"},{default:a(()=>[n("The docker-compose file")]),_:1}),e("ul",null,[e("li",null,[s(t,{to:"#traefik-container"},{default:a(()=>[n("Traefik Container")]),_:1})]),e("li",null,[s(t,{to:"#cisidm-container"},{default:a(()=>[n("cisidm Container")]),_:1})]),e("li",null,[s(t,{to:"#echoserver"},{default:a(()=>[n("Echoserver")]),_:1})])])]),e("li",null,[s(t,{to:"#start-the-example"},{default:a(()=>[n("Start the Example")]),_:1})])])]),I,C,E,e("ul",null,[e("li",null,[e("p",null,[e("a",S,[O,s(i)]),n(":"),A,n(" A flexible and powerful reverse proxy that will handle automatic HTTPS via Let's Encrypt and secure our services by enforcing authentication via cisidm.")])]),e("li",null,[e("p",null,[e("a",L,[D,s(i)]),n(":"),j,n(" The identity management server.")])]),e("li",null,[e("p",null,[e("a",H,[X,s(i)]),n(":"),G,n(" A simple demo application that will be secured using cisidm and traefik.")])])]),U,e("p",null,[e("strong",null,[n("The whole example can be found in the "),e("a",F,[n("GitHub repository"),s(i)]),n(".")])]),N,s(u,null,{default:a(()=>[s(r,{title:"idm.hcl"},{default:a(()=>[P]),_:1})]),_:1}),e("div",V,[W,e("p",null,[n("Refer to the "),s(l,{to:"/architecture/config-reference.html"},{default:a(()=>[n("Configuration File Reference")]),_:1}),n(" for a more detailed explanation of the configuration file.")])]),B,M,s(u,null,{default:a(()=>[s(r,{title:"docker-compose.yml"},{default:a(()=>[Y]),_:1})]),_:1}),z,e("p",null,[n("Open up "),e("a",J,[n("http://account.example.intern/registration"),s(i)]),n(" and create your admin account:")]),s(p,null,{default:a(()=>[Q]),_:1}),K,e("p",null,[n("Awesome, let's try the example application by opening "),e("a",Z,[n("http://app.example.intern"),s(i)]),n(" in your web-browser.")]),$,s(p,null,{default:a(()=>[ee]),_:1}),ne,e("div",se,[ae,te,e("p",null,[s(l,{to:"/guides/setup-oidc.html"},{default:a(()=>[ie]),_:1})])]),oe,le,e("div",ce,[re,ue,e("p",null,[n("Now it's time to check the "),s(l,{to:"/guides/user-role-management.html"},{default:a(()=>[n("User and Role Administration Guide")]),_:1}),n(" or the "),s(l,{to:"/guides/cli-reference.html"},{default:a(()=>[n("Command Line Reference")]),_:1}),n(".")])])])}const ve=h(b,[["render",pe],["__file","getting-started.html.vue"]]);export{ve as default};
