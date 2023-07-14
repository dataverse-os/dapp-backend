package dapp

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func AppBaseDir() string {
	return filepath.Join(lo.Must(os.UserHomeDir()), ".dapp-backend")
}

func ModelsBoltPath() string {
	return filepath.Join(AppBaseDir(), "models.db")
}

func EnsureDir(name string) error {
	s, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(name, os.ModeDir|0750)
			if err != nil {
				return errors.Wrapf(err, "can't create app base dir %s", name)
			}

			return nil
		}

		return errors.Wrap(err, "unexpected error")
	}

	if !s.IsDir() {
		return errors.Wrapf(err, "app base dir %s is not a dir", name)
	}

	return nil
}
