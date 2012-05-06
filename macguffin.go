package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"macguffin/article"
	"macguffin/client"
	"os"
)

func main() {

	addr := ""
	user := ""
	pass := ""

	flag.StringVar(&addr, "addr", "", "Address of usenet server. Example: news.example.com:119")
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&pass, "pass", "", "Password")

	flag.Parse()

	if pass == "" || user == "" || addr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err := mgclient.NewUsenetClient(addr)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("Connected")
	client.Authenticate(user, pass)

	log.Println("Authenticated")
	log.Println("Finding start")
	start, high, err := client.FindStart("alt.binaries.multimedia", 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Start: %d, Num: %d", start, high-start)

	log.Println("Getting overview")
	overview, err := client.OverviewStartingAt("alt.binaries.multimedia", start)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf_overview := bufio.NewReader(overview)

	for counter := 0; ; counter++ {
		line, err := buf_overview.ReadString('\n')
		if err == io.EOF {
			break
		}

		article := mgarticle.ParseArticle(line)
		fmt.Printf("%s: %s %s %d %d %s\n", article.MessageId, article.Subject, article.Filename, article.NumParts, article.PartSequence, article.Date)

		if counter%1000 == 0 {
			log.Println(counter)
		}
	}

	log.Println("Done")
}
