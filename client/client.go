package client

import (
	"fmt"
	"io"
	"log"
	"github.com/peterkeen/macguffin/article"
	"net/textproto"
	"crypto/tls"
	"strings"
	"time"
	"errors"
)

type UsenetClient struct {
	conn *textproto.Conn
}

func (client *UsenetClient) Command(command string, expected int) (int, string, error) {
	err := client.conn.PrintfLine(command)
	if err != nil {
		return 0, "", err
	}

	return client.conn.ReadCodeLine(expected)
}

func (client *UsenetClient) MustCommand(command string, expected int) (results string) {
	_, results, err := client.Command(command, expected)
	if err != nil {
		log.Panic(err)
	}
	return results
}

func (client *UsenetClient) Authenticate(username string, password string) {
	client.MustCommand(fmt.Sprintf("authinfo user %s", username), 381)
	client.MustCommand(fmt.Sprintf("authinfo pass %s", password), 281)
}

func (client *UsenetClient) Group(group string) (int64, int64, int64, error) {
	res := client.MustCommand(fmt.Sprintf("group %s", group), 211)
	parts := strings.Split(res, " ")

	total := article.ParseInt64(parts[0])
	low := article.ParseInt64(parts[1])
	high := article.ParseInt64(parts[2])

	return total, low, high, nil
}

func (client *UsenetClient) FindStart(group string, retention int) (int64, int64, error) {
	_, low, originalHigh, _ := client.Group(group)
	target_date := time.Now().AddDate(0, 0, 0-retention)

	high := originalHigh

	for high > low {
		mid := (low + high) / 2
		art_text, err := client.overviewForArticleId(mid)
		if err != nil {
			return 0, 0, err
		}
		art := article.ParseArticle(art_text)

		parsedDate, err := art.ParsedDate()
		if err != nil {
			log.Fatal(err)
		}
		if parsedDate.Before(target_date) {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return high, originalHigh, nil
}

func (client *UsenetClient) overviewForArticleId(article int64) (string, error) {
	client.MustCommand(fmt.Sprintf("xover %d", article), 224)
	lines, err := client.conn.ReadDotLines()
	if err != nil {
		return "", err
	}
	return lines[0], nil
}

func (client *UsenetClient) OverviewStartingAt(group string, article int64) (io.Reader, error) {

	client.Group(group)

	client.MustCommand(fmt.Sprintf("xover %d-", article), 224)

	return client.conn.DotReader(), nil
}

func buildConnection(addr string, useTls bool) (*textproto.Conn, error) {
	if useTls == true {
		tls, err := tls.Dial("tcp", addr, nil)
		if err != nil {
			return nil, err
		}
		return textproto.NewConn(tls), nil
	} else {
		return textproto.Dial("tcp", addr)
	}
	return nil, errors.New("could not build connection")
}

func NewUsenetClient(addr string, useTls bool) (client *UsenetClient, err error) {
	conn, err := buildConnection(addr, useTls)
	if err != nil {
		return nil, err
	}
	_, _, err = conn.ReadCodeLine(200)
	if err != nil {
		return nil, err
	}

	return &UsenetClient{conn}, nil
}
