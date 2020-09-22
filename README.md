### DNS Cache Resolver (dcr)

[![Build Status](https://travis-ci.org/vkuznet/dcr.svg?branch=master)](https://travis-ci.org/vkuznet/dcr)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkuznet/dcr)](https://goreportcard.com/report/github.com/vkuznet/dcr)
[![GoDoc](https://godoc.org/github.com/vkuznet/dcr?status.svg)](https://godoc.org/github.com/vkuznet/dcr)

dcr is a simple library to resolve and cache URLs from hostbased naming
convention to IP based one. For instance, if you have an URL like
`http://www.google.com` and would like to resolve it into IP based URL
this library will help you, e.g.

```
// set DNSManager
var mgr DNSManager
host := "www.google.com"
rurl := fmt.Sprintf("http://%s/", host)

// resolve given URL
ipUrl := mgr.Resolve(rurl)

// print existing cache content
fmt.Println(mgr.String())

{"http://www.google.com/":["http://172.217.6.228/","http://[2607:f8b0:4006:802::2004]/"]}
```

### Use-cases
In service like applications where code needs to deal with lots of
similart URLs it is handy to define local cache of their IP representations
without going through DNS server. For example, in web server where you need
to handle multiple URLs at high rate and all of those URLs have a small
set of unique hostnames, it is beneficial to resolve and cache theose URLs
to avoid flood of requests to DNS server
