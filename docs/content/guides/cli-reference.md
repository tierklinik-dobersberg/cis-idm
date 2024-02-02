---
next:
  text: Configuration File Reference
  link: ./config-reference.md
---

# CLI Reference

For the time being, `cisidm` does not include a admin web-interface so configuration of roles, permissions as well as user-invitation and creation must be done using the provided command line tools `idmctl`.

The cli tool is based on cli package of [`tierklinik-dobersberg/apis`](https://github.com/tierklinik-dobersberg/apis) and supports a configuration file located at `~/.config/cis/config.yml`.

<CodeGroup>
  <CodeGroupItem title="~/.config/cis/config.yml">

```yaml
default:
  # If self-signed certificates are in use you may set insecure to true to skip certificate
  # validation for HTTPS endpoints.
  insecure: false

  urls:
    # The URL of cisidm:
    idm: https://account.example.com
```

  </CodeGroupItem>
</CodeGroup>

:::tip Multiple Deployments
It's possible to use the `idmctl` cli tool with multiple deployments. Just add multiple settings to `.config/cis/config.yml`, each under a distinct key. When invoking the command line tool, use `--configuration=your-name` to specify which environment to use. If the `--configuration` parameter is not set, `idmctl` will use the configuration under the `default` key.
:::

The cli tool will store any access and refresh tokens next to the configuration file. **It's important to keep those files secret. Do not commit them to a source code repository or otherwise enable public access**.
