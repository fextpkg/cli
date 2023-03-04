//go:build windows && amd64
// +build windows,amd64

package web

func checkCompatibility(platformTag string) bool {
	return platformTag == "win_amd64"
}
