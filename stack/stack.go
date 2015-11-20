/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under lthe License.
*/

package formula

// Stack implementation as directional list
type Stack struct {
	top  *Element
	size int
}

// Element from stack. Has value and points to element behind
type Element struct {
	value interface{}
	next  *Element
}

// Returns number of elements on stack
func (s *Stack) Len() int {
	return s.size
}

// Push value on top of stack
func (s *Stack) Push(value interface{}) {
	s.top = &Element{value: value, next: s.top}
	s.size++
}

// Take value from top of stack
func (s *Stack) Pop() interface{} {
	if s.size > 0 {
		value := s.top.value
		s.top = s.top.next
		s.size--
		return value
	}
	return nil
}

// Pick value (do not take it off from stack)
func (s *Stack) Pick() interface{} {
	if s.size > 0 {
		return s.top.value
	}
	return nil
}

// Checks if stack is empty
func (s *Stack) Empty() bool {
	return s.size == 0
}