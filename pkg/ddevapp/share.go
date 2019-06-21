package ddevapp

import (
	"github.com/drud/ddev/pkg/util"
	"os"
	"os/exec"
	"strings"
	"sync"
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

	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	//process, _ := os.FindProcess(os.Getpid())

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)

	//var ngrokCmd *Cmd
	go func(waitgroup *sync.WaitGroup) {
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

			ngrokErr := ngrokCmd.Run()

			// nil result means ngrok ran and exited normally.
			// It seems to do this fine when hit by SIGTERM or SIGINT
			if ngrokErr == nil {
				break
			}

			exitErr, ok := ngrokErr.(*exec.ExitError)
			if !ok {
				// Normally we'd have an ExitError, but if not, notify
				util.Error("ngrok exited: %v", ngrokErr)
				break
			}

			exitCode := exitErr.ExitCode()
			// In the case of exitCode==1, ngrok seems to have died due to an error,
			// most likely inadequate user permissions/configuration
			if exitCode != 1 {
				util.Error("ngrok exited: %v", exitErr)
				break
			}
			// Otherwise we'll continue and try the next url or exit
		}
		waitgroup.Done()
	}(&waitgroup)

	//s := <-sigs
	//util.Success("Received signal %v", s)

	waitgroup.Wait()
	return nil
}
