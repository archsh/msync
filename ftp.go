package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jlaffaye/ftp"
	"io"
	"net/url"
	"os"
)

type FTPClient struct {
	address  string
	username string
	password string
	conn     *ftp.ServerConn
}

func (c *FTPClient) login() error {
	if conn, e := ftp.Dial(c.address); nil != e {
		return e
	} else if e := conn.Login(c.username, c.password); nil != e {
		return e
	} else {
		c.conn = conn
		return nil
	}
}

func (c *FTPClient) logout() error {
	if c.conn == nil {
		return nil
	}
	if e := c.conn.Quit(); nil != e {
		return e
	} else {
		c.conn = nil
		return nil
	}
}

func (c FTPClient) FileSize(s string) (int64, error) {
	if u, e := url.Parse(s); nil != e {
		log.Errorln("FTPClient::Get> url.Parse failed: ", e, s)
		return 0, e
	} else {
		if nil != u.User {
			c.username = u.User.Username()
			if p, b := u.User.Password(); b {
				c.password = p
			}
		}
		if "" != u.Port() {
			c.address = u.Host //+ ":" + u.Port()
		} else {
			c.address = u.Host + ":21"
		}
		if e := c.login(); nil != e {
			log.Errorln("FTPClient::Get> login failed: ", e, c.address, c.username, c.password)
			return 0, e
		}
		defer func() { _ = c.logout() }()
		return c.conn.FileSize(u.Path)
	}
}

func (c FTPClient) GetFile(d string, s string) (string, error) {
	if fp, e := os.CreateTemp(d, "ftp-download-*"); nil != e {
		return "", e
	} else if _, e := c.Get(s, fp); nil != e {
		return "", e
	} else {
		_ = fp.Close()
		return fp.Name(), nil
	}
}

func (c FTPClient) Get(s string, w io.Writer) (int64, error) {
	if u, e := url.Parse(s); nil != e {
		log.Errorln("FTPClient::Get> url.Parse failed: ", e, s)
		return 0, e
	} else {
		if nil != u.User {
			c.username = u.User.Username()
			if p, b := u.User.Password(); b {
				c.password = p
			}
		}
		if "" != u.Port() {
			c.address = u.Host //+ ":" + u.Port()
		} else {
			c.address = u.Host + ":21"
		}
		if e := c.login(); nil != e {
			log.Errorln("FTPClient::Get> login failed: ", e, c.address, c.username, c.password)
			return 0, e
		}
		defer func() { _ = c.logout() }()
		if resp, e := c.conn.Retr(u.Path); nil != e {
			log.Errorln("FTPClient::Get> conn.Retr failed: ", e, u.Path)
			return 0, e
		} else {
			defer func() { _ = resp.Close() }()
			return io.Copy(w, resp)
		}
	}
}

func FTP_FileSize(s string) (int64, error) {
	return defaultFTPClient.FileSize(s)
}

func FTP_GetFile(d, s string) (string, error) {
	return defaultFTPClient.GetFile(d, s)
}

var defaultFTPClient *FTPClient

func init() {
	defaultFTPClient = &FTPClient{username: "anonymous", password: "anonymous"}
}
