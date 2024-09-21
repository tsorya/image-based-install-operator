package executer

import (
	"bytes"
	"os/exec"
)

//go:generate mockgen --build_flags=--mod=mod -package executer -destination mock_executer.go . Executer
type Executer interface {
	Execute(command string, args ...string) (stdout string, stderr string, exitCode int)
}

type OsExecuter struct{}

func (e *OsExecuter) Execute(command string, args ...string) (stdout string, stderr string, exitCode int) {
	cmd := exec.Command(command, args...)
	return e.execute(cmd)
}

func (e *OsExecuter) execute(cmd *exec.Cmd) (stdout string, stderr string, exitCode int) {
	var stdoutBytes, stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes
	err := cmd.Run()
	return stdoutBytes.String(), getErrorStr(err, &stderrBytes), getExitCode(err)
}

func getExitCode(err error) int {
	if err == nil {
		return 0
	}
	switch value := err.(type) {
	case *exec.ExitError:
		return value.ExitCode()
	default:
		return -1
	}
}

func getErrorStr(err error, stderr *bytes.Buffer) string {
	b := stderr.Bytes()
	if len(b) > 0 {
		return string(b)
	} else if err != nil {
		return err.Error()
	}
	return ""
}
