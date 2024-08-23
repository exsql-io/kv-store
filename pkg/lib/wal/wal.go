package wal

import (
	"github.com/exsql-io/kv-store/pkg/lib/util"
	"io"
	"os"
	"path/filepath"
)

type Wal struct {
	file *os.File
}

const (
	walFileName = "00000000000000000000.log"
)

func Open(path string) (*Wal, error) {
	err := os.MkdirAll(path, 0700)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	file, err := os.OpenFile(filepath.Join(path, walFileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &Wal{file: file}, nil
}

func (w *Wal) Append(command Command) error {
	payload := command.Encode()
	_, err := w.file.Write(util.Join(len(payload)+4, util.UInt32ToBytes(uint32(len(payload))), payload))
	return err
}

func (w *Wal) Load() (map[string]string, error) {
	values := make(map[string]string)
	file, err := os.Open(w.file.Name())
	if err != nil {
		return nil, err
	}

	length := make([]byte, 4)
	for {
		_, err := file.Read(length)
		if err != nil {
			if err == io.EOF {
				return values, file.Close()
			}

			return nil, err
		}

		buffer := make([]byte, util.UInt32FromBytes(length))
		_, err = file.Read(buffer)
		if err != nil {
			return nil, err
		}

		command, err := FromBytes(buffer)
		if err != nil {
			return nil, err
		}

		switch c := command.(type) {
		case *SetCommand:
			values[string(c.Key)] = string(c.Value)
		case *RmCommand:
			delete(values, string(c.Key))
		}
	}
}

func (w *Wal) Close() error {
	return w.file.Close()
}
