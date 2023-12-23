import { defaultTheme, defineUserConfig } from "vuepress";

export default defineUserConfig({
  lang: "en-US",
  title: "Identity Management Service",
  description: "",
  theme: defaultTheme({
    repo: "tierklinik-dobersberg/cis-idm",
    docsRepo: "https://github.com/tierklinik-dobersberg/cis-idm",
    docsBranch: "main",
    docsDir: "docs/content",
    editLinkPattern: ":repo/-/edit/:branch/:path",
    navbar: [
      {
        text: "Get Started",
        link: "/guides/getting-started",
      },
      {
        text: "Architecture",
        link: "architecture",
      },
      {
        text: "Report an issue",
        link: "https://github.com/tierklinik-dobersberg/cis-idm/issues",
      },
    ],
  }),
});
