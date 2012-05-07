package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"github.com/peterkeen/macguffin/article"
	"github.com/peterkeen/macguffin/client"
	"os"
)

func main() {

	addr := ""
	user := ""
	pass := ""
	group := ""
	retention := 0

	flag.StringVar(&addr, "addr", "", "Address of usenet server. Example: news.example.com:119")
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&pass, "pass", "", "Password")
	flag.StringVar(&group, "group", "", "Newsgroup to get headers for")
	flag.IntVar(&retention, "retention", 0, "Number of days to download")

	flag.Parse()

	if pass == "" || user == "" || addr == "" || group == "" || retention == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cl, err := client.NewUsenetClient(addr)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("Connected")
	cl.Authenticate(user, pass)

	log.Println("Authenticated")
	log.Println("Finding start")
	start, high, err := cl.FindStart(group, retention)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Start: %d, Num: %d", start, high-start)

	log.Println("Getting overview")
	overview, err := cl.OverviewStartingAt(group, start)
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

		art := article.ParseArticle(line)
		date, err := art.ParsedDate()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s %s %d %d %s\n", art.MessageId, art.Subject, art.Filename, art.NumParts, art.PartSequence, date)

		if counter%1000 == 0 {
			log.Println(counter)
		}
	}

	log.Println("Done")
}
