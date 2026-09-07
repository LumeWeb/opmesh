package opmesh

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// slicePager implements PageLister[int, int]: the filter type is unused, List
// pages over a slice honoring the Start/Limit cursor exactly like a server.
type slicePager struct {
	items []int
	err   error
}

// recordingPager is a slicePager that records every Start offset it was asked
// to serve, so tests can assert which windows a scan actually requested.
type recordingPager struct {
	items   []int
	err     error
	offsets []int
}

func (s *recordingPager) List(_ context.Context, o ListOptions[int]) ([]int, error) {
	s.offsets = append(s.offsets, o.Start)
	if s.err != nil {
		return nil, s.err
	}
	if o.Start >= len(s.items) {
		return nil, nil
	}
	end := o.Start + o.Limit
	if end > len(s.items) {
		end = len(s.items)
	}
	return s.items[o.Start:end], nil
}

func (s *slicePager) List(_ context.Context, o ListOptions[int]) ([]int, error) {
	if s.err != nil {
		return nil, s.err
	}
	if o.Start >= len(s.items) {
		return nil, nil
	}
	end := o.Start + o.Limit
	if end > len(s.items) {
		end = len(s.items)
	}
	return s.items[o.Start:end], nil
}

func hasValue(v int) MatchPredicate[int] {
	return func(item int) (bool, error) {
		if item == -987 {
			return false, errors.New("predicate failed")
		}
		return item == v, nil
	}
}

func TestScanPages_FindsLaterPage(t *testing.T) {
	// Target far beyond the first defaultPageSize page.
	items := make([]int, 33)
	for i := range items {
		items[i] = i + 1
	}
	got, ok, err := ScanPages(context.Background(), &slicePager{items: items}, hasValue(31))
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 31, got)
}

func TestScanPages_NotFound(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	got, ok, err := ScanPages(context.Background(), &slicePager{items: items}, hasValue(999))
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Zero(t, got)
}

func TestScanPages_StopsAtShortPage(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	got, ok, err := ScanPages(context.Background(), &slicePager{items: items}, hasValue(11))
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 11, got)
}

func TestScanPagesWithOptions_ExplicitPageSize(t *testing.T) {
	items := make([]int, 40)
	for i := range items {
		items[i] = i + 1
	}
	got, ok, err := ScanPagesWithOptions(context.Background(), &slicePager{items: items}, hasValue(40), &ListOptions[int]{}, 25)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 40, got)
}

// TestScanPagesWithOptions_HonorsBaseStart verifies a caller-supplied
// base.Start is the first window actually requested: the scan must pick up
// from the given offset (finding an item on a later page from there) without
// ever re-reading the page(s) before it.
func TestScanPagesWithOptions_HonorsBaseStart(t *testing.T) {
	// 30 items, page size 10: page 0 = [0..10), page 1 = [10..20), where the
	// target (value 16, index 15) lives.
	items := make([]int, 30)
	for i := range items {
		items[i] = i + 1
	}
	pager := &recordingPager{items: items}
	got, ok, err := ScanPagesWithOptions(context.Background(), pager, hasValue(16), &ListOptions[int]{Start: 10}, 10)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 16, got)
	// The first requested offset must be base.Start, and page 0 must never
	// have been re-read.
	assert.Equal(t, []int{10}, pager.offsets)
}

// TestScanPagesWithOptions_StopsOnNonAdvancingWindow verifies the scan
// terminates at end of data even when a pager never returns a short window:
// a page that cannot advance the window (an empty page here) stops the scan
// instead of looping forever.
func TestScanPagesWithOptions_StopsOnNonAdvancingWindow(t *testing.T) {
	// The target is outside the data; the pager returns a full first page
	// (== pageSize, so the short-page rule would not stop the scan) followed
	// by empty pages.
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	pager := &recordingPager{items: items}
	got, ok, err := ScanPagesWithOptions(context.Background(), pager, hasValue(999), &ListOptions[int]{}, 10)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Zero(t, got)
	assert.Equal(t, []int{0, 10}, pager.offsets)
}

func TestScanPages_LoadError(t *testing.T) {
	_, _, err := ScanPages(context.Background(), &slicePager{err: errors.New("list failed")}, hasValue(1))
	require.ErrorContains(t, err, "list failed")
}

func TestScanPages_PredicateError(t *testing.T) {
	_, _, err := ScanPages(context.Background(), &slicePager{items: []int{-987}}, hasValue(1))
	require.ErrorContains(t, err, "predicate failed")
}
