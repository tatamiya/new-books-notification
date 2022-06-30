package main

import (
	"github.com/mmcdole/gofeed"
)

func main() {
	url := "https://www.hanmoto.com/ci/bd/search/hdt/%E6%96%B0%E3%81%97%E3%81%8F%E7%99%BB%E9%8C%B2%E3%81%95%E3%82%8C%E3%81%9F%E6%9C%AC/sdate/today/created/today/order/desc/vw/rss20"
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
}
