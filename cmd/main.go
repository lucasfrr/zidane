package main

import (
	"fmt"
	"lucasfrr/zidane/handlers"
	"lucasfrr/zidane/models"

	"github.com/PuerkitoBio/goquery"
)

const mainUrl = "https://alisports.x.yupoo.com/search/album?uid=1&sort=&q="

// https://royal-sports.x.yupoo.com/
// https://beonestore.x.yupoo.com/search/album?uid=1&sort=&q=real+madrid
// http://aliexpressjoe.x.yupoo.com/search/album?uid=1&sort=&q=real+madrid
// https://kingsemsport2.x.yupoo.com/search/album?uid=1&sort=&q=real+madrid
// https://minkang.x.yupoo.com/search/album?uid=1&sort=&q=real+madrid

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

	var albums []models.Album
	var jerseys []models.Jersey

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

				refererLink := "https://alisports.x.yupoo.com/" + id + "?uid=1"

				fmt.Printf("URI da imagem: %v | %v\n", refererLink, title)

				imageDocument := handlers.MakeRequest(refererLink)

				src, _ := imageDocument.Find(".viewer__img").Attr("src")
				fmt.Printf("Fonte da imagem: https:%v\n", src)

				jerseyLink := "https:" + src

				go handlers.DownloadJersey(jerseyLink, title, refererLink)

				jersey := &models.Jersey{
					Name: title,
					Url:  refererLink,
				}

				jerseys = append(jerseys, *jersey)

			})

		}

		album := &models.Album{
			Title:   albumTitles[i],
			Jerseys: jerseys,
		}

		albums = append(albums, *album)

	}
}
