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
	"encoding/json"
	"fmt"
	"os/exec"
)

type simpleSource struct {
	args    []string
	command string
	rawData []byte
}

// New creates simple source as external program output.
// It takes command name and its arguments as parameters.
// It return pointer to new source object ready to execute.
func New(cmd string, args []string) *simpleSource {
	return &simpleSource{command: cmd, args: args}
}

// Run executes source command and gathers its output in internal structure.
// It returns nil in case of success and returns error if command execution failed.
func (ss *simpleSource) Run() error {

	out, err := exec.Command(ss.command, ss.args...).Output()
	if err != nil {
		return err
	}

	ss.rawData = out
	return nil
}

// Raw returns output of executed command in raw format (as is)
func (ss *simpleSource) Raw() []byte {
	return ss.rawData
}

// OutputMap translates json document data from command execution to a map.
// It returns map literals with its values, possibly nested maps and slices
func (ss *simpleSource) OutputMap() map[string]interface{} {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(ss.rawData, &jsonMap); err != nil {
		fmt.Println(err)
		return nil
	}
	return jsonMap
}
