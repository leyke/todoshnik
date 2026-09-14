package test

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/joho/godotenv"
)

var (
	loadEnvOnce sync.Once
	loadEnvErr  error
)

func loadEnv() error {
	loadEnvOnce.Do(func() {
		_, filename, _, _ := runtime.Caller(0)

		root := filepath.Join(
			filepath.Dir(filename),
			"..",
			"..",
			"..",
			"..",
		)
		loadEnvErr = godotenv.Load(filepath.Join(root, ".env"))
	})

	return loadEnvErr
}
