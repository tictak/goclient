package goclient

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	c        *http.Client
	dialTO   time.Duration
	rwTO     time.Duration
	Header   map[string]string
	redirect bool
	host     string
}

func NewTimeoutClient(timeout time.Duration) *Client {
	client := new(Client)
	client.c = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				return NewTimeoutConn(conn, timeout), nil
			},
		},
	}
	return client
}

func (c *Client) SetTimeout(timeout time.Duration) {
	if c == nil {
		return
	}
	if c.c == nil {
		c.c = &http.Client{}
	}
	c.c.Transport = &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, timeout)
			if err != nil {
				return nil, err
			}
			return NewTimeoutConn(conn, timeout), nil
		},
	}
}

func (c *Client) SetRedirect() {
	c.redirect = true
}

func (c *Client) SetHost(h string) {
	c.host = h
}

func (c *Client) SetHeader(key, value string) {
	if strings.ToLower(key) == "host" {
		c.host = value
		return
	}
	if c.Header == nil {
		c.Header = make(map[string]string)
	}
	c.Header[key] = value
}

func (c *Client) GetData(u string) (int, []byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Host = c.host
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	var resp *http.Response
	if c.redirect {
		resp, err = c.c.Do(req)
	} else {
		resp, err = c.c.Transport.RoundTrip(req)
	}
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, b, err
}

type TimeoutConn struct {
	net.Conn
	timeout time.Duration
}

func NewTimeoutConn(conn net.Conn, timeout time.Duration) *TimeoutConn {
	return &TimeoutConn{
		Conn:    conn,
		timeout: timeout,
	}
}

func (c *TimeoutConn) Read(b []byte) (n int, err error) {
	c.SetReadDeadline(time.Now().Add(c.timeout))
	return c.Conn.Read(b)
}

func (c *TimeoutConn) Write(b []byte) (n int, err error) {
	c.SetWriteDeadline(time.Now().Add(c.timeout))
	return c.Conn.Write(b)
}
