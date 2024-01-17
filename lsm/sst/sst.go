package sst

import (
	"bytes"
	"encoding/gob"
	"os"
)

// Serialize encodes the SstEntry struct into bytes
func Serialize(entry SstEntry) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(entry)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Deserialize decodes the bytes into an SstEntry struct
func Deserialize(data []byte) (SstEntry, error) {
	var entry SstEntry
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&entry)
	if err != nil {
		return SstEntry{}, err
	}
	return entry, nil
}

func ReadFromFile(file os.File) (SstEntry, error) {
	var data []byte
	file.Read(data)
	return Deserialize(data)
}
