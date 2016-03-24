/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package strings

// StringMap is custom wrapper for golang map
// It provides additional capabilities like list of keys, list of values, Empty
type StringMap struct {
	m map[string]string
}

// NewEmptyMap creates new StringMap
func NewEmptyMap() StringMap {
	m := map[string]string{}
	return StringMap{m: m}
}

// NewEmptyMap creates new StringMap
func NewFromMap(m map[string]string) StringMap {
	return StringMap{m: m}
}

// AddMap adds all key/value pairs from map m to StringMap
// If key already exists, it will be overwritten by new value
func (sm StringMap) AddMap(m map[string]string) {
	if sm.Size() == 0 {
		sm.m = m
	} else {
		for k, v := range m {
			sm.m[k] = v
		}
	}
}

// Add key/value pair to map
func (sm StringMap) Add(key, value string) {
	sm.m[key] = value
}

// Remove key/value pair from map
func (sm StringMap) Remove(key string) {
	delete(sm.m, key)
}

// RemoveAll deletes all key/value pairs from map
func (sm StringMap) RemoveAll() {
	sm.m = map[string]string{}
}


// Size return number of key/value pairs
func (sm StringMap) Size() int {
	return len(sm.m)
}

// Empty checks if map contains any key/value pair
func (sm StringMap) Empty() bool {
	return len(sm.m) == 0
}

// Keys returns slice of map keys
func (sm StringMap) Keys() []string {
	keys := []string{}
	for k, _ := range sm.m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns slice of map values
func (sm StringMap) Values() []string {
	values := []string{}
	for _, v := range sm.m {
		values = append(values, v)
	}
	return values
}

// Get returns value for given key
func (sm StringMap) Get(key string) string {
	return sm.m[key]
}

// HasKey checks if key exists in map
func (sm StringMap) HasKey(key string) bool {
	_, found := sm.m[key]
	return found
}
