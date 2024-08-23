package kvstore

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestKVStore_Set(t *testing.T) {
	withKVStore(t, func(t *testing.T, kvs *KVStore) {
		f := func(key string, value string) {
			t.Helper()

			err := kvs.Set(key, value)
			require.NoError(t, err)
		}

		f("key", "value")
	})
}

func TestKVStore_Get(t *testing.T) {
	withKVStore(t, func(t *testing.T, kvs *KVStore) {
		f := func(key string, value string) {
			t.Helper()

			err := kvs.Set(key, value)
			require.NoError(t, err)

			v, ok, err := kvs.Get(key)
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, value, v)
		}

		f("key", "value")
	})
}

func TestKVStore_Delete(t *testing.T) {
	withKVStore(t, func(t *testing.T, kvs *KVStore) {
		f := func(key string, value string) {
			t.Helper()

			err := kvs.Set(key, value)
			require.NoError(t, err)

			err = kvs.Remove(key)
			require.NoError(t, err)
		}

		f("key", "value")
	})
}

func withKVStore(t *testing.T, f func(t *testing.T, kvs *KVStore)) {
	path, err := os.MkdirTemp("", "kv-store-tests-*")
	require.NoError(t, err)

	kvs, err := New(path)
	require.NoError(t, err)

	f(t, kvs)

	defer func() {
		require.NoError(t, os.RemoveAll(path))
	}()
}
