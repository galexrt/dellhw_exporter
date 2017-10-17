package omreport

import (
	"reflect"
	"testing"
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
}

func TestSeverity(t *testing.T) {
	for _, result := range severityTests {
		value := severity(result.Input)
		if !reflect.DeepEqual(value, result.Output) {
			t.Errorf("severity result not equal: %+v - %+v\n", value, result.Output)
		}
	}
}
