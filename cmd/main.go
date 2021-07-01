package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/swarm/faucet/options"
	"github.com/swarm/faucet/utils"
	"io"
	"io/ioutil"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

const (
	// DefaultTimeout TODO [add goroutine timeout logic]
	DefaultTimeout       = time.Second * 10
	DefaultRetryInterval = 5 * time.Second
	DefaultRetryTimes    = 5
)

type Result struct {
	sync.Mutex
	Report
}

type Report struct {
	Message       string   `json:"message"`
	SuccessRate   string   `json:"success_rate"`
	FailedPodList []string `json:"failed_pod_list"`
	Addresses     []Address
}

type Address struct {
	PodName  string `json:"pod_name"`
	Ethereum string `json:"ethereum"`
}

var (
	opt      options.Option
	PodCount int
	result   Result
	wg       sync.WaitGroup
)

func main() {
	flags := pflag.NewFlagSet("main", pflag.ExitOnError)
	flags.StringVar(&opt.LabelSelector, "label", "", "Which label for pod.")
	flags.StringVar(&opt.Container, "container", "", "The name of the container you want to enter.")
	flags.StringVar(&opt.Command, "command", "", "which command be execute in pod.")
	flags.StringVar(&opt.WebhookUrl, "webhook", "", "webhook url")

	pflag.CommandLine = flags
	pflag.Parse()

	for i := 1; i <= DefaultRetryTimes; i++ {
		if err := Run(); err != nil {
			klog.Error(err)
			time.Sleep(DefaultRetryInterval)
			continue
		} else {
			break
		}
	}
}

func Run() error {
	cli := utils.NewKubernetesClient()
	//cli := utils.NewKubernetesClientOutSide("/Users/carson/.kube/config")
	podList, err := cli.Clientset.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: opt.LabelSelector,
	})
	if err != nil {
		return err
	}

	for _, pod := range podList.Items {
		if pod.Status.Phase == "Running" {
			klog.Infof("Exec [%v] in pod %v", opt.Command, pod.Name)
			wg.Add(1)
			PodCount++
			go exec(cli, pod)
		}
	}

	// wait for all pod response
	wg.Wait()

	// Validate result
	if PodCount != len(result.Addresses) {
		klog.Errorf("failed got all faucet addr,runnging pod %d,but got %d faucet addr.\n", PodCount, len(result.Addresses))
		result.Report.Message = "failed got all faucet addr"
		rate := float64(len(result.Addresses)) / float64(PodCount) * 100
		result.Report.SuccessRate = fmt.Sprintf("%f %%", rate)
	} else {
		result.Report.Message = "OK"
		result.Report.SuccessRate = fmt.Sprintf("%d %%", 100)
	}

	marshal, err := json.Marshal(result.Report)
	if err != nil {
		klog.Errorf("failed to encoding result,because of %v", err)
	}

	klog.Infof("Send report %v to %s", string(marshal), opt.WebhookUrl)
	err = utils.SendReport(opt.WebhookUrl, string(marshal))
	if err != nil {
		klog.Error(err)
	}

	return nil
}

// exec is used by get faucet address
func exec(cli *utils.KubernetesClient, pod api.Pod) {
	// 临界区太大
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
	cmd = append(cmd, opt.Command)

	req.VersionedParams(&api.PodExecOptions{
		Container: opt.Container,
		Command:   cmd,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(cli.RestConfig, "POST", req.URL())
	if err != nil {
		klog.Errorf("error when NewSPDYExecutor, err: %s", err)
		result.FailedPodList = append(result.FailedPodList, pod.Name)
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
		result.FailedPodList = append(result.FailedPodList, pod.Name)
		return
	}

	respString := string(buffer)

	var resp Address
	err = json.Unmarshal(buffer, &resp)
	if err != nil {
		klog.Warningf("unmarshal json %s got error: %s", respString, err)
		result.FailedPodList = append(result.FailedPodList, pod.Name)
		return
	}
	resp.PodName = pod.Name
	result.Addresses = append(result.Addresses, resp)
}
