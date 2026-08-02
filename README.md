# Cloudflare DDNS
This project is not related to cloudflare at all. This is a **personal** project for something **I** personally need. Feel free to use it though.

## Configuration & Quick Setup
You can create the default config file running
```
./cloudflare-ddns --init
```

If you're on Windows, it expects the config to be in the same folder as the program (**./cloudflare-ddns.conf**).<br>
If you're on Linux or Mac, it expects the config to be in **/etc/cloudflare-ddns.conf**

All you have to do now is configure the API Token and Zone ID and run
```
./cloudflare-ddns
```

It will only run if your public ip has changed. To force an update, you can use the --force flag.

```
./cloudflare-ddns --force
```

### Default config
```conf
api_token=[APITOKEN]
zone_id=[ZONEID]

# an empty value for last_ipv4 and last_ipv6 will let the script update the records on the first run
# if you set the value to "none", "disabled" or "no", it's gonna be disabled
# if ipv4 is disabled, A records will be ignored
# if ipv6 is disabled, AAAA records will be ignored
last_ipv4=
last_ipv6=none

# supported types: A, AAAA, CNAME, TXT
# ttl supported values = 1min, 2min, 5min, 10min, 15min, 30min, 1h, 2h, 5h, 12h, 1d, auto
# accepted value vars: {public_ipv4} {public_ipv6}

# A record is defined like so:
# record=type,name,value,proxied?,ttl?
#
# proxied and ttl are optional
# proxied's default is true
# ttl's default is auto

# This is an example A record, with the name "ddns", pointing to your public_ipv4. Proxied is defaulted to true. TTL is defaulted to "auto".
record=A,ddns,{public_ipv4}

# In this example below, proxy is disabled
#record=A,ddns,{public_ipv4},false
```