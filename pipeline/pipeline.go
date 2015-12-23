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

// Pipeline processing interface
type Processor interface {
	Run(input, output Pipe)
}

// Pipeline basic block element
type Pipe chan interface{}

// Next sets up following pipeline chain element
// It returns last Pipe in Pipeline
func (p Pipe) Next(proc Processor) Pipe {
	outPipe := make(Pipe)
	go proc.Run(p, outPipe)
	return outPipe
}

func (p Pipe) Destination(pipes ...Pipe) {
	go func() {
		for input := range p {
			for _, output := range pipes {
				output <- input
			}
		}

		for _, output := range pipes {
			close(output)
		}
	}()
}

func (p Pipe) Clone(n int) []Pipe {
	newPipes := make([]Pipe, n)
	for i, _ := range newPipes {
		newPipes[i] = make(Pipe)
	}

	p.Destination(newPipes...)

	return newPipes
}

func Pipeline(input Pipe, processors ...Processor) Pipe {
	lastOutput := input
	for _, proc := range processors {
		lastOutput = lastOutput.Next(proc)
	}
	return lastOutput
}

type WaitForCompletion struct {
}

func (w WaitForCompletion) Run(input, output Pipe) {
	for v := range input {
		_ = v
	}

	output <- true

	close(output)
}

type Filter struct {
	FilterFunc func(item interface{}) bool
}

func (self Filter) Run(input, output Pipe) {
	for v := range input {
		if self.FilterFunc(v) {
			output <- v
		}
	}

	close(output)
}

type Collect struct {
}

func (self Collect) Run(input, output Pipe) {
	group := []interface{}{}
	for v := range input {
		group = append(group, v)
	}

	output <- group

	close(output)
}

const SKIP_ALL = -1

type Skip struct {
	Count int
}

func (self Skip) Run(input, output Pipe) {
	i := 0
	for v := range input {
		if i < self.Count {
			i++
			continue
		}
		output <- v
	}
	close(output)
}

type LastValue struct {
	last interface{}
}

func (self *LastValue) Run(input, output Pipe) {
	for v := range input {
		self.last = v
		output <- v
	}
	close(output)
}

func (self *LastValue) Last() interface{} {
	return self.last
}

type Nonblocking struct {
}

func (self Nonblocking) Run(input, output Pipe) {
	for v := range input {
		select {
		case output <- v:
		default:
		}
	}
	close(output)
}

type StringContains struct {
	Str string
}

func (self StringContains) Run(input, output Pipe) {
	for v := range input {
		if strings.Contains(v.(string), self.Str) {
			output <- v
		}
	}
	close(output)
}
