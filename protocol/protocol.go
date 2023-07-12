package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

// func init() {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)
// }

type Command byte

const (
	CmdNone Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdJoin
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int32
}

type CommandGet struct {
	Key []byte
}

type CommandJoin struct {
}

type Status byte

const (
	StatusNone Status = iota
	StatusOK
	StatusERR
	StatusKeyNotFound
)

func (s Status) String() string {
	switch s {
	case StatusERR:
		return "ERR"
	case StatusOK:
		return "OK"
	case StatusKeyNotFound:
		return "KEYNOTFOUND"
	default:
		return "NONE"
	}
}

type ResponseGet struct {
	Status Status
	Value  []byte
}

type ResponseSet struct {
	Status Status
}

func (r ResponseGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)

	valueLen := int32(len(r.Value))
	binary.Write(buf, binary.LittleEndian, valueLen)
	binary.Write(buf, binary.LittleEndian, r.Value)

	return buf.Bytes()
}

func (r ResponseSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)
	return buf.Bytes()
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdSet)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)

	binary.Write(buf, binary.LittleEndian, int32(c.TTL))
	return buf.Bytes()
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGet)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)
	return buf.Bytes()
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	err := binary.Read(r, binary.LittleEndian, &cmd)
	if err != nil {
		return nil, err
	}
	switch cmd {
	case CmdSet:
		return parseSetCommand(r), nil
	case CmdGet:
		return parseGetCommand(r), nil
	case CmdJoin:
		return &CommandJoin{}, nil
	default:
		return nil, fmt.Errorf("invalid command")
	}

}

func parseSetCommand(r io.Reader) *CommandSet {
	// TODO : find a way of better error handling
	cmd := &CommandSet{}
	var err error

	var keyLen int32
	err = binary.Read(r, binary.LittleEndian, &keyLen)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Key = make([]byte, keyLen)
	err = binary.Read(r, binary.LittleEndian, &cmd.Key)
	if err != nil {
		log.Fatal(err)
	}

	var valueLen int32
	err = binary.Read(r, binary.LittleEndian, &valueLen)
	if err != nil {
		log.Fatal(err)
	}

	cmd.Value = make([]byte, valueLen)
	err = binary.Read(r, binary.LittleEndian, &cmd.Value)
	if err != nil {
		log.Fatal(err)
	}

	err = binary.Read(r, binary.LittleEndian, &cmd.TTL)
	if err != nil {
		log.Fatal(err)
	}
	return cmd
}

func parseGetCommand(r io.Reader) *CommandGet {
	// TODO : find a way of better error handling
	cmd := &CommandGet{}
	var err error

	var keyLen int32
	err = binary.Read(r, binary.LittleEndian, &keyLen)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Key = make([]byte, keyLen)
	err = binary.Read(r, binary.LittleEndian, &cmd.Key)
	if err != nil {
		log.Fatal(err)
	}
	return cmd
}

func ParseSetResponse(r io.Reader) (*ResponseSet, error) {
	resp := &ResponseSet{}
	err := binary.Read(r, binary.LittleEndian, &resp.Status)
	return resp, err
}

func ParseGetResponse(r io.Reader) (*ResponseGet, error) {
	resp := &ResponseGet{}
	binary.Read(r, binary.LittleEndian, &resp.Status)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)

	resp.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &resp.Value)

	return resp, nil
}
