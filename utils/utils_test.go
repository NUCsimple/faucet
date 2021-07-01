package utils

import (
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
