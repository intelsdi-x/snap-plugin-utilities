<!--
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
-->
[![Build Status](https://travis-ci.com/intelsdi-x/snap-plugin-utilities.svg?token=HoxHq3yqBGpySzRd5XUm&branch=master)](https://travis-ci.com/intelsdi-x/snap-plugin-utilities)
# snap Plugin Utilities

Set of packages that provide tools for plugin development

It's used in the [snap framework](http://github.com/intelsdi-x/snap).

1. [Documentation](#documentation)
  * [Features](#features)
  * [Examples](#examples)
2. [Community Support](#community-support)
3. [Contributing](#contributing)
4. [License](#license-and-authors)
5. [Acknowledgements](#acknowledgements)

## Documentation

### Features:

  * [config](#config-package)
  * [logger](#logger-package)
  * [ns](#ns-package)
  * [pipeline](#pipeline-package)
  * [source](#source-package)
  * [stack](#stack-package)

### Examples

[config] package
-------------------------------------------------------------------------------------------

The `config` package provides helpful methods to retrive global config items.

See it in action:

```go
	if interestingValue, err := GetGlobalConfigItem(cfg, "something_interesting"); err != nil {
		DoSomething(interestingValue)
	}

}
```

```go
	interestingItems := []string{"dummy_string", "dummy_bool", "dummy_int", "dummy_float"}

	if interestingValues, err := GetGlobalConfigItems(cfg, interestingItems); err != nil {
		for _, interestingValue := range interestingValues {
			DoSomething(interestingValue)
		}
	}

```

[logger] package
---------------------------------------------------------------------------------------------

The `logger` package wraps logrus package (https://github.com/Sirupsen/logrus).
It sets logging from plugin to separate file. It adds caller function name to each message.

```go
import (
	. "github.com/intelsdi-x/snap-plugin-utilities/logger"
)

func main() {
	
	LogDebug("Some useful information", "varibale", value)
	LogInfo("Some information worth noting", interesting, thing, done)
	LogWarn("Take a look on that", warning)
	LogError("This is bad", "error", err)
	// Exit after
	LogFatal("This is really bad", "error", err, value)
	// Call panic()
	LogPanic("Show me the stacks!")
}

```


[ns]  package
---------------------------------------------------------------------------------------
The `ns` package provides functions to extract namespace from maps, JSON and struct compositions.
It is useful for situations when full knowledge of available metrics is not known at time when GetMetricTypes() is called.

NamespaceFromMap example usage:
```go
	
	Baz := map[string]interface{}{"Bazo": "bazo", "Fazo": "fazo", "Mazo": "mazo"}
	Foo := map[string]interface{}{"Foos": "foos", "Boos": "boos"}
	Bar := []map[string]interface{}{Baz, Baz}
	
	m := map[string]interface{}{
		"Foo": Foo,
		"Bar": Bar,
	}
	
	ns := []string{}
	NamespaceFromMap(m, "root", &ns)

	/*
	ns contains:
	"root/Foo/Foos"
	"root/Foo/Boos"
	"root/Bar/0/Bazo"
	"root/Bar/0/Fazo"
	"root/Bar/0/Mazo"
	"root/Bar/1/Bazo"
	"root/Bar/1/Fazo"
	"root/Bar/1/Mazo"
	*/

```

NamespaceFromJSON example usage:
```go
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
	ns := []string{}
	
	NamespaceFromJSON(&data, "root", &ns)

	/*
	ns contains:
	"root/bar/qaz/0"
	"root/bar/qaz/1"
	"root/bar/faz"
	"root/baz"
	*/
```

NamespaceFromComposition example usage:
```go
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

	ns := []string{}
	NamespaceFromComposition(Foo, "root", &ns)

	/*
	ns contains:
	"root/Bar/Qaz/0"
	"root/Bar/Qaz/1"
	"root/Bar/Faz"
	"root/Baz"
	*/
```
NamespaceFromCompositionTags example usage:
```go
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

	ns := []string{}
	NamespaceFromCompositionTags(Foo, "root", &ns)

	/*
	ns contains:
	"root/bar/qaz/0"
	"root/bar/qaz/1"
	"root/bar/faz"
	"root/baz"
	*/
```

[pipeline] package
----------------------------------------------------------------------------------------
Creates array of Pipes connected by channels. Each Pipe can do single processing on data transmitted by channels

Pipe is interface which consists of single method Process, which takes input channel and returns channel.
It can be used to implement transformations on incoming data

```go
	type DoNothing struct{}

	func (dn DoNothing) Process(in chan interface{}) chan interface{} {
		out := make(chan interface{})
		go func() {
			for i := range in {
				out <- i
			}
			close(out)
		}()
		return out
	}	
	
	type FilterString struct{}

	func (dn FilterString) Process(in chan interface{}) chan interface{} {
		out := make(chan interface{})
		go func() {
			for i := range in {
				if !strings.Contains(i.(string), "foo") {
					out <- i
				}
			}
			close(out)
		}()
		return out
	}
	
	pipeline := NewPipeline(DoNothing{}, FilterString{}) 

```

[source] package
-----------------------------------------------------------------------------------------
The `source` package provides handy way of dealing with external command output. 
It can be used for continous command execution (EMON or PCM like), or for single command calls.


```go
	ech := make(chan error)
	out := make(chan interface{})
	s := Source{"du", []string{"/some/path", "-h", "--max-depth=10"}}
	go s.Generate(out, ech)

	LOOP:
	for {
		select {
		case data := <-out:
			fmt.Printf(">>> Recieving {%v}\n", data)
		case e := <-ech:
			fmt.Printf("ERRROR {%v}\n", e)
			break LOOP	
		case <-time.After(time.Second * 2):
			fmt.Printf("No activity\n")
			break LOOP 
		}
	}
```

[stack] package
-----------------------------------------------------------------------------------------
The `stack` package provides simple implementation of stack.

```go
	stack := new(Stack)
	
	for i := 0; i < 5; i++ {
		stack.Push(i)
	}
	
	for !stack.Empty() {
		fmt.Printf("%v\n", stack.Pop())
	}
	
	/*
	prints:
	5
	4
	3
	2
	1
	*/
```

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions! 

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Izabella Raulin](https://github.com/IzabellaRaulin)
* Author: [Marcin Krolik](https://github.com/marcin-krolik)

**Thank you!** Your contribution is incredibly important to us.
