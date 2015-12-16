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
limitations under the License.
*/

package strings

func ForEach(slice []string, handler func(string) string) {
	for i, element := range slice {
		slice[i] = handler(element)
	}
}

func Contains(slice []string, lookup string) bool {
	for _, element := range slice {
		if element == lookup {
			return true
		}
	}
	return false
}

func Filter(slice []string, handler func(string) bool) []string {
	filtered := []string{}
	for _, element := range slice {
		if handler(element) {
			filtered = append(filtered, element)
		}
	}
	return filtered
}

func Any(slice []string, handler func(string) bool) bool {
	for _, element := range slice {
		if handler(element) {
			return true
		}
	}
	return false
}

func All(slice []string, handler func(string) bool) bool {
	for _, element := range slice {
		if !handler(element) {
			return false
		}
	}
	return true
}

