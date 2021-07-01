package utils

import (
	"fmt"
	"testing"
)

func TestSendReport(t *testing.T) {
	msg := "test message"
	correctWebHook := "https://bin.webhookrelay.com/v1/webhooks/"
	testWebHook := "http://localhost:80"

	t.Run("correct webhook", func(t *testing.T) {
		err := SendReport(correctWebHook, msg)
		if err != nil {
			t.Errorf("send report error: %v", err)
		}
	})

	t.Run("wrong webhook", func(t *testing.T) {
		err := SendReport(testWebHook, msg)
		if err == nil {
			t.Errorf("send report logic is wrong")
		}
	})
}

func TestNewKubernetesClientOutSide(t *testing.T) {
	t.Run("ss", func(t *testing.T) {
		//fmt.Printf("%d %%", v)
		rate := float64(60) / float64(70) * 100
		I := fmt.Sprintf("%f %%", rate)
		fmt.Println(I)
	})
}
