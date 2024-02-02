import { defaultTheme, defineUserConfig } from "vuepress";

export default defineUserConfig({
  lang: "en-US",
  title: "CIS-IDM",
  base: "/cis-idm/",
  description: "",
  theme: defaultTheme({
    repo: "tierklinik-dobersberg/cis-idm",
    docsRepo: "https://github.com/tierklinik-dobersberg/cis-idm",
    docsBranch: "main",
    docsDir: "docs/content",
    editLinkPattern: ":repo/edit/:branch/:path",
    sidebar: {
      "/guides/": [
        {
          text: "Introduction",
          link: "/guides/intro.md",
        },
        {
          text: "How Tos",
          children: [
            {
              text: "Getting Started",
              link: "/guides/getting-started.md",
            },
            {
              text: "User and Role Management",
              link: "/guides/user-role-management.md",
            },
            {
              text: "Additional User Fields",
              link: "/guides/extra-user-fields.md",
            },
            {
              text: "Policies",
              link: "/guides/policies.md"
            },
            {
              text: "CLI Reference",
              link: "/guides/cli-reference.md",
            },
            {
              text: "Configuration File Reference",
              link: "/guides/config-reference.md"
            },
          ],
        },
        {
          text: "Setup",
          link: "/guides/README.md",
          children: [
            {
              text: "WebAuthN",
              link: "/guides/setup-webauthn.md"
            },
            {
              text: "E-Mail",
              link: "/guides/setup-email.md"
            },
            {
              text: "SMS (Twilio)",
              link: "/guides/setup-sms.md"
            },
            {
              text: "WebPush Notifications",
              link: "/guides/setup-webpush.md"
            },
            {
              text: "OpenID Connect (OIDC)",
              link: "/guides/setup-oidc.md"
            },
          ],
        },
        {
          text: "Developer Documentation",
          collapsible: true,
          children: [
            {
              text: "Setup",
            },
            {
              text: "Architecture",
            },
            {
              text: "Connect RPCs and APIs",
              collapsible: true,
              children: [
                {
                  text: "with Go",
                },
                {
                  text: "with Javascript/Typescript",
                },
                {
                  text: "with CURL",
                },
              ],
            },
          ],
        },
      ],
    },
    navbar: [
      {
        text: "Documentation",
        link: "/guides/",
      },
      {
        text: "Report an issue",
        link: "https://github.com/tierklinik-dobersberg/cis-idm/issues",
      },
    ],
  }),
});
