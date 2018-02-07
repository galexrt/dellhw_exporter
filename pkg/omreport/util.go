package omreport

// Command executes the named program with the given arguments. If it does not
import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/prometheus/common/log"
)

var (
	// ErrPath is returned by Command if the program is not in the PATH.
	ErrPath = errors.New("program not in PATH")
	// ErrTimeout is returned by Command if the program timed out.
	ErrTimeout = errors.New("program killed after timeout")
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

// severity returns 1 if s is not "Ok" or "Non-Critical", elif is "Non-Critical" 2 else 0.
func severity(s string) string {
	if s != "Ok" && s != "Non-Critical" {
		return "1"
	}
	if s == "Non-Critical" {
		return "2"
	}
	return "0"
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
func Command(timeout time.Duration, stdin io.Reader, name string, arg ...string) (io.Reader, error) {
	if _, err := exec.LookPath(name); err != nil {
		return nil, ErrPath
	}
	log.Debug("executing command: ", name, arg)
	c := exec.Command(name, arg...)
	b := &bytes.Buffer{}
	c.Stdout = b
	c.Stdin = stdin
	if err := c.Start(); err != nil {
		return nil, err
	}
	timedOut := false
	intTimer := time.AfterFunc(timeout, func() {
		log.Error("Process taking too long. Interrupting: ", name, strings.Join(arg, " "))
		c.Process.Signal(os.Interrupt)
		timedOut = true
	})
	killTimer := time.AfterFunc(timeout, func() {
		log.Error("Process taking too long. Killing: ", name, strings.Join(arg, " "))
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

// ReadCommand runs command name with args and calls line for each line from its
// stdout. Command is interrupted (if supported by Go) after 10 seconds and
// killed after 20 seconds.
func readCommand(line func(string) error, name string, arg ...string) error {
	return readCommandTimeout(time.Second*10, line, nil, name, arg...)
}

// ReadCommandTimeout is the same as ReadCommand with a specifiable timeout.
// It can also take a []byte as input (useful for chaining commands).
func readCommandTimeout(timeout time.Duration, line func(string) error, stdin io.Reader, name string, arg ...string) error {
	b, err := Command(timeout, stdin, name, arg...)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		if err := line(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error(name, " : ", err)
	}
	return nil
}
