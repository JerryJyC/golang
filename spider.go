//在尝试第一步获取动态页面html时发现chromedp无法访问目标网站，而访问其他网站正常，至此陷入僵局，许久未能解决
package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	html := get_htmlcontent("https://www.cnvd.org.cn/flaw/typelist?typeId=33", `div[class="mw Main clearfix"]`)
	log.Println(html)
}

func get_htmlcontent(url string, selector string) string {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("blink-settings", "imagesEnabled=false"), //设置不显示图片
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36`),
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//创建一个上下文，超时时间为5s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 5*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector), //等待某个特定元素出现
		chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath), //获取html文本
	)
	if err != nil {
		log.Fatal(err)
	}

	return htmlContent
}
