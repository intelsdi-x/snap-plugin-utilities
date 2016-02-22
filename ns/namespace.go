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

package ns

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/oleiade/reflections"
)

// FromMap constructs list of namespaces from multilevel map using map keys as namespace entries.
// 'Current' value is prefixed to all namespace elements.
// It returns nil in case of success or error if building namespaces failed.
func FromMap(m map[string]interface{}, current string, namespace *[]string) error {

	for mkey, mval := range m {

		val := reflect.ValueOf(mval)
		typ := reflect.TypeOf(mval)
		cur := strings.Join([]string{current, mkey}, "/")
		switch val.Kind() {

		case reflect.Map:
			err := FromMap(
				mval.(map[string]interface{}),
				cur,
				namespace)
			if err != nil {
				return err
			}

		case reflect.Slice, reflect.Array:
			if typ.Elem().Kind() == reflect.Map {
				for i := 0; i < val.Len(); i++ {
					err := FromMap(
						val.Index(i).Interface().(map[string]interface{}),
						strings.Join([]string{cur, strconv.Itoa(i)}, "/"),
						namespace)
					if err != nil {
						return err
					}
				}
			} else {
				for i := 0; i < val.Len(); i++ {
					*namespace = append(*namespace, strings.Join([]string{cur, strconv.Itoa(i)}, "/"))
				}
			}

		default:
			*namespace = append(*namespace, cur)
		}
	}

	if len(*namespace) == 0 {
		return fmt.Errorf("Namespace empty!")
	}

	return nil
}

// FromJSON constructs list of namespaces from json document using json literals as namespace entries.
// 'Current' value is prefixed to all namespace elements.
// It returns nil in case of success or error if building namespaces failed.
func FromJSON(data *[]byte, current string, namespace *[]string) error {

	var m map[string]interface{}
	err := json.Unmarshal(*data, &m)

	if err != nil {
		return err
	}

	return FromMap(m, current, namespace)
}

// FromComposition constructs list of namespaces from multilevel struct compositions using field names as namespace entries.
// 'Current' value is prefixed to all namespace elements.
// It returns nil in case of success or error if building namespaces failed.
func FromComposition(object interface{}, current string, namespace *[]string) error {

	fields, err := reflections.Fields(object)

	if err != nil {
		return err

	}

	for _, field := range fields {
		f, err := reflections.GetField(object, field)

		if err != nil {
			return err
		}

		val := reflect.ValueOf(f)
		typ := reflect.TypeOf(f)
		cur := filepath.Join(current, field)

		switch reflect.ValueOf(f).Kind() {

		case reflect.Struct:
			err := FromComposition(f, cur, namespace)
			if err != nil {
				return err
			}

		case reflect.Slice, reflect.Array:
			if typ.Elem().Kind() == reflect.Struct {
				for i := 0; i < val.Len(); i++ {
					err := FromComposition(
						val.Index(i).Interface(),
						filepath.Join(cur, strconv.Itoa(i)),
						namespace)

					if err != nil {
						return err
					}
				}
			} else {
				for i := 0; i < val.Len(); i++ {
					*namespace = append(*namespace, filepath.Join(cur, strconv.Itoa(i)))
				}
			}

		default:
			*namespace = append(*namespace, cur)
		}
	}

	if len(*namespace) == 0 {
		return fmt.Errorf("Namespace empty!")
	}

	return nil
}

// FromCompositionTags constructs list of namespaces from multilevel struct composition using field tags as namespace entries.
// 'Current' value is prefixed to all namespace elements.
// It returns nil in case of success or error if building namespaces failed.
func FromCompositionTags(object interface{}, current string, namespace *[]string) error {

	data, err := json.Marshal(object)

	if err != nil {
		return err
	}

	var jmap map[string]interface{}
	err = json.Unmarshal(data, &jmap)

	if err != nil {
		return err
	}

	return FromMap(jmap, current, namespace)
}

// GetValueByNamespace returns value stored in struct composition.
// It requires filed tags on each struct field which may be represented as namespace component.
// It iterates over fields recursively, checks tags until it finds leaf value.
func GetValueByNamespace(object interface{}, ns []string) interface{} {
	// current level of namespace
	current := ns[0]
	fields, err := reflections.Fields(object)
	if err != nil {
		fmt.Printf("Could not return fields for object{%v}\n", object)
		return nil
	}

	for _, field := range fields {
		tag, err := reflections.GetFieldTag(object, field, "json")
		if err != nil {
			fmt.Printf("Could not find tag for field{%s}\n", field)
			return nil
		}
		// remove omitempty from tag
		tag = strings.Replace(tag, ",omitempty", "", -1)
		if tag == current {
			val, err := reflections.GetField(object, field)
			if err != nil {
				fmt.Printf("Could not retrieve field{%s}\n", field)
				return nil
			}
			// handling of special cases for slice and map
			switch reflect.TypeOf(val).Kind() {
			case reflect.Slice:
				idx, _ := strconv.Atoi(ns[1])
				val := reflect.ValueOf(val)
				if val.Index(idx).Kind() == reflect.Struct {
					return GetValueByNamespace(val.Index(idx).Interface(), ns[2:])
				} else {
					return val.Index(idx).Interface()
				}
			case reflect.Map:
				key := ns[1]

				if vi, ok := val.(map[string]uint64); ok {
					return vi[key]
				}

				val := reflect.ValueOf(val)
				kval := reflect.ValueOf(key)
				if reflect.TypeOf(val.MapIndex(kval).Interface()).Kind() == reflect.Struct {
					return GetValueByNamespace(val.MapIndex(kval).Interface(), ns[2:])
				}
			default:
				// last ns, return value found
				if len(ns) == 1 {
					return val
				} else {
					// or go deeper
					return GetValueByNamespace(val, ns[1:])
				}
			}
		}
	}
	return nil
}
