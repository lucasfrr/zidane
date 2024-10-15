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

	// data := make([]map[string]string, )
	var albums []map[string]string
	var images []map[string]string

	for i, link := range pageLinks {
		fmt.Printf("PAGINA %v - %v\n", i+1, link)

		pageDocument := handlers.MakeRequest(link)

		albumLinks := pageDocument.Find(".showindex__children").Find("a").Map(func(i int, s *goquery.Selection) string {
			link, _ := s.Attr("href")

			return link
		})

		albumTitles := pageDocument.Find(".album__main").Map(func(i int, s *goquery.Selection) string {
			title, _ := s.Attr("title")

			return title
		})

		for j, album := range albumLinks {
			urlAlbum := "https://alisports.x.yupoo.com" + album
			fmt.Printf("ACESSANDO ALBUM %v - %v | %v\n", j, urlAlbum, albumTitles[j])

			albumDocument := handlers.MakeRequest(urlAlbum)

			albumDocument.Find(".showalbum__children").Each(func(i int, s *goquery.Selection) {
				id, _ := s.Attr("data-id")
				title, _ := s.Find(".text_overflow").Attr("title")

				image := map[string]string{
					"uri":   "https://alisports.x.yupoo.com/" + id + "?uid=1",
					"title": title,
				}

				images = append(images, image)

				fmt.Printf("URI da imagem: https://alisports.x.yupoo.com/%v?uid=1 | %v\n", id, title)
			})

			album := map[string]string{
				"title":  albumTitles[j],
				"uri":    urlAlbum,
				"images": images,
			}

			albums = append(albums, album)
			// data = append(data, albums)
		}

	}

}
