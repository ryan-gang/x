package vt

import "slices"

// DefaultTabInterval is the default tab interval.
const DefaultTabInterval = 8

// TabStops represents horizontal line tab stops.
type TabStops []int

// NewTabStops creates a new set of tab stops from a number of columns and an
// interval.
func NewTabStops(cols, interval int) TabStops {
	ts := make(TabStops, 0, cols/interval)
	for i := interval; i < cols; i += interval {
		ts = append(ts, i)
	}
	return ts
}

// DefaultTabStops creates a new set of tab stops with the default interval.
func DefaultTabStops(cols int) TabStops {
	return NewTabStops(cols, DefaultTabInterval)
}

// Next returns the next tab stop after the given column.
func (ts TabStops) Next(col int) int {
	// Use col+1 to ensure we get the next tab stop after the current column if
	// one exists.
	i, _ := binarySearch(ts, col+1)
	if i < len(ts) {
		return ts[i]
	}
	return col
}

// Prev returns the previous tab stop before the given column.
func (ts TabStops) Prev(col int) int {
	i, _ := binarySearch(ts, col)
	// Ensure we get the previous tab stop before the current column if one
	// exists.
	for i > 0 && ts[i-1] >= col {
		i--
	}
	if i > 0 {
		return ts[i-1]
	}
	return col
}

// Set adds a tab stop at the given column.
func (ts *TabStops) Set(col int) {
	i, ok := binarySearch(*ts, col)
	if ok {
		return
	}

	*ts = slices.Insert(*ts, i, col)
}

// Reset removes the tab stop at the given column.
func (ts *TabStops) Reset(col int) {
	i, ok := binarySearch(*ts, col)
	if !ok {
		return
	}

	*ts = slices.Delete(*ts, i, i+1)
}

// Clear removes all tab stops.
func (ts *TabStops) Clear() {
	*ts = (*ts)[:0]
}

// resetTabStops resets the terminal tab stops to the default set.
func (t *Terminal) resetTabStops() {
	t.tabstops = DefaultTabStops(t.Width())
}
