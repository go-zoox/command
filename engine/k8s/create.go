package k8s

import (
	"context"
	"fmt"
	"os"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// create creates the Job (and stores clientset/config for later use).
func (k *k8s) create() error {
	var restConfig *rest.Config
	var err error

	if k.cfg.Kubeconfig != "" {
		restConfig, err = clientcmd.BuildConfigFromFlags("", k.cfg.Kubeconfig)
	} else {
		restConfig, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	}
	if err != nil {
		restConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return fmt.Errorf("k8s: build config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("k8s: create clientset: %w", err)
	}

	k.clientset = clientset
	k.restConfig = restConfig

	jobName := k.cfg.ID
	if jobName == "" {
		jobName = "go-zoox-command"
	}
	// Ensure valid DNS label (lowercase, alphanumeric, hyphen)
	jobName = strings.ToLower(strings.ReplaceAll(jobName, "_", "-"))
	if len(jobName) > 52 {
		jobName = jobName[:52]
	}

	envVars := []corev1.EnvVar{}
	for key, val := range k.cfg.Environment {
		envVars = append(envVars, corev1.EnvVar{Name: key, Value: val})
	}
	for _, key := range k.cfg.AllowedSystemEnvKeys {
		if val, ok := os.LookupEnv(key); ok {
			envVars = append(envVars, corev1.EnvVar{Name: key, Value: val})
		}
	}

	args := []string{"-c", k.cfg.Command}
	if k.cfg.Command == "" {
		args = []string{"-c", "sleep 0"}
	}

	backoffLimit := int32(0)
	activeDeadlineSeconds := k.cfg.JobTimeoutSeconds
	if activeDeadlineSeconds <= 0 {
		activeDeadlineSeconds = 3600
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: k.cfg.Namespace,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            &backoffLimit,
			ActiveDeadlineSeconds:   &activeDeadlineSeconds,
			TTLSecondsAfterFinished: ptr(int32(300)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:       "cmd",
							Image:      k.cfg.Image,
							Command:    []string{k.cfg.Shell},
							Args:       args,
							Env:        envVars,
							WorkingDir: k.cfg.WorkDir,
							Stdin:      true,
							StdinOnce:  true,
							TTY:        true,
						},
					},
				},
			},
		},
	}

	_, err = k.clientset.BatchV1().Jobs(k.cfg.Namespace).Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("k8s: create job: %w", err)
	}

	k.jobName = jobName
	k.jobNamespace = k.cfg.Namespace
	return nil
}

func ptr(i int32) *int32 { return &i }
