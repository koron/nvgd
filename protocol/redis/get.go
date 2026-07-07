package redis

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/koron/nvgd/internal/httperror"
	"github.com/koron/nvgd/resource"
	"github.com/redis/go-redis/v9"
)

type getHandler func(context.Context, *redis.Client, string, []string) (*resource.Resource, error)

var getHandlers = map[string]getHandler{
	"string": getString,
	"list":   getList,
	"set":    getSet,
	"zset":   getZset,
	"hash":   getHash,
	"none":   getNone,
}

func get(ctx context.Context, c *redis.Client, args []string) (*resource.Resource, error) {
	if len(args) < 1 {
		return nil, errors.New("require a key at least")
	}
	key, err := url.PathUnescape(args[0])
	if err != nil {
		return nil, fmt.Errorf("key contains invalid sequence: %s", err)
	}

	typ, err := c.Type(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to TYPE: %s", err)
	}
	h, ok := getHandlers[strings.ToLower(typ)]
	if !ok {
		return nil, fmt.Errorf("unsupported redis value type: %s", typ)
	}
	return h(ctx, c, key, args[1:])
}

func getString(ctx context.Context, c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// GET
	if len(args) == 0 {
		s, err := c.Get(ctx, k).Result()
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
		n, err := c.GetBit(ctx, k, off).Result()
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
		s, err := c.GetRange(ctx, k, start, end).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(s), nil
	}

	return nil, errors.New("too many arguments")
}

func getList(ctx context.Context, c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// LLEN
	if len(args) == 0 {
		n, err := c.LLen(ctx, k).Result()
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
		s, err := c.LIndex(ctx, k, index).Result()
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
		ss, err := c.LRange(ctx, k, start, stop).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strings.Join(ss, "\n")), nil
	}

	return nil, errors.New("too many arguments")
}

func getSet(ctx context.Context, c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// SCARD
	if len(args) == 0 {
		n, err := c.SCard(ctx, k).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatInt(n, 10)), nil
	}

	// SISMEMBER
	if len(args) == 1 {
		member := args[0]
		b, err := c.SIsMember(ctx, k, member).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatBool(b)), nil
	}

	return nil, errors.New("too many arguments")
}

func getZset(ctx context.Context, c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// ZCARD
	if len(args) == 0 {
		n, err := c.ZCard(ctx, k).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatInt(n, 10)), nil
	}

	// ZRANK
	if len(args) == 1 {
		member := args[0]
		n, err := c.ZRank(ctx, k, member).Result()
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
		ss, err := c.ZRange(ctx, k, start, stop).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strings.Join(ss, "\n")), nil
	}

	return nil, errors.New("too many arguments")
}

func getHash(ctx context.Context, c *redis.Client, k string, args []string) (*resource.Resource, error) {
	// HLEN
	if len(args) == 0 {
		n, err := c.HLen(ctx, k).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(strconv.FormatInt(n, 10)), nil
	}

	// HGET
	if len(args) == 1 {
		member := args[0]
		s, err := c.HGet(ctx, k, member).Result()
		if err != nil {
			return nil, err
		}
		return resource.NewString(s), nil
	}

	return nil, errors.New("too many arguments")
}

func getNone(ctx context.Context, c *redis.Client, k string, args []string) (*resource.Resource, error) {
	return nil, httperror.Newf(404, "not found a key: %s", k)
}
