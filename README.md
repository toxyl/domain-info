# domain-info
Simple tool to query info about one or more domains.

## Install
```bash
git clone https://www.github.com/Toxyl/domain-info
cd domain-info
go build .
cp domain-info /usr/local/bin/
domain-info
```

## Get domain details
```bash
domain-info microsoft.com google.com
```

## Get list of domains
Will list all domains and subdomains found for the given input. (Sub)Domains that do not have a CNAME will be prefixed with #.
```bash
domain-info list microsoft.com google.com
```

