// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sort provides primitives for sorting slices and user-defined
// collections.
package bools

// Sort sorts data.
// It makes one call to data.Len to determine n, and Operator(n*log(n)) calls to
// data.Less and data.Swap. The sort is not guaranteed to be stable.
func Sort(vs []bool, os []int64) {
	n := len(os)
	quickSort(vs, os, 0, n, maxDepth(n))
}

// maxDepth returns a threshold at which quicksort should switch
// to heapsort. It returns 2*ceil(lg(n+1)).
func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func BoolLess(x, y bool) bool {
	if x != y {
		return !x
	}
	return false
}

func BoolGreatEq(x, y bool) bool {
	if x != y {
		return x
	}
	return true
}

func quickSort(vs []bool, os []int64, a, b, maxDepth int) {
	for b-a > 12 { // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 {
			heapSort(vs, os, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(vs, os, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(vs, os, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(vs, os, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ {
			if BoolLess(vs[os[i]], vs[os[i-6]]) {
				os[i], os[i-6] = os[i-6], os[i]
			}
		}
		insertionSort(vs, os, a, b)
	}
}

// Insertion sort
func insertionSort(vs []bool, os []int64, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && BoolLess(vs[os[j]], vs[os[j-1]]); j-- {
			os[j], os[j-1] = os[j-1], os[j]
		}
	}
}

// siftDown implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown(vs []bool, os []int64, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && BoolLess(vs[os[first+child]], vs[os[first+child+1]]) {
			child++
		}
		if BoolGreatEq(vs[os[first+root]], vs[os[first+child]]) {
			return
		}
		os[first+root], os[first+child] = os[first+child], os[first+root]
		root = child
	}
}

func heapSort(vs []bool, os []int64, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(vs, os, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		os[first], os[first+i] = os[first+i], os[first]
		siftDown(vs, os, lo, i, first)
	}
}

// Quicksort, loosely following Bentley and McIlroy,
// ``Engineering a Sort Function,'' SP&E November 1993.

// medianOfThree moves the median of the three values data[m0], data[m1], data[m2] into data[m1].
func medianOfThree(vs []bool, os []int64, m1, m0, m2 int) {
	// sort 3 elements
	if BoolLess(vs[os[m1]], vs[os[m0]]) {
		os[m1], os[m0] = os[m0], os[m1]
	}
	// data[m0] <= data[m1]
	if BoolLess(vs[os[m2]], vs[os[m1]]) {
		os[m2], os[m1] = os[m1], os[m2]
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if BoolLess(vs[os[m1]], vs[os[m0]]) {
			os[m1], os[m0] = os[m0], os[m1]
		}
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func swapRange(vs []bool, os []int64, a, b, n int) {
	for i := 0; i < n; i++ {
		os[a+i], os[b+i] = os[b+i], os[a+i]
	}
}

func doPivot(vs []bool, os []int64, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(vs, os, lo, lo+s, lo+2*s)
		medianOfThree(vs, os, m, m-s, m+s)
		medianOfThree(vs, os, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree(vs, os, lo, m, hi-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && BoolLess(vs[os[a]], vs[os[pivot]]); a++ {
	}
	b := a
	for {
		for ; b < c && BoolGreatEq(vs[os[pivot]], vs[os[b]]); b++ { // data[b] <= pivot
		}
		for ; b < c && BoolLess(vs[os[pivot]], vs[os[c-1]]); c-- { // data[c-1] > pivot
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		os[b], os[c-1] = os[c-1], os[b]
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let's be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if BoolGreatEq(vs[os[pivot]], vs[os[hi-1]]) { // data[hi-1] = pivot
			os[c], os[hi-1] = os[hi-1], os[c]
			c++
			dups++
		}
		if BoolGreatEq(vs[os[b-1]], vs[os[pivot]]) { // data[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if BoolGreatEq(vs[os[m]], vs[os[pivot]]) { // data[m] = pivot
			os[m], os[b-1] = os[b-1], os[m]
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			for ; a < b && BoolGreatEq(vs[os[b-1]], vs[os[pivot]]); b-- { // data[b] == pivot
			}
			for ; a < b && BoolLess(vs[os[a]], vs[os[pivot]]); a++ { // data[a] < pivot
			}
			if a >= b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			os[a], os[b-1] = os[b-1], os[a]
			a++
			b--
		}
	}
	// Swap pivot into middle
	os[pivot], os[b-1] = os[b-1], os[pivot]
	return b - 1, c
}
