# domain-info
Simple tool to query info about one or more domains.

## Install
```bash
git clone https://www.github.com/Toxyl/domain-info
cd domain-info
CGO_ENABLED=0 go build .
cp domain-info /usr/local/bin/
mkdir -p /etc/domain-info/
cp config.example.yaml /etc/domain-info/config.yaml
domain-info
```

## Setting up DNS servers
The default config contains the DNS servers of Google, Level 3, Verisign, DNS Watch, Comodo, OpenDNS and Norton. These are used to spread DNS requests, i.e. adding more servers results in more spread. The config file is expected to reside in `/etc/domain-info/config.yaml`, `$HOME/.domain-info/config.yaml` or `.` (i.e. the current directory). 

## Get domain details
```bash
domain-info microsoft.com google.com
```

Or combine it with a bruteforce check using ~100k subdomains from [TheRook's subbrute list](https://github.com/TheRook/subbrute).
```bash
domain-info brute microsoft.com google.com
```

## Get list of (sub)domains
Will list all domains and subdomains found for the given input. (Sub)Domains that do not have a CNAME will be prefixed with #.
```bash
domain-info list microsoft.com google.com
```

Or combine it with a bruteforce check using ~100k subdomains from [TheRook's subbrute list](https://github.com/TheRook/subbrute).
```bash
domain-info brute microsoft.com google.com
```

## Note about bruteforce
If the domain is setup as catch-all, chances are that every entry on the list will be a match. In that case the results are, obviously, useless.