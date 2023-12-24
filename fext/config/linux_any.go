//go:build linux && !(arm64 || amd64)

package config

import "runtime"

const MarkerArch = runtime.GOARCH // platform_machine (platform.machine())
