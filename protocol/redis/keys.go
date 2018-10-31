package redis

import (
	"strings"

	"github.com/go-redis/redis"
	"github.com/koron/nvgd/resource"
)

func hasKeysMeta(s string) bool {
	return strings.IndexAny(s, "?*[") >= 0
}

func keys(c *redis.Client, args []string) (*resource.Resource, error) {
	q := strings.Join(args, "/")
	if !hasKeysMeta(q) {
		q += "*"
	}
	r, err := c.Keys(q).Result()
	if err != nil {
		return nil, err
	}
	return resource.NewString(strings.Join(r, "\n")), nil
}
