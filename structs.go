package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type groupData struct {
	size uint32
	name string
}

type item struct {
	name string
	grp  uint16
	id   uint32
	pos  uint64
	size uint32
	hash uint64
}

type listFile struct {
	groups []groupData
	items  []item
}

// nolint: funlen
func readListFile(data []byte) (*listFile, error) {
	var groupCount uint32
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, &groupCount)
	if err != nil {
		return nil, err
	}
	fmt.Printf("groups: %d\n", groupCount)

	listFileInfo := &listFile{}
	listFileInfo.groups = make([]groupData, groupCount)

	// Read groups
	for i := uint32(0); i < groupCount; i++ {
		var size uint32
		err := binary.Read(r, binary.LittleEndian, &size) // nolint: govet
		if err != nil {
			return nil, err
		}

		var name string
		err = readString(r, binary.LittleEndian, &name)
		if err != nil {
			return nil, err
		}

		listFileInfo.groups[i] = groupData{
			size: size,
			name: name,
		}
	}

	var itemCount uint32
	err = binary.Read(r, binary.LittleEndian, &itemCount)
	if err != nil {
		return nil, err
	}
	fmt.Printf("items: %d\n", itemCount)

	listFileInfo.items = make([]item, itemCount)

	for i := uint32(0); i < itemCount; i++ {
		var name string
		err = readString(r, binary.LittleEndian, &name)
		if err != nil {
			return nil, err
		}

		var grp uint16
		err = binary.Read(r, binary.LittleEndian, &grp)
		if err != nil {
			return nil, err
		}

		var id uint32
		err = binary.Read(r, binary.LittleEndian, &id)
		if err != nil {
			return nil, err
		}

		var pos uint64
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return nil, err
		}

		var size uint32
		err = binary.Read(r, binary.LittleEndian, &size)
		if err != nil {
			return nil, err
		}

		var hash uint64
		err = binary.Read(r, binary.LittleEndian, &hash)
		if err != nil {
			return nil, err
		}

		listFileInfo.items[i] = item{
			name: name,
			grp:  grp,
			id:   id,
			pos:  pos,
			size: size,
			hash: hash,
		}
	}

	return listFileInfo, nil
}

// Comeon golang.. why do i need to define this
func readString(f io.Reader, order binary.ByteOrder, str *string) error {

	for {
		var c byte
		err := binary.Read(f, order, &c)
		if err != nil {
			return err
		}

		if c == 0x0 {
			break
		}

		*str += string(c)
	}

	return nil
}
