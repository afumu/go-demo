package main

import (
	"encoding/json"
	"fmt"
	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/factory"
	"github.com/boltdb/bolt"
	"log"
	"testing"
)

type WinVuln struct {
	TableName         gobatis.TableName `win_vuln`
	Id                string            `column:"id"json:"id"`
	DatePosted        string            `column:"date_posted"json:"datePosted"`
	Cve               string            `column:"cve"json:"cve"`
	BulletinKb        string            `column:"bulletin_kb"json:"bulletinKb"`
	Title             string            `column:"title"json:"title"`
	AffectedProduct   string            `column:"affected_product"json:"affectedProduct"`
	AffectedComponent string            `column:"affected_component"json:"affectedComponent"`
	Severity          string            `column:"severity"json:"severity"`
	Impact            string            `column:"impact"json:"impact"`
	Supersedes        string            `column:"supersedes"json:"supersedes"`
	Exploits          string            `column:"exploits"json:"exploits"`
}

var globalDb *bolt.DB

func init() {
	db, err := bolt.Open("winVuln.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	globalDb = db
}

func connect() factory.Factory {
	return gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.MysqlDataSource{
			Host:     "192.168.3.31",
			Port:     3306,
			DBName:   "trivy",
			Username: "root",
			Password: "root-abcd-1234",
			Charset:  "utf8",
		}))
}

func Test_connect(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := mapper.BaseMapper[WinVuln]{SessMgr: mgr}
	queryWrapper := &mapper.QueryWrapper[WinVuln]{}
	queryWrapper.Eq("bulletin_kb", "4013198")
	list, err := userMapper.SelectList(queryWrapper)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range list {
		marshal, _ := json.Marshal(v)
		globalDb.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("winVuln"))
			product := b.Bucket([]byte(v.AffectedProduct))
			if product == nil {
				product, err = b.CreateBucket([]byte(v.AffectedProduct))
			}
			err := product.Put([]byte(v.AffectedComponent), marshal)
			return err
		})
	}

	// Start a writable transaction.

}

func Test_get(t *testing.T) {
	globalDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("winVuln"))
		// Get()函数不会返回错误，如果key存在，则返回byte slice值，如果不存在就会返回nil。
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%v\n", k, v)
			return nil
		})
		return nil
	})
}

func TestCreateBucket(t *testing.T) {
	// Start a writable transaction.
	tx, err := globalDb.Begin(true)
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucket([]byte("winVuln"))
	if err != nil {
		fmt.Println(err)
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}

}
