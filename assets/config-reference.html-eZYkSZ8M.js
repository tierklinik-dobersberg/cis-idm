import{_ as n,o as s,c as e,e as a}from"./app-wNEdNvye.js";const i={},t=a(`<h1 id="configuration-file-reference" tabindex="-1"><a class="header-anchor" href="#configuration-file-reference" aria-hidden="true">#</a> Configuration File Reference</h1><p>Welcome to the Configuration File Reference for cisidm.</p><div class="custom-container tip"><p class="custom-container-title">Example Configuration</p><p>If you just want a full example configuration, explained with comments you can skip below to the <a href="#full-example-configuration-file">Full Example Configuration File</a></p></div><h2 id="full-example-configuration-file" tabindex="-1"><a class="header-anchor" href="#full-example-configuration-file" aria-hidden="true">#</a> Full Example Configuration File</h2><div class="language-hcl line-numbers-mode" data-ext="hcl"><pre class="language-hcl"><code><span class="token comment"># General Configuration Values</span>
<span class="token comment"># -------------------------------------------------------------------------------------------------------</span>


<span class="token comment"># Configures the logging level to use. Valid values for this setting are:</span>
<span class="token comment">#  - debug</span>
<span class="token comment">#  - info</span>
<span class="token comment">#  - warn</span>
<span class="token comment">#  - error</span>
<span class="token property">log_level</span> <span class="token punctuation">=</span> <span class="token string">&quot;info&quot;</span>

<span class="token comment"># The URL for the SQLite3 database that stores any user related information.</span>
<span class="token property">database_url</span> <span class="token punctuation">=</span> <span class="token string">&quot;file:/data/idm.db&quot;</span>

<span class="token comment"># Whehter or not roles can be created, modified or deleted via the</span>
<span class="token comment"># tkd.idm.v1.RoleService API.</span>
<span class="token comment"># If unset, the default value for enable_dynamic_roles depends on the presence</span>
<span class="token comment"># of one or more role blocks (see below).</span>
<span class="token comment">#</span>
<span class="token comment"># If at least one role is defined in the configuration file,</span>
<span class="token comment"># enable_dynamic_roles defaults to false since cisidm expects a static role</span>
<span class="token comment"># configuration. </span>
<span class="token comment"># If no roles are configures, enable_dynamic_roles defaults to true.</span>
<span class="token comment">#</span>
<span class="token comment"># Note that it&#39;s possible to disable roles altogheter by not specifying a role</span>
<span class="token comment"># block and by explicitly setting enable_dynamic_roles to false. Counterwise, </span>
<span class="token comment"># it&#39;s possible to have some static roles configured (which are delete/modify</span>
<span class="token comment"># protected) while still allowing roles to be created dynamically.</span>
<span class="token comment">#</span>
<span class="token property">enable_dynamic_roles</span> <span class="token punctuation">=</span> <span class="token boolean">false</span>

<span class="token comment"># Whether or not users are allowed to change their username.</span>
<span class="token comment"># If setting this to true, make sure all connected/protected services use the unique</span>
<span class="token comment"># user-ID for identification instead of the username. Otherwise enabling this may pose</span>
<span class="token comment"># a security risk.</span>
<span class="token property">allow_username_change</span> <span class="token punctuation">=</span> <span class="token boolean">false</span>

<span class="token comment"># Whether or not user addresses are enabled and can be managed by a user.</span>
<span class="token comment"># Set to true if you don&#39;t need user addresses.</span>
<span class="token comment">#</span>
<span class="token comment"># Defaults to false.</span>
<span class="token property">disable_user_addresses</span> <span class="token punctuation">=</span> <span class="token boolean">false</span>

<span class="token comment"># Whether or not user phone numbers are enabled and can be managed by a user.</span>
<span class="token comment"># Set to true if you don&#39;t need phone number support.</span>
<span class="token comment">#</span>
<span class="token comment"># Defaults to false.</span>
<span class="token property">disable_phone_numbers</span> <span class="token punctuation">=</span> <span class="token boolean">false</span>

<span class="token comment"># Configures the user registration mode:</span>
<span class="token comment">#</span>
<span class="token comment">#  - disabled (default): users are not allowed to register themself. In this </span>
<span class="token comment">#                        mode the administrator must create user accounts for</span>
<span class="token comment">#                        all users. User invitations are not possible!</span>
<span class="token comment">#</span>
<span class="token comment">#  - token: users must provide a token for registration. In this mode, it&#39;s possible to send invitation</span>
<span class="token comment">#           mails to your users and also assign them to specific roles upon successful registration.</span>
<span class="token comment">#</span>
<span class="token comment">#  - public: Anyone can create a user account on cisidm. Note that it&#39;s still possible to use token registration</span>
<span class="token comment">#            and user invitations as specified above.</span>
<span class="token comment">#</span>
<span class="token property">registration</span> <span class="token punctuation">=</span> <span class="token string">&quot;disabled&quot;</span>

<span class="token comment"># The server block configures the built-in HTTP/2 servers.</span>
<span class="token keyword">server</span> <span class="token punctuation">{</span>
    <span class="token comment"># Whether or not cookies issued by the server should enforce </span>
    <span class="token comment"># HTTPS. If unset, the default will be set based on the protocol</span>
    <span class="token comment"># in ui.public_url</span>
    <span class="token property">secure_cookies</span> <span class="token punctuation">=</span> <span class="token boolean">true</span>

    <span class="token comment"># Domain configures the cookie domain and the JWT issuer for access and</span>
    <span class="token comment"># refresh tokens.</span>
    <span class="token property">domain</span> <span class="token punctuation">=</span> <span class="token string">&quot;example.com&quot;</span>

    <span class="token comment"># The listen address of the HTTP server which requires authentication</span>
    <span class="token comment"># for API endpoints.</span>
    <span class="token property">public_listener</span> <span class="token punctuation">=</span> <span class="token string">&quot;:8080&quot;</span>

    <span class="token comment"># The listen address of the admin HTTP server. Requests on this endpoint</span>
    <span class="token comment"># will always be authenticated with the idm_superuser role assigned.</span>
    <span class="token comment">#</span>
    <span class="token comment"># SECURITY: You must make sure that this port is not publically accessible!</span>
    <span class="token property">admin_listener</span> <span class="token punctuation">=</span> <span class="token string">&quot;:8081&quot;</span>

    <span class="token comment"># Path to the static files to serve the web-interface on the public_listener.</span>
    <span class="token comment"># Valid values are:</span>
    <span class="token comment">#    - an empty string (or if the setting is omitted) serves the built-in</span>
    <span class="token comment">#      web-interface.</span>
    <span class="token comment">#    - a path to folder which contains web assets to serve.</span>
    <span class="token comment">#    - a HTTP or HTTPS URL. In this case, cisidm will setup a single host </span>
    <span class="token comment">#      reverse proxy and forward all UI/web related requests to that server.</span>
    <span class="token property">static_files</span> <span class="token punctuation">=</span> <span class="token string">&quot;&quot;</span>

    <span class="token comment"># In addition to the static_files setting above, it is also possible to</span>
    <span class="token comment"># expose additional assets under the /files path. Specify a path to a local</span>
    <span class="token comment"># folder to expose all folder content on the web.</span>
    <span class="token comment">#</span>
    <span class="token comment"># You may include additional assets like a logo or brand image here and set</span>
    <span class="token comment"># ui.logo_url = &quot;/files/my-logo.png&quot;</span>
    <span class="token property">extra_assets</span> <span class="token punctuation">=</span> <span class="token string">&quot;&quot;</span>

    <span class="token comment"># Configures additional origins that are allowed to perform</span>
    <span class="token comment"># cross-origin-resource requests (CORS).</span>
    <span class="token comment"># Note that ui.public_url and server.domain (http and https) are always</span>
    <span class="token comment"># added to this list.</span>
    <span class="token property">allowed_origins</span> <span class="token punctuation">=</span> <span class="token punctuation">[</span>
        <span class="token string">&quot;https://example.com&quot;</span>,
        <span class="token string">&quot;https://*.example.com&quot;</span>
    <span class="token punctuation">]</span>

    <span class="token comment"># A list of CIDR network addresses or hostnames which are considered trusted.</span>
    <span class="token comment"># For any requests originating from one of the specified networks, cisidm</span>
    <span class="token comment"># will trust the X-Forwarded-For header to determine the actual client IP.</span>
    <span class="token property">trusted_networks</span> <span class="token punctuation">=</span> <span class="token punctuation">[</span>
        <span class="token string">&quot;10.1.1.1/32&quot;</span>, <span class="token comment"># A single host</span>
        <span class="token string">&quot;10.1.2.1/24&quot;</span>, <span class="token comment"># A whole sub-net in CIDR notation</span>
        <span class="token string">&quot;traefik&quot;</span>,     <span class="token comment"># A hostname, cisidm will resolve the hostname on a</span>
                       <span class="token comment"># regular basis to detect changes to the IP.</span>
                       <span class="token comment"># In containerized environments with changing container</span>
                       <span class="token comment"># IPs it&#39;s best to use this method.</span>
    <span class="token punctuation">]</span>

    <span class="token comment"># A list of application domains to which cisidm will permit a redirect after</span>
    <span class="token comment"># a successful login/access token refresh. If this is unset, cisidm will</span>
    <span class="token comment"># refuse to redirect the user to prevent open-redirect vulnerabilitites.</span>
    <span class="token property">allowed_redirects</span> <span class="token punctuation">=</span> <span class="token punctuation">[</span>
        <span class="token string">&quot;example.com&quot;</span>, <span class="token comment"># Allow redirects only to example.com</span>
        <span class="token string">&quot;.example.com&quot;</span>, <span class="token comment"># Allow redirects to all subdomains of example.com</span>
    <span class="token punctuation">]</span>
<span class="token punctuation">}</span>

<span class="token comment"># The JWT block configures addtional settings for signing access and refresh</span>
<span class="token comment"># tokens.</span>
<span class="token keyword">jwt</span> <span class="token punctuation">{</span>
    <span class="token comment"># The audience for JWT tokens. This defaults to server.domain</span>
    <span class="token property">audience</span> <span class="token punctuation">=</span> <span class="token string">&quot;&quot;</span>

    <span class="token comment"># The secret used to sign various data and tokens. Rotating this secret will</span>
    <span class="token comment"># invalidate any access and refresh tokens.</span>
    <span class="token comment">#</span>
    <span class="token comment"># SECURITY: for the time being, cisidm only supports signing tokens in HS512</span>
    <span class="token comment">#           mode. Support for public-private keypairs with JWSK</span>
    <span class="token comment">#           (JSON-Web-Signing-Keys) is planned but not yet implemented.</span>
    <span class="token property">secret</span> <span class="token punctuation">=</span> <span class="token string">&quot;some-secure-random-string&quot;</span>

    <span class="token comment"># Configures the time-to-live for all access tokens issued by cisidm. This</span>
    <span class="token comment"># defaults to 1h.</span>
    <span class="token property">access_token_ttl</span> <span class="token punctuation">=</span> <span class="token string">&quot;1h&quot;</span>
    
    <span class="token comment"># cis_idm_access</span>
    <span class="token comment"># The name of the cookie that will hold the access token. This defaults to</span>
    <span class="token property">access_token_cookie_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;cis_idm_access&quot;</span>

    <span class="token comment"># Configures the time-to-live for all refresh tokens issued by cisidm. This</span>
    <span class="token comment"># defaults to 720h.</span>
    <span class="token property">refresh_token_ttl</span> <span class="token punctuation">=</span> <span class="token string">&quot;720h&quot;</span>

    <span class="token comment"># The name of the cookie that will hold the refresh token. This defaults to</span>
    <span class="token comment"># cis_idm_refresh. Note that the refresh cookie is limited to the refresh</span>
    <span class="token comment"># API endpoint /tkd.idm.v1.AuthService/RefreshToken on server.domain</span>
    <span class="token property">refresh_token_cookie_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;cis_idm_refresh&quot;</span>
<span class="token punctuation">}</span>

<span class="token comment"># The UI block configures settings for all user-facing interface like the web-ui</span>
<span class="token comment"># or any mail or SMS templates.</span>
<span class="token keyword">ui</span> <span class="token punctuation">{</span>
    <span class="token comment"># The name of your deployment. This will be visible in the built-in user</span>
    <span class="token comment"># interface and also in compiled mail templates.</span>
    <span class="token property">site_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;Example Inc&quot;</span>

    <span class="token comment"># The public address at which the cisidm server is reachable (the web-ui).</span>
    <span class="token property">public_url</span> <span class="token punctuation">=</span> <span class="token string">&quot;https://account.example.com&quot;</span>

    <span class="token comment"># A optional URL that is used on the web-ui and in mail templates when</span>
    <span class="token comment"># clicking the brand/deployment logo or the name of the deployment. If</span>
    <span class="token comment"># omitted, this defaults to ui.public_url</span>
    <span class="token property">site_name_url</span> <span class="token punctuation">=</span> <span class="token string">&quot;https://example.com&quot;</span>

    <span class="token comment"># Configures the resource where to find the the deployment/brand logo. This</span>
    <span class="token comment"># might be a fully-qualified external URL or a web-resource on the cisidm</span>
    <span class="token comment"># public server. See server.extra_assets for more information.</span>
    <span class="token property">logo_url</span> <span class="token punctuation">=</span> <span class="token string">&quot;/files/logo.png&quot;</span>

    <span class="token comment"># Configures the URL template to build the redirect URL for forward</span>
    <span class="token comment"># authentication. This value should not be set if the built-in web-ui is</span>
    <span class="token comment"># used.</span>
    <span class="token property">login_url</span> <span class="token punctuation">=</span> <span class="token string">&quot;&quot;</span>

    <span class="token comment"># Configures the URL template to build the redirect URL for forward</span>
    <span class="token comment"># authentication if the access token is expired. This value should not be</span>
    <span class="token comment"># set if the built-in web-ui is used. The built-in web-ui will automatically</span>
    <span class="token comment"># try to refresh the access token and redirect the user back to the</span>
    <span class="token comment"># application. If the refresh token is invalid (i.e. expired), the web-ui</span>
    <span class="token comment"># will render the login screen before redirecting the user back.</span>
    <span class="token property">refresh_url</span> <span class="token punctuation">=</span> <span class="token string">&quot;&quot;</span>

    <span class="token comment"># Configures the URL template to build the mail verification URL.</span>
    <span class="token comment"># This value should not be set if the built-in web-ui is used.</span>
    <span class="token property">verify_mail_url</span> <span class="token punctuation">=</span> <span class="token string">&quot;&quot;</span>

    <span class="token comment"># Configures the URL template to build URL used in user invitations.</span>
    <span class="token comment"># This value should not be set if the built-in web-ui is used.</span>
    <span class="token property">registration_url</span> <span class="token punctuation">=</span> <span class="token string">&quot;&quot;</span>
<span class="token punctuation">}</span>

<span class="token comment"># The (single) twilio block configures the Twilio integration which allows</span>
<span class="token comment"># sending SMS messages to your users. This is required for phone-number</span>
<span class="token comment"># verification to work.</span>
<span class="token keyword">twilio</span> <span class="token punctuation">{</span>
    <span class="token property">from</span> <span class="token punctuation">=</span> <span class="token string">&quot;Example Inc&quot;</span>
    <span class="token property">sid</span> <span class="token punctuation">=</span> <span class="token string">&quot;your-twilio-account-sid&quot;</span>
    <span class="token property">token</span> <span class="token punctuation">=</span> <span class="token string">&quot;your-twilio-account-token&quot;</span>
<span class="token punctuation">}</span>

<span class="token comment"># The (single) mail block configures the SMTP relay server that is used to send</span>
<span class="token comment"># mail messages. This is required for email-verification, user-invitation and</span>
<span class="token comment"># password reset mails to work.</span>
<span class="token keyword">mail</span> <span class="token punctuation">{</span>
    <span class="token comment"># The hostname (FQDN) of your mail server</span>
    <span class="token property">host</span> <span class="token punctuation">=</span> <span class="token string">&quot;smtp.example.com&quot;</span>

    <span class="token comment"># The port of your mail server</span>
    <span class="token property">port</span> <span class="token punctuation">=</span> <span class="token number">456</span>

    <span class="token comment"># The user and password to authenticate</span>
    <span class="token property">user</span> <span class="token punctuation">=</span> <span class="token string">&quot;noreply&quot;</span>
    <span class="token property">password</span> <span class="token punctuation">=</span> <span class="token string">&quot;a-secure-password&quot;</span>

    <span class="token comment"># The sender name used in e-mails. Make sure the value specified here is</span>
    <span class="token comment"># actually allowed as a sender on your SMTP server.</span>
    <span class="token property">from</span> <span class="token punctuation">=</span> <span class="token string">&quot;Example Inc &lt;noreply@example.com&gt;&quot;</span>

    <span class="token comment"># Whether or not SSL/TLS should be used when connecting to the mail server.</span>
    <span class="token property">use_tls</span> <span class="token punctuation">=</span> <span class="token boolean">true</span>

    <span class="token comment"># allow_insecure may be set to true to disable certificate validation when</span>
    <span class="token comment"># use_tls = true.</span>
    <span class="token comment">#</span>
    <span class="token comment"># SECURITY: make sure you know what you&#39;re doing before setting this value</span>
    <span class="token comment"># to true.</span>
    <span class="token property">allow_insecure</span> <span class="token punctuation">=</span> <span class="token boolean">false</span>
<span class="token punctuation">}</span>

<span class="token comment"># The (single) webpush block can be used to configure the Voluntary Application</span>
<span class="token comment"># Server Identification (VAPID) for end-to-end encrypted WebPush support.</span>
<span class="token keyword">webpush</span> <span class="token punctuation">{</span>
    <span class="token comment"># The email address of the person responsible for this cisidm deployment.</span>
    <span class="token comment"># This value is sent to web-push gateways upon web-push subscription.</span>
    <span class="token property">admin</span> <span class="token punctuation">=</span> <span class="token string">&quot;admin@example.com&quot;</span>

    <span class="token comment"># The public VAPID key</span>
    <span class="token property">vapid_public_key</span> <span class="token punctuation">=</span> <span class="token string">&quot;...&quot;</span>

    <span class="token comment"># The private VAPID key</span>
    <span class="token property">vapid_private_key</span> <span class="token punctuation">=</span> <span class="token string">&quot;...&quot;</span>
<span class="token punctuation">}</span>

<span class="token comment"># The dry_run block may be used to set cisidm into &quot;dry-run&quot; mode. In this mode,</span>
<span class="token comment"># any outgoing mail or SMS notification will be redirected to addresses</span>
<span class="token comment"># specified here. This setting is mainly for developers that want to test</span>
<span class="token comment"># delivery or new templates.</span>
<span class="token comment">#</span>
<span class="token comment"># dry_run { main = &quot;testing@example.com&quot; sms = &quot;+4312312312&quot;</span>
<span class="token comment"># }</span>

<span class="token comment"># Custom User Field Configuration</span>
<span class="token comment"># -------------------------------------------------------------------------------------------------------</span>
<span class="token comment">#</span>
<span class="token comment"># cisidm already supports storing private user information like email addresses,</span>
<span class="token comment"># phone numbers and addresses (like for delivery/billing). Though, in most</span>
<span class="token comment"># use-cases where cisidm provides authentication and authorization in</span>
<span class="token comment"># micro-service environments there will be some need for additional user data.</span>
<span class="token comment"># Some use cases might include:</span>
<span class="token comment">#  - storing user settings (like notification preferences)</span>
<span class="token comment">#  - adding additional per-user data like company internal phone extensions,</span>
<span class="token comment">#</span>
<span class="token comment"># For such cases, cisidm supports definition additional user fields that are</span>
<span class="token comment"># stored as JSON blobs in the cisidm database. For more information on</span>
<span class="token comment"># additional user fields please refer to the documentation.</span>
<span class="token comment">#</span>
<span class="token comment"># Take the following example:</span>

field <span class="token string">&quot;string&quot;</span> <span class="token string">&quot;internal-phone-extension&quot;</span> <span class="token punctuation">{</span>
    <span class="token comment"># This field may be seen by any authenticated user</span>
    <span class="token property">visibility</span> <span class="token punctuation">=</span> <span class="token string">&quot;public&quot;</span>
    
    <span class="token comment"># An optional description, this is for documentation purposes only</span>
    <span class="token property">description</span> <span class="token punctuation">=</span> <span class="token string">&quot;The company internal phone extension&quot;</span>

    <span class="token comment"># This field is populated by an administrator and cannot be changed by the</span>
    <span class="token comment"># user themself.</span>
    <span class="token property">writeable</span> <span class="token punctuation">=</span> <span class="token boolean">false</span>

    <span class="token comment"># A display name for the self-service web-ui. This is currently unused.</span>
    <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;Notification settings&quot;</span>
<span class="token punctuation">}</span>

field <span class="token string">&quot;object&quot;</span> <span class="token string">&quot;notification-settings&quot;</span> <span class="token punctuation">{</span>
    <span class="token comment"># Only the user is able to see this field</span>
    <span class="token property">visibility</span> <span class="token punctuation">=</span> <span class="token string">&quot;private&quot;</span>

    <span class="token comment"># This field may be written by the user</span>
    <span class="token property">writeable</span> <span class="token punctuation">=</span> <span class="token boolean">true</span>

    property <span class="token string">&quot;bool&quot;</span> <span class="token string">&quot;sms&quot;</span> <span class="token punctuation">{</span>
        <span class="token property">description</span> <span class="token punctuation">=</span> <span class="token string">&quot;Wether the user wants to receive notifications via SMS&quot;</span>
    <span class="token punctuation">}</span>

    property <span class="token string">&quot;bool&quot;</span> <span class="token string">&quot;email&quot;</span> <span class="token punctuation">{</span>
        <span class="token property">description</span> <span class="token punctuation">=</span> <span class="token string">&quot;Wether the user wants to receive notifications via EMail&quot;</span>
    <span class="token punctuation">}</span>
<span class="token punctuation">}</span>

<span class="token comment"># Role, Permissions, Policies and user/role overwrites</span>
<span class="token comment"># -------------------------------------------------------------------------------------------------------</span>
    
<span class="token keyword">forward_auth</span> <span class="token punctuation">{</span>
    <span class="token comment"># This configures the rego query to evaluate forward auth expressions. You</span>
    <span class="token comment"># only ever need to change this if your rego policies do not live under the</span>
    <span class="token comment"># cisidm.forward_auth package. Note that the query is expected to return a</span>
    <span class="token comment"># single result with a single expressions that contains the following keys:</span>
    <span class="token comment">#</span>
    <span class="token comment">#  - allow: a boolean to indicate if the request should be allowed or not</span>
    <span class="token comment">#  - headers: a map of headers to add to the upstream request in case it was</span>
    <span class="token comment">#    accepted.</span>
    <span class="token comment">#  - status_code: If set and the request is not allowed, this status code is</span>
    <span class="token comment">#                 sent to the user instead of letting cisidm figure out if</span>
    <span class="token comment">#                 the user should be redirected to the login/refresh page.</span>
    <span class="token comment">#</span>
    <span class="token property">rego_query</span> <span class="token punctuation">=</span> <span class="token string">&quot;data.cisidm.forward_auth&quot;</span>

    <span class="token comment"># The default policy for forward_auth queries.</span>
	<span class="token comment"># This may either be set to &quot;allow&quot; or &quot;deny&quot; (default).</span>
	<span class="token comment">#</span>
	<span class="token comment"># Depending on the value  cisidm will look for different rules</span>
	<span class="token comment"># when evaluating policies.</span>
    <span class="token comment">#</span>
	<span class="token comment"># If set to &quot;allow&quot;, cisidm will evaluate any &quot;deny&quot; rule.</span>
	<span class="token comment"># If set to &quot;deny&quot;, cisidm will evaluate any &quot;allow&quot; rule.</span>
    <span class="token property">default</span> <span class="token punctuation">=</span> <span class="token string">&quot;deny&quot;</span>

    <span class="token comment"># Whether or not Cross-Origin-Resource-Sharing (CORS) Preflight should</span>
    <span class="token comment"># always be allowed. This defaults to true.</span>
    <span class="token comment">#</span>
    <span class="token comment"># This is a performance optimization to avoid evaluating rego policies</span>
    <span class="token comment"># for preflight requests.</span>
    <span class="token comment">#</span>
    <span class="token comment"># When disabled, your rego policies likely need to account for</span>
    <span class="token comment"># CORS Preflight requests. A rule like the following could be used:</span>
    <span class="token comment">#</span>
    <span class="token comment"># \`\`\`</span>
    <span class="token comment">#    package cisidm.forward_auth</span>
    <span class="token comment">#</span>
    <span class="token comment">#    import rego.v1</span>
    <span class="token comment">#</span>
    <span class="token comment">#    allow if {</span>
    <span class="token comment">#        input.method = &quot;OPTIONS&quot;</span>
    <span class="token comment">#        &quot;Origin&quot; in input.headers</span>
    <span class="token comment">#    }</span>
    <span class="token comment">#</span>
    <span class="token comment"># \`\`\`</span>
    <span class="token property">allow_cors_preflight</span> <span class="token punctuation">=</span> <span class="token boolean">true</span>   

    <span class="token comment"># Below are the default values for headers set when the request is allowed to be</span>
    <span class="token comment"># forwarded to the protected upstream service.</span>
    <span class="token comment"># To disable a header you must explicitly set the configuration stanza to an empty</span>
    <span class="token comment"># string (for example: user_id_header = &quot;&quot;).</span>
    <span class="token comment">#</span>
    <span class="token comment"># Note that forwarding of those headers to the upstream service must likely be configured</span>
    <span class="token comment"># in your reverse proxy.</span>
    
    <span class="token comment"># The header that holds the unique user id</span>
    <span class="token property">user_id_header</span> <span class="token punctuation">=</span> <span class="token string">&quot;X-Remote-User-ID&quot;</span>

    <span class="token comment"># The header that holds the username</span>
    <span class="token property">username_header</span> <span class="token punctuation">=</span> <span class="token string">&quot;X-Remote-User&quot;</span>

    <span class="token comment"># The header that holds the user&#39;s primary email address</span>
    <span class="token property">mail_header</span> <span class="token punctuation">=</span> <span class="token string">&quot;X-Remote-Mail&quot;</span>

    <span class="token comment"># The header that holds all assigned role IDs</span>
    <span class="token property">role_header</span> <span class="token punctuation">=</span> <span class="token string">&quot;X-Remote-Role&quot;</span>

    <span class="token comment"># The header that holds the URL to retrieve the user avatar.</span>
    <span class="token property">avatar_header</span> <span class="token punctuation">=</span> <span class="token string">&quot;X-Remote-Avatar-URL&quot;</span>

    <span class="token comment"># The header that holds the user&#39;s choosen display_name, if any.</span>
    <span class="token property">display_name_header</span> <span class="token punctuation">=</span> <span class="token string">&quot;X-Remote-User-Display-Name&quot;</span>

    <span class="token comment"># The header that holds all resolved user permissions (that is, a distinct set</span>
    <span class="token comment"># of all permissions of all user roles)</span>
    <span class="token property">permission_header</span> <span class="token punctuation">=</span> <span class="token string">&quot;X-Remote-Permission&quot;</span>
<span class="token punctuation">}</span>

<span class="token comment"># The (single) policies block configures Open-Policy-Agent / Rego policies for</span>
<span class="token comment"># cisidm.</span>
<span class="token comment">#</span>
<span class="token comment"># Please refer to the documentation of cisidm on how to write policies and to</span>
<span class="token comment"># the offical OPA documentation for the REGO policy language.</span>
<span class="token comment">#</span>
<span class="token keyword">policies</span> <span class="token punctuation">{</span>

    <span class="token comment"># Enable rego policy debugging.</span>
    <span class="token property">debug</span> <span class="token punctuation">=</span> <span class="token boolean">false</span>

    <span class="token comment"># A list of directories where cisidm should (recursively) load .rego</span>
    <span class="token comment"># policies.</span>
    <span class="token property">directories</span> <span class="token punctuation">=</span> <span class="token punctuation">[</span>
        <span class="token string">&quot;./policies&quot;</span>
    <span class="token punctuation">]</span>

    <span class="token comment"># Instead of/Additional to loading policy files from a directory, it&#39;s also</span>
    <span class="token comment"># possible to specify policies inline by using one or more policy blocks.</span>
    <span class="token comment">#</span>
    <span class="token comment"># In the following example we define a policy named &quot;superuser&quot; that will</span>
    <span class="token comment"># permit any requests if the user has the \`idm_superuser\` role assigned.</span>
    policy <span class="token string">&quot;superuser&quot;</span> <span class="token punctuation">{</span>
        <span class="token property">content</span> <span class="token punctuation">=</span> <span class="token heredoc string">&lt;&lt;EOT
        package cisidm.forward_auth

        import future.keywords.in

        allow {
            user_is_superuser
        }

        user_is_superuser {
            input.subject
            input.subject.roles

            some role in input.subject.roles
            role.ID = &quot;idm_superuser&quot;
        }
        EOT</span>
    <span class="token punctuation">}</span>
<span class="token punctuation">}</span>

<span class="token comment"># A list of permissions that are stored in cisidm. Note that it&#39;s completely</span>
<span class="token comment"># fine to use the permission feature without specifying some in the</span>
<span class="token comment"># configuration. cisidm does treat permissions as &quot;unknown&quot; strings since</span>
<span class="token comment"># evaluation of privileges is either done using rego policies or by other</span>
<span class="token comment"># micro-service applications that just request a list of user permissions.</span>


<span class="token comment"># Though, for permissions specified here, cisidm can build a hierarchical tree so</span>
<span class="token comment"># it&#39;s easier to assign multiple permissions by using a shared prefix.</span>
<span class="token comment"># Hierarchical levels of a permission string are seperated using colons (&quot;:&quot;).</span>
<span class="token comment">#</span>
<span class="token comment"># The example values below will build the following tree of permissions:</span>
<span class="token comment">#</span>
<span class="token comment">#    - roster</span>
<span class="token comment">#      - write</span>
<span class="token comment">#         - create</span>
<span class="token comment">#         - approve</span>
<span class="token comment">#      - read</span>
<span class="token comment">#    - calendar</span>
<span class="token comment">#      - write</span>
<span class="token comment">#         - create</span>
<span class="token comment">#         - delete</span>
<span class="token comment">#         - move</span>
<span class="token comment">#      - read</span>
<span class="token comment">#</span>
<span class="token comment"># With this tree, it&#39;s possible to assign the permission string &quot;roster&quot; and</span>
<span class="token comment"># &quot;calendar:write&quot; to roles which cisidm will resolve to the following</span>
<span class="token comment"># permission set:</span>
<span class="token comment">#</span>
<span class="token comment">#    roster:write:create roster:write:approve roster:read calendar:write:create</span>
<span class="token comment">#    calendar:write:delete calendar:write:move</span>
<span class="token comment">#</span>
<span class="token comment"># Note: permissions themself do not have any authorizational meaning for cisidm</span>
<span class="token comment"># per-se but the can be used in rego policies to implement</span>
<span class="token comment"># permission/attribute-based access control (ABAC) rather than</span>
<span class="token comment"># role-based-access-control (RBAC).</span>
<span class="token comment">#</span>
<span class="token comment"># To enable permission trees set the following to true:</span>
<span class="token property">permission_trees</span> <span class="token punctuation">=</span> <span class="token boolean">false</span>

<span class="token comment"># Here are some examples for permissions.</span>
<span class="token property">permissions</span> <span class="token punctuation">=</span> <span class="token punctuation">[</span>
    <span class="token string">&quot;roster:write:create&quot;</span>,
    <span class="token string">&quot;roster:write:approve&quot;</span>,
    <span class="token string">&quot;roster:read&quot;</span>,
    <span class="token string">&quot;calendar:write:create&quot;</span>,
    <span class="token string">&quot;calendar:write:delete&quot;</span>,
    <span class="token string">&quot;calendar:write:move&quot;</span>,
    <span class="token string">&quot;calendar:read&quot;</span>
<span class="token punctuation">]</span>

<span class="token comment"># Role blocks can be used to configure static roles which cannot be modified or</span>
<span class="token comment"># deleted via the tkd.idm.v1.RoleService API.</span>
role <span class="token string">&quot;computer-accounts&quot;</span> <span class="token punctuation">{</span>
    <span class="token comment"># The name of the role. This is just for human representation but must be</span>
    <span class="token comment"># set.</span>
    <span class="token property">name</span> <span class="token punctuation">=</span> <span class="token string">&quot;Computer Accounts&quot;</span>

    <span class="token comment"># An optional description of the role. This is just for human</span>
    <span class="token comment"># representation.</span>
    <span class="token property">description</span> <span class="token punctuation">=</span> <span class="token string">&quot;Accounts for shared office computers&quot;</span>

    <span class="token comment"># A list of permissions that are assigned to this role.</span>
    <span class="token comment"># See permissions above.</span>
    <span class="token property">permissions</span> <span class="token punctuation">=</span> <span class="token punctuation">[</span>
        <span class="token string">&quot;roster:read&quot;</span>
    <span class="token punctuation">]</span>
<span class="token punctuation">}</span>

<span class="token comment"># Overwrite blocks allow to overwrite certain settings on a per role or per-user</span>
<span class="token comment"># basis.</span>

overwrite <span class="token string">&quot;role&quot;</span> <span class="token string">&quot;idm_superuser&quot;</span> <span class="token punctuation">{</span>
    <span class="token comment"># This block will set the access_token_ttl and the refresh_token_ttl to a</span>
    <span class="token comment"># much lower value for all user accounts with the idm_superuser role. Note</span>
    <span class="token comment"># that roles are matched by ID!</span>
    <span class="token property">access_token_ttl</span> <span class="token punctuation">=</span> <span class="token string">&quot;10m&quot;</span>
    <span class="token property">refresh_token_ttl</span> <span class="token punctuation">=</span> <span class="token string">&quot;2h&quot;</span>
<span class="token punctuation">}</span>

overwrite <span class="token string">&quot;user&quot;</span> <span class="token string">&quot;computer-account-1&quot;</span> <span class="token punctuation">{</span>
    <span class="token comment"># Set access and refresh token TTLs for a user account with ID</span>
    <span class="token comment"># &quot;computer-account-1&quot;. Though, overwriting on per-user basis is less common</span>
    <span class="token comment"># since one must first figure out the ID of the user.</span>
    <span class="token property">access_token_ttl</span> <span class="token punctuation">=</span> <span class="token string">&quot;1h&quot;</span>
    <span class="token property">refresh_token_ttl</span> <span class="token punctuation">=</span> <span class="token string">&quot;1480h&quot;</span>
<span class="token punctuation">}</span>

</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,5),l=[t];function o(c,p){return s(),e("div",null,l)}const d=n(i,[["render",o],["__file","config-reference.html.vue"]]);export{d as default};