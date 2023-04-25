/*
 * Copyright (C) 2023 Stefan Kühnel
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

package book

import (
	"errors"

	"golang.org/x/net/html"

	"stefanco.de/bookprint/internal/util/parsetree"
	"stefanco.de/bookprint/internal/util/slices"
)

type Page struct {
	Id          int
	Level       int
	Path        string
	Title       *Title
	Content     *Content
	Next        *Chapter
	HasNext     bool
	Previous    *Chapter
	HasPrevious bool
	Parents     []*Chapter
	HasParents  bool
	Children    []*Chapter
	HasChildren bool
}

func Pages(body *html.Node) ([]*Page, error) {
	var pages []*Page

	if !parsetree.IsBody(body) {
		return pages, errors.New("passed HTML node is not a 'body' element")
	}

	chapters, err := Chapters(body)
	if err != nil {
		return pages, err
	}

	for _, chapter := range chapters {
		id := chapter.Id
		level := chapter.Level
		path := chapter.Path
		title := chapter.Title
		content := chapter.Content

		next := getNext(chapter, chapters)
		hasNext := next != nil

		previous := getPrevious(chapter, chapters)
		hasPrevious := previous != nil

		parents := getParents(chapter, chapters)
		hasParents := !slices.IsEmpty(parents)

		children := getChildren(chapter, chapters)
		hasChildren := !slices.IsEmpty(children)

		page := &Page{
			Id:          id,
			Level:       level,
			Path:        path,
			Title:       title,
			Content:     content,
			Next:        next,
			HasNext:     hasNext,
			Previous:    previous,
			HasPrevious: hasPrevious,
			Parents:     parents,
			HasParents:  hasParents,
			Children:    children,
			HasChildren: hasChildren,
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func getParents(chapter *Chapter, chapters []*Chapter) []*Chapter {
	var parents []*Chapter

	for index := 0; index < slices.Index(chapters, chapter); index++ {
		current := chapters[index]

		// The for loop iterates through the chapter slice from top to bottom,
		// with the last visited chapter being the one for which parent chapters
		// are being searched. This last visited chapter also has the highest
		// chapter level. Hence, all parent chapter levels must be lower than
		// the chapter level of the target chapter for which parent chapters are
		// being searched.
		isParent := current.Level < chapter.Level

		if isParent {
			// slice indices start at 0
			// chapter levels start at 1 (chapter level <==> HTML heading level)
			// therefore decrement the chapter level by 1 to stay in line with the index
			parents = slices.Insert(parents, current, current.Level-1)
		}
	}

	return parents
}

func getNext(chapter *Chapter, chapters []*Chapter) *Chapter {
	nextIndex := slices.Index(chapters, chapter) + 1

	for index := nextIndex; index < len(chapters); index++ {
		next := chapters[index]

		if next.Level == chapter.Level {
			return next
		}

		isNoLongerInCurrentChapter := next.Level < chapter.Level

		if isNoLongerInCurrentChapter {
			return nil
		}
	}

	return nil
}

func getPrevious(chapter *Chapter, chapters []*Chapter) *Chapter {
	previousIndex := slices.Index(chapters, chapter) - 1

	for index := previousIndex; index >= 0; index-- {
		previous := chapters[index]

		if previous.Level == chapter.Level {
			return previous
		}

		isNoLongerInCurrentChapter := previous.Level < chapter.Level

		if isNoLongerInCurrentChapter {
			return nil
		}
	}

	return nil
}

func getChildren(chapter *Chapter, chapters []*Chapter) []*Chapter {
	var children []*Chapter

	nextChapterIndex := slices.Index(chapters, chapter) + 1

	for index := nextChapterIndex; index < len(chapters); index++ {
		current := chapters[index]

		isChild := current.Level == chapter.Level+1
		isConsidered := current.Level <= 3 // only h1, h2, h3 are considered

		if isChild && isConsidered {
			children = append(children, current)
		}

		isNoLongerInCurrentChapter := current.Level <= chapter.Level

		if isNoLongerInCurrentChapter {
			return children
		}
	}

	return children
}
