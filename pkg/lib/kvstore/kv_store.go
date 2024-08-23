package kvstore

import "github.com/exsql-io/kv-store/pkg/lib/wal"

var (
	Version string
)

type (
	KVStore struct {
		w    *wal.Wal
		data map[string]string
	}
)

func New(path string) (*KVStore, error) {
	w, err := wal.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := w.Load()
	if err != nil {
		return nil, err
	}

	return &KVStore{
		w:    w,
		data: data,
	}, nil
}

func (kvs *KVStore) Set(key string, value string) error {
	err := kvs.w.Append(wal.NewSetCommand([]byte(key), []byte(value)))
	if err != nil {
		return err
	}

	kvs.data[key] = value
	return nil
}

func (kvs *KVStore) Get(key string) (string, bool, error) {
	value, ok := kvs.data[key]
	return value, ok, nil
}

func (kvs *KVStore) Remove(key string) error {
	err := kvs.w.Append(wal.NewRmCommand([]byte(key)))
	if err != nil {
		return err
	}

	delete(kvs.data, key)
	return nil
}

func (kvs *KVStore) Close() error {
	return kvs.w.Close()
}
