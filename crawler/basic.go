package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gocolly/colly"
)

func main() {
	url := "https://books.studygolang.com/gopl-zh/"
	c := colly.NewCollector(
		colly.MaxDepth(2),
	)

	converter := md.NewConverter("", true, nil)

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited: ", r.Request.URL)
		if r.Request.Depth == 2 {
			sub := strings.Split(r.Request.URL.Path, "/")
			localPath := "./" + sub[1] + "/" + sub[2]
			fileName := strings.Split(sub[3], ".")[0]

			if _, err := os.Stat(localPath); os.IsNotExist(err) {
				os.MkdirAll(localPath, os.ModePerm)
			}

			f, err := os.OpenFile(localPath+"/"+fileName+".md", os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			reg := regexp.MustCompile(`((ch[0-9]*-(([0-9][0-9])|([0-9]*)))|(ch[0-9]*)).html`)
			body := reg.ReplaceAllString(string(r.Body), "${1}.md")
			markdown, err := converter.ConvertString(body)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			_, err = f.Write([]byte(markdown))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			f.Close()
		}
	})
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		match, _ := regexp.MatchString("ch[0-9]*/((ch[0-9]*-[0-9]*)|(ch[0-9]*)).html", e.Attr("href"))
		if match {
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.Visit(url)
}
