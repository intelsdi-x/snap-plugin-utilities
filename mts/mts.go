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

package mts

import (
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
)

// ConvertDynamicElements removes dynamic `NamespaceElements` from `core.Namespace`
// and turn them into tags. Useful when metrics with unwanted namespace elements are gathered
// like `/intel/mock/host1/foo`, `/intel/mock/host2/foo` ... `/intel/mock/hostN/foo`
// where `hostN` is dynamic `NamespaceElement` (has Name) and one wants to publish them
// as `/intel/mock/foo` with additional tags as like host_name: hostN
func ConvertDynamicElements(metrics []plugin.MetricType) {
	for j, metric := range metrics {
		if isDynamic, indexes := metric.Namespace().IsDynamic(); isDynamic {
			static := core.Namespace{}
			tags := metric.Tags()
			if tags == nil {
				tags = map[string]string{}
			}
			for i, nse := range metric.Namespace() {
				if contains(indexes, i) {
					tags[nse.Name] = nse.Value
				} else {
					static = append(static, nse)
				}
			}
			metrics[j].Namespace_ = static
			metrics[j].Tags_ = tags
		}
	}
}

func contains(slice []int, lookup int) bool {
	for _, element := range slice {
		if element == lookup {
			return true
		}
	}
	return false
}