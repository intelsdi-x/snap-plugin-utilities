// +build unit

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
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"fmt"
	_ "fmt"
)

func TestSimpleMap(t *testing.T) {
	Convey("Given flat string to string map", t, func() {
		m := map[string]interface{}{
			"Foo": "foo",
			"Bar": "bar",
			"Baz": "baz",
		}

		Convey("When NamespaceFromMap is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromMap(m, current, &ns)

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
			FromMap(m, current, &ns)

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
			FromMap(m, current, &ns)

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
			FromMap(m, current, &ns)

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
			FromMap(m, current, &ns)

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
	Convey("Given flat struct", t, func() {

		Foo := struct {
			Bar int    `json:"bar"`
			Baz string `json:"baz"`
		}{
			Bar: 1,
			Baz: "1",
		}

		data, _ := json.Marshal(Foo)

		Convey("When NamespaceFromJSON is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromJSON(&data, current, &ns)

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

func TestComplexJson(t *testing.T) {
	Convey("Given composition ofd structs", t, func() {

		Foo := struct {
			Bar struct {
				Qaz int `json:"qaz"`
				Faz int `json:"faz"`
			} `json:"bar"`
			Baz string `json:"baz"`
		}{
			struct {
				Qaz int `json:"qaz"`
				Faz int `json:"faz"`
			}{
				1,
				2,
			},
			"baz_val",
		}

		data, _ := json.Marshal(Foo)

		Convey("When NamespaceFromJSON is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromJSON(&data, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 3)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/bar/qaz")
				So(ns, ShouldContain, "root/bar/faz")
				So(ns, ShouldContain, "root/baz")
			})
		})
	})
}

func TestComplexJsonSlice(t *testing.T) {
	Convey("Given composition of structs with slice", t, func() {

		Foo := struct {
			Bar struct {
				Qaz int `json:"qaz"`
				Faz int `json:"faz"`
			} `json:"bar"`
			Baz []string `json:"baz"`
		}{
			struct {
				Qaz int `json:"qaz"`
				Faz int `json:"faz"`
			}{
				1,
				2,
			},
			[]string{"baz_val_1", "baz_val_2", "baz_val_3"},
		}

		data, _ := json.Marshal(Foo)

		Convey("When NamespaceFromJSON is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromJSON(&data, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 5)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/bar/qaz")
				So(ns, ShouldContain, "root/bar/faz")
				So(ns, ShouldContain, "root/baz/0")
				So(ns, ShouldContain, "root/baz/1")
				So(ns, ShouldContain, "root/baz/2")
			})
		})
	})
}

func TestComplexJsonSliceNested(t *testing.T) {
	Convey("Given composition of structs with nested slice", t, func() {

		Foo := struct {
			Bar struct {
				Qaz []int `json:"qaz"`
				Faz int   `json:"faz"`
			} `json:"bar"`
			Baz string `json:"baz"`
		}{
			struct {
				Qaz []int `json:"qaz"`
				Faz int   `json:"faz"`
			}{
				[]int{1, 2},
				2,
			},
			"baz_val",
		}

		data, _ := json.Marshal(Foo)

		Convey("When NamespaceFromJSON is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromJSON(&data, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 4)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/bar/qaz/0")
				So(ns, ShouldContain, "root/bar/qaz/1")
				So(ns, ShouldContain, "root/bar/faz")
				So(ns, ShouldContain, "root/baz")
			})
		})
	})
}

func TestSimpleStruct(t *testing.T) {
	Convey("Given flat struct", t, func() {

		Foo := struct {
			Bar int
			Baz string
		}{
			Bar: 1,
			Baz: "1",
		}

		Convey("When NamespaceFromComposition is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromComposition(Foo, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 2)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/Bar")
				So(ns, ShouldContain, "root/Baz")
			})
		})
	})
}

func TestComplexStruct(t *testing.T) {
	Convey("Given composition of structs", t, func() {
		fmt.Printf("\nTU\n")
		Foo := struct {
			Bar struct {
				Qaz int
				Faz int
			}
			Baz string
		}{
			struct {
				Qaz int
				Faz int
			}{
				1,
				2,
			},
			"baz_val",
		}

		Convey("When NamespaceFromComposition is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromComposition(Foo, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 3)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/Bar/Qaz")
				So(ns, ShouldContain, "root/Bar/Faz")
				So(ns, ShouldContain, "root/Baz")
			})
		})
	})
}

func TestComplexStructSlice(t *testing.T) {
	Convey("Given composition of structs with slice", t, func() {

		Foo := struct {
			Bar struct {
				Qaz int
				Faz int
			}
			Baz []string
		}{
			struct {
				Qaz int
				Faz int
			}{
				1,
				2,
			},
			[]string{"baz_val_1", "baz_val_2", "baz_val_3"},
		}

		Convey("When NamespaceFromComposition is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromComposition(Foo, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 5)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/Bar/Qaz")
				So(ns, ShouldContain, "root/Bar/Faz")
				So(ns, ShouldContain, "root/Baz/0")
				So(ns, ShouldContain, "root/Baz/1")
				So(ns, ShouldContain, "root/Baz/2")
			})
		})
	})
}

func TestComplexCompositionSliceNested(t *testing.T) {
	Convey("Given composition of structs with nested slice", t, func() {

		Foo := struct {
			Bar struct {
				Qaz []int
				Faz int
			}
			Baz string
		}{
			struct {
				Qaz []int
				Faz int
			}{
				[]int{1, 2},
				2,
			},
			"baz_val",
		}

		Convey("When NamespaceFromComposition is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromComposition(Foo, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 4)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/Bar/Qaz/0")
				So(ns, ShouldContain, "root/Bar/Qaz/1")
				So(ns, ShouldContain, "root/Bar/Faz")
				So(ns, ShouldContain, "root/Baz")
			})
		})
	})
}

func TestSimpleCompositionTags(t *testing.T) {
	Convey("Given flat struct with json tags", t, func() {

		Foo := struct {
			Bar int    `json:"bar"`
			Baz string `json:"baz"`
		}{
			Bar: 1,
			Baz: "1",
		}

		Convey("When NamespaceFromJSON is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromCompositionTags(Foo, current, &ns)

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

func TestComplexCompositionTags(t *testing.T) {
	Convey("Given composition of structs with json tags", t, func() {

		Foo := struct {
			Bar struct {
				Qaz int `json:"qaz"`
				Faz int `json:"faz"`
			} `json:"bar"`
			Baz string `json:"baz"`
		}{
			struct {
				Qaz int `json:"qaz"`
				Faz int `json:"faz"`
			}{
				1,
				2,
			},
			"baz_val",
		}

		Convey("When NamespaceFromJSON is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromCompositionTags(Foo, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 3)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/bar/qaz")
				So(ns, ShouldContain, "root/bar/faz")
				So(ns, ShouldContain, "root/baz")
			})
		})
	})
}

func TestComplexCompositionSliceTags(t *testing.T) {
	Convey("Given composition of structs with slice and json tags", t, func() {

		Foo := struct {
			Bar struct {
				Qaz int `json:"qaz"`
				Faz int `json:"faz"`
			} `json:"bar"`
			Baz []string `json:"baz"`
		}{
			struct {
				Qaz int `json:"qaz"`
				Faz int `json:"faz"`
			}{
				1,
				2,
			},
			[]string{"baz_val_1", "baz_val_2", "baz_val_3"},
		}

		Convey("When NamespaceFromJSON is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromCompositionTags(Foo, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 5)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/bar/qaz")
				So(ns, ShouldContain, "root/bar/faz")
				So(ns, ShouldContain, "root/baz/0")
				So(ns, ShouldContain, "root/baz/1")
				So(ns, ShouldContain, "root/baz/2")
			})
		})
	})
}

func TestComplexCompositionTagsSliceNested(t *testing.T) {
	Convey("Given composition of structs with nested slice", t, func() {

		Foo := struct {
			Bar struct {
				Qaz []int `json:"qaz"`
				Faz int   `json:"faz"`
			} `json:"bar"`
			Baz string `json:"baz"`
		}{
			struct {
				Qaz []int `json:"qaz"`
				Faz int   `json:"faz"`
			}{
				[]int{1, 2},
				2,
			},
			"baz_val",
		}

		Convey("When NamespaceFromJSON is called with root as current", func() {
			ns := []string{}
			current := "root"
			FromCompositionTags(Foo, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 4)
			})

			Convey("Then namespace should contain correct entries", func() {
				So(ns, ShouldContain, "root/bar/qaz/0")
				So(ns, ShouldContain, "root/bar/qaz/1")
				So(ns, ShouldContain, "root/bar/faz")
				So(ns, ShouldContain, "root/baz")
			})
		})
	})
}

func TestReplaceNotAllowedCharsInNamespacePart(t *testing.T) {
	Convey("Given list of namespace parts with not allowed chars", t, func() {

		incorrectMetricParts := []string{"test(test1)", "test[test2]", "test{test3}",
			"test 4", "test.5", "test,6", "test;7", "test!8", "test?9", "test|10",
			"test\\11", "test/12", "test^13", "test\"14", "test`15", "test'16"}

		Convey("When ReplaceSpecialCharsInNamespacePart is called namespace parts should contain only allowed characters", func() {
			correctMetricParts := []string{"test_test1", "test_test2", "test_test3",
				"test_4", "test_5", "test_6", "test_7", "test_8", "test_9", "test_10",
				"test_11", "test_12", "test_13", "test_14", "test_15", "test_16"}

			for i := range incorrectMetricParts {
				correctedMetricParts := ReplaceNotAllowedCharsInNamespacePart(incorrectMetricParts[i])
				So(correctedMetricParts, ShouldEqual, correctMetricParts[i])
				So(ValidateMetricNamespacePart(correctedMetricParts), ShouldBeNil)
			}
		})
	})
}

func TestValidateMetricNamespacePart(t *testing.T) {
	Convey("Given list of namespace parts with not allowed chars", t, func() {

		incorrectMetricParts := []string{"test1)", "[test2]", "{test3}",
			"test 4", "test.5", "test,6", "test;7", "test!8", "test?9", "test|10",
			"test\\11", "test/12", "test^13", "test\"14", "test`15", "test'16"}

		Convey("When ValidateMetricNamespace is called not allowed chars should be detected", func() {
			for i := range incorrectMetricParts {
				So(ValidateMetricNamespacePart(incorrectMetricParts[i]), ShouldNotBeNil)
			}
		})
	})
}
