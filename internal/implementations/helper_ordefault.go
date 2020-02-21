package implementations

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
)

func tord(in *time.Time) (out time.Time) {
	if in != nil {
		out = *in
	}

	return out
}

func tord1970(in *time.Time) (out time.Time) {
	if in != nil {
		if in.After(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)) {
			out = *in
		} else {
			out = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		}
	} else {
		out = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	return out
}

func sord(in *string) (out string) {
	if in != nil {
		out = *in
	}

	return out
}

func ford(in *float32) (out float32) {
	if in != nil {
		out = *in
	}

	return out
}

func iord(in *int) (out int) {
	if in != nil {
		out = *in
	}

	return out
}

// RemoveHTMLTags strips the input of it's HTML tags
func removeHTMLTags(input string) (output string) {
	output = input
	var err error
	var stringReader = strings.NewReader(input)
	var docReader *goquery.Document
	docReader, err = goquery.NewDocumentFromReader(stringReader)
	if err == nil {
		docReader.Find("script").Each(func(i int, selection *goquery.Selection) {
			selection.Remove()
		})
		output = docReader.Text()
	}

	return output
}
