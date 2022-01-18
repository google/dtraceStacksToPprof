// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/pprof/profile"
)

type pprofHelper struct {
	functions map[string]*profile.Function
	locations map[string]*profile.Location
	nextId    uint64
}

func (helper *pprofHelper) getOrInsertLocation(name string) *profile.Location {
	location, ok := helper.locations[name]
	if !ok {
		location = &profile.Location{
			ID:   helper.nextId,
			Line: []profile.Line{{Function: helper.getOrInsertFunction(name)}},
		}
		helper.locations[name] = location
		helper.nextId++
	}
	return location
}

func (helper *pprofHelper) getOrInsertFunction(name string) *profile.Function {
	function, ok := helper.functions[name]
	if !ok {
		function = &profile.Function{
			ID:   helper.nextId,
			Name: name,
		}
		helper.functions[name] = function
		helper.nextId++
	}
	return function
}

func main() {
	fileNamePtr := flag.String("output", "profile.pb.gz", "Pprof output file name")
	flag.Parse()
	pprof := &profile.Profile{
		SampleType: []*profile.ValueType{{Type: "stacks", Unit: "count"}},
		Sample:     make([]*profile.Sample, 0),
		Location:   make([]*profile.Location, 0),
		Function:   make([]*profile.Function, 0),
	}
	helper := &pprofHelper{
		functions: make(map[string]*profile.Function),
		locations: make(map[string]*profile.Location),
		nextId:    1,
	}
	stackHeaderRe := regexp.MustCompile(`^\S+:$`)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && !stackHeaderRe.MatchString(strings.TrimSpace(scanner.Text())) {
	}
	var stack []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if numStacks, err := strconv.ParseInt(line, 10, 64); err == nil {
			stackTrace := make([]*profile.Location, 0)
			for _, function := range stack {
				stackTrace = append(stackTrace, helper.getOrInsertLocation(function))
			}
			pprof.Sample = append(pprof.Sample, &profile.Sample{
				Location: stackTrace,
				Value:    []int64{numStacks},
			})
			stack = make([]string, 0)
		} else if !stackHeaderRe.MatchString(line) {
			function := line
			if strings.Contains(line, "`") {
				function = strings.Split(strings.Split(line, "`")[1], "+")[0]
			}
			stack = append(stack, function)
		}
	}
	for _, fn := range helper.functions {
		pprof.Function = append(pprof.Function, fn)
	}
	for _, loc := range helper.locations {
		pprof.Location = append(pprof.Location, loc)
	}
	out, err := os.Create(*fileNamePtr)
	if err != nil {
		log.Fatalf("Unable to write output file %s: %q", *fileNamePtr, err)
	}
	defer out.Close()
	pprof.Write(out)
}
