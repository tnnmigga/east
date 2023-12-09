package algorithms

type Set[T comparable] map[T]struct{}

func (s Set[T]) Insert(key T) bool {
	if s.Find(key) {
		return false
	}
	s[key] = struct{}{}
	return true
}

func (s Set[T]) Find(key T) bool {
	if _, has := s[key]; has {
		return true
	}
	return false
}

func (s Set[T]) Delete(key T) {
	delete(s, key)
}

func (s Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s))
	for k := range s {
		slice = append(slice, k)
	}
	return slice
}
