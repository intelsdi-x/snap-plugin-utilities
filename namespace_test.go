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
	"testing"
	"encoding/json"

	. "github.com/smartystreets/goconvey/convey"

	_"fmt"

)

func TestSimpleMap(t *testing.T) {
	Convey("Given flat strig to string map", t, func() {
		m := map[string]interface{}{
			"Foo": "foo",
			"Bar": "bar",
			"Baz": "baz",
		}

		Convey("When NamespaceFromMap is called with root as current", func() {
			ns := []string{}
			current := "root"
			NamespaceFromMap(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 3)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/Foo")
				So(ns, ShouldContain, "root/Bar")
				So(ns, ShouldContain, "root/Baz")
			})
		})
	})
}

func TestMapWithSlice(t *testing.T) {
	Convey("Given two leyer map with slice", t, func() {
		Foo := []string{"foo_0", "foo_1"}
		Bar := []string{"bar_0"}
		m := map[string]interface{}{
			"Foo": Foo,
			"Bar": Bar,
		}

		Convey("When NamespaceFromMap is called with root as current", func() {
			ns := []string{}
			current := "root"
			NamespaceFromMap(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 3)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/Foo/0")
				So(ns, ShouldContain, "root/Foo/1")
				So(ns, ShouldContain, "root/Bar/0")
			})
		})
	})
}

func TestMapWithMap(t *testing.T) {
	Convey("Given two leyer nested map", t, func() {
		Foo := map[string]interface{}{"Foos": "foos", "Boos": "boos"}
		Bar := map[string]interface{}{"Goos": "goos"}
		m := map[string]interface{}{
			"Foo": Foo,
			"Bar": Bar,
		}

		Convey("When NamespaceFromMap is called with root as current", func() {
			ns := []string{}
			current := "root"
			NamespaceFromMap(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 3)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/Foo/Foos")
				So(ns, ShouldContain, "root/Foo/Boos")
				So(ns, ShouldContain, "root/Bar/Goos")
			})
		})
	})
}

func TestMapComposition(t *testing.T) {
	Convey("Given composition map", t, func() {
		Foo := map[string]interface{}{"Foos": "foos", "Boos": "boos"}
		Bar := []string{"1", "2", "3"}
		m := map[string]interface{}{
			"Foo": Foo,
			"Bar": Bar,
		}

		Convey("When NamespaceFromMap is called with root as current", func() {
			ns := []string{}
			current := "root"
			NamespaceFromMap(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 5)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/Foo/Foos")
				So(ns, ShouldContain, "root/Foo/Boos")
				So(ns, ShouldContain, "root/Bar/0")
				So(ns, ShouldContain, "root/Bar/1")
				So(ns, ShouldContain, "root/Bar/2")
			})
		})
	})
}

func TestMapCompositionComplex(t *testing.T) {
	Convey("Given complex composition in map", t, func() {
		Baz := map[string]interface{}{"Bazo": "bazo", "Fazo": "fazo", "Mazo": "mazo"}
		Foo := map[string]interface{}{"Foos": "foos", "Boos": "boos"}
		Bar := []map[string]interface{}{Baz, Baz}
		m := map[string]interface{}{
			"Foo": Foo,
			"Bar": Bar,
		}

		Convey("When NamespaceFromMap is called with root as current", func() {
			ns := []string{}
			current := "root"
			NamespaceFromMap(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 8)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/Foo/Foos")
				So(ns, ShouldContain, "root/Foo/Boos")
				So(ns, ShouldContain, "root/Bar/0/Bazo")
				So(ns, ShouldContain, "root/Bar/0/Fazo")
				So(ns, ShouldContain, "root/Bar/0/Mazo")
				So(ns, ShouldContain, "root/Bar/1/Bazo")
				So(ns, ShouldContain, "root/Bar/1/Fazo")
				So(ns, ShouldContain, "root/Bar/1/Mazo")
			})
		})
	})
}

func TestSimpleJson(t *testing.T) {
	Convey("Given flat strig to string map", t, func() {

		Foo := struct {
			Bar int `json:"bar"`
			Baz string `json:"baz"`
		}{
			Bar: 1,
			Baz: "1",
		}

		data, _ := json.Marshal(Foo)

		Convey("When NamespaceFromMap is called with root as current", func() {
			ns := []string{}
			current := "root"
			NamespaceFromJSON(&data, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 2)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/bar")
				So(ns, ShouldContain, "root/baz")
			})
		})
	})
}
