package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

type runtimeInfo struct {
	crypto   *cryptoInfo
	listInfo *listFile
	dataBin  *os.File
}

var info = &runtimeInfo{}

func process(c *cli.Context) error {
	dataBinPath := c.String("data")
	dataBin, err := os.Open(dataBinPath) // nolint: gosec
	if err != nil {
		return err
	}
	info.dataBin = dataBin

	listBinPath := c.String("list")
	listBytes, err := ioutil.ReadFile(listBinPath) // nolint: gosec
	if err != nil {
		return err
	}

	// Setup crypto
	info.crypto = newCryptoInfo()

	// Decrypt list.bin
	listBytesDecrypted := info.crypto.decrypt(listBytes)

	// If the user specified a hash file, use it to verify the list was correctly decrypted
	hashBinPath := c.String("hash")
	if hashBinPath != "" {
		hashBin, err := os.Open(hashBinPath) // nolint: gosec,govet
		if err != nil {
			return err
		}

		var expectedHash uint64
		err = binary.Read(hashBin, binary.LittleEndian, &expectedHash)
		if err != nil {
			return err
		}

		actualHash := info.crypto.fnv(listBytesDecrypted)
		if expectedHash != actualHash {
			return fmt.Errorf("hash mismatch, expected %v got %v", expectedHash, actualHash)
		}
	}

	// Read list data
	listInfo, err := readListFile(listBytesDecrypted)
	if err != nil {
		return err
	}
	info.listInfo = listInfo

	err = info.dumpItems()
	if err != nil {
		return err
	}

	return nil
}

func (info *runtimeInfo) dumpItems() error {
	var err error
	for _, item := range info.listInfo.items {
		group := info.listInfo.groups[item.grp]
		url := filepath.Join(".", "out", group.name, item.name+".unity3d")
		fmt.Printf("Writing %s\n", url)

		folder := filepath.Dir(url)
		_ = os.MkdirAll(folder, os.ModePerm)

		// Go to item position
		_, err = info.dataBin.Seek(int64(item.pos), 0)
		if err != nil {
			return err
		}

		// Create buffer and read item contents
		contents := make([]byte, item.size)
		_, err = info.dataBin.Read(contents)
		if err != nil {
			return err
		}

		// Decrypt buffer
		info.crypto.key = item.hash
		dec := info.crypto.decrypt(contents)

		f, err := os.Create(url)
		if err != nil {
			return err
		}

		_, err = f.Write(dec)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
