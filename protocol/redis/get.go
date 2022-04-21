package redis

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/koron/nvgd/internal/httperror"
	"github.com/koron/nvgd/resource"
)

type getHandler func(*redis.Client, string, []string) (*resource.Resource, error)

var getHandlers = map[string]getHandler{
	"string": getString,
	"list":   getList,
	"set":    getSet,
	"zset":   getZset,
	"hash":   getHash,
	"none":   getNone,
}

func get(c *redis.Client, args []string) (*resource.Resource, error) {
	if len(args) < 1 {
		return nil, errors.New("require a key at least")
	}
	key, err := url.PathUnescape(args[0])
	if err != nil {
		return nil, fmt.Errorf("key contains invalid sequence: %s", err)
	}

	typ, err := c.Type(key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to TYPE: %s", err)
	}
	h, ok := getHandlers[strings.ToLower(typ)]
	if !ok {
		return nil, fmt.Errorf("unsupported redis value type: %s", typ)
	}
	return h(c, key, args[1:])
}

func getString(c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// GET
	if len(args) == 0 {
		s, err := c.Get(k).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(s), nil
	}

	// GETBIT
	if len(args) == 1 {
		off, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return nil, err
		}
		n, err := c.GetBit(k, off).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatInt(n, 10)), nil
	}

	// GETRANGE
	if len(args) == 2 {
		start, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return nil, err
		}
		end, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return nil, err
		}
		s, err := c.GetRange(k, start, end).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(s), nil
	}

	return nil, errors.New("too many arguments")
}

func getList(c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// LLEN
	if len(args) == 0 {
		n, err := c.LLen(k).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatInt(n, 10)), nil
	}

	// LINDEX
	if len(args) == 1 {
		index, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return nil, err
		}
		s, err := c.LIndex(k, index).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(s), nil
	}

	// LRANGE
	if len(args) == 2 {
		start, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return nil, err
		}
		stop, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return nil, err
		}
		ss, err := c.LRange(k, start, stop).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strings.Join(ss, "\n")), nil
	}

	return nil, errors.New("too many arguments")
}

func getSet(c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// SCARD
	if len(args) == 0 {
		n, err := c.SCard(k).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatInt(n, 10)), nil
	}

	// SISMEMBER
	if len(args) == 1 {
		member := args[0]
		b, err := c.SIsMember(k, member).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatBool(b)), nil
	}

	return nil, errors.New("too many arguments")
}

func getZset(c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// ZCARD
	if len(args) == 0 {
		n, err := c.ZCard(k).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatInt(n, 10)), nil
	}

	// ZRANK
	if len(args) == 1 {
		member := args[0]
		n, err := c.ZRank(k, member).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatInt(n, 10)), nil
	}

	// ZRANGE
	if len(args) == 2 {
		start, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return nil, err
		}
		stop, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return nil, err
		}
		ss, err := c.ZRange(k, start, stop).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strings.Join(ss, "\n")), nil
	}

	return nil, errors.New("too many arguments")
}

func getHash(c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// HLEN
	if len(args) == 0 {
		n, err := c.HLen(k).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatInt(n, 10)), nil
	}

	// HGET
	if len(args) == 0 {
		member := args[0]
		s, err := c.HGet(k, member).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(s), nil
	}

	return nil, errors.New("too many arguments")
}

func getNone(c *redis.Client, k string, args []string) (*resource.Resource, error) {
	return nil, httperror.Newf(404, "not found a key: %s", k)
}
