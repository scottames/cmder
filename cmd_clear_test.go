package cmder

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

const (
	echo      string = "echo"
	foo       string = "foo"
	bar       string = "bar"
	baz       string = "baz"
	envFooBar string = "FOO=bar"
	tmp       string = "/tmp"
)

var (
	env      = os.Environ()
	testCmd1 = cmd{
		ctx:     context.Background(),
		env:     env,
		stdin:   os.Stdin,
		stderr:  os.Stderr,
		stdout:  os.Stdout,
		strings: []string{echo, foo},
	}
	testCmd2 = cmd{
		ctx:     context.Background(),
		env:     env,
		stdin:   os.Stdin,
		stderr:  os.Stderr,
		stdout:  os.Stdout,
		strings: []string{echo, foo, bar},
	}
	testCmd3 = cmd{
		ctx:     context.Background(),
		env:     env,
		stdin:   os.Stdin,
		stderr:  os.Stderr,
		stdout:  os.Stdout,
		strings: []string{echo, foo, bar, baz},
	}
	testCmd4 = cmd{
		ctx:     context.Background(),
		env:     env,
		stdin:   os.Stdin,
		stderr:  os.Stderr,
		stdout:  os.Stdout,
		strings: []string{echo, foo},
	}
	testCmd5 = cmd{
		ctx:     context.Background(),
		dir:     tmp,
		env:     env,
		stdin:   os.Stdin,
		stderr:  os.Stderr,
		stdout:  os.Stdout,
		strings: []string{echo, foo},
	}
	testCmd6 = cmd{
		ctx:     context.Background(),
		env:     append(env, envFooBar),
		stdin:   os.Stdin,
		stderr:  os.Stderr,
		stdout:  os.Stdout,
		strings: []string{echo, foo},
	}
)

func Test_New(t *testing.T) {
	expected := &testCmd1
	actual := New(echo, foo)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_ArgsSingle(t *testing.T) {
	expected := &testCmd2
	actual := New(echo, foo).Args(bar)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_ArgsMultiple(t *testing.T) {
	expected := &testCmd3
	actual := New(echo, foo).Args(bar, baz)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_Clone(t *testing.T) {
	expected := &testCmd1
	actual := New(echo, foo).Clone()

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_Ctx(t *testing.T) {
	expected := &testCmd4
	actual := New(echo, foo).Ctx(context.Background())

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_Dir(t *testing.T) {
	expected := &testCmd5
	actual := New(echo, foo).Dir(tmp)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_DryRunGlobal(t *testing.T) {
	expected := false
	cmd := New(echo, foo)
	cmd.DryRun()

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}

	actual := cmd.Complete()

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}

	dryRun = false
}

func Test_Env(t *testing.T) {
	expected := &testCmd6
	actual := New(echo, foo).Env(envFooBar)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}

func Test_RunFnCmd(t *testing.T) {
	expected := &testCmd1
	fn := New(echo, foo).RunFnCmd()

	actual, err := fn()
	if err != nil {
		t.Error(err)
	}

	msg := fmt.Sprintf("Expected %T. Got %T.", expected, actual)
	assert.IsType(t, expected, actual, msg)
}

func Test_StartFnCmd(t *testing.T) {
	expected := &testCmd1
	fn := New(echo, foo).StartFnCmd()

	actual, err := fn()
	if err != nil {
		t.Error(err)
	}

	msg := fmt.Sprintf("Expected %T. Got %T.", expected, actual)
	assert.IsType(t, expected, actual, msg)
}

func Test_clearStdOutStdErr(t *testing.T) {
	expected := &cmd{
		ctx:     context.Background(),
		env:     env,
		strings: []string{echo, foo},
	}
	actual := &cmd{
		ctx:     context.Background(),
		env:     env,
		stderr:  os.Stderr,
		stdout:  os.Stdout,
		strings: []string{echo, foo},
	}
	actual.stdout = os.Stdout
	actual.stderr = os.Stderr
	actual.clearStdOutStdErr()

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}
