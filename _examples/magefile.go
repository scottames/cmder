//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/scottames/cmder"
	"github.com/scottames/cmder/pkg/log"
)

// Execute command returning the combined output from stdout and stderr
func Combinedoutput() error {
	out, err := cmder.New("bash", "-c", "echo foo | tee /dev/stderr").CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}

// Print the command to be executed but don't actually run it
func Dryrun1() error {
	return cmder.New("echo", "foo").DryRun().Run()
}

// Print the command to be executed but don't actually run it
// in the scope of an individual command
func Dryrun2() error {
	err := cmder.New("echo", "foo").DryRun().Run()
	if err != nil {
		return err
	}
	return cmder.New("echo", "not dryrun").Run()
}

// Print the command to be executed but don't actually run it
// in the scope of the cmd package
func Dryrun3() error {
	cmder.DryRun()
	err := cmder.New("echo", "foo").DryRun().Run()
	if err != nil {
		return err
	}
	return cmder.New("echo", "not dryrun").Run()
}

// Set environment variables
//
// Note: piping supported if passed to a shell
//   - alternatively us cmd.Output + cmd.In
func Env() error {
	return cmder.New("bash", "-c", "env | grep -i foo").Env("FOO=bar", "BAR=foo").Run()
}

// Execute a command and inspect the exit code 0
func ExitCode0() error {
	command := cmder.New("echo")
	err := command.Run()
	if err != nil {
		return err
	}
	fmt.Println("exit code: ", command.ExitCode())
	return nil
}

// Execute a command and inspect the exit code 1
func ExitCode1() {
	command := cmder.New("mkdir", "/tmp/does/not/exist/thank/you")
	// ignore error for demonstration purposes only
	command.Run()
	fmt.Println("exit code: ", command.ExitCode())
}

// Execute a command and pass a string (bytes) to command's stdin
func In() error {
	mystring := "foo"
	return cmder.New("cat").In([]byte(mystring)...).Run()
}

// Execute a command and pipe to the command from the calling process' stdin (echo foo | mage infromstdin)
func Infromstdin() error {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return fmt.Errorf("stdin expected to be passed.. use: echo foo | mage infromstdin")
	}

	return cmder.New("cat").In().Run()
}

// Set a different logger in the scope of the given command
func Logger() error {
	cmder.New("echo", "uno").Logger(log.New()).Run()
	return cmder.New("echo", "dos").Run()
}

// Set a different logger in the scope of cmd package
func Logger2() error {
	cmder.SetLogger(log.New().Key("logger2"))
	return cmder.New("echo").Run()
}

// Customize and externally use the Cmder custom logger
func Logger3() error {
	// run command w/o any modifications
	err := cmder.New("echo", "no mods").Run()
	if err != nil {
		return err
	}

	// set the log color (green is not default) at the package scope
	log.LoggerColor = log.LoggerGreen

	// create a new logger and set the key (replace 'run')
	logger := log.New().Key("foo")

	// log something
	logger.Log("global color set")

	// run another command after the global color has been set at the package scope
	err = cmder.New("echo", "foo with global color set").Run()
	if err != nil {
		return err
	}

	// log without the timestamp and a unique color for this log
	logger.Key("ERR").WithoutTimestamp().Color(log.LoggerMagenta).Log("I'm an error!")

	// set the outside logger for a specific command
	// note: not recommended to log ERR when running a command :D
	err = cmder.New("echo", "foo with logger (do not do this in production)").Logger(logger).Run()
	if err != nil {
		return err
	}

	err = cmder.New("echo", "foo without logger").Run()
	if err != nil {
		return err
	}

	// log with formatting
	logger.Logf("foo: %s", "bar")

	return nil
}

// Write output of command's stdout to file
func Outtofile() error {
	f, err := os.OpenFile("/tmp/exampleStdout.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0750)
	if err != nil {
		return err
	}
	defer f.Close()

	return cmder.New("echo", "Hello world!").Out(f).Run()
}

// Write output of command's stdout and stderr to files
func Outtofiles() error {
	fStdout, err := os.OpenFile("/tmp/exampleStdout.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0750)
	if err != nil {
		return err
	}
	defer fStdout.Close()

	fStderr, err := os.OpenFile("/tmp/exampleStderr.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0750)
	if err != nil {
		return err
	}
	defer fStderr.Close()

	return cmder.New("echo", "Hello world!").Out(fStdout, fStderr).Run()
}

// Execute a command slurping up the stdout for later use
func Output() error {
	out, err := cmder.New("echo", "foo").Output()
	if err != nil {
		return err
	}
	fmt.Println("output: ", string(out))
	return nil
}

func Fn() error {
	// define a command for reuse
	echo := cmder.New("echo").RunFn()

	// use it once
	err := echo("uno")
	if err != nil {
		return err
	}

	// use it twice
	err = echo("dos")
	if err != nil {
		return err
	}

	// use it as many times as you'd like
	return echo("tres")
}

// Execute a command
func Run(s string) error {
	return cmder.New("echo", s).Run()
}

// Execute a command without printing log (similar to appending @ to a shell string in make)
func Run2(s string) error {
	return cmder.New("echo", s).Silent().Run()
}

// Execute a command with additional arguments
func Run3(s string) error {
	return cmder.New("echo", s).
		Args("with", "additional", "args").
		Run()
}

// Execute a command passing in context
// - can be useful for timeouts
func Run4(s string) error {
	ctx := context.Background()
	return cmder.New("echo", s).Ctx(ctx).Run()
}

// Execute a command in a specified directory
func Run6(dir string) error {
	return cmder.New("ls").Dir(dir).Run()
}

// Execute a command without logging the command prior to being run
func Silent() error {
	return cmder.New("echo", "foo").Silent().Run()
}
