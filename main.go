package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"strings"
	"sync"
)

// Categoria de livros
type Categoria struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

// FichaTecnica do livro
type FichaTecnica struct {
	Dimensoes string
	Marca     string
	Peso      string
}

// Produto corresponde ao livro
type Produto struct {
	Key         string
	Categoria    string
	Autor        string
	Editora      string
	Preco        string
	Descricao    string
	FichaTecnica *FichaTecnica
	ImgURL       string
}

type categoriasMap struct {
	sync.RWMutex
	items map[string]*Categoria
}

type produtosMap struct {
	sync.RWMutex
	items map[string]*Produto
}

func main() {
	baseURL := "https://www.thesaurus.com.br/"
	categorias := categoriasMap{
		sync.RWMutex{},
		map[string]*Categoria{},
	}
	produtos := produtosMap{
		sync.RWMutex{},
		map[string]*Produto{},
	}

	c := colly.NewCollector()
	catCollector := colly.NewCollector()
	productCollector := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("nav.categories", func(e *colly.HTMLElement) {
		e.ForEach("a", func(i int, e *colly.HTMLElement) {
			title := e.Attr("title")
			link := e.Attr("href")
			key := strings.Replace(link, "/categorias/", "", -1)
			categoria := &Categoria{Link: link, Title: title}
			categorias.Lock()
			categorias.items[key] = categoria
			categorias.Unlock()
			fmt.Printf("%s %s %s\n", title, link, key)
			catCollector.Visit(e.Request.AbsoluteURL(link))
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	catCollector.OnHTML("ul.listItems", func(e *colly.HTMLElement) {
		catKey := strings.Replace(e.Request.URL.Path, "/categorias/", "", -1)
		e.ForEach("li", func(i int, e *colly.HTMLElement) {
			link := e.ChildAttr("figure a", "href")
			prodKey := strings.Replace(link, "/produto/", "", -1)
			produto := &Produto{Categoria: catKey, Key: prodKey}
			produtos.RWMutex.Lock()
			produtos.items[prodKey] = produto
			produtos.RWMutex.Unlock()
			productCollector.Visit(e.Request.AbsoluteURL(link))
		})
	})

	productCollector.OnHTML("div.detailProduct", func(e *colly.HTMLElement) {
		title := e.ChildText("h1")
		fmt.Println("Scrapping product " + title)
	})

	catCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Category", r.URL)
	})

	productCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Product", r.URL)
	})

	c.Visit(baseURL)
}
