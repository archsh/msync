package main

import (
	"net/url"
	"strings"
)

type MinioOption struct {
	Endpoint string
	Access   string
	Secret   string
	Secure   bool
}

func ParseMinioURL(s string) (*MinioOption, error) {
	var opt MinioOption
	if u, e := url.Parse(s); nil != e {
		return nil, e //log.Fatalln(e)
	} else {
		if strings.ToLower(u.Scheme) == "https" {
			opt.Secure = true
		}
		opt.Endpoint = u.Host
		opt.Access = u.User.Username()
		opt.Secret, _ = u.User.Password()
	}
	return &opt, nil
}
