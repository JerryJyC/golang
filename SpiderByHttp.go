//基本解决了各部分问题，待明天完善组装
package main

import (
	"fmt"
	"io/ioutil"

	"net/http"
	"strings"
	"log"
	"github.com/PuerkitoBio/goquery"
)




func main() {
	s := fetch("https://www.cnvd.org.cn/flaw/typeResult?typeId=33&max=20&offset=0")

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
/*
	u, _ := dom.Find(`a[href="/flaw/show/CNVD-2021-44382"]`).Attr("href")
	fmt.Println("https://www.cnvd.org.cn/" + u)
*/
	dom.Find(`tr`).Each(func(i int, selection *goquery.Selection) {
		u, _ := selection.Find(`td[width="45%"]`).Find("a").Attr("href")
		title, _ := selection.Find(`td[width="45%"]`).Find("a").Attr("title")
		d := selection.Find(`td[width="13%"]`).Last().Text()
		fmt.Println(i, u, title, d)
	})

}

func fetch(url string) string {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")
	req.Header.Set("cookie", "__jsluid_s=935befba431feb67def32ec3f24fa454; __jsl_clearance_s=1626967638.119|0|g4tlXRpQQnwUdrIUvEEPzc5HrTI%3D; JSESSIONID=3D52F64A22842B2A3BF268F399B68FFF")
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
