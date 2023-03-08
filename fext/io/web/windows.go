//go:build windows && amd64

package web

import "github.com/fextpkg/cli/fext/config"

func checkCompatibility(platformTag string) (bool, error) {
	return platformTag == config.PlatformTag, nil
}
