/*
 * Copyright (C) 2023 Stefan Kühnel
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

package book

import (
	"net/url"

	"stefanco.de/bookprint/internal/util/parsetree"
	"stefanco.de/bookprint/internal/util/slices"
)

// ResolveCrossReferences replaces cross-references to other pages with their corresponding paths.
func ResolveCrossReferences(pages []*Page) error {
	for _, page := range pages {
		tree, err := parsetree.New(string(page.Content.Html))
		if err != nil {
			return err
		}

		body := parsetree.Body(tree)
		links := parsetree.ElementsByTagName(body, "a")

		for _, link := range links {
			for index, attribute := range link.Attr {
				if attribute.Key == "href" {
					crossReferencedPage := slices.FindFunc(pages, func(page *Page) bool {
						// Using cross-references in the source document involves
						// adding the complete section heading (with spaces) into
						// the "href" attribute. However, some tools, like Pandoc,
						// replace spaces and other characters with query escape
						// sequences. Thus, it is better to remove them before
						// conducting a lookup, if that cross-reference, i.e.
						// title exists in the pages slice.
						referencedTitle, err := url.QueryUnescape(attribute.Val)
						if err != nil {
							return false // ignore error
						}

						return page.Title.Text == referencedTitle
					})

					if crossReferencedPage != nil {
						// Using `attribute.Val = "..."` does not, as intended,
						// update the attribute value. Thus, it is required to
						// refer to the attribute by index to preserve the pointer.
						// See: https://stackoverflow.com/a/63870840
						link.Attr[index].Val = crossReferencedPage.Path
					}
				}
			}
		}

		html, err := parsetree.Html(parsetree.Children(body)...)
		if err != nil {
			return err
		}

		page.Content.Html = html // update content via pointer
	}

	return nil
}
