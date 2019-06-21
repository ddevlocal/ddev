package ddevapp

import (
	"fmt"
	"github.com/drud/ddev/pkg/util"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Share provides the `ddev share` functionality by running ngrok with correct flags
func (app *DdevApp) Share(useHTTP bool, extraNgrokFlags []string) error {
	ngrokLoc, err := exec.LookPath("ngrok")
	if ngrokLoc == "" || err != nil {
		util.Failed("ngrok not found in path, please install it, see https://ngrok.com/download")
	}
	urls := []string{app.GetWebContainerDirectHTTPSURL(), app.GetWebContainerDirectHTTPURL()}

	if useHTTP {
		urls = []string{app.GetWebContainerDirectHTTPURL()}
	}

	var ngrokErr error
	for _, url := range urls {
		ngrokArgs := []string{"http"}
		if app.NgrokArgs != "" {
			ngrokArgs = append(ngrokArgs, strings.Split(app.NgrokArgs, " ")...)
		}
		ngrokArgs = append(ngrokArgs, url)
		ngrokArgs = append(ngrokArgs, extraNgrokFlags...)

		if strings.Contains(url, "http://") {
			util.Warning("Using local http URL, your data may be exposed on the internet. Create a free ngrok account instead...")
			time.Sleep(time.Second * 3)
		}
		util.Success("Running %s %s", ngrokLoc, strings.Join(ngrokArgs, " "))
		ngrokCmd := exec.Command(ngrokLoc, ngrokArgs...)
		ngrokCmd.Stdout = os.Stdout
		ngrokCmd.Stderr = os.Stderr
		ngrokErr = ngrokCmd.Run()

		// nil result means ngrok ran and exited normally.
		// It seems to do this fine when hit by SIGTERM or SIGINT
		if ngrokErr == nil {
			break
		}

		exitErr, ok := ngrokErr.(*exec.ExitError)
		if !ok {
			// Normally we'd have an ExitError, but if not, notify
			return fmt.Errorf("ngrok exited: %v", ngrokErr)
		}

		exitCode := exitErr.ExitCode()
		// In the case of exitCode==1, ngrok seems to have died due to an error,
		// most likely inadequate user permissions.
		if exitCode != 1 {
			return fmt.Errorf("ngrok exited: %v", exitErr)
		}
		// Otherwise we'll continue and do the next url or exit
	}
	return nil
}
