package main

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"net/url"
	"strconv"
	"strings"
)

func ParseRedisURL(s string) (*redis.Options, error) {
	var opts redis.Options
	if u, e := url.Parse(s); nil != e {
		return nil, e //log.Fatalln(e)
	} else {
		opts.Addr = u.Host
		opts.Password = u.User.Username()
		opts.DB = 1
		var p = strings.TrimRight(strings.TrimLeft(u.Path, "/"), "/")
		if p != "" {
			if n, e := strconv.Atoi(p); nil != e {
				return nil, e //log.Fatalln(e)
			} else {
				opts.DB = n
			}
		}
	}
	return &opts, nil
}
