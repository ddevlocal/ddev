package ddevapp_test

import (
	"bufio"
	"encoding/json"
	"github.com/drud/ddev/pkg/ddevapp"
	"github.com/drud/ddev/pkg/testcommon"
	asrt "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestDdevShare tests ddevappapp.Share()
func TestDdevShare(t *testing.T) {
	assert := asrt.New(t)
	runTime := testcommon.TimeTrack(time.Now(), t.Name())
	defer runTime()

	site := TestSites[0]

	app, err := ddevapp.NewApp(site.Dir, false, "")
	require.NoError(t, err)

	err = app.Start()
	require.NoError(t, err)

	tmpDir := testcommon.CreateTmpDir(t.Name())
	logFile := filepath.Join(tmpDir, "ngrok.out")
	app.NgrokArgs = "-log " + logFile + " -log-format=json"
	var ngrokCmd *exec.Cmd

	go func() {
		err = app.Share(true, nil, ngrokCmd)
	}()

	time.Sleep(time.Second * 10)

	f, err := os.Open(logFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var urlRead = false

	// Read through the ngrok json output until we get the url it has opened
	for scanner.Scan() {
		logLine := scanner.Text()
		logData := make(map[string]string)

		err := json.Unmarshal([]byte(logLine), &logData)
		if err != nil {
			switch err.(type) {
			case *json.SyntaxError:
				continue
			default:
				t.Errorf("failed unmarshalling %v: %v", logLine, err)
				break
			}
		}
		// If URL is provided, try to hit it and look for expected response
		if url, ok := logData["url"]; ok {
			resp, err := http.Get(url + site.Safe200URIWithExpectation.URI)
			assert.NoError(err)
			//nolint: errcheck
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(err)
			assert.Contains(string(body), site.Safe200URIWithExpectation.Expect)
			urlRead = true

			//// The complexity here using github.com/shirou/gopsutil/process
			//// is a result of the need to kill subprocesses in Windows.
			//// Suggested by https://forum.golangbridge.org/t/how-can-i-use-syscall-kill-in-windows/11472/5
			//p, err := process.NewProcess(int32(cmd.Process.Pid))
			//assert.NoError(err)
			//children, err := p.Children()
			//assert.NoError(err)
			//for _, v := range children {
			//	err = v.Kill() // Kill each child
			//	assert.NoError(err)
			//}
			//err = p.Kill()
			//assert.NoError(err)
		}
	}
	assert.True(urlRead)
}
