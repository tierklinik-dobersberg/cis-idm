import{_ as u,r as i,o as r,c as d,a as n,d as s,b as e,w as a,e as c}from"./app-wNEdNvye.js";const m="/cis-idm/assets/notification_settings-jUD5shZw.png",v="/cis-idm/assets/integration_example-CxsnpZeZ.png",k={},b=n("h1",{id:"additional-user-fields",tabindex:"-1"},[n("a",{class:"header-anchor",href:"#additional-user-fields","aria-hidden":"true"},"#"),s(" Additional User Fields")],-1),g=n("p",null,[s("Although cisidm already provides common user information (see "),n("a",{href:"#users"},"Users"),s("), there may be use-cases where additional user metadata needs to be stored. In order to avoid requiring a dedicated service to store this metadata, "),n("code",null,"cisidm"),s(" provides support for custom user fields. Those fields are stored in the database as a "),n("strong",null,"JSON blob"),s(".")],-1),h=n("code",null,"cisidm",-1),f=n("p",null,"To add custom user fields, you need configured them in the configuration file.",-1),y=n("div",{class:"custom-container tip"},[n("p",{class:"custom-container-title"},"Examples"),n("p",null,[s("Head over to "),n("a",{href:"#use-cases-and-examples"},"Use-Cases and Examples"),s(" for some ideas.")])],-1),q=n("br",null,null,-1),_=n("hr",null,null,-1),w=n("p",null,[n("strong",null,"Contents")],-1),x={class:"table-of-contents"},S=c(`<hr><h2 id="field-types" tabindex="-1"><a class="header-anchor" href="#field-types" aria-hidden="true">#</a> Field Types</h2><p><code>cisidm</code> supports the following field types:</p><table><thead><tr><th><strong>Type</strong></th><th><strong>Comment</strong></th></tr></thead><tbody><tr><td><code>string</code></td><td>A text value</td></tr><tr><td><code>number</code></td><td>A number value</td></tr><tr><td><code>bool</code></td><td>A boolean value</td></tr><tr><td><code>date</code></td><td>A date value in the format of <code>YYYY-MM-DD</code></td></tr><tr><td><code>time</code></td><td>A time value in the format of <code>HH:MM</code></td></tr><tr><td><code>object</code></td><td>A JSON object. Allowed properties must be configured in the configuration file using <code>property</code> blocks and may have any of the supported field types.</td></tr><tr><td><code>any</code></td><td>A special field type that instructs <code>cisidm</code> to not perform any validation. <code>any</code> fields cannot be seen or edited in the user-interface.</td></tr><tr><td><code>list</code></td><td>A JSON list. The element-type must be set using the <code>element_type</code> block. List values cannot be seen or edited in the user-interface for now.</td></tr></tbody></table><div class="custom-container warning"><p class="custom-container-title">User-Interface Support</p><p>The user-interface does not yet support <code>list</code> and <code>any</code> types!</p></div><h2 id="field-visibility" tabindex="-1"><a class="header-anchor" href="#field-visibility" aria-hidden="true">#</a> Field Visibility</h2><p>Each field has a visibility specified that instructs <code>cisidm</code> who is allowed to see the field.</p><table><thead><tr><th><strong>Visibility</strong></th><th><strong>Comment</strong></th></tr></thead><tbody><tr><td><code>private</code></td><td>Only administrators/services that integrate with cisidm can see those fields.</td></tr><tr><td><code>self</code></td><td>A user can only see his/her own field</td></tr><tr><td><code>public</code></td><td>Every authenticated user can see the field for each other user</td></tr></tbody></table><div class="custom-container tip"><p class="custom-container-title">Inheritance</p><p>If <code>visibilty</code> or <code>writeable</code> is not defined for <code>object</code> properties than the value of the parent field is inherited. If top-level fields do not specify a visibility it always defaults to <code>self</code>. Fields are writeable by default.</p></div><h2 id="configuration" tabindex="-1"><a class="header-anchor" href="#configuration" aria-hidden="true">#</a> Configuration</h2><p>The definition of a custom user field looks like the following:</p><div class="language-hcl line-numbers-mode" data-ext="hcl"><pre class="language-hcl"><code><span class="token comment"># Possible values for &lt;field-type&gt; are:</span>
<span class="token comment">#  - string</span>
<span class="token comment">#  - number</span>
<span class="token comment">#  - bool</span>
<span class="token comment">#  - date</span>
<span class="token comment">#  - time</span>
<span class="token comment">#  - object</span>
<span class="token comment">#  - list</span>
<span class="token comment">#  - any</span>
field <span class="token string">&quot;&lt;field-type&gt;&quot;</span> <span class="token string">&quot;&lt;field-name&gt;&quot;</span> <span class="token punctuation">{</span>
    <span class="token comment"># Possible values are: &quot;private&quot;, &quot;self&quot;, &quot;public&quot;</span>
    <span class="token property">visibility</span> <span class="token punctuation">=</span> <span class="token string">&quot;&lt;visibility&gt;&quot;</span>

    <span class="token comment"># Whether or not the field is writeable</span>
    <span class="token property">writeable</span> <span class="token punctuation">=</span> <span class="token boolean">true</span>

    <span class="token comment"># A custom display name for the field in the User-Interface.</span>
    <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;&quot;</span>

    <span class="token comment"># A custom description for the field in the User-Interface.</span>
    <span class="token property">description</span> <span class="token punctuation">=</span> <span class="token string">&quot;&quot;</span>

    <span class="token comment"># Property defines an object property.</span>
    <span class="token comment"># Those blocks are only allowed if the field-type is set to &quot;object&quot;.</span>
    <span class="token comment"># A property block may be specified multiple times.</span>
    property <span class="token string">&quot;&lt;property-type&gt;&quot;</span> <span class="token string">&quot;&lt;property-name&gt;&quot;</span> <span class="token punctuation">{</span>
        <span class="token comment"># Defines a object property. This is exactly the same as the &quot;field&quot; block</span>
        <span class="token comment"># so objects, lists and fields may be nested in any way</span>
    <span class="token punctuation">}</span>

    <span class="token comment"># ElementType defines the type of list elements.</span>
    <span class="token comment"># This block is only allowed if the field-type is set to &quot;list&quot;.</span>
    <span class="token comment"># A element_type block must be specified exactly once for &quot;list&quot;s.</span>
    element_type <span class="token string">&quot;&lt;element-type&gt;&quot;</span> <span class="token string">&quot;&lt;element-name&gt;&quot;</span> <span class="token punctuation">{</span>
        <span class="token comment"># Defines a lists element type. This is exactly the same as the &quot;field&quot; block</span>
        <span class="token comment"># so objects, lists and fields may be nested in any way</span>
    <span class="token punctuation">}</span>

    <span class="token comment"># For &quot;string&quot; fields, it is possible to define a set of allowed values using</span>
    <span class="token comment"># multiple \`value\` blocks.</span>
    <span class="token comment"># The &quot;&lt;value&gt;&quot; is the actually stored value for this field while the optional</span>
    <span class="token comment"># display_name stanza may be used to provide a better human value representation.</span>
    <span class="token comment">#</span>
    <span class="token comment"># As soon as at least one \`value\` block is defined \`cisidm\` does not allow</span>
    <span class="token comment"># other values for this field and the user-interface will render a select</span>
    <span class="token comment"># box instead of a simple text field.</span>
    value <span class="token string">&quot;&lt;value&gt;&quot;</span> <span class="token punctuation">{</span>
        <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;&lt;display-name&gt;&quot;</span>
    <span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="use-cases-and-examples" tabindex="-1"><a class="header-anchor" href="#use-cases-and-examples" aria-hidden="true">#</a> Use-Cases and Examples</h2>`,13),A=c('<h3 id="notification-settings" tabindex="-1"><a class="header-anchor" href="#notification-settings" aria-hidden="true">#</a> Notification Settings</h3><p>When integrating your apps directly with <code>cisidm</code>s API you may want to allow your users to update their notification preferences without needing to store those settings directly in your application.</p><p><img src="'+m+`" alt="Settings"></p><p>One can easily add support for notification settings like this:</p><div class="language-hcl line-numbers-mode" data-ext="hcl"><pre class="language-hcl"><code>field <span class="token string">&quot;object&quot;</span> <span class="token string">&quot;notificationSettings&quot;</span> <span class="token punctuation">{</span>
    <span class="token property">visibility</span> <span class="token punctuation">=</span> <span class="token string">&quot;self&quot;</span>
    <span class="token property">writeable</span> <span class="token punctuation">=</span> <span class="token boolean">true</span>
    <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;Notification Settings&quot;</span>
    <span class="token property">description</span> <span class="token punctuation">=</span> <span class="token string">&quot;Manage your notification settings and preferences&quot;</span>

    property <span class="token string">&quot;string&quot;</span> <span class="token string">&quot;newsletter&quot;</span> <span class="token punctuation">{</span>
        <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;Newsletter&quot;</span>
        <span class="token property">description</span> <span class="token punctuation">=</span> <span class="token string">&quot;If and how you want to receive our weekly newsletter&quot;</span>

        value <span class="token string">&quot;never&quot;</span> <span class="token punctuation">{</span>
            <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;Never&quot;</span>
        <span class="token punctuation">}</span>

        value <span class="token string">&quot;email&quot;</span> <span class="token punctuation">{</span>
            <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;E-Mail&quot;</span>
        <span class="token punctuation">}</span>

        value <span class="token string">&quot;sms&quot;</span> <span class="token punctuation">{</span>
            <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;SMS&quot;</span>
        <span class="token punctuation">}</span>

        value <span class="token string">&quot;both&quot;</span> <span class="token punctuation">{</span>
            <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;E-Mail + SMS&quot;</span>
        <span class="token punctuation">}</span>
    <span class="token punctuation">}</span>

    property <span class="token string">&quot;string&quot;</span> <span class="token string">&quot;comments&quot;</span> <span class="token punctuation">{</span>
        <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;Comments &amp; Replies&quot;</span>
        <span class="token property">description</span> <span class="token punctuation">=</span> <span class="token string">&quot;If and how you want to receive notifications about comments and replies&quot;</span>

        value <span class="token string">&quot;never&quot;</span> <span class="token punctuation">{</span>
            <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;Never&quot;</span>
        <span class="token punctuation">}</span>

        value <span class="token string">&quot;email&quot;</span> <span class="token punctuation">{</span>
            <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;E-Mail&quot;</span>
        <span class="token punctuation">}</span>

        value <span class="token string">&quot;sms&quot;</span> <span class="token punctuation">{</span>
            <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;SMS&quot;</span>
        <span class="token punctuation">}</span>

        value <span class="token string">&quot;both&quot;</span> <span class="token punctuation">{</span>
            <span class="token property">display_name</span> <span class="token punctuation">=</span> <span class="token string">&quot;E-Mail + SMS&quot;</span>
        <span class="token punctuation">}</span>
    <span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>Your application than just needs to check the <code>extra</code> property of the user profile. For example, using <code>idmctl</code>:</p><div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code>$ idmctl <span class="token function">users</span> get-extra alice <span class="token string">&quot;notificationSettings.newsletter&quot;</span>
<span class="token string">&quot;never&quot;</span>

$ idmctl <span class="token function">users</span> get-extra alice <span class="token string">&quot;notificationSettings&quot;</span>
<span class="token punctuation">{</span> <span class="token string">&quot;newsletter&quot;</span><span class="token builtin class-name">:</span> <span class="token string">&quot;never&quot;</span>, <span class="token string">&quot;comments&quot;</span><span class="token builtin class-name">:</span> <span class="token string">&quot;email&quot;</span> <span class="token punctuation">}</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>When using <code>Go</code>, an example might look like the following:</p><div class="language-go line-numbers-mode" data-ext="go"><pre class="language-go"><code><span class="token keyword">package</span> newsletter

<span class="token comment">// A JSON-Path under which the newsletter settings are stored.</span>
<span class="token keyword">const</span> newsletterSettingsKey <span class="token operator">=</span> <span class="token string">&quot;notificationSettings.newsletter&quot;</span>

<span class="token keyword">var</span> <span class="token punctuation">(</span>
    notify idmv1<span class="token punctuation">.</span>NotifyServiceClient
    users idmv1<span class="token punctuation">.</span>UserServiceClient
<span class="token punctuation">)</span>

<span class="token keyword">func</span> <span class="token function">init</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
    <span class="token comment">// adjust client init to your liking ...</span>
    <span class="token comment">// this just uses constant URLs</span>
    notify <span class="token operator">:=</span> idmv1connect<span class="token punctuation">.</span><span class="token function">NewNotifyServiceClient</span><span class="token punctuation">(</span><span class="token string">&quot;https://account.example.com&quot;</span><span class="token punctuation">)</span>
    users <span class="token operator">:=</span> idmv1connect<span class="token punctuation">.</span><span class="token function">NewUserServiceClient</span><span class="token punctuation">(</span><span class="token string">&quot;https://account.example.com&quot;</span><span class="token punctuation">)</span>
<span class="token punctuation">}</span>

<span class="token keyword">func</span> <span class="token function">SendNotification</span><span class="token punctuation">(</span>ctx context<span class="token punctuation">.</span>Context<span class="token punctuation">,</span> userId <span class="token builtin">string</span><span class="token punctuation">,</span> message <span class="token builtin">string</span><span class="token punctuation">)</span> <span class="token builtin">error</span> <span class="token punctuation">{</span>
    <span class="token comment">// get user preferences</span>
    res<span class="token punctuation">,</span> err <span class="token operator">:=</span> users<span class="token punctuation">.</span><span class="token function">GetUserExtraKey</span><span class="token punctuation">(</span>ctx<span class="token punctuation">,</span> connect<span class="token punctuation">.</span><span class="token function">NewRequest</span><span class="token punctuation">(</span>idmv1<span class="token punctuation">.</span>GetUserExtraKeyRequest<span class="token punctuation">{</span>
        userId<span class="token punctuation">:</span> userId<span class="token punctuation">,</span>
        path<span class="token punctuation">:</span> newsletterSettingsKey<span class="token punctuation">,</span>
    <span class="token punctuation">}</span><span class="token punctuation">)</span><span class="token punctuation">)</span>

    <span class="token comment">// Handle any request error</span>
    <span class="token keyword">if</span> err <span class="token operator">!=</span> <span class="token boolean">nil</span> <span class="token punctuation">{</span>
        <span class="token keyword">return</span> err
    <span class="token punctuation">}</span>

    <span class="token comment">// Make sure the property value is actually a string.</span>
    val<span class="token punctuation">,</span> ok <span class="token operator">:=</span> res<span class="token punctuation">.</span>Msg<span class="token punctuation">.</span>Value<span class="token punctuation">.</span><span class="token punctuation">(</span><span class="token operator">*</span>structpb<span class="token punctuation">.</span>StringValue<span class="token punctuation">)</span>
    <span class="token keyword">if</span> <span class="token operator">!</span>ok <span class="token punctuation">{</span>
        <span class="token keyword">return</span> errors<span class="token punctuation">.</span><span class="token function">New</span><span class="token punctuation">(</span><span class="token string">&quot;invalid customer user property value, expected string, got %T&quot;</span><span class="token punctuation">,</span> res<span class="token punctuation">.</span>Msg<span class="token punctuation">.</span>Value<span class="token punctuation">)</span>
    <span class="token punctuation">}</span>

    <span class="token comment">// Decide what to do based on the string value</span>
    <span class="token keyword">switch</span> val<span class="token punctuation">.</span>StringValue <span class="token punctuation">{</span>
        <span class="token keyword">case</span> <span class="token string">&quot;never&quot;</span><span class="token punctuation">:</span>
            <span class="token keyword">return</span> <span class="token boolean">nil</span>

        <span class="token keyword">case</span> <span class="token string">&quot;sms&quot;</span><span class="token punctuation">:</span>
            <span class="token comment">// Send an SMS</span>

        <span class="token keyword">case</span> <span class="token string">&quot;email&quot;</span><span class="token punctuation">:</span>
            <span class="token comment">// Send an Email</span>

        <span class="token keyword">case</span> <span class="token string">&quot;&quot;</span><span class="token punctuation">,</span> <span class="token string">&quot;both&quot;</span><span class="token punctuation">:</span> <span class="token comment">// we consider an unset-value as &quot;both&quot;</span>
            <span class="token comment">// Send E-Mail and SMS</span>
    <span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h3 id="external-service-integration" tabindex="-1"><a class="header-anchor" href="#external-service-integration" aria-hidden="true">#</a> External Service Integration</h3><p>Imagine your app having support to send messages to your users on Discord or Reddit but those integrations must be enabled by each user on it&#39;s own:</p><p><img src="`+v+'" alt="Integration Example"></p>',12),C=n("div",{class:"language-hcl line-numbers-mode","data-ext":"hcl"},[n("pre",{class:"language-hcl"},[n("code",null,[s("field "),n("span",{class:"token string"},'"object"'),s(),n("span",{class:"token string"},'"integrations"'),s(),n("span",{class:"token punctuation"},"{"),s(`
    `),n("span",{class:"token property"},"display_name"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Integrations"'),s(`
    `),n("span",{class:"token property"},"description"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Configure integrations with external services"'),s(`

    property `),n("span",{class:"token string"},'"object"'),s(),n("span",{class:"token string"},'"discord"'),s(),n("span",{class:"token punctuation"},"{"),s(`
        `),n("span",{class:"token property"},"display_name"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Discord"'),s(`
        `),n("span",{class:"token property"},"description"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Link your discord account"'),s(`

        property `),n("span",{class:"token string"},'"string"'),s(),n("span",{class:"token string"},'"id"'),s(),n("span",{class:"token punctuation"},"{"),s(`
            `),n("span",{class:"token property"},"display_name"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Discord ID"'),s(`
            `),n("span",{class:"token property"},"description"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Enter your discord user ID to link your account"'),s(`
        `),n("span",{class:"token punctuation"},"}"),s(`

        property `),n("span",{class:"token string"},'"bool"'),s(),n("span",{class:"token string"},'"notify"'),s(),n("span",{class:"token punctuation"},"{"),s(`
            `),n("span",{class:"token property"},"display_name"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Notify on Discord"'),s(`
            `),n("span",{class:"token property"},"description"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Do you want to be notified on discord"'),s(`
        `),n("span",{class:"token punctuation"},"}"),s(`
    `),n("span",{class:"token punctuation"},"}"),s(`

    property `),n("span",{class:"token string"},'"object"'),s(),n("span",{class:"token string"},'"reddit"'),s(),n("span",{class:"token punctuation"},"{"),s(`
        `),n("span",{class:"token property"},"display_name"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Reddit"'),s(`
        `),n("span",{class:"token property"},"description"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Link your Reddit account"'),s(`

        property `),n("span",{class:"token string"},'"string"'),s(),n("span",{class:"token string"},'"handle"'),s(),n("span",{class:"token punctuation"},"{"),s(`
            `),n("span",{class:"token property"},"display_name"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Reddit User-Handle"'),s(`
            `),n("span",{class:"token property"},"description"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Enter your Reddit username to link your account"'),s(`
        `),n("span",{class:"token punctuation"},"}"),s(`

        property `),n("span",{class:"token string"},'"bool"'),s(),n("span",{class:"token string"},'"notify"'),s(),n("span",{class:"token punctuation"},"{"),s(`
            `),n("span",{class:"token property"},"display_name"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Notify on Reddit"'),s(`
            `),n("span",{class:"token property"},"description"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Do you want to be notified on Reddit"'),s(`
        `),n("span",{class:"token punctuation"},"}"),s(`
    `),n("span",{class:"token punctuation"},"}"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`
`)])]),n("div",{class:"line-numbers","aria-hidden":"true"},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),N=n("h3",{id:"attribute-based-access-control-abac",tabindex:"-1"},[n("a",{class:"header-anchor",href:"#attribute-based-access-control-abac","aria-hidden":"true"},"#"),s(" Attribute-Based-Access-Control (ABAC)")],-1),I=n("p",null,"Custom user fields can also be used to implement ABAC:",-1),E=n("div",{class:"language-hcl line-numbers-mode","data-ext":"hcl"},[n("pre",{class:"language-hcl"},[n("code",null,[s("field "),n("span",{class:"token string"},'"string"'),s(),n("span",{class:"token string"},'"department"'),s(),n("span",{class:"token punctuation"},"{"),s(`
    `),n("span",{class:"token comment"},"# Users can see their assigned department"),s(`
    `),n("span",{class:"token property"},"visibility"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"self"'),s(`

    `),n("span",{class:"token comment"},"# But are not allowed to change them (only administrators can set those"),s(`
    `),n("span",{class:"token comment"},"# fields)"),s(`
    `),n("span",{class:"token property"},"writeable"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token boolean"},"false"),s(`

    `),n("span",{class:"token property"},"display_name"),s(),n("span",{class:"token punctuation"},"="),s(),n("span",{class:"token string"},'"Department"'),s(`

    value `),n("span",{class:"token string"},'"HR"'),s(),n("span",{class:"token punctuation"},"{"),n("span",{class:"token punctuation"},"}"),s(`
    value `),n("span",{class:"token string"},'"Support"'),s(),n("span",{class:"token punctuation"},"{"),n("span",{class:"token punctuation"},"}"),s(`
    value `),n("span",{class:"token string"},'"Development"'),s(),n("span",{class:"token punctuation"},"{"),n("span",{class:"token punctuation"},"}"),s(`
    value `),n("span",{class:"token string"},'"Management"'),s(),n("span",{class:"token punctuation"},"{"),n("span",{class:"token punctuation"},"}"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`
`)])]),n("div",{class:"line-numbers","aria-hidden":"true"},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),j=n("p",null,[s("Then one can use those fields in rego policies (See our "),n("a",{href:"./policies"},"Policies Guide"),s("):")],-1),M=n("div",{class:"language-rego line-numbers-mode","data-ext":"rego"},[n("pre",{class:"language-rego"},[n("code",null,[n("span",{class:"token keyword"},"package"),s(" cisidm"),n("span",{class:"token punctuation"},"."),s(`forward_auth

`),n("span",{class:"token keyword"},"import"),s(" rego"),n("span",{class:"token punctuation"},"."),s(`v1

`),n("span",{class:"token keyword"},"default"),s(" allow "),n("span",{class:"token operator"},":="),s(),n("span",{class:"token boolean"},"false"),s(`

allow if `),n("span",{class:"token punctuation"},"{"),s(`
    input`),n("span",{class:"token punctuation"},"."),s(`subject
    input`),n("span",{class:"token punctuation"},"."),s("subject"),n("span",{class:"token punctuation"},"."),s("fields"),n("span",{class:"token punctuation"},"["),n("span",{class:"token string"},'"department"'),n("span",{class:"token punctuation"},"]"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token string"},'"HR"'),s(`
    input`),n("span",{class:"token punctuation"},"."),s("host "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token string"},'"hr.example.com"'),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

allow if `),n("span",{class:"token punctuation"},"{"),s(`
    input`),n("span",{class:"token punctuation"},"."),s(`subject
    input`),n("span",{class:"token punctuation"},"."),s("subject"),n("span",{class:"token punctuation"},"."),s("fields"),n("span",{class:"token punctuation"},"["),n("span",{class:"token string"},'"department"'),n("span",{class:"token punctuation"},"]"),s(" in "),n("span",{class:"token punctuation"},"["),n("span",{class:"token string"},'"Development"'),n("span",{class:"token punctuation"},","),s(),n("span",{class:"token string"},'"Support"'),n("span",{class:"token punctuation"},"]"),s(`

    input`),n("span",{class:"token punctuation"},"."),s("host in "),n("span",{class:"token punctuation"},"["),n("span",{class:"token string"},'"prod.example.com"'),n("span",{class:"token punctuation"},","),s(),n("span",{class:"token string"},'"dev.example.com"'),n("span",{class:"token punctuation"},"]"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`
`)])]),n("div",{class:"line-numbers","aria-hidden":"true"},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),D=c(`<h2 id="setting-fields" tabindex="-1"><a class="header-anchor" href="#setting-fields" aria-hidden="true">#</a> Setting fields</h2><p>While users can update <code>writeable</code> non-<code>private</code> fields via the user-interface, as an administrator, you may set all fields using the idmctl commandline utility:</p><div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code>idmctl <span class="token function">users</span> set-extra <span class="token punctuation">[</span>user<span class="token punctuation">]</span> <span class="token punctuation">[</span>path<span class="token punctuation">]</span> <span class="token punctuation">[</span>value<span class="token punctuation">]</span>

<span class="token comment"># Examples:</span>
idmctl <span class="token function">users</span> set-extra alice <span class="token string">&quot;notification-settings.sms&quot;</span> <span class="token boolean">true</span>
idmctl <span class="token function">users</span> set-extra alice <span class="token string">&quot;internal-phone-extenstion&quot;</span> <span class="token string">&#39;&quot;34&quot;&#39;</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><div class="custom-container tip"><p class="custom-container-title">Note</p><p>Values should be encoded in their JSON representation. For example, the following should also work:</p><div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code>idmctl <span class="token function">users</span> set-extra alice                    <span class="token punctuation">\\</span>
    <span class="token string">&quot;notification-settings&quot;</span>                     <span class="token punctuation">\\</span>
    <span class="token string">&#39;{&quot;sms&quot;: true, &quot;email&quot;: false}&#39;</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div></div>`,4);function R(T,U){const p=i("RouterLink"),t=i("router-link"),o=i("CodeGroupItem"),l=i("CodeGroup");return r(),d("div",null,[b,g,n("p",null,[s("Custom fields may either be used by services that directly integrate with "),h,s(" using the API or may also be used in "),e(p,{to:"/guides/policies.html"},{default:a(()=>[s("Policies")]),_:1}),s(" to implement Attribute-Based-Access-Control (ABAC).")]),f,y,q,_,w,n("nav",x,[n("ul",null,[n("li",null,[e(t,{to:"#field-types"},{default:a(()=>[s("Field Types")]),_:1})]),n("li",null,[e(t,{to:"#field-visibility"},{default:a(()=>[s("Field Visibility")]),_:1})]),n("li",null,[e(t,{to:"#configuration"},{default:a(()=>[s("Configuration")]),_:1})]),n("li",null,[e(t,{to:"#use-cases-and-examples"},{default:a(()=>[s("Use-Cases and Examples")]),_:1}),n("ul",null,[n("li",null,[e(t,{to:"#notification-settings"},{default:a(()=>[s("Notification Settings")]),_:1})]),n("li",null,[e(t,{to:"#external-service-integration"},{default:a(()=>[s("External Service Integration")]),_:1})]),n("li",null,[e(t,{to:"#attribute-based-access-control-abac"},{default:a(()=>[s("Attribute-Based-Access-Control (ABAC)")]),_:1})])])]),n("li",null,[e(t,{to:"#setting-fields"},{default:a(()=>[s("Setting fields")]),_:1})])])]),S,n("p",null,[s("Below are a few examples on how to configure custom fields. Refer to the "),e(p,{to:"/guides/config-reference.html"},{default:a(()=>[s("Configuration File Reference")]),_:1}),s(" for more information.")]),A,e(l,null,{default:a(()=>[e(o,{title:"config.hcl"},{default:a(()=>[C]),_:1})]),_:1}),N,I,e(l,null,{default:a(()=>[e(o,{title:"config.hcl"},{default:a(()=>[E]),_:1})]),_:1}),j,e(l,null,{default:a(()=>[e(o,{title:"policy.rego"},{default:a(()=>[M]),_:1})]),_:1}),D])}const V=u(k,[["render",R],["__file","extra-user-fields.html.vue"]]);export{V as default};
