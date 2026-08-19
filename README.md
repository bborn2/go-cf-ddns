# go-cf-ddns

A lightweight Go-based DDNS client that updates a Cloudflare DNS record to match the current public IPv4 address.

## Why

GoDaddy has disabled API access without prior notice, and the official documentation is difficult to find. In response, this project provides a simple, free alternative: keep a DNS zone on Cloudflare and update the record automatically when the public IP changes.

The typical migration flow is:

1. Export the zone file from GoDaddy
2. Change the nameservers to Cloudflare
3. Disable DNSSEC in GoDaddy if required
4. Import the zone file into Cloudflare

For Cloudflare DNS management guidance, see: https://developers.cloudflare.com/dns/manage-dns-records/how-to/

## Overview

This project is designed for environments where the public IP changes frequently and the DNS record should stay in sync automatically. It retrieves the current public IPv4 address from a public IP service, compares it with the Cloudflare DNS record, and updates the record if needed.

## Features

- Retrieves the current public IPv4 address
- Reads the configured Cloudflare DNS record
- Updates the record when the IP changes
- Supports a simple CLI for token, zone, and domain configuration
- Suitable for cron or periodic scheduled execution

## Requirements

- Go 1.21 or newer
- A Cloudflare API token with DNS edit permissions
- The target Cloudflare zone ID
- A DNS record name or root domain

## Build

```bash
go build -o cfddns .
```

## Usage

```bash
./cfddns --token YOUR_CLOUDFLARE_TOKEN --zoneid YOUR_ZONE_ID --domain example.com
```

### Optional flags

- `--token`: Cloudflare API token
- `--zoneid`: Cloudflare zone ID
- `--domain`: Domain or record name to update; default is `@`
- `-v`: Print build metadata

## Cloudflare Token Setup

Create a Cloudflare API token at:

https://dash.cloudflare.com/profile/api-tokens

Use a token with permission to edit DNS records for the target zone.

## Zone ID Lookup

```bash
curl --request GET \
  --url https://api.cloudflare.com/client/v4/zones \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer YOUR_TOKEN'
```

From the response JSON, copy the `id` value for the target zone.

## Deployment Example

This is best run on a schedule using `cron` or a systemd timer, for example every 5 minutes:

```bash
*/5 * * * * /usr/local/bin/cfddns --token YOUR_TOKEN --zoneid YOUR_ZONE_ID --domain example.com >> /var/log/cfddns.log 2>&1
```

## Notes

This utility is intentionally simple and focused on a single job: keep a Cloudflare A record aligned with the current public IPv4 address.

For more information about Cloudflare DNS management, see:
https://developers.cloudflare.com/dns/manage-dns-records/
