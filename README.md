# Cloudflare Ddns

## Why
Godaddy has “disabled” API access without any prior notification，Official documentation cannot even be found。 WTF！

## Free alternative solutions  （10 mins)
created a free DNS at Cloudflare and then pointed the Godaddy nameservers to it。
1. export zone file. (hub.godaddy.com-domain-Dns-Dns Records-action)
2. change Nameservers 
3. turn off dnssec in godaddy
4. import zone file in cloudflare.

[Learn how to use Cloudflare DNS to manage your DNS records](https://developers.cloudflare.com/dns/manage-dns-records/how-to/)

## Usage
```
cfddns --token xxx --zoneid xxx
```

## How to get token
create api token: https://dash.cloudflare.com/profile/api-tokens, use template for Edit zone DNS

## How to get zoneid
https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-list-dns-records
```
curl --request GET \
  --url https://api.cloudflare.com/client/v4/zones/zone_id/dns_records \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer xxxtoken'
```
Copy the zone_id from the return json
