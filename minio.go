package main

import (
	"fmt"
	"net/url"
	"strings"
)

type MinioOption struct {
	Endpoint string
	Access   string
	Secret   string
	Bucket   string
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
		sss := strings.Split(strings.TrimLeft(u.Path, "/"), "/")
		if len(sss) < 1 {
			return nil, fmt.Errorf("bucket should specified")
		} else {
			opt.Bucket = sss[0]
		}
		opt.Endpoint = u.Host
		opt.Access = u.User.Username()
		opt.Secret, _ = u.User.Password()
	}
	return &opt, nil
}
