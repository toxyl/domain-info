package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(
			`domain-info v%.2f
Usage:   %s      [domain1] .. [domainN]
         %s list [domain1] .. [domainN]
Example: %s      google.com microsoft.com
         %s list google.com microsoft.com
`,
			0.1,
			os.Args[0],
			os.Args[0],
			os.Args[0],
			os.Args[0],
		)
		return
	}

	hosts := os.Args[1:]
	list := false
	if hosts[0] == "list" {
		list = true
		hosts = hosts[1:]
	}

	for _, host := range hosts {
		if !list {
			printStatusLn("[%s] Collecting info...", host)
			res := NewDomain(host, false)
			fmt.Println(res)
			continue
		}
		res := NewDomain(host, true)
		if res.CanonicalName != "" {
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
