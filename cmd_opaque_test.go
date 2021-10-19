package cmder_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/scottames/cmder"
	"github.com/stretchr/testify/assert"
)

const (
	echo       string = "echo"
	cat        string = "cat"
	sleep      string = "sleep"
	foo        string = "foo"
	five       string = "5"
	newLineStr string = "\n"
)

var (
	logCmdStr string
	silentStr string
)

func Test_CombinedOutput(t *testing.T) {
	expected := foo + foo
	out, err := cmder.
		New("bash", "-c", fmt.Sprintf("printf %s | tee /dev/stderr", foo)).
		CombinedOutput()

	if err != nil {
		t.Error(err)
	}

	actual := string(out)
	msg := fmt.Sprintf("Expected %v. Got %v.", expected, actual)
	assert.Equal(t, expected, actual, msg)
}

func Test_DryRunLocal(t *testing.T) {
	expected := false
	cmd := cmder.New(echo, foo).DryRun()

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}

	actual := cmd.Complete()
	msg := fmt.Sprintf("Expected %t. Got %t.", expected, actual)
	assert.Equal(t, expected, actual, msg)
}

func Test_Duration(t *testing.T) {
	cmd := cmder.New(echo, foo)

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}

	duration := cmd.Duration()
	if duration <= 0 {
		t.Errorf("Expected positive duration. Got %d.", duration)
	}
}

func Test_CmdRunExitCode(t *testing.T) {
	expected := 0
	cmd := cmder.New(echo, foo)

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}

	actual := cmd.ExitCode()
	msg := fmt.Sprintf("Expected %v. Got %v.", expected, actual)
	assert.Equal(t, expected, actual, msg)
}

func Test_InWithInput(t *testing.T) {
	expected := []byte(foo)

	actual, err := cmder.New(cat).In([]byte(foo)...).Output()
	if err != nil {
		t.Error(err)
	}

	msg := fmt.Sprintf("Expected %s. Got %s.", expected, actual)
	assert.Equal(t, expected, actual, msg)
}

func Test_InFromStdin(t *testing.T) {
	expected := []byte("foo")

	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}

	_, err = w.Write(expected)
	if err != nil {
		t.Error(err)
	}

	w.Close()

	orig := os.Stdin
	defer func() { os.Stdin = orig }()

	os.Stdin = r

	actual, err := cmder.New("cat").In().Output()
	if err != nil {
		t.Error(err)
	}

	msg := fmt.Sprintf("Expected %s. Got %s.", expected, actual)
	assert.Equal(t, expected, actual, msg)
}

func Test_StartKill(t *testing.T) {
	cmd := cmder.New(echo, foo)

	err := cmd.Start()
	if err != nil {
		t.Error(err)
	}

	err = cmd.Kill()
	if err != nil {
		t.Error(err)
	}
}

type testLogger struct{}

func (t testLogger) Log(v ...interface{}) {
	logCmdStr = fmt.Sprint(v...)
}

func (t testLogger) Logf(format string, v ...interface{}) {
	logCmdStr = fmt.Sprint(v...)
}

func Test_LogCmd(t *testing.T) {
	expected := "[echo foo]"

	cmder.New(echo, foo).Logger(testLogger{}).LogCmd()

	actual := logCmdStr
	msg := fmt.Sprintf("Expected '%s' Got '%s'", expected, actual)
	assert.Equal(t, expected, actual, msg)
}

func Test_OutStdout(t *testing.T) {
	expected := foo + newLineStr

	var buf bytes.Buffer

	err := cmder.New(echo, foo).Out(&buf).Run()
	if err != nil {
		t.Error(err)
	}

	actual := buf.String()
	msg := fmt.Sprintf("Expected '%s' Got '%s'", expected, actual)
	assert.Equal(t, expected, actual, msg)
}

func Test_OutStdoutStderr(t *testing.T) {
	expected := foo

	var bufStdout bytes.Buffer

	var bufStderr bytes.Buffer

	err := cmder.New("bash", "-c", fmt.Sprintf("printf %s | tee /dev/stderr", foo)).
		Out(&bufStdout, &bufStderr).
		Run()

	if err != nil {
		t.Error(err)
	}

	stdout := bufStdout.String()
	stderr := bufStderr.String()
	msg := fmt.Sprintf("Expected stdout to be '%s' Got '%s'", expected, stdout)
	assert.Equal(t, expected, stdout, msg)

	msg = fmt.Sprintf("Expected stderr to be '%s' Got '%s'", expected, stderr)
	assert.Equal(t, expected, stderr, msg)
}

func Test_Output(t *testing.T) {
	expected := []byte(foo + newLineStr)

	actual, err := cmder.New(echo, foo).Output()
	if err != nil {
		t.Error(err)
	}

	msg := fmt.Sprintf("Expected '%s' Got '%s'", expected, actual)
	assert.Equal(t, expected, actual, msg)
}

func Test_Process(t *testing.T) {
	cmd := cmder.New(sleep, five)

	err := cmd.Start()
	if err != nil {
		t.Error(err)
	}

	process := cmd.Process()
	if process == nil {
		t.Errorf("Expected non-nil process. Got nil value.")
	}

	err = cmd.Kill()
	if err != nil {
		t.Error(err)
	}
}

func Test_Pid(t *testing.T) {
	cmd := cmder.New(sleep, five)

	err := cmd.Start()
	if err != nil {
		t.Error(err)
	}

	pid := cmd.Pid()
	if pid == nil || *pid <= 0 {
		msg := fmt.Sprintf("Expected pid non-zero pid. Got %d.", pid)
		t.Errorf(msg, pid)
	}

	err = cmd.Kill()
	if err != nil {
		t.Error(err)
	}
}

func Test_Run(t *testing.T) {
	cmd := cmder.New(echo, foo)

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
}

func Test_RunFn(t *testing.T) {
	fn := cmder.New(echo, foo).RunFn()

	err := fn()
	if err != nil {
		t.Error(err)
	}
}

func Test_Silent(t *testing.T) {
	expected := ""

	cmder.New(echo, foo).Logger(testLogger{}).Silent().LogCmd()

	actual := silentStr
	msg := fmt.Sprintf("Expected '%s' Got '%s'", expected, actual)
	assert.Equal(t, expected, actual, msg)
}

func Test_StartFn(t *testing.T) {
	fn := cmder.New(echo, foo).StartFn()

	err := fn()
	if err != nil {
		t.Error(err)
	}
}

func Test_StartWait(t *testing.T) {
	cmd := cmder.New(echo, foo)

	err := cmd.Start()
	if err != nil {
		t.Error(err)
	}

	err = cmd.Wait()
	if err != nil {
		t.Error(err)
	}
}

func Test_String(t *testing.T) {
	expected := "/usr/bin/echo foo"
	actual := cmder.New(echo, foo).String()
	msg := fmt.Sprintf("Expected %v. Got %v.", expected, actual)
	assert.Equal(t, expected, actual, msg)
}
