# Proxy server for ACME DNS-01 challenges

**dns01proxy** is a server for using DNS-01 challenges to get TLS/SSL
certificates from Let's Encrypt, or any ACME-compatible certificate authority,
without exposing your DNS credentials to every host that needs a certificate.

It acts as a proxy for DNS-01 challenge requests, allowing hosts to delegate
their DNS record updates during ACME validation. This makes it possible to
issue certificates to internal or private hosts that can't (or shouldn't) have
direct access to your DNS provider or API keys.

dns01proxy is designed to work with:
* [acme.sh](https://acme.sh/)'s
  [`acmeproxy`](https://github.com/acmesh-official/acme.sh/wiki/dnsapi2#dns_acmeproxy)
  provider,
* [lego](https://go-acme.github.io/lego/)'s
  [`httpreq`](https://go-acme.github.io/lego/dns/httpreq/index.html) DNS
  provider, and
* [Caddy](https://caddyserver.com/)'s
  [`acmeproxy`](https://caddyserver.com/docs/modules/dns.providers.acmeproxy)
  DNS provider module.

## Features

* Privilege separation. Internal or private hosts can complete DNS-01
  challenges without having direct access to DNS API keys. In turn, the private
  keys for issued certificates stay private to the ACME clients.
* HTTPS built-in, automatic, and always on. dns01proxy uses its configured DNS
  credentials to automatically get and renew its own TLS/SSL certificate.
* Mandatory client authentication using HTTP Basic Authentication.
* Optional per-client policies for limiting which DNS names each client can get
  a certificate for.

## Installing dns01proxy

There are two options for getting dns01proxy.

### Pre-compiled binaries

dns01proxy is built using Caddy, and uses DNS provider modules that are written
by the Caddy community. dns01proxy ships a number of binaries, each built with
a single DNS module. To install, just download a build of the [latest
release](https://github.com/liujed/dns01proxy/releases) that matches your DNS
provider.

> [!CAUTION]
> Always check that you trust the author of the DNS module. The
> [release notes](https://github.com/liujed/dns01proxy/releases) has details
> about the source of the DNS module in each build.

### Caddy module

Alternatively, dns01proxy is also available as a Caddy module, which adds
dns01proxy to [Caddy](https://caddyserver.com/) as a subcommand, app, and HTTP
handler. See the [caddy-dns01proxy](https://github.com/liujed/caddy-dns01proxy)
project for more on this second option.

## Configuring dns01proxy

dns01proxy is configured through a single JSON file. Below is an example
configuration for running at `https://dns01proxy.example.com` with Cloudflare
as a DNS provider. A few things are worth noting about this example:
* The user's password is hashed using `dns01proxy hash-password`.
* Environment variables can be referenced using the `{env.VAR_NAME}` syntax.
* Each DNS provider has a different set of configuration parameters. See the
  documentation link for your provider in the [release
  notes](https://github.com/liujed/dns01proxy/releases).

```json
{
  "hostnames": ["dns01proxy.example.com"],
  "listen": [":443"],
  "dns": {
    "provider": {
      "name": "cloudflare",
      "api_token": "{env.CF_API_TOKEN}"
    }
  },
  "accounts": [
    {
      "username": "AzureDiamond",
      "password": "$2a$14$N5bGBXf7zwAW9Ym7IQ/mxOHTGsvFNOTEAiN4/r1LnvfzYCpiWcHOa",
      "allow_domains": ["private.example.com"]
    }
  ]
}
```

<details>
<summary>Full JSON structure</summary>

```jsonc
{
  // The server's hostnames. Used for obtaining TLS/SSL certificates.
  "hostnames": ["<hostname>"],

  // The sockets on which to listen.
  "listen": ["<ip_addr:port>"],

  // Configures the set of trusted proxies, for accurate logging of client IP
  // addresses.
  "trusted_proxies": {
    // An `http.ip_sources` Caddy module.
    "source": "<module_name>",
    // ...
  },

  "dns": {
    // The DNS provider for publishing DNS-01 responses.
    "provider": {
      // A `dns.providers` Caddy module.
      "name": "<provider_name>",

      // Module configuration. See the documentation link for your provider in
      // the release notes: https://github.com/liujed/dns01proxy/releases

      // ... 
    },

    // The TTL to use in DNS TXT records. Optional. Not usually needed.
    "ttl": "<ttl>",  // e.g., "2m"

    // Custom DNS resolvers to prefer over system or built-in defaults. Set
    // this to a public resolver if you are using split-horizon DNS.
    "resolvers": ["<resolver>"]
  },

  // Configures HTTP basic authentication and the domains for which each user
  // can get TLS/SSL certificates.
  "accounts": [
    {
      "user_id": "<userID>",

      // To hash passwords, use `dns01proxy hash-password`.
      "password": "<hashed_password>",

      // These largely follow Smallstep's domain name rules:
      //
      //   https://smallstep.com/docs/step-ca/policies/#domain-names
      //
      // Due to a limitation in ACME and DNS-01, allowing a domain also allows
      // wildcard certificates for that domain.
      "allow_domains": ["<domain>"],
      "deny_domains": ["<domain>"]
    }
  ]
}
```

</details>

## Running dns01proxy

To run dns01proxy, use the `run` subcommand. For example,
```
dns01proxy run --config /usr/local/etc/dns01proxy.json
```

## Integrating with acme.sh

dns01proxy works with [acme.sh](https://acme.sh/)'s
[`acmeproxy`](https://github.com/acmesh-official/acme.sh/wiki/dnsapi2#dns_acmeproxy)
provider:
```sh
export ACMEPROXY_ENDPOINT='https://dns01proxy.example.com'
export ACMEPROXY_USERNAME='AzureDiamond'
export ACMEPROXY_PASSWORD='hunter2'
acme.sh --issue --dns dns_acmeproxy -d example.com
```

## Integrating with lego

dns01proxy works with [lego](https://go-acme.github.io/lego/)'s
[`httpreq`](https://go-acme.github.io/lego/dns/httpreq/index.html) DNS
provider:
```sh
export HTTPREQ_ENDPOINT='https://dns01proxy.example.com'
export HTTPREQ_USERNAME='AzureDiamond'
export HTTPREQ_PASSWORD='hunter2'
lego --email you@example.com --dns httpreq -d example.com run
```

## Integrating with Caddy

dns01proxy works with [Caddy](https://caddyserver.com/)'s
[`acmeproxy`](https://caddyserver.com/docs/modules/dns.providers.acmeproxy) DNS
provider module:
```json
{
  "endpoint": "https://dns01proxy.example.com",
  "username": "AzureDiamond",
  "password": "hunter2"
}
```

## Acknowledgements

dns01proxy is a reimplementation of
[acmeproxy](https://github.com/mdbraber/acmeproxy/), which is no longer being
developed. Whereas acmeproxy was built on top of lego, dns01proxy uses
[libdns](https://github.com/libdns/libdns) under the hood, which allows for
better compatibility with acme.sh.

[acmeproxy.pl](https://github.com/madcamel/acmeproxy.pl) is another
reimplementation of acmeproxy, written in Perl.
