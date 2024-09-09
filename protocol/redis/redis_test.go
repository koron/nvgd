package redis

import (
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/internal/protocoltest"
	"github.com/koron/nvgd/protocol"
)

func init() {
	cfg.Stores = map[string]*Store{}
}

func TestRegistered(t *testing.T) {
	protocoltest.CheckRegistered(t, "redis", protocol.ProtocolFunc(open))

	want := &cfg
	got, ok := config.Root().Protocols["redis"].(*Config)
	if !ok {
		t.Fatalf("config for redis not registered")
	}
	if got != want {
		t.Fatalf("differ &Config found: want=%p got=%p", want, got)
	}
}

func TestUnexistStore(t *testing.T) {
	got := protocoltest.OpenFail(t, "redis://_unexist_store_").Error()
	want := `unknown redis store: _unexist_store_`
	if got != want {
		t.Errorf("unexpected failure:\nwant=%s\n got=%s", want, got)
	}
}

func startRedis(t *testing.T) *miniredis.Miniredis {
	m := miniredis.RunT(t)
	addr := m.Addr()
	name := t.Name()
	cfg.Stores[name] = &Store{
		URL: "redis://" + addr + "/0",
	}
	t.Cleanup(func() {
		delete(cfg.Stores, name)
	})
	//t.Logf("start redis: addr=%s name=%s stores=%+v", addr, name, cfg.Stores)
	return m
}

func TestInvalidCommand(t *testing.T) {
	startRedis(t)
	got := protocoltest.OpenFail(t, "redis://TestInvalidCommand/_invalid_command_/foo/bar").Error()
	want := `unsupported command: _invalid_command_`
	if got != want {
		t.Fatalf("unexpected error:\nwant=%s\n got=%s", want, got)
	}
}

func TestKeys(t *testing.T) {
	redisSrv := startRedis(t)
	for k, v := range map[string]string{
		"aaa": "v1",
		"aab": "v2",
		"abc": "v3",
		"bbb": "v4",
		"bbc": "v5",
		"bcc": "v6",
		"ccc": "v7",
	} {
		err := redisSrv.Set(k, v)
		if err != nil {
			t.Fatalf("miniredis.Set failed: k=%s v=%s: %s", k, v, err)
		}
	}
	testKeys := func(req, want string) {
		t.Helper()
		got := protocoltest.OpenString(t, req)
		if got != want {
			t.Errorf("wrong result for query %q:\nwant=%q\n got=%q", req, want, got)
		}
	}
	testKeys("redis://TestKeys/keys/a", "aaa\naab\nabc")
	testKeys("redis://TestKeys/keys/*bc", "abc\nbbc")
}

func TestKeysForm(t *testing.T) {
	startRedis(t)
	got := protocoltest.OpenString(t, "redis://TestKeysForm")
	if !strings.HasPrefix(got, "<!DOCTYPE html>") {
		t.Errorf("don't start with DOCTYPE: %q", got[:min(16, len(got))])
	}
}

func TestGet(t *testing.T) {
	redisSrv := startRedis(t)

	testGet := func(key, want string) {
		t.Helper()
		got := protocoltest.OpenString(t, "redis://TestGet/get/"+key)
		if got != want {
			t.Errorf("get %q failed: want=%q got=%q", key, want, got)
		}
	}
	testGetFail := func(key, wantErr string) {
		t.Helper()
		got := protocoltest.OpenFail(t, "redis://TestGet/get/"+key).Error()
		if got != wantErr {
			t.Errorf("get %q unexpected fail: want=%q got=%q", key, wantErr, got)
		}
	}

	t.Run("string", func(t *testing.T) {
		redisSrv.Set("string0", "value0")
		// GET
		testGet("string0", "value0")
		// GETBIT
		testGet("string0/0", "0")
		testGet("string0/1", "1")
		testGet("string0/41", "0")
		testGet("string0/42", "1")
		testGet("string0/43", "1")
		testGet("string0/44", "0")
		// GETRANGE
		testGet("string0/2/4", "lue")
		// failures
		testGetFail("string0/2/4/3", "too many arguments")
	})

	t.Run("list", func(t *testing.T) {
		redisSrv.RPush("list0", "item0", "item1", "item2", "item3", "item4", "item5")
		// LLEN
		testGet("list0", "6")
		// LINDEX
		testGet("list0/0", "item0")
		testGet("list0/1", "item1")
		testGet("list0/4", "item4")
		testGet("list0/5", "item5")
		// LRANGE
		testGet("list0/2/4", "item2\nitem3\nitem4")
		// failures
		testGetFail("list0/2/4/3", "too many arguments")
	})

	t.Run("set", func(t *testing.T) {
		redisSrv.SAdd("set0", "member0", "member1", "member2", "member3")
		// SCARD
		testGet("set0", "4")
		// SISMEMBER
		testGet("set0/member0", "true")
		testGet("set0/member1", "true")
		testGet("set0/member2", "true")
		testGet("set0/member3", "true")
		testGet("set0/member", "false")
		testGet("set0/member4", "false")
		// failures
		testGetFail("set0/foo/bar", "too many arguments")
	})

	t.Run("zset", func(t *testing.T) {
		for k, v := range map[string]float64{
			"zmember0": 4,
			"zmember1": 2,
			"zmember2": 5,
			"zmember3": 1,
			"zmember4": 3,
		} {
			redisSrv.ZAdd("zset0", v, k)
		}
		// ZCARD
		testGet("zset0", "5")
		// ZRANK
		testGet("zset0/zmember0", "3")
		testGet("zset0/zmember1", "1")
		testGet("zset0/zmember2", "4")
		testGet("zset0/zmember3", "0")
		testGet("zset0/zmember4", "2")
		// ZRANGE
		testGet("zset0/1/3", "zmember1\nzmember4\nzmember0")
		// failures
		testGetFail("zset0/foo/bar/baz", "too many arguments")
	})

	t.Run("hash", func(t *testing.T) {
		redisSrv.HSet("hash0",
			"field0", "value0",
			"field1", "value1",
			"field2", "value2",
			"field3", "value3",
			"field4", "value4",
		)
		// HLEN
		testGet("hash0", "5")
		// HGET
		testGet("hash0/field0", "value0")
		testGet("hash0/field1", "value1")
		testGet("hash0/field4", "value4")
		// failures
		testGetFail("hash0/field5", "redis: nil")
		testGetFail("hash0/foo/bar", "too many arguments")
	})

	t.Run("none", func(t *testing.T) {
		testGetFail("_unexist_key_0_", "not found a key: _unexist_key_0_")
	})
}

func TestGetFailures(t *testing.T) {
	startRedis(t)
	testFail := func(subpath, wantErr string) {
		t.Helper()
		got := protocoltest.OpenFail(t, "redis://TestGetFailures/get"+subpath).Error()
		if got != wantErr {
			t.Errorf("get %q unexpected fail: want=%q got=%q", subpath, wantErr, got)
		}
	}
	testFail("", "require a key at least")
}
