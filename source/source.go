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

package source

import (
	"bufio"
	"os/exec"
	"syscall"
	"sync"
)

type Sourcer interface {
	Generate(out chan interface{}, ech chan error)
}

type Source struct {
	Command string
	Args    []string
}

func (s *Source) Generate(out chan interface{}, ech chan error) {
	cmd := exec.Command(s.Command, s.Args...)
	reader, err := cmd.StdoutPipe()

	if err != nil {
		ech <- err
		close(ech)
		return
	}

	scanner := bufio.NewScanner(reader)

	var done sync.WaitGroup
	done.Add(1)
	go func() {
		for scanner.Scan() {
			out <- scanner.Text()
		}
		close(out)
		done.Done()
	}()

	if err = cmd.Start(); err != nil {
		ech <- err
		close(ech)
		return
	}

	done.Wait()
	status := cmd.Wait()

	var waitStatus syscall.WaitStatus
	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	if waitStatus.ExitStatus() > 0 {
		ech <- status
		close(ech)
		return
	}
}
