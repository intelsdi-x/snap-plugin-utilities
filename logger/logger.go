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

package logger

import (
	"os"
	"runtime"

	log "github.com/Sirupsen/logrus"
)

func init() {
	// Log as default ASCII formatter
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func LogInfo(message string, args ...interface{}) {
	setEntry(args...).Info(message)
}

func LogWarn(message string, args ...interface{}) {
	setEntry(args...).Warn(message)
}

func LogDebug(message string, args ...interface{}) {
	setEntry(args...).Debug(message)
}

func LogError(message string, args ...interface{}) {
	setEntry(args...).Error(message)
}

func LogFatal(message string, args ...interface{}) {
	setEntry(args...).Fatal(message)
}

func LogPanic(message string, args ...interface{}) {
	setEntry(args...).Panic(message)
}

func setEntry(args ...interface{}) *log.Entry {
	caller := getFunctionName()
	if len(args) > 1 {

		return log.WithFields(log.Fields{
			"_func":          caller,
			args[0].(string): args[1],
		})
	} else {
		return log.WithFields(log.Fields{
			"_func": caller,
		})
	}
}

func getFunctionName() string {
	pc := make([]uintptr, 10)

	if runtime.Callers(4, pc) == 0 {
		return "!<not_accessible>"
	}

	return runtime.FuncForPC(pc[0]).Name()
}
