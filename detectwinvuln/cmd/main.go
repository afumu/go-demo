package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/acmestack/gobatis"
	"github.com/boltdb/bolt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// go版本扫描window漏洞
// 通过解析window命令systeminfo输出的信息获取系统的版本的补丁
// 然后去数据库匹配当前的window系统是否产生漏洞

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

var buildnumberMap = make(map[int]string)

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
	Relevant          bool
}

var globalDb *bolt.DB

func init() {
	db, err := bolt.Open("D:\\workplace-go\\local\\go-demo\\winvuln\\winVuln.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	globalDb = db
}

func init() {
	buildnumberMap[10240] = "1507"
	buildnumberMap[10586] = "1511"
	buildnumberMap[14393] = "1607"
	buildnumberMap[15063] = "1703"
	buildnumberMap[16299] = "1709"
	buildnumberMap[17134] = "1803"
	buildnumberMap[17763] = "1809"
	buildnumberMap[18362] = "1903"
	buildnumberMap[18363] = "1909"
	buildnumberMap[19041] = "2004"
	buildnumberMap[19042] = "20H2"
	buildnumberMap[19043] = "21H1"
	buildnumberMap[19044] = "21H2"
	buildnumberMap[20348] = "21H2"
	buildnumberMap[22000] = "21H2"
}

func main() {
	version, build, win, arch, product, kbResults := determineProduct()
	fmt.Println(product, win, build, version, arch, kbResults)

	var kbmap = make(map[string]string)
	for _, v := range kbResults {
		kbmap[v] = ""
	}
	var newKbResults []string
	for key := range kbmap {
		newKbResults = append(newKbResults, key)
	}

	allCve := getAllCve()
	filtered, found := determineMissingPatches(product, allCve, newKbResults)
	fmt.Println(filtered, found)
}

func determineMissingPatches(product string, list []*WinVuln, kbResults []string) ([]*WinVuln, []*WinVuln) {
	var filtered []*WinVuln

	if strings.Contains(product, "Service Pack") {
		for _, cve := range list {
			if !strings.Contains(cve.AffectedProduct, product) {
				continue
			}
			cve.Relevant = true
			filtered = append(filtered, cve)
			if cve.Supersedes != "" {
				kbResults = append(kbResults, cve.Supersedes)
			}
		}
	} else {
		productSp := product + " Service Pack"
		for _, cve := range list {
			// 判断当前cve漏洞是否影响当前操作系统
			if !strings.Contains(cve.AffectedProduct, product) || strings.Contains(cve.AffectedProduct, productSp) {
				continue
			}

			// 把当前漏洞软件标记为相关
			cve.Relevant = true
			filtered = append(filtered, cve)

			// 为什么需要把这个Supersedes添加到kbResults补丁中呢？
			// Supersedes的意思代表是当前软件包会包含的补丁号，如果当前操作系统有这个软件包，说明当前操作系统也就存在这个补丁号
			// Supersedes 指的是当前包软件存在的补丁包的编号
			if cve.Supersedes != "" {
				kbResults = append(kbResults, cve.Supersedes)
			}
		}
	}

	// 合并补丁包

	join := strings.Join(kbResults, ";")
	markSuperseededHotfix(filtered, join)

	// 获取存在漏洞的软件
	var check []*WinVuln
	for _, v := range filtered {
		if v.Relevant {
			check = append(check, v)
		}
	}

	// 这里其实是做兜底的，应该不会起作用
	var supersedes = make(map[string]*WinVuln)
	for _, v := range check {
		supersedes[v.Supersedes] = v
	}
	var checked []*WinVuln
	for _, v := range check {
		if supersedes[v.BulletinKb] != nil {
			v.Relevant = false
			checked = append(checked, v)
		}
	}

	var found []*WinVuln
	for _, v := range filtered {
		if v.Relevant {
			found = append(found, v)
		}
	}
	return filtered, found
}

//

func markSuperseededHotfix(filterd []*WinVuln, superseeded string) {
	hotfixes := strings.Split(superseeded, ";")
	// 遍历所有的系统补丁
	for _, kb := range hotfixes {
		// 获取到操作系统已经存在补丁的漏洞软件包
		// 然后把他们设置为不相关
		for _, cve := range filterd {
			if cve.Relevant && cve.BulletinKb == kb {
				cve.Relevant = false
			}
		}
	}
}
func markContain(marked []string, str string) bool {
	for _, v := range marked {
		if str == v {
			return true
		}
	}
	return false
}

func getAllCve() []*WinVuln {
	var list []*WinVuln
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
					winVuln := &WinVuln{}
					err := json.Unmarshal(v, winVuln)
					if err != nil {
						fmt.Println(err)
					}
					list = append(list, winVuln)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return list
}

func determineProduct() (string, int, string, string, string, []string) {
	info := getSystemInfo()

	// 系统版本
	versionMatches := GetValueStringByRegex(info, ".*?((\\d+\\.?){3}) ((Service Pack (\\d)|N\\/\\w|.+) )?[ -\\xa5]+ (\\d+).*")
	servicePack := versionMatches[5]
	osBuild := versionMatches[6]
	fmt.Println("servicePack:", servicePack)
	fmt.Println("osbuild:", osBuild)

	var version string
	build, _ := strconv.Atoi(osBuild)
	for key, value := range buildnumberMap {
		if build == key {
			version = value
			break
		}
		if build > key {
			version = value
		} else {
			break
		}
	}
	fmt.Println("version:", version)

	// 系统名称版本
	winMatches := GetValueStringByRegex(info, ".*?Microsoft[\\(R\\)]{0,3} Windows[\\(R\\)?]{0,3} ?(Serverr? )?(\\d+\\.?\\d?( R2)?|XP|VistaT).*")
	win := winMatches[2]
	fmt.Println("win:", win)

	// 系统架构
	osArchs := GetValueStringByRegex(info, ".*?([\\w\\d]+?)-based PC.*")
	fmt.Println("osArchs:", osArchs[1])
	arch := osArchs[1]

	if !isProducts(win) {
		if arch == "X86" {
			arch = "32-bit"
		} else if arch == "x64" {
			arch = "x64-based"
		}
	}

	var product string
	if win == "XP" {
		product = "Microsoft Windows XP"
		if arch != "X86" {
			product += fmt.Sprintf(" Professional %s Edition", arch)
		}
		if servicePack != "" {
			product += fmt.Sprintf(" Service Pack %s", servicePack)
		}
	} else if win == "VistaT" {
		product = "Windows Vista"
		if arch != "X86" {
			product += fmt.Sprintf(" %s Edition", arch)
		}
		if servicePack != "" {
			product += fmt.Sprintf(" Service Pack %s", servicePack)
		}
	} else if win == "7" {
		product = fmt.Sprintf("Windows %s for %s Systems", win, arch)
		if servicePack != "" {
			product += fmt.Sprintf(" Service Pack %s", servicePack)
		}
	} else if win == "8" {
		product = fmt.Sprintf("Windows %s for %s Systems", win, arch)
	} else if win == "8.1" {
		product = fmt.Sprintf("Windows %s for %s Systems", win, arch)
	} else if win == "10" {
		product = fmt.Sprintf("Windows %s Version %s for %s Systems", win, version, arch)
	} else if win == "11" {
		product = fmt.Sprintf("Windows %s for %s Systems", win, arch)
	} else if win == "2003" {
		if arch == "X86" {
			arch = ""
		} else if arch == "x64" {
			arch = " x64 Edition"
		}
		var pversion = " "
		if version != "" {
			pversion += version
		}
		product = fmt.Sprintf("Microsoft Windows Server %s%s%s", win, arch, pversion)
	} else if win == "2008" {
		var pversion = " "
		if version != "" {
			pversion += version
		}
		product = fmt.Sprintf("Windows Server %s for %s Systems%s", win, arch, pversion)
	} else if win == "2008 R2" {
		var pversion = " "
		if version != "" {
			pversion += version
		}
		product = fmt.Sprintf("Windows Server %s for %s Systems%s", win, arch, pversion)
	} else if win == "2012" || win == "2012 R2" || win == "2016" || win == "2019" || win == "2022" {
		product = fmt.Sprintf("Windows Server %s", win)
	}

	// 补丁
	compile := regexp.MustCompile(".*KB\\d+.*")
	all := compile.FindAll([]byte(info), -1)
	var kbResults []string
	for _, v := range all {
		kbs := GetValueStringByRegex(string(v), ".*KB(\\d+).*")
		kbResults = append(kbResults, kbs[1])
	}
	fmt.Printf("kbResults%+v\n:", kbResults)
	return version, build, win, arch, product, kbResults
}

func isProducts(win string) bool {
	var products = []string{"XP", "VistaT", "2003", "2003 R2"}
	for _, v := range products {
		if win == v {
			return true
		}
	}
	return false
}

func GetValueStringByRegex(str, rule string) []string {
	reg, err := regexp.Compile(rule)
	if reg == nil || err != nil {
		return nil
	}
	//提取关键信息
	result := reg.FindStringSubmatch(str)
	if len(result) < 2 {
		return nil
	}
	return result
}

// 获取window的系统信息
func getSystemInfo() string {
	cmd := exec.Command("systeminfo")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return ""
	}

	in := bufio.NewScanner(stdout)
	builder := strings.Builder{}
	for in.Scan() {
		cmdRe := ConvertByte2String(in.Bytes(), "GB18030")
		builder.Write([]byte(cmdRe + "\n"))
	}
	err = cmd.Wait()
	if err != nil {
		log.Println(err)
		return ""
	}
	systemInfo := builder.String()

	//fmt.Println(systemInfo)
	return systemInfo
}

// 解码
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
