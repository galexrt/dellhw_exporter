// This file has been taken from the https://github.com/coreos/pkg repository
// The https://github.com/coreos/pkg repositroy and it's code is licensed under
// tha Apache 2.0 License, see https://github.com/coreos/pkg/blob/master/LICENSE.

package flagutil

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// SetFlagsFromEnv parses all registered flags in the given flagset,
// and if they are not already set it attempts to set their values from
// environment variables. Environment variables take the name of the flag but
// are UPPERCASE, and any dashes are replaced by underscores. Environment
// variables additionally are prefixed by the given string followed by
// and underscore. For example, if prefix=PREFIX: some-flag => PREFIX_SOME_FLAG
func SetFlagsFromEnv(fs *flag.FlagSet, prefix string) (err error) {
	alreadySet := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})
	fs.VisitAll(func(f *flag.Flag) {
		if !alreadySet[f.Name] {
			key := prefix + "_" + strings.ToUpper(strings.Replace(strings.Replace(f.Name, ".", "_", -1), "-", "_", -1))
			val := os.Getenv(key)
			if val != "" {
				if serr := fs.Set(f.Name, val); serr != nil {
					err = fmt.Errorf("invalid value %q for %s: %v", val, key, serr)
				}
			}
		}
	})
	return err
}
