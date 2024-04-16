import{_ as d,r as l,o as u,c as p,a as e,b as s,w as a,d as n,e as o}from"./app-wNEdNvye.js";const m={},v=o('<h1 id="user-and-role-management" tabindex="-1"><a class="header-anchor" href="#user-and-role-management" aria-hidden="true">#</a> User and Role Management</h1><p>In <code>cisidm</code> identity management is split up into <a href="#users">Users</a>, <a href="#roles">Roles</a> and <a href="#permissions">Permissions</a>.</p><p>Each user might be assigned multiple roles (or even none) with each role having a set associated permissions.</p><p>Note that cisidm itself does not define any permissions itself since this is up to the administrator to define. <a href="#permissions">Permissions</a> are just opaque strings in cisidm. Instead, the cisidm APIs are already protected by either just requiring valid authentication or an administrator account with the built-in <code>idm_superuser</code> role. See <a href="#roles">Roles</a> for more information.</p><br><hr><p><strong>Contents</strong></p>',7),h={class:"table-of-contents"},b=e("hr",null,null,-1),k=e("h2",{id:"users",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#users","aria-hidden":"true"},"#"),n(" Users")],-1),g=e("p",null,[n("Users are the main identities used for authentication in cisidm. Users are identified by a unique ID that is not allowed to change. Note that by default, "),e("code",null,"cisidm"),n(" permit users to change their username as long as it's not already take. The unique ID ensures that services that integrate with cisidm have a stable identifier.")],-1),f={class:"custom-container tip"},_=e("p",{class:"custom-container-title"},"TIP",-1),y=e("p",null,[n("When using cisidm just for forward authentication using a supported reverse proxy like Traefik or Caddy, it's best to disable the username-change feature using the configuration. This way the username is also expected to be stable and can thus safely be sent to protected upstream services using the "),e("code",null,"X-Remote-User"),n(" header.")],-1),w=o('<p><code>cisidm</code> provides a default set of fields that users can change on their own using the <code>tkd.idm.v1.SelfService</code> API endpoint. Those fields include:</p><ul><li>First and Lastname</li><li>A Display name that should/can be used instead of the username</li><li>Birthday</li><li>Avatar</li><li>Multiple addresses: for example billing, delivery, ...</li><li>Multiple phone numbers</li><li>Multiple email addresses</li></ul><p>Each user can also specify a primary phone and email address that is used for password-reset codes and any other notification.</p><p>In addition, <code>cisidm</code> maintains metadata for two-factor authentication (2FA), Backup-Codes, WebauthN/Passkeys.</p><h3 id="creating-users" tabindex="-1"><a class="header-anchor" href="#creating-users" aria-hidden="true">#</a> Creating Users</h3><p>There are multiple ways to setup user accounts in cisidm. If <code>registration</code> is set to <code>&quot;public&quot;</code> in the configuration, any person can create a user account on your deployment using the built-in web interface.</p><p>If <code>registration</code> is set to <code>&quot;token&quot;</code>, only users with a valid registration token can register themself. An administrator can either send an invitation mail to a user or generate a token and distribute manually. Using token based registration allows the administrator to immediately assign one or more roles as soon as the user completes the registration.</p><p>An administrator may also manually create user accounts and can then send an account-creation-notice that contains a password reset link so your users can choose their initial password.</p>',8),q=e("div",{class:"language-bash line-numbers-mode","data-ext":"sh"},[e("pre",{class:"language-bash"},[e("code",null,[n("idmctl "),e("span",{class:"token function"},"users"),n(" invite "),e("span",{class:"token punctuation"},"["),n("email-addresses "),e("span",{class:"token punctuation"},".."),n("."),e("span",{class:"token punctuation"},"]"),n(),e("span",{class:"token parameter variable"},"--roles"),n(),e("span",{class:"token punctuation"},"["),n("role-name/id "),e("span",{class:"token punctuation"},".."),n("."),e("span",{class:"token punctuation"},"]"),n(`

`),e("span",{class:"token comment"},"# For example:"),n(`
idmctl `),e("span",{class:"token function"},"users"),n(" invite alice@example.com "),e("span",{class:"token parameter variable"},"--roles"),n(" help-desk "),e("span",{class:"token parameter variable"},"--roles"),n(` support
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),x=e("div",{class:"language-bash line-numbers-mode","data-ext":"sh"},[e("pre",{class:"language-bash"},[e("code",null,[n("idmctl "),e("span",{class:"token function"},"users"),n(` generate-registration-token
    --max-usage `),e("span",{class:"token punctuation"},"["),n("number"),e("span",{class:"token punctuation"},"]"),n("                    "),e("span",{class:"token comment"},"# How often the token can be used"),n(`
    `),e("span",{class:"token parameter variable"},"--ttl"),n(),e("span",{class:"token punctuation"},"["),n("duration"),e("span",{class:"token punctuation"},"]"),n("                        "),e("span",{class:"token comment"},"# How long the token remains valid "),n(`
    `),e("span",{class:"token parameter variable"},"--roles"),n(),e("span",{class:"token punctuation"},"["),n("role-name/id "),e("span",{class:"token punctuation"},".."),n("."),e("span",{class:"token punctuation"},"]"),n("              "),e("span",{class:"token comment"},"# One or more roles to assign upon"),n(`
                                            `),e("span",{class:"token comment"},"# registration"),n(`

`),e("span",{class:"token comment"},"# For example:"),n(`
idmctl `),e("span",{class:"token function"},"users"),n(` generate-registration-token
    --max-usage `),e("span",{class:"token number"},"3"),n(` 
    `),e("span",{class:"token parameter variable"},"--ttl"),n(` 24h                                          
    `),e("span",{class:"token parameter variable"},"--roles"),n(` help-desk                                  
    `),e("span",{class:"token parameter variable"},"--roles"),n(` support

`),e("span",{class:"token comment"},"# Output:"),n(`
`),e("span",{class:"token comment"},"#   token: some-random-token-to-distribute"),n(`
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),A=e("div",{class:"language-bash line-numbers-mode","data-ext":"sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token comment"},"# Manually create a user"),n(`
idmctl `),e("span",{class:"token function"},"users"),n(` create
    `),e("span",{class:"token parameter variable"},"--name"),n(),e("span",{class:"token punctuation"},"["),n("username"),e("span",{class:"token punctuation"},"]"),n(`
    --display-name `),e("span",{class:"token punctuation"},"["),n("display-name"),e("span",{class:"token punctuation"},"]"),n(` 
    --first-name `),e("span",{class:"token punctuation"},"["),n("first-name"),e("span",{class:"token punctuation"},"]"),n(`
    --last-name `),e("span",{class:"token punctuation"},"["),n("last-name"),e("span",{class:"token punctuation"},"]"),n(`
    `),e("span",{class:"token parameter variable"},"--phone"),n(),e("span",{class:"token punctuation"},"["),n("phone-numbers"),e("span",{class:"token punctuation"},".."),n("."),e("span",{class:"token punctuation"},"]"),n(`
    `),e("span",{class:"token parameter variable"},"--email"),n(),e("span",{class:"token punctuation"},"["),n("email-addresses"),e("span",{class:"token punctuation"},".."),n("."),e("span",{class:"token punctuation"},"]"),n(`
    `),e("span",{class:"token parameter variable"},"--role"),n(),e("span",{class:"token punctuation"},"["),n("roles"),e("span",{class:"token punctuation"},".."),n("."),e("span",{class:"token punctuation"},"]"),n(`
    `),e("span",{class:"token parameter variable"},"--password"),n(),e("span",{class:"token punctuation"},"["),n("plain-text-or-bcrypt-password"),e("span",{class:"token punctuation"},"]"),n(`

`),e("span",{class:"token comment"},"# Send account creation notice to one or more users containing a password-reset"),n(`
`),e("span",{class:"token comment"},"# link."),n(`
idmctl `),e("span",{class:"token function"},"users"),n(" send-account-notice "),e("span",{class:"token punctuation"},"["),n("username"),e("span",{class:"token punctuation"},"]"),n(`

`),e("span",{class:"token comment"},"# Example:"),n(`
`),e("span",{class:"token comment"},"# Note that the first phone number and the first email address will be marked as"),n(`
`),e("span",{class:"token comment"},"# the primary one."),n(`
idmctl `),e("span",{class:"token function"},"users"),n(` create
    `),e("span",{class:"token parameter variable"},"--name"),n(` alice
    --first-name Alice
    --last-name Mustermann
    `),e("span",{class:"token parameter variable"},"--phone"),n(` +4312341234                     
    `),e("span",{class:"token parameter variable"},"--phone"),n(` +49987987987
    `),e("span",{class:"token parameter variable"},"--email"),n(` alice@example.com
    `),e("span",{class:"token parameter variable"},"--email"),n(` alice@help-desk.example.com
    `),e("span",{class:"token parameter variable"},"--role"),n(` help-desk
    `),e("span",{class:"token parameter variable"},"--role"),n(` support

idmctl `),e("span",{class:"token function"},"users"),n(` send-account-notice alice
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),I=e("h3",{id:"additional-user-fields",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#additional-user-fields","aria-hidden":"true"},"#"),n(" Additional User Fields")],-1),C=e("h2",{id:"roles",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#roles","aria-hidden":"true"},"#"),n(" Roles")],-1),P=o(`<p><code>cisidm</code> itself does not care about the defined roles but rather allows an administrator to configure roles based on their requirements.</p><div class="custom-container warning"><p class="custom-container-title">Note</p><p>There is one special role in <code>cisidm</code> called the <code>idm_superuser</code> role. This role actually does imply a set of permissions: Any user with this role can perform any action on any API endpoint of cisidm and is thus considered an administrative account.</p><p>It&#39;s <strong>strongly advised</strong> to only use a <code>idm_superuser</code> account for administrative tasks with multi-factor authentication (TOTP or SMS/E-Mail codes) enabled and use a separate user account for daily work/authentication.</p></div><h3 id="static-configuration" tabindex="-1"><a class="header-anchor" href="#static-configuration" aria-hidden="true">#</a> Static configuration</h3><p>Roles may already be specified in the configuration file. In this case, <code>cisidm</code> will prevent modifications to those roles using the API (they are considered system-roles).</p><p>To define a role in the configuration file, you just need to add a <code>role</code> block:</p><div class="language-hcl line-numbers-mode" data-ext="hcl"><pre class="language-hcl"><code>role <span class="token string">&quot;computer-accounts&quot;</span> <span class="token punctuation">{</span>
    <span class="token comment"># The name of the role. This is just for human representation but must be</span>
    <span class="token comment"># set.</span>
    <span class="token property">name</span> <span class="token punctuation">=</span> <span class="token string">&quot;Computer Accounts&quot;</span>

    <span class="token comment"># An optional description of the role. This is just for human</span>
    <span class="token comment"># representation.</span>
    <span class="token property">description</span> <span class="token punctuation">=</span> <span class="token string">&quot;Accounts for shared office computers&quot;</span>

    <span class="token comment"># A list of permissions that are assigned to this role.</span>
    <span class="token comment"># See Permissions below.</span>
    <span class="token property">permissions</span> <span class="token punctuation">=</span> <span class="token punctuation">[</span>
        <span class="token string">&quot;roster:read&quot;</span>
    <span class="token punctuation">]</span>
<span class="token punctuation">}</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h3 id="dynamic-configuration" tabindex="-1"><a class="header-anchor" href="#dynamic-configuration" aria-hidden="true">#</a> Dynamic configuration</h3><p>When no roles are configured in the configuration file or if <code>enable_dynamic_roles</code> is set, then roles may be created dynamically using the <code>tkd.idm.v1.RolesService</code> API.</p>`,8),T=o(`<div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code><span class="token comment"># Create the role, note if --id is not set, a random ID will be generated by</span>
<span class="token comment"># cisidm.</span>
idmctl roles create
    <span class="token parameter variable">--id</span> <span class="token string">&quot;computer-accounts&quot;</span>
    <span class="token parameter variable">--name</span> <span class="token string">&quot;Computer Accounts&quot;</span>
    <span class="token parameter variable">--description</span> <span class="token string">&quot;Accounts for shared office computers&quot;</span>

<span class="token comment"># Optionally, assign permissions to the role</span>
<span class="token comment"># FIXME(ppacher): not yet implemented</span>
idmctl roles add-permissions <span class="token string">&quot;computer-accounts&quot;</span> <span class="token parameter variable">--permission</span> <span class="token string">&quot;roster:read&quot;</span>

<span class="token comment"># Inspect role and it&#39;s assigned permissions</span>
idmctl roles <span class="token string">&quot;computer-account&quot;</span>
idmctl roles get-permissions <span class="token string">&quot;computer-accounts&quot;</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="permissions" tabindex="-1"><a class="header-anchor" href="#permissions" aria-hidden="true">#</a> Permissions</h2>`,2),S=e("code",null,"cisidm",-1),R=e("code",null,"cisidm",-1),M=o(`<p>For simple setups, one may enable <code>permission_trees = true</code> in the configuration file. If enabled, permssions strings will be parsed as trees and <code>cisidm</code> will resolve all user permissions based on that tree.</p><p>For example, consider the following configuration file:</p><div class="language-hcl line-numbers-mode" data-ext="hcl"><pre class="language-hcl"><code><span class="token property">permission_trees</span> <span class="token punctuation">=</span> <span class="token boolean">true</span>

<span class="token property">permissions</span> <span class="token punctuation">=</span> <span class="token punctuation">[</span>
    <span class="token string">&quot;calendar:events:read&quot;</span>
    <span class="token string">&quot;calendar:events:write&quot;</span>
    <span class="token string">&quot;calendar:events:write:delete&quot;</span>,
    <span class="token string">&quot;calendar:events:write:create&quot;</span>,
    <span class="token string">&quot;calendar:events:write:move&quot;</span>,
    <span class="token string">&quot;calendar:events:write:update&quot;</span>
<span class="token punctuation">]</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>Taken the above permission configuration, <code>cisidm</code> will build the following permission tree:</p><div class="language-text line-numbers-mode" data-ext="text"><pre class="language-text"><code>- calendar
    - events
        - read
        - write
            - delete
            - create
            - move
            - update
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>When resolving user permissions, cisidm will try to detect inherited child permissions automatically:</p><ul><li><p>assigned: <code>calendar:events:read</code><br> resolved: <code>calendar:events:read</code></p></li><li><p>assigned: <code>calendar:events:write</code><br> resolved: <code>calendar:events:write:delete</code>, <code>calendar:events:write:create</code>, <code>calendar:events:write:move</code>, <code>calendar:events:write:update</code></p></li><li><p>assigned: <code>calendar</code><br> resolved: <code>calendar:events.read</code>, <code>calendar:events:write:delete</code>, <code>calendar:events:write:create</code>, <code>calendar:events:write:move</code>, <code>calendar:events:write:update</code></p></li></ul>`,7);function U(F,N){const t=l("router-link"),i=l("RouterLink"),r=l("CodeGroupItem"),c=l("CodeGroup");return u(),p("div",null,[v,e("nav",h,[e("ul",null,[e("li",null,[s(t,{to:"#users"},{default:a(()=>[n("Users")]),_:1}),e("ul",null,[e("li",null,[s(t,{to:"#creating-users"},{default:a(()=>[n("Creating Users")]),_:1})]),e("li",null,[s(t,{to:"#additional-user-fields"},{default:a(()=>[n("Additional User Fields")]),_:1})])])]),e("li",null,[s(t,{to:"#roles"},{default:a(()=>[n("Roles")]),_:1}),e("ul",null,[e("li",null,[s(t,{to:"#static-configuration"},{default:a(()=>[n("Static configuration")]),_:1})]),e("li",null,[s(t,{to:"#dynamic-configuration"},{default:a(()=>[n("Dynamic configuration")]),_:1})])])]),e("li",null,[s(t,{to:"#permissions"},{default:a(()=>[n("Permissions")]),_:1})])])]),b,k,g,e("div",f,[_,y,e("p",null,[n("See "),s(i,{to:"/guides/policies.html"},{default:a(()=>[n("Policies")]),_:1}),n(" and the "),s(i,{to:"/guides/cli-reference.html"},{default:a(()=>[n("Configuration File Reference")]),_:1}),n(" for more information on forward authentication.")])]),w,s(c,null,{default:a(()=>[s(r,{title:"Invitation Mail"},{default:a(()=>[q]),_:1}),s(r,{title:"Generating Tokens"},{default:a(()=>[x]),_:1}),s(r,{title:"Manual Account Creation"},{default:a(()=>[A]),_:1})]),_:1}),I,e("p",null,[n("When integrating cisidm with your own services / applications you might need additional metadata for users. Refer to the "),s(i,{to:"/guides/extra-user-fields.html"},{default:a(()=>[n("Additional User Fields")]),_:1}),n(" guide for more information.")]),C,e("p",null,[n("As mentioned at the beginning, each user may be assigned multiple roles. Those roles may either be used in "),s(i,{to:"/guides/policies.html"},{default:a(()=>[n("Policies")]),_:1}),n(" to implement Role-Based-Access-Control (RBAC) or by services that integrate directly with cisidm using it's exposed API.")]),P,e("p",null,[n("For example, to create the same role as above using the "),s(i,{to:"/guides/cli-reference.html"},{default:a(()=>[n("idmctl utility")]),_:1}),n(":")]),T,e("p",null,[n("Last but not least, each role may be assigned multiple permissions. Permissions are (mostly) opaque text values without any special meaning in "),S,n(" but enable an administrator to implement Permission-Based-Access-Control (PBAC). Since "),R,n(" does not inspect the value of permissions, it's also possible to implement more sophisticated authorization systems like AWS IAM Statements by using JSON encoded objects as permission strings and parse/evaluate them in your "),s(i,{to:"/guides/policies.html"},{default:a(()=>[n("rego policies")]),_:1}),n(".")]),M])}const j=d(m,[["render",U],["__file","user-role-management.html.vue"]]);export{j as default};