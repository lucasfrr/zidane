package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const realMadrid = "https://alisports.x.yupoo.com/search/album?uid=1&sort=&q=real+madrid"

func main() {
	res, err := http.Get(realMadrid)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		panic(err)
	}

	pages := 0
	doc.Find(".pagination__number").Each(func(i int, s *goquery.Selection) {
		pages++
	})

	fmt.Println("Pages:", pages)

	links := doc.Find(".showindex__children").Find("a").Map(func(i int, s *goquery.Selection) string {
		link, _ := s.Attr("href")

		return link
	})

	fmt.Println(links)

}
