package utils

import (
	"k8s.io/klog/v2"
	"net/http"
	"strings"
)

func SendReport(webhookUrl string, msg string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", webhookUrl, strings.NewReader(msg))
	if err != nil {
		klog.Errorf("failed to new http request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "Content-Type=application/json")

	resp, err := client.Do(req)
	if err != nil {
		klog.Errorf("failed to send http request: %v", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
