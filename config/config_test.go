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
	"testing"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/core/cdata"
	"github.com/intelsdi-x/pulse/core/ctypes"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	// dummy config value
	dummy_str   = "asdf"
	dummy_bool  = true
	dummy_int   = 1
	dummy_float = 1.11
)

func TestGetGlobalConfigItem(t *testing.T) {

	Convey("Get a value of item from Global Config with no error", t, func() {
		//create global config and add item (different types)
		cfg := plugin.NewPluginConfigType()
		cfg.AddItem("dummy_string", ctypes.ConfigValueStr{Value: dummy_str})
		cfg.AddItem("dummy_bool", ctypes.ConfigValueBool{Value: dummy_bool})
		cfg.AddItem("dummy_int", ctypes.ConfigValueInt{Value: dummy_int})
		cfg.AddItem("dummy_float", ctypes.ConfigValueFloat{Value: dummy_float})

		Convey("string type of item", func() {
			result, err := GetGlobalConfigItem(cfg, "dummy_string")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, dummy_str)
		})

		Convey("bool type of item", func() {
			result, err := GetGlobalConfigItem(cfg, "dummy_bool")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, dummy_bool)
		})

		Convey("int type of item", func() {
			result, err := GetGlobalConfigItem(cfg, "dummy_int")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, dummy_int)
		})

		Convey("float type of item", func() {
			result, err := GetGlobalConfigItem(cfg, "dummy_float")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, dummy_float)
		})
	})

	Convey("Try to get value of item not defined in Global Config", t, func() {
		cfg := plugin.NewPluginConfigType()
		cfg.AddItem("foo", ctypes.ConfigValueStr{Value: "foo"})
		result, err := GetGlobalConfigItem(cfg, "foo_not_exist")
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("No item defined in Global Config", t, func() {
		cfg_empty := plugin.NewPluginConfigType()
		result, err := GetGlobalConfigItem(cfg_empty, "foo")
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})
}

func TestGetGlobalConfigItems(t *testing.T) {

	Convey("Get values of items from Global Config with no error", t, func() {
		//create global config and add item (different types)
		cfg := plugin.NewPluginConfigType()
		cfg.AddItem("dummy_string", ctypes.ConfigValueStr{Value: dummy_str})
		cfg.AddItem("dummy_bool", ctypes.ConfigValueBool{Value: dummy_bool})
		cfg.AddItem("dummy_int", ctypes.ConfigValueInt{Value: dummy_int})
		cfg.AddItem("dummy_float", ctypes.ConfigValueFloat{Value: dummy_float})

		names := []string{"dummy_string", "dummy_bool", "dummy_int", "dummy_float"}

		result, err := GetGlobalConfigItems(cfg, names)
		So(err, ShouldBeNil)
		for _, name := range names {
			So(result[name], ShouldNotBeEmpty)
		}
	})

	Convey("Try to get values of items not defined in Global Config", t, func() {
		cfg := plugin.NewPluginConfigType()
		cfg.AddItem("foo", ctypes.ConfigValueStr{Value: "foo"})
		names := []string{"foo1", "foo2"}
		result, err := GetGlobalConfigItems(cfg, names)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("No item defined in Global Config", t, func() {
		cfg_empty := plugin.NewPluginConfigType()
		names := []string{"foo", "bar"}
		result, err := GetGlobalConfigItems(cfg_empty, names)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})
}

func TestGetMetricConfigItem(t *testing.T) {

	Convey("Get a value of items from Metrics Config with no error", t, func() {
		// create config
		config := cdata.NewNode()
		config.AddItem("dummy_string", ctypes.ConfigValueStr{Value: dummy_str})
		config.AddItem("dummy_bool", ctypes.ConfigValueBool{Value: dummy_bool})
		config.AddItem("dummy_int", ctypes.ConfigValueInt{Value: dummy_int})
		config.AddItem("dummy_float", ctypes.ConfigValueFloat{Value: dummy_float})

		// create metric and set config
		metric := plugin.PluginMetricType{}
		metric.Config_ = config

		Convey("string type of item", func() {
			result, err := GetMetricConfigItem(metric, "dummy_string")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, dummy_str)
		})

		Convey("bool type of item", func() {
			result, err := GetMetricConfigItem(metric, "dummy_bool")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, dummy_bool)
		})

		Convey("int type of item", func() {
			result, err := GetMetricConfigItem(metric, "dummy_int")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, dummy_int)
		})

		Convey("float type of item", func() {
			result, err := GetMetricConfigItem(metric, "dummy_float")
			So(err, ShouldBeNil)
			So(result, ShouldEqual, dummy_float)
		})
	})

	Convey("Try to get a value of items not defined in Metrics Config", t, func() {
		config := cdata.NewNode()
		config.AddItem("foo", ctypes.ConfigValueStr{Value: "foo_val"})
		metric := plugin.PluginMetricType{}
		metric.Config_ = config

		result, err := GetMetricConfigItem(metric, "foo_not_exist")
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("No item defined in Metrics Config", t, func() {
		metric := plugin.PluginMetricType{}
		metric.Config_ = cdata.NewNode()

		result, err := GetMetricConfigItem(metric, "foo")
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})
}

func TestGetMetricConfigItems(t *testing.T) {

	Convey("Get values of items from Metrics Config with no error", t, func() {

		// create config
		config := cdata.NewNode()
		config.AddItem("dummy_string", ctypes.ConfigValueStr{Value: dummy_str})
		config.AddItem("dummy_bool", ctypes.ConfigValueBool{Value: dummy_bool})
		config.AddItem("dummy_int", ctypes.ConfigValueInt{Value: dummy_int})
		config.AddItem("dummy_float", ctypes.ConfigValueFloat{Value: dummy_float})

		// create metric and set config
		metric := plugin.PluginMetricType{}
		metric.Config_ = config

		names := []string{"dummy_string", "dummy_bool", "dummy_int", "dummy_float"}

		result, err := GetMetricConfigItems(metric, names)
		So(err, ShouldBeNil)
		for _, name := range names {
			So(result[name], ShouldNotBeEmpty)
		}
	})

	Convey("Try to get values of items not defined in Metrics config", t, func() {
		config := cdata.NewNode()
		config.AddItem("foo", ctypes.ConfigValueStr{Value: "foo_val"})
		metric := plugin.PluginMetricType{}
		metric.Config_ = config

		names := []string{"foo1", "foo2"}
		result, err := GetMetricConfigItems(metric, names)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

	Convey("No item defined in Metrics Config", t, func() {
		metric := plugin.PluginMetricType{}
		metric.Config_ = cdata.NewNode()
		names := []string{"foo", "bar"}
		result, err := GetMetricConfigItems(metric, names)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})
}

func TestConfigItem(t *testing.T) {

	Convey("Get a value of item with no error", t, func() {

		Convey("Source: global config", func() {
			//create global config and add item (different types)
			cfg := plugin.NewPluginConfigType()
			cfg.AddItem("dummy_string", ctypes.ConfigValueStr{Value: dummy_str})
			cfg.AddItem("dummy_bool", ctypes.ConfigValueBool{Value: dummy_bool})
			cfg.AddItem("dummy_int", ctypes.ConfigValueInt{Value: dummy_int})
			cfg.AddItem("dummy_float", ctypes.ConfigValueFloat{Value: dummy_float})

			Convey("string type of item", func() {
				result, err := GetConfigItem(cfg, "dummy_string")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, dummy_str)
			})

			Convey("bool type of item", func() {
				result, err := GetConfigItem(cfg, "dummy_bool")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, dummy_bool)
			})

			Convey("int type of item", func() {
				result, err := GetConfigItem(cfg, "dummy_int")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, dummy_int)
			})

			Convey("float type of itemr", func() {
				result, err := GetConfigItem(cfg, "dummy_float")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, dummy_float)
			})
		})

		Convey("Source: metrics config", func() {
			// create metric's config
			config := cdata.NewNode()
			config.AddItem("dummy_string", ctypes.ConfigValueStr{Value: dummy_str})
			config.AddItem("dummy_bool", ctypes.ConfigValueBool{Value: dummy_bool})
			config.AddItem("dummy_int", ctypes.ConfigValueInt{Value: dummy_int})
			config.AddItem("dummy_float", ctypes.ConfigValueFloat{Value: dummy_float})

			// create metric and set config
			metric := plugin.PluginMetricType{}
			metric.Config_ = config

			Convey("string type of item", func() {
				result, err := GetConfigItem(metric, "dummy_string")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, dummy_str)
			})

			Convey("bool type of item", func() {
				result, err := GetConfigItem(metric, "dummy_bool")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, dummy_bool)
			})

			Convey("int type of item", func() {
				result, err := GetConfigItem(metric, "dummy_int")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, dummy_int)
			})

			Convey("float type of item", func() {
				result, err := GetConfigItem(metric, "dummy_float")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, dummy_float)
			})
		})
	})

	Convey("Try to get a value of item not defined in config", t, func() {

		Convey("Source: metrics config", func() {
			config := cdata.NewNode()
			config.AddItem("foo", ctypes.ConfigValueStr{Value: "foo_val"})
			metric := plugin.PluginMetricType{}
			metric.Config_ = config
			result, err := GetMetricConfigItem(metric, "foo_not_exist")
			So(err, ShouldNotBeNil)
			So(result, ShouldBeNil)
		})

		Convey("Source: global config", func() {
			cfg := plugin.NewPluginConfigType()
			cfg.AddItem("foo", ctypes.ConfigValueStr{Value: "foo"})
			names := []string{"foo1", "foo2"}
			result, err := GetGlobalConfigItems(cfg, names)
			So(err, ShouldNotBeNil)
			So(result, ShouldBeNil)
		})
	})

	Convey("No item defined in config", t, func() {

		Convey("Source: metrics config", func() {
			metric := plugin.PluginMetricType{}
			metric.Config_ = cdata.NewNode()
			result, err := GetMetricConfigItem(metric, "foo")
			So(err, ShouldNotBeNil)
			So(result, ShouldBeNil)
		})

		Convey("Source: global config", func() {
			cfg_empty := plugin.NewPluginConfigType()
			names := []string{"foo", "bar"}
			result, err := GetGlobalConfigItems(cfg_empty, names)
			So(err, ShouldNotBeNil)
			So(result, ShouldBeNil)
		})
	})

	Convey("Try to get a value of item from invalid config (unsupported type)", t, func() {
		invalid_cfg := []string{"invalid", "config", "source"}
		result, err := GetConfigItem(invalid_cfg, "foo")
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})
}

func TestConfigItems(t *testing.T) {
	names := []string{"dummy_string", "dummy_bool", "dummy_int", "dummy_float"}

	Convey("Get values of configuration items with no error", t, func() {

		Convey("Source: global config", func() {
			//create global config and add item (different types)
			cfg := plugin.NewPluginConfigType()
			cfg.AddItem("dummy_string", ctypes.ConfigValueStr{Value: dummy_str})
			cfg.AddItem("dummy_bool", ctypes.ConfigValueBool{Value: dummy_bool})
			cfg.AddItem("dummy_int", ctypes.ConfigValueInt{Value: dummy_int})
			cfg.AddItem("dummy_float", ctypes.ConfigValueFloat{Value: dummy_float})

			result, err := GetConfigItems(cfg, names)
			So(err, ShouldBeNil)
			for _, name := range names {
				So(result[name], ShouldNotBeEmpty)
			}
		})

		Convey("Source: metrics config", func() {
			// create metrics config
			config := cdata.NewNode()
			config.AddItem("dummy_string", ctypes.ConfigValueStr{Value: dummy_str})
			config.AddItem("dummy_bool", ctypes.ConfigValueBool{Value: dummy_bool})
			config.AddItem("dummy_int", ctypes.ConfigValueInt{Value: dummy_int})
			config.AddItem("dummy_float", ctypes.ConfigValueFloat{Value: dummy_float})

			// create metric and set config
			metric := plugin.PluginMetricType{}
			metric.Config_ = config

			result, err := GetConfigItems(metric, names)
			So(err, ShouldBeNil)
			for _, name := range names {
				So(result[name], ShouldNotBeEmpty)
			}
		})
	})

	Convey("Try to get values of items from invalid config (unsupported type)", t, func() {
		invalid_cfg := []string{"invalid", "config", "source"}
		result, err := GetConfigItems(invalid_cfg, names)
		So(err, ShouldNotBeNil)
		So(result, ShouldBeNil)
	})

}
