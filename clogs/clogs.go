package clogs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/manifoldco/promptui"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"

	// auth needed for proxy
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type Clogser interface {
	Do() error
}

type Config struct {
	Namespace       string
	PodFilter       string
	ContainerFilter string
}

type Clogs struct {
	restConfig *rest.Config
	config     *Config
}

func NewClogs(restConfig *rest.Config, config *Config) *Clogs {

	log.WithFields(log.Fields{
		"containerFilter": config.ContainerFilter,
		"podFilter":       config.PodFilter,
		"Namespace":       config.Namespace,
	}).Debug("clogs config values...")

	return &Clogs{restConfig: restConfig, config: config}
}

func selectPod(pods []v1.Pod, config Config) (v1.Pod, error) {

	if len(pods) == 1 {
		return pods[0], nil
	}

	templates := podTemplate

	podsPrompt := promptui.Select{
		Label:     "Select Pod",
		Items:     pods,
		Templates: templates,
	}

	i, _, err := podsPrompt.Run()
	if err != nil {
		return pods[i], err
	}

	return pods[i], nil
}

func containerPrompt(containers []v1.Container, config Config) (v1.Container, error) {

	if len(containers) == 1 {
		return containers[0], nil
	}

	templates := containerTemplates

	containersPrompt := promptui.Select{
		Label:     "Select Container",
		Items:     containers,
		Templates: templates,
	}

	i, _, err := containersPrompt.Run()
	if err != nil {
		return containers[i], err
	}

	return containers[i], nil
}

func (r *Clogs) Do() error {
	client, err := kubernetes.NewForConfig(r.restConfig)
	if err != nil {
		return err
	}

	pods, err := getAllPods(client, r.config.Namespace)
	if err != nil {
		return err
	}

	filteredPods, err := r.matchPods(pods)
	if err != nil {
		return err
	}

	pod, err := selectPod(filteredPods.Items, *r.config)
	if err != nil {
		return err
	}

	containers, err := matchContainers(pod, *r.config)
	if err != nil {
		return err
	}

	container, err := containerPrompt(containers, *r.config)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"pod":       pod.GetName(),
		"container": container.Name,
		"namespace": r.config.Namespace,
	}).Info("Logs from pod...")

	str := logs(r.restConfig, pod, container)
	fmt.Print(str)

	return nil
}

func logs(restCfg *rest.Config, pod v1.Pod, container v1.Container) string {
	client, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return "error in getting config"
	}
	logOptions := corev1.PodLogOptions{}
	logOptions.Container = container.Name

	req := client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &logOptions)

	podLogs, err := req.Stream()
	if err != nil {
		log.Info(err)
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()
	return str
}
