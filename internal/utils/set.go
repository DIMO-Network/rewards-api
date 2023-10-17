package utils

import "golang.org/x/exp/maps"

type Set[A comparable] map[A]struct{}

func NewSet[A comparable](as ...A) Set[A] {
	s := make(Set[A])
	for _, a := range as {
		s.Add(a)
	}
	return s
}

func (s Set[A]) Add(a A) {
	s[a] = struct{}{}
}

func (s Set[A]) Remove(a A) {
	delete(s, a)
}

func (s Set[A]) Contains(a A) bool {
	_, ok := s[a]
	return ok
}

func (s Set[A]) Slice() []A {
	return maps.Keys(s)
}

func (s Set[A]) Len() int {
	return len(s)
}

func (s Set[A]) Union(other Set[A]) Set[A] {
	union := NewSet[A]()
	for a := range s {
		union.Add(a) // make copy of original set
	}
	for a := range other {
		union.Add(a) // add elements from other set
	}
	return union
}

func (s Set[A]) Intersection(other Set[A]) Set[A] {
	intersection := NewSet[A]()
	for a := range s {
		if other.Contains(a) {
			intersection.Add(a)
		}
	}
	return intersection
}

func (s Set[A]) Difference(other Set[A]) Set[A] {
	difference := NewSet[A]()
	for a := range s {
		if !other.Contains(a) {
			difference.Add(a)
		}
	}
	return difference
}
