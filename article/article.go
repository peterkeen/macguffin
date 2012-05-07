package article

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	ArticleId    int64
	Subject      string
	From         string
	Date         string
	MessageId    string
	Bytes        int64
	Lines        int64
	Xref         string
	Comment      string
	Filename     string
	NumParts     int64
	PartSequence int64
}

func ParseInt64(in string) int64 {
	parsed, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return parsed
}

func (article *Article) ParsedDate() (time.Time, error) {
	formats := []string{
		"02 Jan 2006 15:04:05 MST",
		"Mon, _2 Jan 2006 15:04:05 MST",
		"Mon, _2 Jan 2006 15:04:05 -0500",
	}

	for _, format := range formats {
		d, err := time.Parse(format, article.Date)
		if err != nil {
			continue
		}
		return d, nil
	}
	return time.Now(), errors.New(fmt.Sprintf("Could not parse date: %s", article.Date))
}

func ParseArticle(line string) *Article {
	parts := strings.Split(line, "\t")

	re := regexp.MustCompile("(.*)\\s+\\\"([^\\\"]+)\\\" (yEnc)? \\((\\d+)\\/(\\d+)\\)")

	subject := parts[1]
	matches := re.FindStringSubmatch(subject)
	comment := ""
	filename := ""
	num_parts := int64(0)
	part_sequence := int64(0)

	if matches != nil {
		comment = matches[1]
		filename = matches[2]
		part_sequence = ParseInt64(matches[4])
		num_parts = ParseInt64(matches[5])
	}

	return &Article{
		ParseInt64(parts[0]),
		subject,
		parts[2],
		parts[3],
		parts[4],
		ParseInt64(parts[6]),
		ParseInt64(parts[7]),
		parts[8],
		comment,
		filename,
		num_parts,
		part_sequence,
	}
}
