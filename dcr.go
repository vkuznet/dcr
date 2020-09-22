package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

// DNS cache resolver (dcr)
//
// Copyright (c) 2020 - Valentin Kuznetsov <vkuznet AT gmail dot com>
//

// DNSCache defines DNS cache map
type DNSCache map[string][]string

// DNSManager handlers DNS cache for period of time defined by TTL
type DNSManager struct {
	Cache         DNSCache // dns cache map
	TTL           int64    // time-to-live of current cache snapshot
	RenewInterval int64    // renew interval for cache
}

// String method of DBSManager defines DNSCache string representation
func (d *DNSManager) String() string {
	data, err := json.Marshal(d.Cache)
	if err != nil {
		log.Fatalf("Unable to marshal DNS cache, error: %+v\n", err)
	}
	return string(data)
}

// Resolve function resolves given request URL into request URL with IP address of the host
func (d *DNSManager) Resolve(rurl string) string {
	if d.Cache == nil {
		d.Cache = make(DNSCache)
	}
	if d.TTL < time.Now().Unix() {
		d.update()
	}
	if vals, ok := d.Cache[rurl]; ok {
		idx := rand.Intn(len(vals))
		return vals[idx]
	}
	ips := GetIPs(rurl)
	d.Cache[rurl] = ips
	idx := rand.Intn(len(ips))
	return ips[idx]
}

// helper function to update DNSManager cache
func (d *DNSManager) update() {
	if d.Cache == nil {
		d.Cache = make(DNSCache)
	}
	rand.Seed(12345)
	for r := range d.Cache {
		d.Cache[r] = GetIPs(r)
	}
	log.Println("update DNSManager", d.String())
	d.TTL = time.Now().Unix() + d.RenewInterval
}

// NewDNSManager method properly initialize DNSManager
func NewDNSManager(renew ...int64) *DNSManager {
	dcr := DNSManager{RenewInterval: 10} // by default renew cache every 10 seconds
	if len(renew) > 0 {
		dcr.RenewInterval = renew[0]
	}
	return &dcr
}

// GetIPs helper function to resolve url hostname into IP
func GetIPs(rurl string) []string {
	var urls []string
	host := GetHostname(rurl)
	ips, err := ResolveHost(host)
	if err != nil {
		// in case of error we'll use host name itself
		log.Printf("unable to resolve host %s, error %v\n", host, err)
		urls = append(urls, rurl)
	} else {
		for _, ip := range ips {
			urls = append(urls, strings.Replace(rurl, host, ip, -1))
		}
	}
	return urls
}

// GetHostname helper function to extract hostname from given url
func GetHostname(rurl string) string {
	var path string
	if strings.Contains(rurl, "https://") {
		path = strings.Split(rurl, "https://")[1]
	} else {
		path = strings.Split(rurl, "http://")[1]
	}
	arr := strings.Split(path, "/")
	return strings.Replace(arr[0], "/", "", -1)
}

// ResolveHost helper function to resolve given Host name into set of IP addresses
func ResolveHost(host string) ([]string, error) {
	var out []string
	addrs, err := net.LookupIP(host)
	if err != nil {
		log.Printf("Unable to resolve host %s into IP addresses, error %v\n", host, err)
		return out, err
	}
	for _, addr := range addrs {
		if strings.Contains(addr.String(), ":") { // IPv6 address
			out = append(out, fmt.Sprintf("[%s]", addr))
		} else {
			out = append(out, fmt.Sprintf("%s", addr))
		}
	}
	return out, nil
}
