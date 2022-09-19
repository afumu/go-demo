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
	"sync"
	"testing"
)

type WinVuln struct {
	TableName         gobatis.TableName `win_vuln`
	Id                string            `column:"id" json:"id"`
	DatePosted        string            `column:"date_posted" json:"datePosted"`
	Cve               string            `column:"cve" json:"cve"`
	BulletinKb        string            `column:"bulletin_kb" json:"bulletinKb"`
	Title             string            `column:"title" json:"title"`
	AffectedProduct   string            `column:"affected_product" json:"affectedProduct"`
	AffectedComponent string            `column:"affected_component" json:"affectedComponent"`
	Severity          string            `column:"severity" json:"severity"`
	Impact            string            `column:"impact" json:"impact"`
	Supersedes        string            `column:"supersedes" json:"supersedes"`
	Exploits          string            `column:"exploits" json:"exploits"`
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

var wg sync.WaitGroup

func Test_connect(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := mapper.BaseMapper[WinVuln]{SessMgr: mgr}
	queryWrapper := &mapper.QueryWrapper[WinVuln]{}
	list, err := userMapper.SelectList(queryWrapper)

	if err != nil {
		fmt.Println(err.Error())
	}

	var end = 0
	var count = 500
	// 30万
	// 300 个协程，每个 1000
	// map[AffectedProduct,map[AffectedComponent][value]]
	// 0-300
	// 300-600
	for i := 0; i < len(list); i++ {
		if i%count == 0 {
			end += count
			//fmt.Println(i, end)
			wg.Add(1)
			go putData(i, end, list)
		}
	}
	wg.Wait()
	fmt.Println("全部完成")
}

func putData(start int, end int, list []WinVuln) {
	//fmt.Printf("%d-%d:开始填充数据\n", start, end)
	newList := list[start:end]
	tx, err := globalDb.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	winVulnBucket := tx.Bucket([]byte("winVuln"))
	for _, v := range newList {
		//fmt.Printf("%d-%d:处理到了第%d条数据\n", start, end, i)
		marshal, _ := json.Marshal(v)
		product := winVulnBucket.Bucket([]byte(v.AffectedProduct))
		if product == nil {
			product, err = winVulnBucket.CreateBucket([]byte(v.AffectedProduct))
			if err != nil {
				fmt.Printf("创建桶 %s 出现异常", v.AffectedProduct)
			}
		}
		if product != nil {
			err = product.Put([]byte(v.Id), marshal)
		}

	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("%d-%d:填充数据失败", start, end)
	}

	//fmt.Printf("%d-%d:填充数据完毕\n", start, end)
	wg.Done()
}

// 查询嵌套数据库
func TestView(t *testing.T) {
	err := globalDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("winVuln"))
		// Get()函数不会返回错误，如果key存在，则返回byte slice值，如果不存在就会返回nil。
		c := b.Cursor()
		var keys []string
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			//fmt.Printf("key = %s,value = %s\n", k, v)
			keys = append(keys, fmt.Sprintf("%s", k))
		}
		for _, k := range keys {
			bucket := b.Bucket([]byte(k))
			if bucket != nil {
				c := bucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					fmt.Printf("key = %s,value = %s\n", k, v)
				}
			}
		}
		return nil
	})

	fmt.Println(err)
}

// 搜索数据
func TestSearch(t *testing.T) {
	err := globalDb.View(func(tx *bolt.Tx) error {
		// 先获取到桶
		b := tx.Bucket([]byte("winVuln"))
		// Get()函数不会返回错误，如果key存在，则返回byte slice值，如果不存在就会返回nil。

		c := b.Cursor()
		var keys []string
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			kstr := fmt.Sprintf("%s", k)
			if kstr == ".NET 5.0" {
				keys = append(keys, kstr)
			}
		}

		for _, k := range keys {
			bucket := b.Bucket([]byte(k))
			if bucket != nil {
				bucket.ForEach(func(k, v []byte) error {
					fmt.Printf("key = %s,value = %s\n", k, v)
					return nil
				})
				//c := bucket.Cursor()
				//for k, v := c.First(); k != nil; k, v = c.Next() {
				//	fmt.Printf("key = %s,value = %s\n", k, v)
				//}
			}
		}
		return nil
	})

	fmt.Println(err)
}

func TestCreate(t *testing.T) {
	create()
}

func create() {
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
