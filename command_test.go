package command

import (
	"io"
	"testing"

	"github.com/go-zoox/testify"
)

func TestCommand(t *testing.T) {
	cmd := &Command{
		Script: `echo "hello world"`,
	}

	cfg := `{
 "content": "echo \"hello world\"",
 "context": "",
 "environment": null,
 "shell": ""
}`

	testify.Equal(t, cfg, cmd.MustConfig())

	// cmd.Run()
}

func TestCommandWithContext(t *testing.T) {
	cmd := &Command{
		Script:  `echo "PWD: $PWD"`,
		Context: `/tmp`,
	}

	cfg := `{
 "content": "echo \"PWD: $PWD\"",
 "context": "/tmp",
 "environment": null,
 "shell": ""
}`

	testify.Equal(t, cfg, cmd.MustConfig())

	// cmd.Run()
}

func TestCommandWithEnvironment(t *testing.T) {
	cmd := &Command{
		Script:  `echo "PWD: $PWD"`,
		Context: `/tmp`,
		Environment: map[string]string{
			"FOO":  "BAR",
			"FOO1": "NAR1",
		},
	}

	cfg := `{
 "content": "echo \"PWD: $PWD\"",
 "context": "/tmp",
 "environment": {
  "FOO": "BAR",
  "FOO1": "NAR1"
 },
 "shell": ""
}`

	testify.Equal(t, cfg, cmd.MustConfig())

	// cmd.Run()
}

func TestCommandWithShell(t *testing.T) {
	cmd := &Command{
		Script:  `echo "PWD: $PWD"`,
		Context: `/tmp`,
		Environment: map[string]string{
			"FOO":  "BAR",
			"FOO1": "NAR1",
		},
		Shell: "/bin/bash",
	}

	cfg := `{
 "content": "echo \"PWD: $PWD\"",
 "context": "/tmp",
 "environment": {
  "FOO": "BAR",
  "FOO1": "NAR1"
 },
 "shell": "/bin/bash"
}`

	testify.Equal(t, cfg, cmd.MustConfig())

	// cmd.Run()
}

func TestCommandWithMultilines(t *testing.T) {
	output := &Output{}
	cmd := &Command{
		Script: `
echo 1
echo 2
echo 3
`,
		Stdout: output,
		Stderr: output,
	}

	testify.Assert(t, cmd.Run() == nil)
	testify.Equal(t, output.String(), "1\n2\n3\n")
}

type Output struct {
	io.Writer

	data []byte
}

func (o *Output) Write(b []byte) (n int, err error) {
	o.data = append(o.data, b...)
	return len(b), nil
}

func (o *Output) String() string {
	return string(o.data)
}
