# Cloudflare DDNS
This project is not related to cloudflare at all. This is a **personal** project for something **I** personally need. Feel free to use it though.

## Configuration & Quick Setup
The name of the file must be cloudflare-ddns.conf

If you're on Windows, it expects the config to be in the same folder as the program (**./cloudflare-ddns.conf**).<br>
If you're on Linux or Mac, it expects the config to be in **/etc/cloudflare-ddns.conf**

### Example config
```conf
api_token=[APITOKEN]
zone_id=[ZONEID]

# supported types: A, AAAA, CNAME, TXT
# ttl supported values = 1min, 2min, 5min, 10min, 15min, 30min, 1h, 2h, 5h, 12h, 1d, auto
# accepted value vars: {public_ipv4} {public_ipv6}

# record=type,name,value,ttl,proxied
record=A,ddns,{public_ipv4},auto,false
```