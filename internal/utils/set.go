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
