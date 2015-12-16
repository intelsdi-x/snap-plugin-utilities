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

package ns

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"reflect"
	"strconv"

	"github.com/oleiade/reflections"
	"github.com/vektra/errors"

)

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
		return errors.New("Namespace empty!\n")
	}

	return nil
}

func FromJSON(data *[]byte, current string, namespace *[]string) error {

	var m map[string]interface{}
	err := json.Unmarshal(*data, &m)

	if err != nil {
		return err
	}

	return FromMap(m, current, namespace)
}

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
		return errors.New("Namespace empty!\n")
	}

	return nil
}

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

// TODO - Value getters (GetValueByTag, GetValueByNamespace etc)
