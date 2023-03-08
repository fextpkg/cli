//go:build linux && !(arm64 || amd64)

package config

import "runtime"

const PythonArch = runtime.GOARCH
