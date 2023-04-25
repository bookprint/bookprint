/*
 * Copyright (C) 2023 Stefan Kühnel
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

package parsetree

import (
	"bytes"
	"errors"
	"html/template"
	"strings"

	"golang.org/x/net/html"

	"stefanco.de/bookprint/internal/util/slices"
)

// New returns the HTML parse tree from the given HTML string.
func New(string string) (*html.Node, error) {
	reader := strings.NewReader(string)

	parsetree, err := html.Parse(reader)
	if err != nil {
		return parsetree, err
	}

	return parsetree, nil
}

// IsNil checks if the given HTML node is nil.
func IsNil(node *html.Node) bool {
	return node == nil
}

// IsText checks if the given HTML node represents a text node.
func IsText(node *html.Node) bool {
	return !IsNil(node) && node.Type == html.TextNode
}

// IsElement checks if the given HTML node represents an element node.
func IsElement(node *html.Node) bool {
	return !IsNil(node) && node.Type == html.ElementNode
}

// IsElementFunc checks if the given HTML node represents an HTML element node and satisfies a custom predicate
// function.
func IsElementFunc(node *html.Node, predicate func(*html.Node) bool) bool {
	if predicate == nil {
		return IsElement(node)
	}

	return IsElement(node) && predicate(node)
}

// IsHeading checks if the given HTML node represents an HTML heading element.
func IsHeading(node *html.Node) bool {
	return IsElementFunc(node, func(n *html.Node) bool {
		_, ok := HeadingMap()[node.Data]

		return ok
	})
}

// IsHead checks if the given HTML node represents the HTML "head" element.
func IsHead(node *html.Node) bool {
	return IsElementFunc(node, func(n *html.Node) bool {
		return n.Data == "head"
	})
}

// IsBody checks if the given HTML node represents the HTML "body" element.
func IsBody(node *html.Node) bool {
	return IsElementFunc(node, func(n *html.Node) bool {
		return n.Data == "body"
	})
}

// Head returns the first occurrence of a HTML "head" element in the subtree of the
// given HTML node, or nil if not found.
func Head(node *html.Node) *html.Node {
	var head *html.Node

	if IsNil(node) {
		return nil
	}

	if IsHead(node) {
		return node
	}

	// ToDo: Check if returning nil or error is favorable.
	// ToDo: Check if IsDocument(node) check is favorable.

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if IsHead(child) {
			return child
		}

		head = Head(child)
	}

	return head
}

// Body returns the first occurrence of a HTML "body" element in the subtree of the
// given HTML node, or nil if not found.
func Body(node *html.Node) *html.Node {
	var body *html.Node

	if IsNil(node) {
		return nil
	}

	if IsBody(node) {
		return node
	}

	// ToDo: Check if returning nil or error is favorable.
	// ToDo: Check if IsDocument(node) check is favorable.

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if IsBody(child) {
			return child
		}

		body = Body(child)
	}

	return body
}

// AttributeMap returns a map of attribute key-value pairs for a given HTML node.
// If the HTML node is nil or the HTML node contains no attributes, an empty map is returned.
func AttributeMap(node *html.Node) map[string]string {
	attributeMap := make(map[string]string)

	if IsNil(node) {
		return attributeMap
	}

	for _, attribute := range node.Attr {
		attributeMap[attribute.Key] = attribute.Val
	}

	return attributeMap
}

// Headings returns a slice of HTML nodes representing the HTML heading elements
// (h1 to h6) that are direct children of the given HTML node.
func Headings(node *html.Node) []*html.Node {
	var headings []*html.Node

	if IsNil(node) {
		return headings
	}

	headings = ChildrenFunc(node, IsHeading)

	return headings
}

// HeadingMap returns a map of HTML heading tags (h1, h2, ..., h6) as keys and their corresponding
// HTML heading levels (1, 2, ..., 6) as values.
func HeadingMap() map[string]int {
	headingMap := make(map[string]int) // key: heading tag (h1, h2, ..., h6), value: heading level (1, 2, ..., 6)

	headingMap["h1"] = 1
	headingMap["h2"] = 2
	headingMap["h3"] = 3
	headingMap["h4"] = 4
	headingMap["h5"] = 5
	headingMap["h6"] = 6

	return headingMap
}

// HeadingLevel returns the HTML heading level (numeric value of "h" tag) for a given HTML heading node.
// If the given HTML node does not contain a valid HTML heading tag (h1, h2, ..., h6), it returns -1 and an error.
func HeadingLevel(node *html.Node) (int, error) {
	headingMap := HeadingMap()

	if headingLevel, ok := headingMap[node.Data]; ok {
		return headingLevel, nil
	}

	return -1, errors.New("the given node contains no valid heading tag")
}

// Children returns a slice of child HTML nodes of the provided HTML node that are either an element or a text node.
// If none of the child HTML nodes satisfies this condition, an empty slice is returned.
func Children(node *html.Node) []*html.Node {
	return ChildrenFunc(node, func(n *html.Node) bool {
		return IsElement(n) || IsText(n)
	})
}

// ChildrenFunc returns a slice of child HTML nodes of the given HTML node that satisfy the given predicate.
// The predicate function is used to filter child HTML nodes based on the provided condition.
// The returned slice contains all child HTML nodes that satisfy the predicate, in the order they appear in the HTML node.
func ChildrenFunc(node *html.Node, predicate func(*html.Node) bool) []*html.Node {
	var children []*html.Node

	if IsNil(node) {
		return children
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if predicate != nil && predicate(child) {
			children = append(children, child)
		}
	}

	return children
}

// SiblingsUntilFunc returns a slice of HTML nodes representing the siblings
// of the given HTML node, until a HTML sibling node that satisfies the provided predicate function
// is encountered. If the predicate function is nil or if no HTML node satisfies the condition
// in the predicate function, all siblings of the given HTML node are returned.
func SiblingsUntilFunc(node *html.Node, predicate func(*html.Node) bool) []*html.Node {
	var siblings []*html.Node

	if IsNil(node) {
		return siblings
	}

	for sibling := node.NextSibling; sibling != nil; sibling = sibling.NextSibling {
		if predicate != nil && predicate(sibling) {
			return siblings
		}

		siblings = append(siblings, sibling)
	}

	return siblings
}

// ElementsByTagName recursively returns a slice of HTML nodes with the given HTML tag name(s).
// If the provided HTML node is nil or if no tags are provided, an empty slice is returned.
func ElementsByTagName(node *html.Node, tags ...string) []*html.Node {
	var elements []*html.Node

	if IsNil(node) {
		return elements
	}

	if slices.IsEmpty(tags) {
		return elements
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		tag := child.Data

		if slices.Contains(tags, tag) {
			elements = append(elements, child)
		}

		elements = append(elements, ElementsByTagName(child, tags...)...)
	}

	return elements
}

// TagName returns the tag name of the given HTML node.
// If the given HTML node is nil or does not represent an HTML element, an empty string is returned.
func TagName(node *html.Node) string {
	if IsElement(node) {
		return node.Data
	}

	return ""
}

// Text recursively returns the rendered text content of an HTML node and its children.
func Text(node *html.Node) string {
	var stringBuilder strings.Builder

	if IsNil(node) {
		return ""
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if IsText(child) {
			stringBuilder.WriteString(child.Data)
		}

		if IsElement(child) {
			stringBuilder.WriteString(Text(child))
		}
	}

	return stringBuilder.String()
}

// Html returns the serialized HTML content of the given HTML node.
func Html(nodes ...*html.Node) (template.HTML, error) {
	var buffer bytes.Buffer

	for _, node := range nodes {
		if node != nil {
			err := html.Render(&buffer, node)
			if err != nil {
				return "", err
			}
		}
	}

	return template.HTML(buffer.String()), nil
}
