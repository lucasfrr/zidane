package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const urlSearch = "https://alisports.x.yupoo.com/search/album?uid=1&sort=&q=real+madrid"
const mainUrl = "https://alisports.x.yupoo.com"

func main() {
	fmt.Println("PESQUISA INICIAL...")
	res, err := http.Get(urlSearch)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		panic(err)
	}

	pageLinks := getPages(doc)

	fmt.Println("NÚMERO DE PÁGINAS:", len(pageLinks))
	fmt.Println("PRÓXIMOS LINKS DA PAGINAÇÃO:", pageLinks)

	links := doc.Find(".showindex__children").Find("a").Map(func(i int, s *goquery.Selection) string {
		link, _ := s.Attr("href")

		return link
	})

	titles := doc.Find(".showindex__children").Find("a").Map(func(i int, s *goquery.Selection) string {
		title, _ := s.Attr("title")

		return title
	})

	fmt.Println(titles)
	fmt.Println("LINKS DE TODOS OS MODELOS ENCONTRADOS:", links)

	for i, link := range links {
		fmt.Printf("ACESSANDO MODELO NÚMERO %v: %v\n", i+1, link)

		fmt.Println()
		jerseyRes, err := http.Get(mainUrl + link)

		if err != nil {
			fmt.Println(err)
			continue
		}

		defer jerseyRes.Body.Close()

		jerseyDoc, err := goquery.NewDocumentFromReader(jerseyRes.Body)

		if err != nil {
			fmt.Println(err)
			continue
		}

		jerseyDoc.Find(".showalbum__children").Each(func(i int, s *goquery.Selection) {
			id, _ := s.Attr("data-id")
			fmt.Println("EXIBINDO OS LINKS DAS IMAGENS DO MODELO DA CAMISA")
			fmt.Println(mainUrl + id + "?uid=1")
		})

	}

}

func getPages(document *goquery.Document) []string {
	pages := 0
	var pageLinks []string

	document.Find(".pagination__number").Each(func(i int, s *goquery.Selection) {
		pages++
	})

	for i := 1; i <= pages; i++ {
		newLink := urlSearch + "&page=" + strconv.Itoa(i)
		pageLinks = append(pageLinks, newLink)
	}

	return pageLinks
}
