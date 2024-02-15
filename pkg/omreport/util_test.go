/*
Copyright 2021 The dellhw_exporter Authors. All rights reserved.

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

package omreport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type SeverityTestResult struct {
	Input  string
	Output string
}

var severityTests = []SeverityTestResult{
	{
		Input:  "Ok",
		Output: "0",
	},
	{
		Input:  "Unknown",
		Output: "1",
	},
	{
		Input:  "Non-Critical",
		Output: "2",
	},
	{
		Input:  "What if a high severity issue walks into a bar?",
		Output: "1",
	},
}

func TestSeverity(t *testing.T) {
	for _, result := range severityTests {
		value := severity(result.Input)
		assert.Equal(t, result.Output, value)
	}
}
