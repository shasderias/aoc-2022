package grid

import "sort"

type Range struct {
	Min, Max int
}

func (r1 Range) Width() int {
	return r1.Max - r1.Min + 1
}

func (r1 Range) Overlap(r2 Range) bool {
	return r1.Min <= r2.Max && r2.Min <= r1.Max
}

func (r1 Range) Union(r2 Range) Range {
	if !r1.Overlap(r2) {
		panic("cannot union non-overlapping ranges")
	}
	return Range{Min: min(r1.Min, r2.Min), Max: max(r1.Max, r2.Max)}
}

type RangeSet []Range

func (rs *RangeSet) Add(r Range) {
	*rs = append(*rs, r)
	rs.sort()
	rs.consolidate()
}

func (rs *RangeSet) consolidate() {
	for {
		changed := false
		for i := 0; i < len(*rs)-1; i++ {
			if (*rs)[i].Overlap((*rs)[i+1]) {
				(*rs)[i] = (*rs)[i].Union((*rs)[i+1])
				*rs = append((*rs)[:i+1], (*rs)[i+2:]...)
				changed = true
				break
			}
		}
		if !changed {
			break
		}
	}
}

func (rs *RangeSet) sort() {
	sort.Slice(*rs, func(i, j int) bool {
		switch {
		case (*rs)[i].Min < (*rs)[j].Min:
			return true
		case (*rs)[i].Min == (*rs)[j].Min:
			return (*rs)[i].Max < (*rs)[j].Max
		}
		return false
	})

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
