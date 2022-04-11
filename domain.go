package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
)

type Domain struct {
	Name       string
	CNAME      string
	PTR        []*PTRRecord
	NS         []string
	MX         []*net.MX
	Subdomains struct {
		Active   []*Domain
		Inactive []*Domain
	}
	Active bool
}

func (d *Domain) status(format string, a ...interface{}) {
	if !listMode {
		printStatusLn(format, a...)
	}
}

func (d *Domain) getCNAME() {
	d.Active = false
	cname, err := resolvers.Next().ResolveCNAME(d.Name)
	if err == nil {
		d.CNAME = cname
		d.Active = true
	}
}

func (d *Domain) getNSRecords() {
	nsrecords, err := resolvers.Next().ResolveNS(d.Name)
	if err == nil {
		for _, ns := range nsrecords {
			d.NS = append(d.NS, ns.Host)
		}
	}
}

func (d *Domain) getMXRecords() {
	mxrecords, err := resolvers.Next().ResolveMX(d.Name)
	if err == nil {
		d.MX = append(d.MX, mxrecords...)
	}
}

func (d *Domain) getPTRRecords() {
	ips, err := resolvers.Next().ResolveIP(d.Name)
	if err == nil {
		for _, ip := range ips {
			ipi := &PTRRecord{
				IP:    ip,
				Names: []string{},
			}
			ptrs, _ := resolvers.Next().ResolveAddr(ip.String())
			ipi.Names = append(ipi.Names, ptrs...)
			d.PTR = append(d.PTR, ipi)
		}
	}
}

func (d *Domain) getSubdomains(list []string) {
	l := len(list)

	for i := 0; i < l; i++ {
		var wg sync.WaitGroup
		j := 0
		for j = 0; j < Conf.ConcurrentRequests && i+j < l; j++ {
			wg.Add(1)
			go func(subdomain string) {
				defer wg.Done()
				sd := NewDomain(subdomain, false, false) // don't lookup subdomains of subdomains
				if sd.Active {
					d.Subdomains.Active = append(d.Subdomains.Active, sd)
				} else if !bruteforceMode {
					d.Subdomains.Inactive = append(d.Subdomains.Inactive, sd)
				}
			}(list[i+j] + "." + d.Name)
		}
		i += j
		wg.Wait()
	}
}

func (d *Domain) getSubdomainsFromSSLCerts() {
	res := Do("GET", fmt.Sprintf("https://crt.sh/?q=%s", d.Name))
	res = strings.ReplaceAll(res, "<BR>", "\n")
	res = strings.ReplaceAll(res, "<TD>", "\n")
	re, _ := regexp.Compile(fmt.Sprintf("([^\\s]+)\\.%s", regexp.QuoteMeta(d.Name)))
	matches := re.FindAllString(res, 9999)
	for i, m := range matches {
		matches[i] = re.ReplaceAllString(m, "$1")
	}
	matches = uniqueStrings(matches)
	d.getSubdomains(matches)
}

func (d *Domain) getSubdomainsFromList() {
	d.getSubdomains(subdomainList)
}

func (d *Domain) ToString() string {
	csv := NewCSV()
	csv.AddHeaders(CSVHeader{
		Name:   "Type",
		Format: "%-s",
	}, CSVHeader{
		Name:   "Active",
		Format: "%s",
	}, CSVHeader{
		Name:   "Domain",
		Format: "%s",
	}, CSVHeader{
		Name:   "CNAME",
		Format: "%s",
	}, CSVHeader{
		Name:   "PTR",
		Format: "%s",
	}, CSVHeader{
		Name:   "Prio.",
		Format: "%s",
	}, CSVHeader{
		Name:   "IP",
		Format: "%-s",
	})

	if len(d.NS) > 0 {
		for _, ns := range d.NS {
			csv.AddRow("NS", "x", "", ns, "", "", "")
		}
	}

	if len(d.MX) > 0 {
		for _, mx := range d.MX {
			csv.AddRow("MX", "x", "", mx.Host, "", fmt.Sprint(mx.Pref), "")
		}
	}

	if len(d.PTR) > 0 {
		for _, ipi := range d.PTR {
			for _, name := range ipi.Names {
				csv.AddRow("Domain", "x", d.Name, d.CNAME, name, "", ipi.IP.String())
			}
		}
	}

	if len(d.Subdomains.Active) > 0 {
		for _, sd := range d.Subdomains.Active {
			for _, ipi := range sd.PTR {
				for _, name := range ipi.Names {
					csv.AddRow("Subdomain", "x", sd.Name, sd.CNAME, name, "", ipi.IP.String())
				}
			}
		}
	}

	if len(d.Subdomains.Inactive) > 0 {
		for _, sd := range d.Subdomains.Inactive {
			csv.AddRow("Subdomain", "", sd.Name, sd.CNAME, "", "", "")
		}
	}

	return csv.String()
}

func NewDomain(name string, lookupCerts, lookupList bool) *Domain {
	d := &Domain{
		Name:  name,
		CNAME: "",
		PTR:   []*PTRRecord{},
		NS:    []string{},
		MX:    []*net.MX{},
		Subdomains: struct {
			Active   []*Domain
			Inactive []*Domain
		}{},
		Active: false,
	}
	d.status("[%s] Collecting info...", d.Name)
	d.getCNAME()
	if d.Active {
		d.status("[%s] Domain is active", d.Name)
		if !listMode {
			d.status("[%s] Looking up NS records...", d.Name)
			d.getNSRecords()
			d.status("[%s] Looking up MX records...", d.Name)
			d.getMXRecords()
			d.status("[%s] Looking up PTR records...", d.Name)
			d.getPTRRecords()
		}
		if lookupCerts {
			d.status("[%s] Looking up subdomains from SSL certs...", d.Name)
			d.getSubdomainsFromSSLCerts()
		}
		if lookupList {
			d.status("[%s] Looking up subdomains from list...", d.Name)
			d.getSubdomainsFromList()
		}
	} else {
		d.status("[%s] Domain is inactive", d.Name)
	}
	d.status("[%s] Finished resolving.", d.Name)
	return d
}
