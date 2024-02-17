package kakapo

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutGetInMemory(t *testing.T) {
	db := Open()

	db.Put("1", "2")
	val := db.Get("1")

	assert.Equal(t, "2", val)

	db.Put("1", "3")

	val = db.Get("1")

	assert.Equal(t, "3", val)
}

func TestScanInMemory(t *testing.T) {
	db := Open()

	db.Put("1", "2")
	val := db.Get("1")

	assert.Equal(t, "2", val)

	db.Put("1", "3")

	val = db.Get("1")

	assert.Equal(t, "3", val)

	// get the scan result and assert it
	pairs := db.Scan()

	t.Log("pairs", pairs)

}

func TestFlushToDisk(t *testing.T) {
	db := Open()

	// put ten key value pair
	for i := 0; i < 9; i++ {
		db.Put(strconv.Itoa(i), strconv.Itoa(i*2))
	}

	// get the scan result and assert it
	pairs := db.Scan()

	t.Log("pairs", pairs)

}

func TestSearchOnDisk(t *testing.T) {
	db := Open()

	// put ten key vialue pair
	for i := 0; i < 10; i++ {
		db.Put(strconv.Itoa(i), strconv.Itoa(i*2))
	}

	// get the scan result and assert it
	pairs := db.Scan()

	t.Log("pairs", pairs)

	val := db.Get("1")

	t.Log("val", val)

	assert.Equal(t, "2", val)

}

func TestChineseString(t *testing.T) {
	db := Open()

	db.Put("1", "你好")
	val := db.Get("1")

	assert.Equal(t, "你好", val)

	db.Put("1", "再见")

	val = db.Get("1")

	assert.Equal(t, "再见", val)

	db.Close()
}

func TestDelete(t *testing.T) {
	db := Open()

	db.Put("1", "2")
	val := db.Get("1")

	assert.Equal(t, "2", val)

	db.Delete("1")

	val = db.Get("1")

	assert.Equal(t, "", val)

	db.Close()
}

func TestCrash(t *testing.T) {
	db := Open()

	// put ten key vialue pair
	for i := 0; i < 9; i++ {
		db.Put(strconv.Itoa(i), strconv.Itoa(i*2))
	}

	// get the scan result and assert it
	pairs := db.Scan()

	t.Log("pairs", pairs)

	val := db.Get("1")

	t.Log("val", val)

	assert.Equal(t, "2", val)

	db.Close()

	db = Open()

	val = db.Get("1")

	t.Log("val", val)

	assert.Equal(t, "2", val)

	db.Put("10", "20")

	val = db.Get("10")

	t.Log("val", val)

	assert.Equal(t, "20", val)

	db.Close()
}
