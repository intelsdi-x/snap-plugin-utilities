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

package config

import (
	"fmt"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/ctypes"
)

// getConfigItemValue returns value of configuration item
func getConfigItemValue(item ctypes.ConfigValue) (interface{}, error) {

	var value interface{}

	switch item.Type() {
	case "string":
		value = item.(ctypes.ConfigValueStr).Value
		break

	case "float":
		value = item.(ctypes.ConfigValueFloat).Value
		break

	case "integer":
		value = item.(ctypes.ConfigValueInt).Value
		break

	case "bool":
		value = item.(ctypes.ConfigValueBool).Value
		break

	default:
		return nil, fmt.Errorf("Unsupported type of configuration item, type=%+v", item.Type())
	}

	return value, nil
}

// GetGlobalConfigItem returns value of config item specified by `name` defined in Plugin Global Config
// Notes: GetGlobalConfigItem() will be helpful to access and get configuration item's value in GetMetricTypes()
func GetGlobalConfigItem(cfg plugin.ConfigType, name string) (interface{}, error) {

	if item, ok := cfg.Table()[name]; ok {
		return getConfigItemValue(item)
	}

	return nil, fmt.Errorf("Cannot find %+v in Global Config", name)
}

// GetGlobalConfigItems returns map to values of multiple configuration items defined in Plugin Global Config and specified in 'names' slice
// Notes: GetGlobalConfigItems() will be helpful to access and get multiple configuration items' values in GetMetricTypes()
func GetGlobalConfigItems(cfg plugin.ConfigType, names []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, name := range names {
		item, ok := cfg.Table()[name]
		if !ok {
			return nil, fmt.Errorf("Cannot find %+v in Global Config", name)
		}

		val, err := getConfigItemValue(item)
		if err != nil {
			return nil, err
		}
		result[name] = val
	}

	return result, nil
}

// GetMetricConfigItem returns value of configuration item specified by `name` defined in Metrics Config
// Notes: GetMetricConfigItem() will be helpful to access and get configuration item's value in CollectMetrics()
// (Plugin Global Config is merged into Metric Config)
func GetMetricConfigItem(metric plugin.MetricType, name string) (interface{}, error) {

	if item, ok := metric.Config().Table()[name]; ok {
		return getConfigItemValue(item)
	}

	return nil, fmt.Errorf("Cannot find %+v in Metrics Config", name)
}

// GetMetricConfigItems returns map to values of multiple configuration items defined in Metric Config and specified in 'names' slice
// Notes: GetMetricConfigItems() will be helpful to access and get multiple configuration items' values in CollectMetrics()
// (Plugin Global Config is merged into Metric Config)
func GetMetricConfigItems(metric plugin.MetricType, names []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, name := range names {
		item, ok := metric.Config().Table()[name]
		if !ok {
			return nil, fmt.Errorf("Cannot find %+v in Metrics Config", name)
		}

		val, err := getConfigItemValue(item)
		if err != nil {
			return nil, err
		}
		result[name] = val
	}

	return result, nil
}

// GetConfigItem returns value of configuration item specified by `name` defined in Global Config or Metrics Config
// Notes: GetConfigItem() will be helpful to access and get configuration item'a value, both in GetMetricTypes() and CollectMetrics()
func GetConfigItem(config interface{}, name string) (interface{}, error) {

	switch config.(type) {
	case plugin.ConfigType:
		return GetGlobalConfigItem(config.(plugin.ConfigType), name)

	case plugin.MetricType:
		return GetMetricConfigItem(config.(plugin.MetricType), name)
	}

	return nil, fmt.Errorf("Unsupported type of config. Input 'config' needs to be PluginConfigType or PluginMetricType")
}

// GetConfigItems returns map to values of multiple  configuration items defined in Global Config or Metrics Config
// Notes: GetConfigItems() will be helpful to access and and get multiple configuration items' values, both in GetMetricTypes() and CollectMetrics()
func GetConfigItems(config interface{}, names ...string) (map[string]interface{}, error) {

	switch config.(type) {
	case plugin.ConfigType:
		return GetGlobalConfigItems(config.(plugin.ConfigType), names)

	case plugin.MetricType:
		return GetMetricConfigItems(config.(plugin.MetricType), names)
	}

	return nil, fmt.Errorf("Unsupported type of config. Input 'config' needs to be PluginConfigType or PluginMetricType")
}
