package nodeps

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"runtime"
	"time"
)

// ArrayContainsString returns true if slice contains element
func ArrayContainsString(slice []string, element string) bool {
	return !(PosString(slice, element) == -1)
}

// PosString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func PosString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// RemoveItemFromSlice returns a slice with item removed
// If the item does not exist, the slice is unchanged
// This is quite slow in the scheme of things, so shouldn't
// be used without examination
func RemoveItemFromSlice(slice []string, item string) []string {
	pos := PosString(slice, item)
	if pos != -1 {
		// Remove the element at index i from a.
		copy(slice[pos:], slice[pos+1:]) // Shift slice[pos+1:] left one index.
		slice[len(slice)-1] = ""         // Erase last element (write zero value).
		slice = slice[:len(slice)-1]     // Truncate slice.
	}
	return slice
}

// From https://www.calhoun.io/creating-random-strings-in-go/
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// RandomString creates a random string with a set length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// IsWSL2() returns true if running WSL2
func IsWSL2() bool {
	return GetWSLDistro() != ""
}

// GetWSLDistro returns the WSL2 distro name if on Linux
func GetWSLDistro() string {
	wslDistro := ""
	if runtime.GOOS == "linux" {
		wslDistro = os.Getenv("WSL_DISTRO_NAME")
	}
	return wslDistro
}

// ClearDockerEnv unsets env vars set in platform DockerEnv() so that
// they can be set by another test run.
func ClearDockerEnv() {
	envVars := []string{
		"COMPOSE_PROJECT_NAME",
		"COMPOSE_CONVERT_WINDOWS_PATHS",
		"DDEV_SITENAME",
		"DDEV_DBIMAGE",
		"DDEV_WEBIMAGE",
		"DDEV_APPROOT",
		"DDEV_HOST_WEBSERVER_PORT",
		"DDEV_HOST_HTTPS_PORT",
		"DDEV_DOCROOT",
		"DDEV_HOSTNAME",
		"DDEV_PHP_VERSION",
		"DDEV_WEBSERVER_TYPE",
		"DDEV_PROJECT_TYPE",
		"DDEV_ROUTER_HTTP_PORT",
		"DDEV_ROUTER_HTTPS_PORT",
		"DDEV_HOST_DB_PORT",
		"DDEV_HOST_WEBSERVER_PORT",
		"DDEV_PHPMYADMIN_PORT",
		"DDEV_PHPMYADMIN_HTTPS_PORT",
		"DDEV_MAILHOG_PORT",
		"COLUMNS",
		"LINES",
		"DDEV_XDEBUG_ENABLED",
		"IS_DDEV_PROJECT",
	}
	for _, env := range envVars {
		err := os.Unsetenv(env)
		if err != nil {
			logrus.Printf("failed to unset %s: %v\n", env, err)
		}
	}
}
