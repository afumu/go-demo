package bolttest

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
	"testing"
)

var globalDb *bolt.DB

func init() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	globalDb = db
}

// 创建数据库
func TestCreate(t *testing.T) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

// 创建桶
func TestCreateBucket(t *testing.T) {
	// Start a writable transaction.
	tx, err := globalDb.Begin(true)
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucket([]byte("MyBucket"))
	if err != nil {
		fmt.Println(err)
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}

}

// 更新桶
func TestUpdate(t *testing.T) {
	globalDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		err := b.Put([]byte("answer2"), []byte("43"))
		return err
	})
}

// 查看桶
func TestView(t *testing.T) {
	globalDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		// Get()函数不会返回错误，如果key存在，则返回byte slice值，如果不存在就会返回nil。
		v := b.Get([]byte("answer"))
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})
}

// 自动增长
func TestSequence(t *testing.T) {
	for i := 5; i < 10; i++ {
		if err := globalDb.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			id, _ := b.NextSequence()
			idStr := fmt.Sprintf("%d", id)
			fmt.Println("idStr", idStr)
			err := b.Put([]byte(idStr), []byte("abc"+strconv.Itoa(i)))
			return err
		}); err != nil {
			fmt.Println(err)
		}

	}
}

// 遍历所有的key和value
func TestCursor(t *testing.T) {
	globalDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key = %s,value = %s\n", k, v)
		}
		return nil

	})
}

// 搜索前缀
func TestSeek(t *testing.T) {
	globalDb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket([]byte("MyBucket")).Cursor()

		prefix := []byte("ans")
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})
}

// 范围搜索
func TestSearchRange(t *testing.T) {
	globalDb.View(func(tx *bolt.Tx) error {
		// Assume our events bucket exists and has RFC3339 encoded time keys.
		c := tx.Bucket([]byte("Events")).Cursor()

		// Our time range spans the 90's decade.
		min := []byte("1990-01-01T00:00:00Z")
		max := []byte("2022-01-01T00:00:00Z")

		// Iterate over the 90's.
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			fmt.Printf("%s: %s\n", k, v)
		}

		return nil
	})
}

// 遍历所有的key和value
func TestForEach(t *testing.T) {
	globalDb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("MyBucket"))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
}

// 桶嵌套创建，遍历
func TestBucketS(t *testing.T) {
	globalDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("keys"))
		fmt.Println(err)
		return err
	})

	if err := globalDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("keys"))
		bkt, err := b.CreateBucketIfNotExists([]byte("user"))
		if err != nil {
			return err
		}
		bkt.Put([]byte("name"), []byte("zouchangfu"))
		return err
	}); err != nil {
		fmt.Println(err)
	}

	globalDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("keys"))
		bucketUser := b.Bucket([]byte("user"))
		bucketUser.ForEach(func(k, v []byte) error {
			fmt.Printf("for each key = %s,value = %s\n", k, v)
			return nil
		})
		return nil
	})

}
