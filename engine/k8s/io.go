package k8s

import "io"

// SetStdin sets the stdin for the k8s engine.
func (k *k8s) SetStdin(stdin io.Reader) error {
	k.stdin = stdin
	return nil
}

// SetStdout sets the stdout for the k8s engine.
func (k *k8s) SetStdout(stdout io.Writer) error {
	k.stdout = stdout
	return nil
}

// SetStderr sets the stderr for the k8s engine.
func (k *k8s) SetStderr(stderr io.Writer) error {
	k.stderr = stderr
	return nil
}
