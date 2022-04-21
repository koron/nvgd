// Package redis provides redis protocol for NVGD.
package redis

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/go-redis/redis/v7"
	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

func open(u *url.URL) (*resource.Resource, error) {
	c, err := getClient(u.Hostname())
	if err != nil {
		return nil, err
	}
	cmd, args, err := parseCommand(u)
	if err != nil {
		return nil, err
	}
	h, ok := handlers[strings.ToUpper(cmd)]
	if !ok {
		return nil, fmt.Errorf("unsupported command: %s", cmd)
	}
	return h(c, args)
}

type handler func(*redis.Client, []string) (*resource.Resource, error)

var handlers = map[string]handler{
	"":     keysForm,
	"GET":  get,
	"KEYS": keys,
}

func parseCommand(u *url.URL) (cmd string, args []string, err error) {
	raw := u.Path
	if len(raw) > 0 && raw[0] == '/' {
		raw = raw[1:]
	}
	args = strings.SplitN(raw, "/", 10)
	if len(args) == 0 {
		return "", nil, nil
	}
	return args[0], args[1:], nil
}

var (
	clients = map[string]*redis.Client{}
	mu      sync.Mutex
)

func getClient(name string) (*redis.Client, error) {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := clients[name]; ok {
		return c, nil
	}
	s, ok := cfg.Stores[name]
	if !ok {
		return nil, fmt.Errorf("unknown redis store: %s", name)
	}
	o, err := redis.ParseURL(s.URL)
	if err != nil {
		return nil, err
	}
	c := redis.NewClient(o)
	clients[name] = c
	return c, nil
}

// Config provides configuration for redis protocol.
type Config struct {
	Stores map[string]*Store `yaml:"stores"`
}

// Store represents a data store of redis.
type Store struct {
	// URL store redis URL.
	URL string `yaml:"url"`
}

var cfg Config

func init() {
	protocol.Register("redis", protocol.ProtocolFunc(open))
	config.RegisterProtocol("redis", &cfg)
}
