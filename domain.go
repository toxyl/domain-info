package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
)

type Domain struct {
	Name          string
	CanonicalName string
	Addresses     []*Address
	NameServers   []string
	MailServers   []*net.MX
	Subdomains    struct {
		Active   []*Domain
		Inactive []*Domain
	}
	Simple bool
}

func (d *Domain) getCanonicalName() {
	cname, err := net.LookupCNAME(d.Name)
	if err == nil {
		d.CanonicalName = cname
	}
}

func (d *Domain) getNameServers() {
	nameservers, err := net.LookupNS(d.Name)
	if err == nil {
		for _, ns := range nameservers {
			d.NameServers = append(d.NameServers, ns.Host)
		}
	}
}

func (d *Domain) getMailServers() {
	mxrecords, err := net.LookupMX(d.Name)
	if err == nil {
		d.MailServers = append(d.MailServers, mxrecords...)
	}
}

func (d *Domain) getAddresses() {
	ips, err := net.LookupIP(d.Name)
	if err == nil {
		for _, ip := range ips {
			ipi := &Address{
				IP:    ip,
				Names: []string{},
			}
			ptrs, _ := net.LookupAddr(ip.String())
			ipi.Names = append(ipi.Names, ptrs...)
			d.Addresses = append(d.Addresses, ipi)
		}
	}
}

func (d *Domain) getSubdomains() {
	res := Do("GET", fmt.Sprintf("https://crt.sh/?q=%s", d.Name))
	res = strings.ReplaceAll(res, "<BR>", "\n")
	res = strings.ReplaceAll(res, "<TD>", "\n")
	re, _ := regexp.Compile(fmt.Sprintf("[^\\s]+\\.%s", regexp.QuoteMeta(d.Name)))
	matches := uniqueStrings(re.FindAllString(res, 9999))
	var wg sync.WaitGroup
	for _, m := range matches {
		if m != d.Name {
			wg.Add(1)
			go func(subdomain string) {
				defer wg.Done()
				sd := NewDomain(subdomain, d.Simple)
				if sd.CanonicalName != "" {
					d.Subdomains.Active = append(d.Subdomains.Active, sd)
				} else {
					d.Subdomains.Inactive = append(d.Subdomains.Inactive, sd)
				}
			}(m)
		}
	}
	wg.Wait()
}

func (d *Domain) String() string {
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

	if len(d.NameServers) > 0 {
		for _, ns := range d.NameServers {
			csv.AddRow("NS", "x", "", ns, "", "", "")
		}
	}

	if len(d.MailServers) > 0 {
		for _, mx := range d.MailServers {
			csv.AddRow("MX", "x", "", mx.Host, "", fmt.Sprint(mx.Pref), "")
		}
	}

	if len(d.Addresses) > 0 {
		for _, ipi := range d.Addresses {
			for _, name := range ipi.Names {
				csv.AddRow("Domain", "x", d.Name, d.CanonicalName, name, "", ipi.IP.String())
			}
		}
	}

	if len(d.Subdomains.Active) > 0 {
		for _, sd := range d.Subdomains.Active {
			for _, ipi := range sd.Addresses {
				for _, name := range ipi.Names {
					csv.AddRow("Subdomain", "x", sd.Name, sd.CanonicalName, name, "", ipi.IP.String())
				}
			}
		}
	}

	if len(d.Subdomains.Inactive) > 0 {
		for _, sd := range d.Subdomains.Inactive {
			csv.AddRow("Subdomain", "", sd.Name, sd.CanonicalName, "", "", "")
		}
	}

	return csv.String()
}

func NewDomain(name string, simple bool) *Domain {
	d := &Domain{
		Name:          name,
		CanonicalName: "",
		Addresses:     []*Address{},
		NameServers:   []string{},
		MailServers:   []*net.MX{},
		Subdomains: struct {
			Active   []*Domain
			Inactive []*Domain
		}{},
		Simple: simple,
	}
	d.getCanonicalName()
	if d.CanonicalName != "" {
		if !d.Simple {
			d.getNameServers()
			d.getMailServers()
			d.getAddresses()
		}
		d.getSubdomains()
	}
	return d
}
