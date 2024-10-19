/*
Copyright 2024 The dellhw_exporter Authors. All rights reserved.

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

// Command executes the named program with the given arguments. If it does not
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"sync/atomic"
	"time"
	"unicode"
	"unicode/utf8"

	"go.uber.org/zap"
)

type Output = []Report

type Report struct {
	Title       string
	Description string

	Lines []Line
}

type Line = map[string]string

var (
	// ErrPath is returned by Command if the program is not in the PATH.
	ErrPath = errors.New("program not in PATH")
	// ErrTimeout is returned by Command if the program timed out.
	ErrTimeout = errors.New("program killed after timeout")

	// cmdTimeout configurable timeout for commands.
	cmdTimeout int64 = 10

	logger = zap.NewNop()
)

// clean concatenates arguments with a space and removes extra whitespace.
func clean(ss ...string) string {
	v := strings.Join(ss, " ")
	fs := strings.Fields(v)
	return strings.Join(fs, " ")
}

// extract tries to return a parsed number from s with given suffix. A space may
// be present between number ond suffix.
func extract(s, suffix string) (string, error) {
	if !strings.HasSuffix(s, suffix) {
		return "0", fmt.Errorf("extract: suffix not found")
	}
	s = s[:len(s)-len(suffix)]
	return strings.TrimSpace(s), nil
}

// severity returns 1 if s is not "Ok" or "Non-Critical" (should be "Critical" then in most cases)
// elif is "Non-Critical" 2 else 0.
func severity(s string) string {
	if s != "Ok" && s != "Non-Critical" {
		return "1"
	}
	if s == "Non-Critical" {
		return "2"
	}
	return "0"
}

func pdiskState(s string) string {
	states := map[string]string{
		"Unknown":              "0",
		"Ready":                "1",
		"Online":               "2",
		"Degraded":             "3",
		"Failed":               "4",
		"Offline":              "5",
		"Rebuilding":           "6",
		"Incompatible":         "7",
		"Removed":              "8",
		"Clear":                "9",
		"SMART Alert Detected": "10",
		"Foreign":              "11",
		"Unsupported":          "12",
		"Replacing":            "13",
		"Non-RAID":             "14",
	}

	return states[s]
}

func vdiskState(s string) string {
	states := map[string]string{
		"Ready":                     "1",
		"Degraded":                  "2",
		"Resynching":                "3",
		"Resynching Paused":         "4",
		"Regenerating":              "5",
		"Reconstructing":            "6",
		"Failed":                    "7",
		"Failed Redundancy":         "8",
		"Background Initialization": "9",
		"Formatting":                "10",
		"Initializing":              "11",
		"Degraded Redundancy":       "12",
	}

	return states[s]
}

func vdiskReadPolicy(s string) string {
	policies := map[string]string{
		"Not Applicable":      "0",
		"Read Ahead":          "1",
		"No Read Ahead":       "2",
		"Read Cache Enabled":  "3",
		"Read Cache Disabled": "4",
		"Adaptive Read Ahead": "5",
	}

	return policies[s]
}

func vdiskWritePolicy(s string) string {
	policies := map[string]string{
		"Not Applicable":                "0",
		"Write Ahead":                   "1",
		"Force Write Back":              "2",
		"Write Back Enabled":            "3",
		"Write Through":                 "4",
		"Write Cache Enabled Protected": "5",
		"Write Cache Disabled":          "6",
		"Write Back":                    "7",
	}

	return policies[s]
}

func vdiskCachePolicy(s string) string {
	policies := map[string]string{
		"Not Applicable": "0",
		"Cache I/O":      "1",
		"Direct I/O":     "2",
	}

	return policies[s]
}

// yesNoToBool returns "1" for "Yes" and "0" for "No"
func yesNoToBool(s string) string {
	if s == "Yes" {
		return "1"
	}
	return "0"
}

var getNumberFromStringRegex = regexp.MustCompile("[0-9]+")

func getNumberFromString(s string) string {
	result := getNumberFromStringRegex.FindString(s)
	if result != "" {
		return result
	}
	return "-1"
}

func replace(name string) string {
	r, _ := Replace(name, "_")
	return r
}

// Replace certain chars in a string
func Replace(s, replacement string) (string, error) {
	var c string
	replaced := false
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' || r == '/' {
			c += string(r)
			replaced = false
		} else if !replaced {
			c += replacement
			replaced = true
		}
		s = s[size:]
	}
	if len(c) == 0 {
		return "", fmt.Errorf("clean result is empty")
	}
	return c, nil
}

// Command exit within timeout, it is sent SIGINT (if supported by Go). After
// another timeout, it is killed.
func Command(timeout time.Duration, stdin io.Reader, name string, args ...string) (io.Reader, error) {
	if _, err := exec.LookPath(name); err != nil {
		return nil, ErrPath
	}
	logger.Debug("executing command", zap.String("command", name), zap.Strings("args", args))
	c := exec.Command(name, args...)
	b := &bytes.Buffer{}
	c.Stdout = b
	c.Stdin = stdin
	if err := c.Start(); err != nil {
		return nil, err
	}
	timedOut := false
	intTimer := time.AfterFunc(timeout, func() {
		logger.Error("process taking too long, interrupting: ", zap.String("command", name), zap.Strings("args", args))
		c.Process.Signal(os.Interrupt)
		timedOut = true
	})
	killTimer := time.AfterFunc(timeout, func() {
		logger.Error("process taking too long, killing", zap.String("command", name), zap.Strings("args", args))
		c.Process.Signal(os.Interrupt)
		timedOut = true
	})
	err := c.Wait()
	intTimer.Stop()
	killTimer.Stop()
	if timedOut {
		return nil, ErrTimeout
	}
	return b, err
}

// ReadCommand runs command name with args and calls fn for the output from
// stdout. Command is interrupted (if supported by Go) after 10 seconds and
// killed after 20 seconds.
func readCommand(fn func(string) error, name string, arg ...string) error {
	timeout := time.Duration(int(atomic.LoadInt64(&cmdTimeout)))
	return readCommandTimeout(timeout*time.Second, fn, nil, name, arg...)
}

// ReadCommandTimeout is the same as ReadCommand with a specifiable timeout.
// It can also take a []byte as input (useful for chaining commands).
func readCommandTimeout(timeout time.Duration, fn func(string) error, stdin io.Reader, name string, args ...string) error {
	b, err := Command(timeout, stdin, name, args...)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Skip exit code 255, it should indicate that no devices have been found in some cases
			if exitErr.ExitCode() == 255 {
				return nil
			}
		}
		return fmt.Errorf("failed to execute command (\"%s %s\"). %w", name, args, err)
	}

	out, err := io.ReadAll(b)
	if err != nil {
		logger.Error("failed to read command output", zap.String("command", name), zap.Strings("args", args), zap.Error(err))
	}

	if err := fn(string(out[:])); err != nil {
		return fmt.Errorf("failed to process command (\"%s %s\") output. %w", name, args, err)
	}

	return nil
}

// SetCommandTimeout this function can be used to atomically set the command execution timeout
func SetCommandTimeout(timeout int64) {
	atomic.StoreInt64(&cmdTimeout, timeout)
}

func hasKeys(in map[string]string, fields ...string) bool {
	for _, field := range fields {
		field = normalizeName(field)
		if _, ok := in[field]; !ok {
			return false
		}
	}

	return true
}

func parseOutput(mode ReaderMode, input string) Output {
	output := Output{
		{},
	}
	ri := 0
	gotTitle := false
	kvSeparated := false

	keyLine := ""
	keys := []string{}

	prevLine := ""
	nextLine := ""

	spl := strings.Split(input, "\n")
	for i, line := range spl {
		if i > 0 {
			prevLine = spl[i-1]
		}
		if len(spl) > i+1 {
			nextLine = spl[i+1]
		} else {
			nextLine = ""
		}

		line = clean(line)

		if line == "" {
			if strings.Contains(prevLine, ";") {
				output = append(output, Report{})
				ri++
				gotTitle = false
				kvSeparated = false
				keyLine = ""
				keys = []string{}
			}
			continue
		}

		if strings.HasPrefix(line, "For further help") {
			continue
		}

		if !strings.Contains(line, ";") {
			if !gotTitle {
				output[ri].Title = line
				gotTitle = true
			} else {
				output[ri].Description = line
			}
		} else {
			sp := strings.Split(line, ";")

			// Handle special cases..
			if output[ri].Description == "Version Information" {
				keyLine = line
				keys = []string{"component", "version"}
				kvSeparated = true
			} else if output[ri].Title == "Amperage" {
				keys = []string{"psu", "amperage"}
			} else if output[ri].Title == "BIOS Information" {
				kvSeparated = true
			} else if strings.Count(line, ";") > 1 {
				kvSeparated = true
			}

			if len(keys) == 0 {
				keyLine = line
				keys = append(keys, sp...)

				// Normalize keys to lower case
				for i := 0; i < len(keys); i++ {
					keys[i] = normalizeName(keys[i])
				}

				if prevLine == "" && strings.Contains(nextLine, ";") && keyLine != line {
					continue
				}
			}

			if (mode == KeyValueReaderMode || kvSeparated) && strings.Count(line, ";") == 1 {
				output[ri].Lines = append(output[ri].Lines, Line{
					normalizeName(sp[0]): sp[1],
				})
			} else if mode <= TableReaderMode && keyLine != line {
				l := Line{}
				for i, s := range sp {
					if i > len(keys)-1 {
						continue
					}

					l[keys[i]] = s
				}

				output[ri].Lines = append(output[ri].Lines, l)
			}
		}
	}

	if len(output[ri].Lines) == 0 {
		output = slices.Delete(output, ri, ri+1)
	}

	return output
}

func normalizeName(in string) string {
	return strings.Replace(strings.ToLower(in), " ", "_", -1)
}
