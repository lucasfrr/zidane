package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const mainUrl = "https://alisports.x.yupoo.com/search/album?uid=1&sort=&q="

func main() {
	fmt.Println("PESQUISA INICIAL...")
	fmt.Print("Digite o nome do clube/seleção: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	term := scanner.Text()

	formattedTerm := formatSearchTerm(term)

	urlSearch := mainUrl + formattedTerm

	fmt.Println(urlSearch)
	res, err := http.Get(urlSearch)

	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		panic(err)
	}

	res.Body.Close()

	pageLinks := getPages(doc, urlSearch)

	fmt.Println("NÚMERO DE PÁGINAS:", len(pageLinks))
	fmt.Println("PRÓXIMOS LINKS DA PAGINAÇÃO:", pageLinks)

	for page, link := range pageLinks {
		fmt.Printf("ACESSANDO PÁGINA %v - link: %v", page+1, link)

		response, err := http.Get(link)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()

		modelLinks := doc.Find(".showindex__children").Find("a").Map(func(i int, s *goquery.Selection) string {
			link, _ := s.Attr("href")

			return link
		})

		// modelTitles := doc.Find(".showindex__children").Find("a").Map(func(i int, s *goquery.Selection) string {
		// 	title, _ := s.Attr("title")

		// 	return title
		// })

		for i, link := range modelLinks {
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

}

func getPages(document *goquery.Document, url string) []string {
	pages := 0
	var pageLinks []string

	document.Find(".pagination__number").Each(func(i int, s *goquery.Selection) {
		pages++
	})

	if pages > 0 {
		for i := 1; i <= pages; i++ {
			newLink := url + "&page=" + strconv.Itoa(i)
			pageLinks = append(pageLinks, newLink)
		}
	} else {
		pageLinks = append(pageLinks, url)
	}

	return pageLinks
}

func formatSearchTerm(term string) string {
	var newTerm string

	s := strings.ToLower(term)

	for i := 0; i < len(s); i++ {
		if string(s[i]) == " " {
			newTerm += "+"
		} else {
			newTerm += string(s[i])
		}
	}

	return newTerm
}
