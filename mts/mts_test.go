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

package mts

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
)

func TestConvertDynamicElements(t *testing.T) {
	Convey("Given list of metrics", t, func() {
		metrics := []plugin.MetricType{
			plugin.MetricType{
				Namespace_: core.NewNamespace("intel", "mock", "foo"),
				Tags_:      map[string]string{},
			},
			plugin.MetricType{
				Namespace_: core.NewNamespace("intel", "mock").AddDynamicElement("foos", "foos contains many foo").AddStaticElement("bar"),
				Tags_:      map[string]string{},
			},
			plugin.MetricType{
				Namespace_: core.NewNamespace("intel", "mock").AddDynamicElement("foos", "foos contains many foo").AddStaticElement("bar"),
				Tags_:      map[string]string{core.STD_TAG_PLUGIN_RUNNING_ON: "127.0.0.1"},
			},
			plugin.MetricType{
				Namespace_: core.NewNamespace("intel", "mock").AddDynamicElement("foos", "foos contains many foo").AddStaticElements("bar", "tar"),
				Tags_:      map[string]string{},
			},
			plugin.MetricType{
				Namespace_: core.NewNamespace("intel").AddDynamicElement("foos", "foos contains many foo").AddStaticElements("bar", "tar"),
				Tags_:      map[string]string{},
			},
			plugin.MetricType{
				Namespace_: core.NewNamespace("intel").AddDynamicElement("foos", "foos contains many foo").AddDynamicElement("bars", "bars contains many bar").AddStaticElement("tar"),
				Tags_:      map[string]string{},
			},
		}

		Convey("When conversion of dynamic namespaces is triggered", func() {
			mockCollectorActionOnDynamicMetric(metrics)
			ConvertDynamicElements(metrics)

			Convey("Then metric namespaces are converted to static and proper tags are added", func() {
				So(metrics[0].Namespace().String(), ShouldEqual, core.NewNamespace("intel", "mock", "foo").String())
				So(metrics[0].Tags_, ShouldResemble, map[string]string{})
				So(metrics[1].Namespace().String(), ShouldEqual, core.NewNamespace("intel", "mock", "bar").String())
				So(metrics[1].Tags_, ShouldResemble, map[string]string{"foos": "value_1_2"})
				So(metrics[2].Namespace().String(), ShouldEqual, core.NewNamespace("intel", "mock", "bar").String())
				So(metrics[2].Tags_, ShouldResemble, map[string]string{"foos": "value_2_2", core.STD_TAG_PLUGIN_RUNNING_ON: "127.0.0.1"})
				So(metrics[3].Namespace().String(), ShouldEqual, core.NewNamespace("intel", "mock", "bar", "tar").String())
				So(metrics[3].Tags_, ShouldResemble, map[string]string{"foos": "value_3_2"})
				So(metrics[4].Namespace().String(), ShouldEqual, core.NewNamespace("intel", "bar", "tar").String())
				So(metrics[4].Tags_, ShouldResemble, map[string]string{"foos": "value_4_1"})
				So(metrics[5].Namespace().String(), ShouldEqual, core.NewNamespace("intel", "tar").String())
				So(metrics[5].Tags_, ShouldResemble, map[string]string{"foos": "value_5_1", "bars": "value_5_2"})
			})
		})
	})
}

func mockCollectorActionOnDynamicMetric(metrics []plugin.MetricType) {
	for i, metric := range metrics {
		if isDynamic, indexes := metric.Namespace().IsDynamic(); isDynamic {
			for _, index := range indexes {
				metrics[i].Namespace()[index].Value = fmt.Sprintf("value_%d_%d", i, index)
			}
		}
	}
}
