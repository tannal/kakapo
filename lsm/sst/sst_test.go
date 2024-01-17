package sst

import (
	"os"
	"reflect"
	"testing"
)

// TestSerializationToFile tests the serialization of SstEntry, writing to a file,
// and then reading from the file for deserialization.
func TestSerializationToFile(t *testing.T) {
	// Create a test entry
	testEntry := SstEntry{
		Key:     "ABCD",
		Value:   []byte("abcd"),
		Deleted: true,
	}

	// Serialize the entry
	serializedData, err := Serialize(testEntry)
	if err != nil {
		t.Fatalf("Failed to serialize: %s", err)
	}

	// Create a temporary file
	tmpfile, err := os.Create("./test.bin")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	// defer os.Remove(tmpfile.Name()) // Clean up

	// Write serialized data to temp file
	if _, err := tmpfile.Write(serializedData); err != nil {
		tmpfile.Close()
		t.Fatalf("Failed to write to temporary file: %s", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %s", err)
	}

	// Read the data back from the temp file
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read from temporary file: %s", err)
	}

	// Deserialize the data
	deserializedEntry, err := Deserialize(data)
	if err != nil {
		t.Fatalf("Failed to deserialize: %s", err)
	}

	// Compare the original and deserialized entries
	if !reflect.DeepEqual(testEntry, deserializedEntry) {
		t.Errorf("Original and deserialized entries are not equal. Original: %+v, Deserialized: %+v", testEntry, deserializedEntry)
	}
}

func TestReadFromFile(t *testing.T) {
	var bytes []byte
	bytes, err := os.ReadFile("./test.bin")
	if err != nil {
	}

	entry, err := Deserialize(bytes)
	t.Log(entry.Key)
	t.Log(string(entry.Value))
}
