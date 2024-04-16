import{_ as n,o as s,c as e,e as a}from"./app-wNEdNvye.js";const i={},t=a(`<h1 id="e-mail-setup" tabindex="-1"><a class="header-anchor" href="#e-mail-setup" aria-hidden="true">#</a> E-Mail Setup</h1><p>To enalbe email verification and support for password-reset links, account invitations etc it&#39;s require to configure an outgoing E-Mail (SMTP) server in your <code>cisidm</code> configuration:</p><div class="language-hcl line-numbers-mode" data-ext="hcl"><pre class="language-hcl"><code>
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

</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,3),l=[t];function o(c,r){return s(),e("div",null,l)}const u=n(i,[["render",o],["__file","setup-email.html.vue"]]);export{u as default};
