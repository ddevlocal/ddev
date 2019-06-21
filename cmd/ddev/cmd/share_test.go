package cmd

import (
	"bufio"
	"encoding/json"
	"github.com/shirou/gopsutil/process"
	asrt "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
	"testing"
)

// TestShareCmd tests `ddev share`
func TestShareCmd(t *testing.T) {
	assert := asrt.New(t)
	urlRead := false

	site := TestSites[0]
	defer site.Chdir()()

	// Configure ddev/ngrok to use json output to stdout
	cmd := exec.Command(DdevBin, "config", "--ngrok-args", "-log stdout -log-format=json")
	err := cmd.Start()
	require.NoError(t, err)
	err = cmd.Wait()
	require.NoError(t, err)

	cmd = exec.Command(DdevBin, "share", "--use-http")

	cmdReader, err := cmd.StdoutPipe()
	require.NoError(t, err)
	scanner := bufio.NewScanner(cmdReader)

	// Read through the ngrok json output until we get the url it has opened
	go func() {
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

				// The complexity here using https://github.com/shirou/gopsutil/process
				// is a result of the need to kill subprocesses in Windows.
				// Suggested by https://forum.golangbridge.org/t/how-can-i-use-syscall-kill-in-windows/11472/5
				ddevProc, err := process.NewProcess(int32(cmd.Process.Pid))
				assert.NoError(err)
				children, err := ddevProc.Children()
				assert.NoError(err)
				// Kill off the child ngrok process first.
				for _, v := range children {
					name, _ := v.Name()
					status, _ := v.Status()
					t.Logf("Killing %v with status %v", name, status)
					err = v.Terminate()
					assert.NoError(err)
					err = v.Kill()
					assert.NoError(err)
				}
				name, _ := ddevProc.Name()
				status, _ := ddevProc.Status()
				t.Logf("Killing %v with status %v", name, status)

				err = ddevProc.Terminate()
				assert.NoError(err)
				return
			}
		}
		return
	}()
	err = cmd.Start()
	require.NoError(t, err)
	err = cmd.Wait()
	t.Logf("cmd.Wait() err: %v", err)
	assert.True(urlRead)
	_ = cmdReader.Close()
	t.Logf("goprocs: %v", runtime.NumGoroutine())
}
