/*
 * Copyright (C) 2023 Stefan Kühnel
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

package slices

// Contains reports whether an element is present in a slice.
func Contains[Element comparable](slice []Element, element Element) bool {
	// Copyright (c) 2021 The Go Authors. All rights reserved.
	//
	// Redistribution and use in source and binary forms, with or without
	// modification, are permitted provided that the following conditions are
	// met:
	//
	//    * Redistributions of source code must retain the above copyright
	// notice, this list of conditions and the following disclaimer.
	//    * Redistributions in binary form must reproduce the above
	// copyright notice, this list of conditions and the following disclaimer
	// in the documentation and/or other materials provided with the
	// distribution.
	//    * Neither the name of Google Inc. nor the names of its
	// contributors may be used to endorse or promote products derived from
	// this software without specific prior written permission.
	//
	// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
	// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
	// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
	// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
	// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
	// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
	// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
	// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
	// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
	// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
	//
	// See: https://cs.opensource.google/go/x/exp/+/10a50721:slices/slices.go;l=127
	return Index(slice, element) >= 0
}

// FindFunc returns the first element in a slice that satisfies predicate(element).
func FindFunc[Element comparable](slice []Element, predicate func(Element) bool) Element {
	var empty Element

	for _, element := range slice {
		if predicate(element) {
			return element
		}
	}

	return empty
}

// Index returns the index of the first occurrence of element in slice,
// or -1 if not present.
func Index[Element comparable](slice []Element, element Element) int {
	// Copyright (c) 2021 The Go Authors. All rights reserved.
	//
	// Redistribution and use in source and binary forms, with or without
	// modification, are permitted provided that the following conditions are
	// met:
	//
	//    * Redistributions of source code must retain the above copyright
	// notice, this list of conditions and the following disclaimer.
	//    * Redistributions in binary form must reproduce the above
	// copyright notice, this list of conditions and the following disclaimer
	// in the documentation and/or other materials provided with the
	// distribution.
	//    * Neither the name of Google Inc. nor the names of its
	// contributors may be used to endorse or promote products derived from
	// this software without specific prior written permission.
	//
	// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
	// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
	// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
	// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
	// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
	// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
	// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
	// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
	// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
	// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
	//
	// See: https://cs.opensource.google/go/x/exp/+/10a50721:slices/slices.go;l=106
	for index, current := range slice {
		if element == current {
			return index
		}
	}
	return -1
}

// Insert inserts an element at an index into a slice, returning a shallow copy.
func Insert[Element comparable](slice []Element, element Element, index int) []Element {
	size := 0

	// no increase of size required
	isIndexWithinBounds := index < len(slice)
	if isIndexWithinBounds {
		size = len(slice)
	}

	// increase of size required
	isIndexOutOfBounds := index >= len(slice)
	if isIndexOutOfBounds {
		boundsDifference := (index + 1) - len(slice)
		size = len(slice) + boundsDifference
	}

	// create a new slice
	values := make([]Element, size)

	// create a shallow copy with possible zero values
	copy(values, slice)

	// insert or replace element
	values[index] = element

	return values
}

// IsEmpty returns true if the slice is empty, false otherwise.
func IsEmpty[Element comparable](slice []Element) bool {
	return len(slice) == 0
}
