import { defaultTheme, defineUserConfig } from "vuepress";

export default defineUserConfig({
  lang: "en-US",
  title: "Identity Management Service",
  base: "/cis-idm/",
  description: "",
  theme: defaultTheme({
    repo: "tierklinik-dobersberg/cis-idm",
    docsRepo: "https://github.com/tierklinik-dobersberg/cis-idm",
    docsBranch: "main",
    docsDir: "docs/content",
    editLinkPattern: ":repo/-/edit/:branch/:path",
    sidebar: {
      "/guides/": [
        {
          text: "Introduction",
          link: "/intro.md",
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
              text: "CLI Reference",
              link: "/guides/cli-reference.md",
            },
          ],
        },
        {
          text: "Setup",
          link: "/guides/README.md",
          children: [
            {
              text: "WebAuthN",
            },
            {
              text: "E-Mail",
            },
            {
              text: "SMS (Twilio)",
            },
            {
              text: "WebPush Notifications",
            },
            {
              text: "OpenID Connect (OIDC)",
            },
            {
              text: "Configuration File Reference",
            },
          ],
        },
        {
          text: "Developer Docs",
          children: [
            {
              text: "ConnectRPC APIs",
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
        text: "Guides",
        link: "/guides/",
      },
      {
        text: "Report an issue",
        link: "https://github.com/tierklinik-dobersberg/cis-idm/issues",
      },
    ],
  }),
});
