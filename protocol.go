package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
)

// func init() {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)
// }

type Command byte

const (
	CmdNonce Command = iota
	CmdSet
	CmdGet
	CmdDel
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int32
}

type CommandGet struct {
	Key []byte
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

func ParseCommand(r io.Reader) any {
	var cmd Command
	binary.Read(r, binary.LittleEndian, &cmd)
	switch cmd {
	case CmdSet:
		return parseSetCommand(r)
	case CmdGet:
		return parseGetCommand(r)
	default:
		panic("invalid command")
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
