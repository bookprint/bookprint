/*
 * Copyright (C) 2023 Stefan Kühnel
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

package book

import (
	"errors"
	"html/template"

	"golang.org/x/net/html"

	"stefanco.de/bookprint/internal/util/parsetree"
)

type Book struct {
	MetaData *MetaData
	Pages    []*Page
}

type MetaData struct {
	Title   string
	Author  string
	Date    string
	Preface template.HTML
}

func New(file []byte) (*Book, error) {
	tree, err := parsetree.New(string(file))
	if err != nil {
		return nil, err
	}

	head := parsetree.Head(tree)
	body := parsetree.Body(tree)

	if !parsetree.IsHead(head) {
		return nil, errors.New("passed HTML node is not a 'head' element")
	}

	if !parsetree.IsBody(body) {
		return nil, errors.New("passed HTML node is not a 'body' element")
	}

	title, err := getTitle(head)
	if err != nil {
		return nil, err
	}

	author, err := getAuthor(head)
	if err != nil {
		return nil, err
	}

	date, err := getDate(head)
	if err != nil {
		return nil, err
	}

	preface, err := getPreface(body)
	if err != nil {
		return nil, err
	}

	pages, err := Pages(body)
	if err != nil {
		return nil, err
	}

	err = ResolveCrossReferences(pages)
	if err != nil {
		return nil, err
	}

	book := &Book{
		MetaData: &MetaData{
			Title:   title,
			Author:  author,
			Date:    date,
			Preface: preface,
		},
		Pages: pages,
	}

	return book, nil
}

func getTitle(head *html.Node) (string, error) {
	for child := head.FirstChild; child != nil; child = child.NextSibling {
		tag := child.Data

		if tag == "title" {
			return parsetree.Text(child), nil
		}
	}

	return "", errors.New("passed HTML 'head' node contains no 'title'")
}

func getAuthor(head *html.Node) (string, error) {
	metas := parsetree.ElementsByTagName(head, "meta")

	for _, meta := range metas {
		attributes := parsetree.AttributeMap(meta)

		if name, hasNameAttribute := attributes["name"]; hasNameAttribute {
			if name == "author" {
				if content, hasContentAttribute := attributes["content"]; hasContentAttribute {
					return content, nil
				}
			}
		}
	}

	return "", nil
}

func getDate(head *html.Node) (string, error) {
	metas := parsetree.ElementsByTagName(head, "meta")

	for _, meta := range metas {
		attributes := parsetree.AttributeMap(meta)

		if name, hasNameAttribute := attributes["name"]; hasNameAttribute {
			if name == "dcterms.date" {
				if content, hasContentAttribute := attributes["content"]; hasContentAttribute {
					return content, nil
				}
			}
		}
	}

	return "", nil
}

func getPreface(body *html.Node) (template.HTML, error) {
	preface, err := parsetree.Html(parsetree.SiblingsUntilFunc(body.FirstChild, parsetree.IsHeading)...)
	if err != nil {
		return preface, err
	}

	return preface, nil
}
