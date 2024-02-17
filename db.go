package kakapo

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/google/btree"
)

type DB struct {
	btree      *btree.BTree
	max        int
	activeFile os.File
	walFile    os.File
}

type node struct {
	Key     string
	Val     string
	Deleted bool
}

func (n node) toJson() (string, error) {
	jsonData, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (k node) Less(item btree.Item) bool {
	otherKey := item.(node)
	return k.Key < otherKey.Key
}

func Open() DB {
	btree := btree.New(2)

	file, err := os.OpenFile("data.json", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	walfile, err := os.OpenFile("wal.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	// read wal file and replay changes
	walfile.Seek(0, 0)
	reader := bufio.NewReader(walfile)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var n node
		err = json.Unmarshal([]byte(line), &n)
		if err != nil {
			panic(err)
		}

		btree.ReplaceOrInsert(n)
	}

	return DB{
		btree:      btree,
		max:        10,
		activeFile: *file,
		walFile:    *walfile,
	}

}

func (db DB) Get(k string) string {
	key1 := node{
		Key:     k,
		Deleted: false,
	}
	item := db.btree.Get(key1)

	// read from active file line by line
	if item == nil {
		db.activeFile.Seek(0, 0)
		reader := bufio.NewReader(&db.activeFile)

		var val = ""
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			var n node
			err = json.Unmarshal([]byte(line), &n)
			if err != nil {
				panic(err)
			}
			if n.Key == k && !n.Deleted {
				val = n.Val
			}
		}
		return val
	}
	return item.(node).Val

}

func (db DB) Put(k string, v string) {
	key := node{
		Key: k,
		Val: v,
	}

	db.btree.ReplaceOrInsert(key)

	// write to wal file
	jsonData, err := key.toJson()
	if err != nil {
		panic(err)
	}
	db.walFile.WriteString(jsonData + "\n")

	if db.btree.Len() >= db.max {
		node2 := db.btree.DeleteMin()
		for node2 != nil {
			// write json to active file line 1

			jsonData, err := node2.(node).toJson()
			if err != nil {
				panic(err)
			}

			// write to the begginning of the file
			db.activeFile.WriteString(jsonData + "\n")
			node2 = db.btree.DeleteMin()
		}
		db.walFile.Truncate(0)
	}
}

func (db DB) Delete(k string) {
	key := node{
		Key:     k,
		Deleted: true,
	}

	db.btree.ReplaceOrInsert(key)

	// write to wal file
	jsonData, err := key.toJson()
	if err != nil {
		panic(err)
	}
	db.walFile.WriteString(jsonData + "\n")

	if db.btree.Len() >= db.max {
		node2 := db.btree.DeleteMin()
		for node2 != nil {
			// write json to active file line 1

			jsonData, err := node2.(node).toJson()
			if err != nil {
				panic(err)
			}

			// write to the begginning of the file
			db.activeFile.WriteString(jsonData + "\n")
			node2 = db.btree.DeleteMin()
		}
		db.walFile.Truncate(0)
	}
}

func (db DB) Scan() []node {
	var pairs []node
	db.btree.Ascend(func(item btree.Item) bool {
		pairs = append(pairs, item.(node))
		return true
	})

	// read from active file line by line and append to pairs array
	db.activeFile.Seek(0, 0)
	reader := bufio.NewReader(&db.activeFile)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var n node
		err = json.Unmarshal([]byte(line), &n)
		if err != nil {
			panic(err)
		}
		pairs = append(pairs, n)
	}

	// remove duplicates
	seen := make(map[string]struct{}, len(pairs))
	j := 0
	for _, v := range pairs {
		if _, ok := seen[v.Key]; ok {
			continue
		}
		seen[v.Key] = struct{}{}
		pairs[j] = v
		j++
	}
	pairs = pairs[:j]

	return pairs
}

// close
func (db DB) Close() {
	db.activeFile.Seek(2, 0)
	node2 := db.btree.DeleteMin()
	for node2 != nil {
		// write json to active file line 1

		jsonData, err := node2.(node).toJson()
		if err != nil {
			panic(err)
		}

		// write to the begginning of the file
		db.activeFile.WriteString(jsonData + "\n")
		node2 = db.btree.DeleteMin()
	}
	db.activeFile.Close()
	db.walFile.Close()
}
