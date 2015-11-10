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

package pipeline

import (
	"strings"
)

type Pipe interface {
	Process(in chan interface{}) chan interface{}
}

type Pipeline struct {
	head chan interface{}
	tail chan interface{}
}

func (p *Pipeline) Enqueue(item interface{}) {
	p.head <- item
}

func (p *Pipeline) Dequeue(handler func(interface{})) {
	for i := range p.tail {
		handler(i)
	}
}

func (p *Pipeline) Output() chan interface{} {
	return p.tail
}

func (p *Pipeline) Close() {
	close(p.head)
}

func NewPipeline(pipes ...Pipe) Pipeline {
	head := make(chan interface{})
	var next_chan chan interface{}
	for _, pipe := range pipes {
		if next_chan == nil {
			next_chan = pipe.Process(head)
		} else {
			next_chan = pipe.Process(next_chan)
		}
	}
	return Pipeline{head: head, tail: next_chan}
}

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

type StringContains struct {
	str string
}

func (t StringContains) Process(in chan interface{}) chan interface{} {
	out := make(chan interface{})
	go func() {
		for i := range in {
			if !strings.Contains(i.(string), t.str) {
				out <- i.(string)
			}
		}
		close(out)
	}()
	return out
}
