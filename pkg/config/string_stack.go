package config

// stringStack is an implementation of stack mainly to support `jsonLoader`.
//
// TODO: write tests stringStack?
type stringStack []string

func (s *stringStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *stringStack) Push(str string) {
	*s = append(*s, str)
}

func (s *stringStack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

func (s *stringStack) Peek() string {
	if s.IsEmpty() {
		return ""
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		return element
	}
}

func (s *stringStack) Count() int {
	return len(*s)
}
