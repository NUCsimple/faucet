package main

import (
	"context"
	"encoding/json"
	"github.com/spf13/pflag"
	"github.com/swarm/faucet/utils"
	"io"
	"io/ioutil"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
	"log"
	"sync"
	"time"
)

const DEFAULT_TIMEOUT = time.Second * 10

type Results struct {
	sync.Mutex
	Addresses []Address
}

type Address struct {
	PodName  string `json:"pod_name"`
	Ethereum string `json:"ethereum"`
}

type Option struct {
	labelSelector string
	container     string
	command       string
	webhookUrl    string
}

var (
	opt      Option
	PodCount int
	result   Results
	wg       sync.WaitGroup
)

func main() {
	flags := pflag.NewFlagSet("main", pflag.ExitOnError)
	flags.StringVar(&opt.labelSelector, "label", "", "pod label")
	flags.StringVar(&opt.container, "container", "", "container name")
	flags.StringVar(&opt.command, "command", "", "which command be exec in pod")
	flags.StringVar(&opt.webhookUrl, "webhook", "", "webhook url")

	pflag.CommandLine = flags
	pflag.Parse()

	Run()
}

func Run() {
	cli := utils.NewKubernetesClientOutSide("/Users/carson/.kube/mulan-new")

	podList, err := cli.Clientset.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: opt.labelSelector,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range podList.Items {
		if pod.Status.Phase == "Running" {
			klog.Infof("Exec [%v] in pod %v", opt.command, pod.Name)
			wg.Add(1)
			PodCount++
			go exec(cli, pod)
		}
	}

	// wait for all pod response
	wg.Wait()

	// Check
	if PodCount != len(result.Addresses) {
		log.Fatalf("failed got all faucet addr,runnging pod %d,but got %d faucet addr.\n", PodCount, len(result.Addresses))
	}

	marshal, err := json.Marshal(result.Addresses)
	if err != nil {
		klog.Errorf("failed to encoding result,because of %v", err)
	}

	klog.Infof("Send report %v to %s", string(marshal), opt.webhookUrl)
	err = utils.SendReport(opt.webhookUrl, string(marshal))
	if err != nil {
		klog.Error(err)
	}
}

func exec(cli *utils.KubernetesClient, pod api.Pod) {
	result.Lock()
	defer result.Unlock()
	defer wg.Done()

	restClient := cli.Clientset.CoreV1().RESTClient()
	req := restClient.Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")

	cmd := make([]string, 0)
	cmd = append(cmd, "sh")
	cmd = append(cmd, "-c")
	cmd = append(cmd, opt.command)

	req.VersionedParams(&api.PodExecOptions{
		Container: opt.container,
		Command:   cmd,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(cli.RestConfig, "POST", req.URL())
	if err != nil {
		klog.Errorf("error when NewSPDYExecutor, err: %s", err)
		return
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
		return
	}

	respString := string(buffer)

	var resp Address
	klog.Infof("pod %s response is %s", pod.Name, respString)
	err = json.Unmarshal(buffer, &resp)
	if err != nil {
		klog.Warningf("unmarshal json %s got error: %s", respString, err)
		return
	}
	resp.PodName = pod.Name
	result.Addresses = append(result.Addresses, resp)

}
