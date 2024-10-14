package main

import (
	"fmt"
	"lucasfrr/zidane/handlers"

	"github.com/PuerkitoBio/goquery"
)

const mainUrl = "https://alisports.x.yupoo.com/search/album?uid=1&sort=&q="

func main() {
	fmt.Println("PESQUISA INICIAL...")
	fmt.Print("Digite o nome do clube/seleção: ")

	term := handlers.InputSearch()

	urlSearch := mainUrl + term

	fmt.Println(urlSearch)

	document := handlers.MakeRequest(urlSearch)

	pageLinks := handlers.GetPages(document, urlSearch)
	fmt.Println("NÚMERO DE PÁGINAS:", len(pageLinks))
	fmt.Println("PRÓXIMOS LINKS DA PAGINAÇÃO:", pageLinks)

	for i, link := range pageLinks {
		fmt.Printf("PAGINA %v - %v\n", i+1, link)

		pageDocument := handlers.MakeRequest(link)

		albumLinks := pageDocument.Find(".showindex__children").Find("a").Map(func(i int, s *goquery.Selection) string {
			link, _ := s.Attr("href")

			return link
		})

		for j, album := range albumLinks {
			urlAlbum := "https://alisports.x.yupoo.com" + album + "?uid=1"
			fmt.Printf("ACESSANDO ALBUM %v - %v\n", j, urlAlbum)

			albumDocument := handlers.MakeRequest(urlAlbum)

			albumDocument.Find(".showalbum__children").Each(func(i int, s *goquery.Selection) {
				id, _ := s.Attr("data-id")
				fmt.Printf("ID da imagem: %v\n", id)
			})
		}

	}

}
