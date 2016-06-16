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
	"reflect"
	"strings"
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

func TestFromCompositeObject(t *testing.T) {

	//// checking the effect of options, first - nil pointers
	Convey("Given struct with nil pointers", t, func() {
		type delta struct {
			Eins int
		}
		type uno struct {
			Alpha bool
			Delta *delta
		}
		m := struct {
			First *uno
		}{}

		Convey("When NamespaceFromCompositeObject is called with default options", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns)

			Convey("Then nil pointers should be expanded by default at all levels", func() {
				So(ns, ShouldContain, "root/First/Alpha")
				So(ns, ShouldContain, "root/First/Delta/Eins")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with nil pointer expansion disabled at all levels", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, InspectNilPointers(AlwaysFalse))

			Convey("Then no nil pointers should be expanded at any level", func() {
				So(ns, ShouldNotContain, "root/First/Alpha")
				So(ns, ShouldNotContain, "root/First/Delta/Eins")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with nil pointer expansion disabled at second level", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, InspectNilPointers(func(current string, typ reflect.Type) bool {
				return strings.Count(current, "/") < 2
			}))

			Convey("Then nil pointer should be expanded only at first level", func() {
				So(ns, ShouldContain, "root/First/Alpha")
				So(ns, ShouldNotContain, "root/First/Delta/Eins")
			})
		})
	})

	//// checking effect on empty containers
	Convey("Given struct with empty containers", t, func() {
		type dos struct {
			Alpha int
		}
		type first struct {
			Uno bool
			Dos []dos
		}
		m := struct {
			First []first
		}{}

		Convey("When NamespaceFromCompositeObject is called with default options", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns)

			Convey("Then empty containers should be expanded at all levels", func() {
				So(ns, ShouldContain, "root/First/*/Uno")
				So(ns, ShouldContain, "root/First/*/Dos/*/Alpha")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with disabled empty container expansion at all levels", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, InspectEmptyContainers(AlwaysFalse))

			Convey("Then empty containers should be expanded at all levels", func() {
				So(ns, ShouldNotContain, "root/First/*/Uno")
				So(ns, ShouldNotContain, "root/First/*/Dos/*/Alpha")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with empty container expansion disabled at second level", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, InspectEmptyContainers(func(current string, typ reflect.Type) bool {
				return strings.Count(current, "/") < 3
			}))

			Convey("Then nil pointer should be expanded only at first level", func() {
				So(ns, ShouldContain, "root/First/*/Uno")
				So(ns, ShouldNotContain, "root/First/*/Dos/*/Alpha")
			})
		})
	})

	Convey("Given struct with some containers", t, func() {
		type dos struct {
			Alpha int
		}
		type first struct {
			Uno bool
			Dos []dos
		}
		m := struct {
			First map[string]first
		}{}

		Convey("When NamespaceFromCompositeObject is called with default options", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns)

			Convey("Then no empty container should be expanded at any levels", func() {
				So(ns, ShouldNotContain, "root/First")
				So(ns, ShouldNotContain, "root/First/*/Dos")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with enabled entries for container roots at all levels", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, EntryForContainersRoot(AlwaysTrue))

			Convey("Then containers should have own entries at all levels", func() {
				So(ns, ShouldContain, "root/First")
				So(ns, ShouldContain, "root/First/*/Dos")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with entries for container roots disabled at second level", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, EntryForContainersRoot(func(current string, typ reflect.Type) bool {
				return strings.Count(current, "/") < 3
			}))

			Convey("Then entry for container root should be added only at first level", func() {
				So(ns, ShouldContain, "root/First")
				So(ns, ShouldNotContain, "root/First/*/Dos")
			})
		})
	})

	//// checking json tag expansion
	Convey("Given struct with some structs annotated with json tags", t, func() {
		type dos struct {
			Alpha int `json:""`
		}
		type first struct {
			Uno  bool `json:"uno_f,omitempty"`
			Dos  dos
			Tres string `json:"-,omitempty"`
		}
		m := struct {
			First  first
			Second string `json:"second_f"`
			Third  bool   `json:"-"`
			Fourth int    `json:",omitempty"`
		}{}

		Convey("When NamespaceFromCompositeObject is called with default options", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns)

			Convey("Then fields at all levels should be exported basing on json tags", func() {
				So(ns, ShouldContain, "root/First/uno_f")
				So(ns, ShouldContain, "root/First/Dos/Alpha")
				So(ns, ShouldContain, "root/second_f")
				So(ns, ShouldContain, "root/Fourth")
			})

			Convey("Then fields hidden by json tag should not be exported at any level", func() {
				So(ns, ShouldNotContain, "root/Third")
				So(ns, ShouldNotContain, "root/First/Tres")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with json export disabled at all levels", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, ExportJsonFieldNames(AlwaysFalse))

			Convey("Then fields at all levels should be exported basing on real field names", func() {
				So(ns, ShouldContain, "root/First/Uno")
				So(ns, ShouldContain, "root/First/Dos/Alpha")
				So(ns, ShouldContain, "root/First/Tres")
				So(ns, ShouldContain, "root/Second")
				So(ns, ShouldContain, "root/Third")
				So(ns, ShouldContain, "root/Fourth")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with json export disabled at second level", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, ExportJsonFieldNames(func(current string, typ reflect.Type) bool {
				return strings.Count(current, "/") == 0
			}))

			Convey("Then fields at first level should be exported basing on json tags", func() {
				So(ns, ShouldContain, "root/second_f")
				So(ns, ShouldNotContain, "root/Third")
				So(ns, ShouldContain, "root/Fourth")
			})

			Convey("Then fields at second level and deeper should be exported basing on actual field names", func() {
				So(ns, ShouldContain, "root/First/Uno")
				So(ns, ShouldContain, "root/First/Dos/Alpha")
				So(ns, ShouldContain, "root/First/Tres")
			})
		})
	})

	//// testing wildcard entries for containers
	Convey("Given struct with some containers", t, func() {
		type dos struct {
			Alpha int
		}
		type first struct {
			Uno bool
			Dos []dos
		}
		m := struct {
			First []first
		}{
			First: []first{
				first{},
				first{Dos: []dos{
					dos{},
					dos{},
				}},
				first{}}}

		Convey("When NamespaceFromCompositeObject is called with default options", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns)

			Convey("Then discovered entries should be added for containers at any level", func() {
				So(ns, ShouldContain, "root/First/0/Uno")
				So(ns, ShouldContain, "root/First/1/Dos/0/Alpha")
			})

			Convey("Then NO wildcard entry should be added for non-empty container at any level", func() {
				So(ns, ShouldNotContain, "root/First/*/Uno")
				So(ns, ShouldNotContain, "root/First/1/Dos/*/Alpha")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with  WildcardEntryInContainer enabled at all levels", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, WildcardEntryInContainer(AlwaysTrue))

			Convey("Then wildcard entries should be added for non-empty container at any level", func() {
				So(ns, ShouldContain, "root/First/*/Uno")
				So(ns, ShouldContain, "root/First/1/Dos/*/Alpha")
			})

			Convey("Then discovered entries should still be present at any level", func() {
				So(ns, ShouldContain, "root/First/0/Uno")
				So(ns, ShouldContain, "root/First/1/Dos/0/Alpha")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with  WildcardEntryInContainer disabled at 4th level", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, WildcardEntryInContainer(func(current string, typ reflect.Type) bool {
				res := strings.Count(current, "/")+1 < 4
				return res
			}))

			Convey("Then wildcard entries should be added only for containers at level <4", func() {
				So(ns, ShouldContain, "root/First/*/Uno")
				So(ns, ShouldNotContain, "root/First/1/Dos/*/Alpha")
			})

			Convey("Then discovered entries should still be present at any level", func() {
				So(ns, ShouldContain, "root/First/0/Uno")
				So(ns, ShouldContain, "root/First/1/Dos/0/Alpha")
			})
		})
	})

	//// checking namespace sanitization
	Convey("Given struct with some namespace-unsafe characters", t, func() {
		type first struct {
			Uno  bool   `json:"uno[f]"`
			Tres string `json:"tres(f)"`
		}
		m := struct {
			First  first
			Second string `json:"second(f)"`
			Third  bool   `json:"third[f]"`
		}{}

		Convey("When NamespaceFromCompositeObject is called with default options", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns)

			Convey("Then all parts of namespace at all levels should be free from unsafe characters", func() {
				So(ns, ShouldContain, "root/second_f")
				So(ns, ShouldContain, "root/third_f")
				So(ns, ShouldContain, "root/First/uno_f")
				So(ns, ShouldContain, "root/First/tres_f")
			})
		})

		Convey("When NamespaceFromCompositeObject is called with  SanitizeNamespaceParts disabled for all levels", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns, SanitizeNamespaceParts(AlwaysFalse))

			Convey("Then all parts of namespace at all levels should retain original contents", func() {
				So(ns, ShouldContain, "root/second(f)")
				So(ns, ShouldContain, "root/third[f]")
				So(ns, ShouldContain, "root/First/uno[f]")
				So(ns, ShouldContain, "root/First/tres(f)")
			})
		})
	})

	//// checking default behavior of all options
	Convey("Given struct with nil pointers and empty maps", t, func() {
		type delta struct {
			Eins int
		}
		type uno struct {
			Alpha bool           `json:"alpha_f"`
			Beta  map[string]int `json:"beta_f"`
			Gamma bool           `json:"-"`
			Delta *delta         `json:"delta_f"`
		}
		type zwei struct {
			Foo bool
		}
		type cuatro struct {
			Kappa int
			Zeta  map[string]zwei
		}
		m := struct {
			First  *uno           `json:"first_f"`
			Second map[string]int `json:"second_f"`
			Third  int            `json:"-"`
			Fourth map[string]cuatro
		}{}

		Convey("When NamespaceFromCompositeObject is called with default options", func() {
			ns := []string{}
			FromCompositeObject(m, "root", &ns)

			Convey("Then nil pointers should be expanded by default at all levels", func() {
				So(ns, ShouldContain, "root/first_f/alpha_f")
				So(ns, ShouldContain, "root/first_f/delta_f/Eins")
			})

			Convey("Then empty containers should be expanded by default at all levels", func() {
				So(ns, ShouldContain, "root/Fourth/*/Kappa")
				So(ns, ShouldContain, "root/Fourth/*/Zeta/*/Foo")
			})

			Convey("Then there should not be entries for containers on their own at any level", func() {
				So(ns, ShouldNotContain, "root/first_f/beta_f")
				So(ns, ShouldNotContain, "root/second_f")
				So(ns, ShouldNotContain, "root/Fourth")
				So(ns, ShouldNotContain, "root/Fourth/*/Zeta")
			})

			Convey("Then namespace should contain json names for fields rather than those from structs", func() {
				So(ns, ShouldContain, "root/first_f/alpha_f")
				So(ns, ShouldContain, "root/first_f/beta_f/*")
				So(ns, ShouldContain, "root/first_f/delta_f/Eins")
				So(ns, ShouldContain, "root/second_f/*")
				So(ns, ShouldNotContain, "root/First/Alpha")
				So(ns, ShouldNotContain, "root/First/Beta/*")
				So(ns, ShouldNotContain, "root/First/Delta/Eins")
			})

			Convey("Then namespace should not contain entries hidden by json tags", func() {
				So(ns, ShouldNotContain, "root/first_f/Gamma")
				So(ns, ShouldNotContain, "root/Third")
			})
		})
	})

	//// check the correctness
	Convey("Given nested map of maps", t, func() {
		m := map[string]interface{}{
			"First":  "first",
			"Second": 13,
			"Third": map[string]interface{}{
				"Uno": false,
				"Dos": map[string]interface{}{
					"Alpha": "alpha"}}}

		Convey("When NamespaceFromCompositeObject is called with  current being root", func() {
			ns := []string{}
			current := "root"
			FromCompositeObject(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 4)
			})

			Convey("Then namespace should contain correct entries for all levels", func() {
				So(ns, ShouldContain, "root/First")
				So(ns, ShouldContain, "root/Second")
				So(ns, ShouldContain, "root/Third/Uno")
				So(ns, ShouldContain, "root/Third/Dos/Alpha")
			})
		})
	})
	Convey("Given map of structs", t, func() {
		m := map[string]interface{}{
			"First": struct {
				Uno string
			}{},
			"Second": struct {
				Tres struct {
					Alpha string
				}
			}{}}

		Convey("When NamespaceFromCompositeObject is called with  current being root", func() {
			ns := []string{}
			current := "root"
			FromCompositeObject(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 2)
			})

			Convey("Then namespace should contain entries for nested structs", func() {
				So(ns, ShouldContain, "root/First/Uno")
				So(ns, ShouldContain, "root/Second/Tres/Alpha")
			})
		})
	})

	Convey("Given map with pointer to struct", t, func() {
		m := map[string]interface{}{
			"Third": &(struct {
				Dos struct {
					Alpha int
				}
			}{})}

		Convey("When NamespaceFromCompositeObject is called with  current being root", func() {
			ns := []string{}
			current := "root"
			FromCompositeObject(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 1)
			})

			Convey("Then namespace should contain entries for nested structs", func() {
				So(ns, ShouldContain, "root/Third/Dos/Alpha")
			})
		})
	})

	Convey("Given typed empty map", t, func() {
		m := map[string]struct{ Uno int }{}

		Convey("When NamespaceFromCompositeObject is called with  current being root", func() {
			ns := []string{}
			current := "root"
			FromCompositeObject(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 1)
			})

			Convey("Then namespace should contain entry for potential nested struct", func() {
				So(ns, ShouldContain, "root/*/Uno")
			})
		})
	})

	Convey("Given typed map with nil pointer", t, func() {
		type first struct {
			Uno int
		}
		m := map[string]*first{"First": nil}

		Convey("When NamespaceFromCompositeObject is called with  current being root", func() {
			ns := []string{}
			current := "root"
			//NamespaceFromCompositeObject(m, current, &ns, EntryForContainersRoot(AlwaysTrue))
			FromCompositeObject(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 1)
			})

			Convey("Then namespace should contain expanded entry for nil pointer's declared type", func() {
				So(ns, ShouldContain, "root/First/Uno")
			})
		})
	})

	Convey("Given struct of maps", t, func() {
		type uno struct {
			Alpha string
		}
		type dos struct {
			Beta int
		}
		m := struct {
			First  map[string]*uno
			Second map[string]dos
		}{}

		Convey("When NamespaceFromCompositeObject is called with  current being root", func() {
			ns := []string{}
			current := "root"
			FromCompositeObject(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 2)
			})

			Convey("Then namespace should contain expanded entries for potential entries in both maps", func() {
				So(ns, ShouldContain, "root/First/*/Alpha")
				So(ns, ShouldContain, "root/Second/*/Beta")
			})
		})
	})

	Convey("Given struct of pointers", t, func() {
		type uno struct {
			Alpha bool
		}
		m := struct {
			First  ***int
			Second **string
			Third  ***uno
		}{}

		Convey("When NamespaceFromCompositeObject is called with  current being root", func() {
			ns := []string{}
			current := "root"
			FromCompositeObject(m, current, &ns)

			Convey("Then namespace should be populated", func() {
				So(len(ns), ShouldEqual, 3)
			})

			Convey("Then namespace should contain expanded entries for pointed types", func() {
				So(ns, ShouldContain, "root/First")
				So(ns, ShouldContain, "root/Second")
				So(ns, ShouldContain, "root/Third/Alpha")
			})
		})
	})

	Convey("Given struct of non-nil and nil pointers", t, func() {
		type first struct {
			Uno int
		}
		type dos struct {
			Alpha int
		}
		type second struct {
			Dos *dos
		}
		makeRefRef_first := func() **first {
			ref_first := &first{}
			return &ref_first
		}
		makeRefRefRef_second := func() ***second {
			ref_second := &second{}
			refref_second := &ref_second
			return &refref_second
		}
		m := struct {
			First  **first
			Second ***second
		}{
			First:  makeRefRef_first(),
			Second: makeRefRefRef_second()}

		Convey("When NamespaceFromCompositeObject is called with  current being root", func() {
			ns := []string{}
			current := "root"
			FromCompositeObject(m, current, &ns)

			Convey("Then namespace should be populated with correct entries", func() {
				So(ns, ShouldContain, "root/First/Uno")
				So(ns, ShouldContain, "root/Second/Dos/Alpha")
			})
		})
	})

	Convey("Given struct with some pointers and empty maps and empty slices", t, func() {
		type first struct {
			Uno bool `json:"uno_f"`
		}
		type fourth struct {
			Cuatro bool
		}
		m := struct {
			First  *first         `json:"first_f"`
			Second map[string]int `json:"second_f"`
			Third  int            `json:"-"`
			Fourth []*fourth
		}{
			First:  nil,
			Fourth: []*fourth{&fourth{}, &fourth{}}}

		ns := []string{}
		FromCompositeObject(m, "root", &ns, WildcardEntryInContainer(AlwaysTrue))
		Convey("When NamespaceFromCompositeObject is called with default options plus wildcard entry", func() {

			Convey("Then namespace should contain entries for all potential exported fields", func() {
				So(ns, ShouldContain, "root/first_f/uno_f")
				So(ns, ShouldContain, "root/second_f/*")
				So(ns, ShouldContain, "root/Fourth/*/Cuatro")
				So(ns, ShouldContain, "root/Fourth/0/Cuatro")
				So(ns, ShouldContain, "root/Fourth/1/Cuatro")
			})
		})
	})
}

func TestGetValueByNamespace(t *testing.T) {
	Convey("When GetValueByNamespace is given some nested structs", t, func() {
		dataTwo := new(float64)
		nestedDataTwo := new(float64)
		nestedDataThree := new(float32)
		*dataTwo = 3.14
		*nestedDataTwo = 99.99
		*nestedDataThree = 20.0

		type bar struct {
			NestedDataOne   uint32   `json:"nested_data_one,omitempty"`
			NestedDataTwo   *float64 `json:"nested_data_two,omitempty"`
			NestedDataThree *float32 `json:"nested_data_three,omitempty"`
		}
		type foo struct {
			DataOne   uint32   `json:"data_one,omitempty"`
			DataTwo   *float64 `json:"data_two,omitempty"`
			DataThree *bar     `json:"data_three,omitempty"`
		}
		m := struct {
			ID      string      `json:"id"`
			Invalid interface{} `json:"invalid"`
			Data    *foo        `json:"data"`
		}{
			ID:      "FooID-12345",
			Invalid: new(interface{}),
			Data: &foo{
				DataOne: 127,
				DataTwo: dataTwo,
				DataThree: &bar{
					NestedDataOne:   254,
					NestedDataTwo:   nestedDataTwo,
					NestedDataThree: nestedDataThree,
				},
			},
		}

		Convey("Should return values contained in the struct by a given namespace", func() {
			namespaces := [][]string{
				[]string{"data_one"},
				[]string{"data_two"},
				[]string{"data_three", "nested_data_one"},
				[]string{"data_three", "nested_data_two"},
				[]string{"data_three", "nested_data_three"},
			}

			So(GetValueByNamespace(m.Data, namespaces[0]), ShouldEqual, 127)
			So(GetValueByNamespace(m.Data, namespaces[1]), ShouldEqual, 3.14)
			So(GetValueByNamespace(m.Data, namespaces[2]), ShouldEqual, 254)
			So(GetValueByNamespace(m.Data, namespaces[3]), ShouldEqual, 99.99)
			So(GetValueByNamespace(m.Data, namespaces[4]), ShouldEqual, 20.0)
		})

		Convey("Should not attempt to interface a zero value", func() {
			So(func() { GetValueByNamespace(m, []string{"invalid"}) }, ShouldNotPanic)
		})
	})
}
