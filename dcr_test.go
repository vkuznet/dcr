package dcr

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

// TestGetIPs function
func TestGetIPs(t *testing.T) {
	host := "www.yahoo.com"
	rurl := fmt.Sprintf("http://%s:80/", host)
	rhost := GetHostname(rurl)
	if rhost != host {
		t.Fatalf("unable to get hostname, rhost=%s host=%s\n", rhost, host)
	}
	ips := GetIPs(rurl)
	addrs, err := net.LookupIP(host)
	if err != nil {
		t.Fatalf("unable to resolve %s, error %v\n", host, err)
	}
	if len(ips) != len(addrs) {
		t.Fatalf("unable to resolve %s, expected %v got %v\n", rurl, addrs, ips)
	}
}

// TestResolveHost function
func TestResolveHost(t *testing.T) {
	host := "www.yahoo.com"
	ips, err := ResolveHost(host)
	if err != nil {
		t.Fatalf("ResolveHost error %v\n", err)
	}
	for _, ip := range ips {
		if ip == host {
			t.Fatalf("host and IP address are the same, %s\n", host)
		}
	}
}

// TestDNSCache resolver function
func TestDNSCache(t *testing.T) {
	mgr := NewDNSManager(1)
	host := "www.yahoo.com"
	rurl := fmt.Sprintf("http://%s/", host)
	ipUrl := mgr.Resolve(rurl)
	addrs, err := net.LookupIP(host)
	if err != nil {
		t.Fatalf("unable to resolve %s, error %v\n", host, err)
	}
	match := ""
	for _, addr := range addrs {
		if strings.Contains(ipUrl, fmt.Sprintf("%s", addr)) {
			match = fmt.Sprintf("%s", addr)
			break
		}
	}
	if match == "" {
		t.Fatalf("unable to match any IP from %v to url %s ipUrl %s\n", addrs, rurl, ipUrl)
	}
	fmt.Println("DNSCache", mgr.String())
	// this time we should see update call
	ipUrl = mgr.Resolve(rurl)
	fmt.Println("no update call", rurl, ipUrl)
	// let's clear up TTL
	mgr.TTL = 0
	// this time we'll see the update call
	ipUrl = mgr.Resolve(rurl)
	fmt.Println("update call", rurl, ipUrl)
}
