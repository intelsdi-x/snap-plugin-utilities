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

package utilities

import (
	"reflect"
	"strconv"
	"encoding/json"
	"path/filepath"

	"github.com/vektra/errors"
	"github.com/oleiade/reflections"
)

func NamespaceFromMap(m map[string]interface{}, current string, namespace *[]string) error {

	for mkey, mval := range m {

		val := reflect.ValueOf(mval)
		typ := reflect.TypeOf(mval)
		cur := filepath.Join(current, mkey)

		switch val.Kind() {

		case reflect.Map:
			NamespaceFromMap(
				mval.(map[string]interface{}),
				cur,
				namespace)

		case reflect.Slice, reflect.Array:
			if typ.Elem().Kind() == reflect.Map {
				for i := 0; i < val.Len(); i++ {
					NamespaceFromMap(
						val.Index(i).Interface().(map[string]interface{}),
						filepath.Join(cur, strconv.Itoa(i)),
						namespace)
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
		return errors.New("Namespace empty!\n")
	}

	return nil
}

func NamespaceFromJSON(data *[]byte, current string, namespace *[]string) error {

	var m map[string]interface{}
	err := json.Unmarshal(*data, &m)

	if err != nil {
		return err
	}

	return NamespaceFromMap(m, current, namespace)
}

func NamespaceFromComposition(object interface{}, current string, namespace *[]string) error {

	fields, err := reflections.Fields(object)

	if err != nil {
		return err

	}

	for _, field := range fields {
		tag, err := reflections.GetFieldTag(object, field, "json")

		if err != nil {
			return err
		}


	}

	return nil
}
