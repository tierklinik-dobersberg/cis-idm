import{_ as n,o as s,c as e,e as i}from"./app-wNEdNvye.js";const o={},a=i(`<h1 id="sms-twilio-setup" tabindex="-1"><a class="header-anchor" href="#sms-twilio-setup" aria-hidden="true">#</a> SMS (Twilio) Setup</h1><p><code>cisidm</code> supports phone number verifications and two-factor authentication using SMS codes. For this to work you need to configure Twilio support in <code>cisidm</code> by adding the <code>twilio</code> configuration block to your configuration file:</p><div class="language-hcl line-numbers-mode" data-ext="hcl"><pre class="language-hcl"><code><span class="token comment"># The (single) twilio block configures the Twilio integration which allows</span>
<span class="token comment"># sending SMS messages to your users. This is required for phone-number</span>
<span class="token comment"># verification to work.</span>
<span class="token keyword">twilio</span> <span class="token punctuation">{</span>
    <span class="token property">from</span> <span class="token punctuation">=</span> <span class="token string">&quot;Example Inc&quot;</span>
    <span class="token property">sid</span> <span class="token punctuation">=</span> <span class="token string">&quot;your-twilio-account-sid&quot;</span>
    <span class="token property">token</span> <span class="token punctuation">=</span> <span class="token string">&quot;your-twilio-account-token&quot;</span>
<span class="token punctuation">}</span>

</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,3),t=[a];function c(l,r){return s(),e("div",null,t)}const u=n(o,[["render",c],["__file","setup-sms.html.vue"]]);export{u as default};
