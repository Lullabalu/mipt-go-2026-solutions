//go:build !solution

package fileleak

import (
	"os"
	"path/filepath"
)

type testingT interface {
	Errorf(msg string, args ...any)
	Cleanup(func())
}

func VerifyNone(t testingT) {
	mp := make(map[string]int)

	files, _ := os.ReadDir("/proc/self/fd")

	for _, file := range files {
		fd := file.Name()

		path := filepath.Join("/proc/self/fd", fd)

		name, err := os.Readlink(path)
		if err != nil {
			continue
		}

		mp[name] += 1
	}

	t.Cleanup(func() {
		files, _ := os.ReadDir("/proc/self/fd")

		for _, file := range files {
			fd := file.Name()

			path := filepath.Join("/proc/self/fd", fd)

			name, err := os.Readlink(path)
			if err != nil {
				continue
			}

			mp[name] -= 1
		}
		for k, v := range mp {
			if v < 0 {
				t.Errorf("file %k appeared", k)
			}
		}
	})
}
