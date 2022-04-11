package main

import (
	"fmt"
	"os"
)

var listMode bool = false
var bruteforceMode bool = false
var resolvers *DNSResolverRing

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(
			`domain-info v%.2f
Usage:   %s            [domain1] .. [domainN]
         %s brute      [domain1] .. [domainN]
         %s       list [domain1] .. [domainN]
         %s brute list [domain1] .. [domainN]
Example: %s            google.com microsoft.com
         %s brute      google.com microsoft.com
         %s       list google.com microsoft.com
         %s brute list google.com microsoft.com

Use "brute" to enable bruteforce search (slow). 
When using this mode the output will only contain active (sub)domains. 
This might yield a lot of false positives if a catch-all is defined for the target.

Use "list" to output a list of (sub)domains, records that do not resolve to an IP are prefixed with "#".
When combined with "brute" the list will only contain active (sub)domains.
`,
			0.1,
			os.Args[0],
			os.Args[0],
			os.Args[0],
			os.Args[0],
			os.Args[0],
			os.Args[0],
			os.Args[0],
			os.Args[0],
		)
		return
	}

	hosts := os.Args[1:]
	if hosts[0] == "brute" {
		bruteforceMode = true
		hosts = hosts[1:]
	}

	if hosts[0] == "list" {
		listMode = true
		hosts = hosts[1:]
	}

	initConfig()
	r := []*DNSResolver{}
	for _, ds := range Conf.DNSServers {
		r = append(r, NewDNSResolver(ds.Host, ds.Port, ds.Timeout))
	}
	resolvers = NewDNSResolverRing(r...)

	for _, host := range hosts {
		res := NewDomain(host, true, bruteforceMode)
		if !listMode {
			fmt.Println(res.ToString())
			continue
		}
		if res.Active {
			fmt.Println(res.Name)
		} else {
			fmt.Println("# " + res.Name)
		}
		for _, sd := range res.Subdomains.Active {
			fmt.Println(sd.Name)
		}
		for _, sd := range res.Subdomains.Inactive {
			fmt.Println("# " + sd.Name)
		}
	}
}
