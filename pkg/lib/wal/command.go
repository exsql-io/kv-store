package wal

import (
	"fmt"
	"github.com/exsql-io/kv-store/pkg/lib/util"
)

type CommandFlag byte

const (
	SET CommandFlag = iota
	RM
)

type Command interface {
	GetFlag() CommandFlag
	Encode() []byte
}

func FromBytes(bytes []byte) (Command, error) {
	flag := CommandFlag(bytes[0])
	switch flag {
	case SET:
		return SetCommandFromBytes(bytes[1:]), nil
	case RM:
		return NewRmCommand(bytes[1:]), nil
	}

	return nil, fmt.Errorf("unknown flag 0x%x", bytes[0])
}

type SetCommand struct {
	Key   []byte
	Value []byte
}

func NewSetCommand(key []byte, value []byte) Command {
	return &SetCommand{
		Key:   key,
		Value: value,
	}
}

func SetCommandFromBytes(bytes []byte) *SetCommand {
	keyLen := util.UInt32FromBytes(bytes[0:4])
	key := bytes[4 : 4+keyLen]
	return &SetCommand{
		Key:   key,
		Value: bytes[4+keyLen:],
	}
}

func (s *SetCommand) GetFlag() CommandFlag {
	return SET
}

func (s *SetCommand) Encode() []byte {
	size := len(s.Key) + len(s.Value) + 5
	return util.Join(size, []byte{byte(SET)}, util.UInt32ToBytes(uint32(len(s.Key))), s.Key, s.Value)
}

type RmCommand struct {
	Key []byte
}

func NewRmCommand(key []byte) Command {
	return &RmCommand{
		Key: key,
	}
}

func (s *RmCommand) GetFlag() CommandFlag {
	return RM
}

func (s *RmCommand) Encode() []byte {
	size := len(s.Key) + 1
	return util.Join(size, []byte{byte(RM)}, s.Key)
}
