package db

import (
	"fmt"
	"io"
	"os"
)

func isEmpty(name string) (bool, error) {
	//nolint:gosec // The function is package-only and input is checked in caller.
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	//nolint:errcheck,gosec // This runs as one-shot program, so worst case, system cleans this up.
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func checkDataDirectory(datadir string) error {
	fi, err := os.Stat(datadir)

	if os.IsNotExist(err) {
		if err = os.Mkdir(datadir, datadirPermissions); err != nil {
			return err
		}
		fi, err = os.Stat(datadir)
	}
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("'%s' exists and is not a directory", datadir)
	}

	if fi.Mode().Perm() != datadirPermissions {
		return fmt.Errorf("permissions for '%s' are incorrect (%o != %o)", datadir, fi.Mode(), datadirPermissions)
	}

	return nil
}
