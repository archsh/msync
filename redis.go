package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	//log "github.com/sirupsen/logrus"
	"github.com/go-redis/redis"
)

type RedisOption struct {
	redis.Options
	Key string
}

func ParseRedisURL(s string) (*RedisOption, error) {
	var opts RedisOption
	if u, e := url.Parse(s); nil != e {
		return nil, e //log.Fatalln(e)
	} else {
		opts.Addr = u.Host
		opts.Password = u.User.Username()
		opts.DB = 1
		var sss = strings.Split(strings.TrimLeft(u.Path, "/"), "/")
		if len(sss) < 2 {
			return nil, fmt.Errorf("DB and key name should provided, eg: /1/KEY ")
		} else {
			if n, e := strconv.Atoi(sss[0]); nil != e {
				return nil, e //log.Fatalln(e)
			} else {
				opts.DB = n
			}
			opts.Key = sss[1]
		}
	}
	return &opts, nil
}
