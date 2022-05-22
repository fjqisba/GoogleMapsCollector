package ChromeBrowser

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

var(
	browser context.Context
	browserCancel context.CancelFunc
)

func GetPage(targetUrl string)string  {
	taskCtx, cancel := chromedp.NewContext(browser)
	defer cancel()
	var html string
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(targetUrl),
		chromedp.OuterHTML("body",&html,chromedp.ByQuery))
	if err != nil{
		log.Println(err)
	}
	return html
}

func init()  {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless",false),
	)
	browser, browserCancel = chromedp.NewExecAllocator(context.Background(), opts...)
}

