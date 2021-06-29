package main

import (
	"context"
	"encoding/json"
	"github.com/spf13/pflag"
	"github.com/swarm/faucet/utils"
	"io"
	"io/ioutil"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
	"log"
	"sync"

	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Results struct {
	sync.Mutex
	Addresses []Address
}

type Address struct {
	PodName string `json:"pod_name"`
	Ethereum string `json:"ethereum"`
}

var (
	labelSelector string
	container     string
	command       string
	result        Results
	webhookUrl    string
)

func main() {
	flags := pflag.NewFlagSet("main", pflag.ExitOnError)
	flags.StringVar(&labelSelector, "label", "", "app label")
	flags.StringVar(&container, "container", "", "container name")
	flags.StringVar(&command, "command", "", "exec in pod")
	flags.StringVar(&webhookUrl, "webhook", "", "webhook Url")

	pflag.CommandLine = flags
	pflag.Parse()

	ExecCommand()
}

func ExecCommand() {
	cli := utils.NewKubernetesClient()
	wg := &sync.WaitGroup{}

	podList, err := cli.Clientset.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range podList.Items {
		wg.Add(1)
		go execInPod(cli, pod, wg)
	}
	wg.Wait()

	marshal, err := json.Marshal(result.Addresses)
	if err != nil {
		klog.Error(err)
	}

	klog.Infof("send report %v to %s", string(marshal), webhookUrl)
	err = utils.SendReport(webhookUrl, string(marshal))
	if err != nil {
		klog.Error(err)
	}
}

func execInPod(cli *utils.KubernetesClient, pod api.Pod, wg *sync.WaitGroup) {
Exec:
	result.Lock()
	defer result.Unlock()

	if pod.Status.Phase != api.PodRunning {
		klog.Errorf("Pod %s status is %s", pod.Name, pod.Status.Phase)
		goto Exec
	}
	restClient := cli.Clientset.CoreV1().RESTClient()
	req := restClient.Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")

	cmd := make([]string, 0)
	cmd = append(cmd, "sh")
	cmd = append(cmd, "-c")
	cmd = append(cmd, command)

	req.VersionedParams(&api.PodExecOptions{
		Container: container,
		Command:   cmd,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(cli.RestConfig, "POST", req.URL())
	if err != nil {
		klog.Errorf("error when NewSPDYExecutor, err: %s", err)
		goto Exec
	}

	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		err = exec.Stream(remotecommand.StreamOptions{
			Stdout: writer,
			Stderr: writer,
			Tty:    false,
		})
	}()

	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		klog.Warningf("read resp got error: %s", err)
		goto Exec
	}

	respString := string(buffer)

	var resp Address
	err = json.Unmarshal(buffer, &resp)
	if err != nil {
		klog.Warningf("unmarshal json %s got error: %s", respString, err)
		goto Exec
	}
	resp.PodName = pod.Name
	result.Addresses = append(result.Addresses, resp)
	wg.Done()
}
