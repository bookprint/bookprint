/*
 * Copyright (C) 2023 Stefan Kühnel
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

package book

import (
	"errors"
	"fmt"
	"html/template"
	"strings"

	"golang.org/x/net/html"

	"stefanco.de/bookprint/internal/util/parsetree"
)

type Chapter struct {
	Id      int
	Level   int
	Path    string
	Title   *Title
	Content *Content
}

type Title struct {
	Prefix string
	Html   template.HTML
	Text   string
}

type Content struct {
	Html template.HTML
}

func Chapters(body *html.Node) ([]*Chapter, error) {
	var chapters []*Chapter

	if !parsetree.IsBody(body) {
		return chapters, errors.New("passed HTML node is not a 'body' element")
	}

	headings := parsetree.Headings(body)

	prefix := getPrefix()

	for index, heading := range headings {
		id := index + 1
		path := fmt.Sprintf("page%d.html", id)

		level, err := parsetree.HeadingLevel(heading)
		if err != nil {
			return chapters, err
		}

		headingText := parsetree.Text(heading)
		headingText = strings.Join(strings.Fields(headingText), " ") // See: https://stackoverflow.com/a/42251527
		headingHtml, err := parsetree.Html(parsetree.Children(heading)...)
		if err != nil {
			return chapters, err
		}

		content := parsetree.SiblingsUntilFunc(heading, parsetree.IsHeading)
		contentHtml, err := parsetree.Html(content...)
		if err != nil {
			return chapters, err
		}

		chapter := &Chapter{
			Id:    id,
			Level: level,
			Path:  path,
			Title: &Title{
				Prefix: prefix(heading), // heading prefix: 1, 1.1, 1.1.1, etc.
				Html:   headingHtml,
				Text:   headingText,
			},
			Content: &Content{
				Html: contentHtml,
			},
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func getPrefix() func(*html.Node) string {
	var h1Counter int
	var h2Counter int
	var h3Counter int

	return func(heading *html.Node) string {
		tagName := parsetree.TagName(heading)

		if tagName == "h1" {
			h1Counter++
			h2Counter = 0

			return fmt.Sprintf("%d", h1Counter)
		}

		if tagName == "h2" {
			h2Counter++
			h3Counter = 0

			return fmt.Sprintf("%d.%d", h1Counter, h2Counter)
		}

		if tagName == "h3" {
			h3Counter++

			return fmt.Sprintf("%d.%d.%d", h1Counter, h2Counter, h3Counter)
		}

		return ""
	}
}

func getPath() {

}
