//尚未添加重试机制
package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"bufio"
	"os"
	"regexp"
	"time"
	"net/http"
	"strings"
	"log"
	"github.com/PuerkitoBio/goquery"
)


var (
	url = "https://www.cnvd.org.cn/flaw/typeResult?typeId=33&max=20&offset="
	url1 = "https://www.cnvd.org.cn/"
	reg = regexp.MustCompile("\\s+")
)



func main() {
	for i:=0; i < 3; i++ {
		u := url + strconv.Itoa(i*20)
		s1 := fetch(u)
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(s1))
		if err != nil {
			log.Fatal(err)
		}

		dom.Find(`tr`).Each(func(i int, selection *goquery.Selection) {
			dic := make(map[string]string)
			u1, _ := selection.Find(`td[width="45%"]`).Find("a").Attr("href")
			u1 = url1 + u1
			title, _ := selection.Find(`td[width="45%"]`).Find("a").Attr("title")
			d := selection.Find(`td[width="13%"]`).Last().Text()
			dic["文章标题"] = title
			dic["文章链接"] = u1
			dic["发布日期"] = d
			
			//抓取每篇文章内部内容
			//抓取间隔为3s
			time.Sleep(3*time.Second)
			s2 := fetch(u1)
			dom1, err := goquery.NewDocumentFromReader(strings.NewReader(s2))
			if err != nil {
				log.Fatal(err)
			}
			dom1.Find(`tr`).EachWithBreak(func(i int, selection *goquery.Selection) bool {
				if i > 13 {
					return false
				}
				t := selection.Find("td").First().Text()
				c := selection.Find("td").Last().Text()
				c = reg.ReplaceAllString(c, "")
				dic[t] = c
				return true
			})
			
			//将爬取结果写入本地文件
			err = mapWriteToFile(dic, `D:\新建文件夹\爬取结果.txt`)
			if err != nil {
				log.Fatal(err)
			}

		})
		time.Sleep(3*time.Second)
	} 
	
	
}

//抓取html文档
func fetch(url string) string {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")
	req.Header.Set("cookie", "__jsluid_s=935befba431feb67def32ec3f24fa454; __jsl_clearance_s=1627010209.942|0|1G5X7xkkfIoX0pLKsOINgLXMCJA%3D; JSESSIONID=4698431905DA02C10D3EC2A55A6D8D00")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}

	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}

	return string(body)
}

//将map写入文件
func mapWriteToFile(dic map[string]string, filename string) error {
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		f, err = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		if err != nil {
			fmt.Println("打开文件失败", err)
			return err
		}
	} else {
		f, err = os.Create(filename) //创建文件
		if err != nil {
			fmt.Println("创建文件失败", err)
			return err
		}
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	writer.WriteString("{\n")
	for k, v := range dic {
		writer.WriteString(k)
		writer.WriteString(": ")
		writer.WriteString(v)
		writer.WriteString(",\n")
	}
	writer.WriteString("}\n")
	writer.Flush()
	return nil
}

//检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}