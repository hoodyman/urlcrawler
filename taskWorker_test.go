package UrlCrawler

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

func TestTaskWorker(t *testing.T) {

	startingNode := "example.com"
	t_data := "www.example.com"

	config := NewConfig()

	// check normal
	fakeUrlFunc := func(node_name nodeName, timeout int) (*io.ReadCloser, error) {
		data := io.NopCloser(strings.NewReader("https://" + t_data))
		return &data, nil
	}
	taskWorkerUrlFunction = fakeUrlFunc

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	taskData := make(chan nodeName)
	taskResult := make(chan [2]nodeName)
	taskDone := make(chan nodeName)
	errCh := make(chan errorData)

	x := NewUrlCrawler(&config)

	go taskWorker(ctx, taskData, taskDone, taskResult, errCh, &x)
	taskData <- nodeName(startingNode)

	for {
		select {
		case <-ctx.Done():
			t.Error("Operation timeout")
			return
		case node := <-taskDone:
			if node != nodeName(startingNode) {
				t.Log("Result is:", node)
				t.Error("Must be:", startingNode)
			}
			goto L0
		case r := <-taskResult:
			if r[0] != nodeName(startingNode) {
				t.Log("Result is:", r[0])
				t.Error("Must be:", startingNode)
			}
			if r[1] != nodeName(t_data) {
				t.Log("Result is:", r[1])
				t.Error("Must be:", t_data)
			}
		case <-errCh:
			t.Error("Error occurred")
		}
	}

L0:

	// check error returned fakeUrlFunc
	urlError := errors.New("urlError")

	fakeUrlFunc = func(node_name nodeName, timeout int) (*io.ReadCloser, error) {
		data := io.NopCloser(strings.NewReader("https://www.example.com"))
		return &data, urlError
	}
	taskWorkerUrlFunction = fakeUrlFunc

	ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(5))
	defer cancel()

	taskData = make(chan nodeName)
	taskResult = make(chan [2]nodeName)
	taskDone = make(chan nodeName)
	errCh = make(chan errorData)

	go taskWorker(ctx, taskData, taskDone, taskResult, errCh, &x)
	taskData <- nodeName(startingNode)

	for {
		select {
		case <-ctx.Done():
			t.Error("Operation timeout")
		case node := <-taskDone:
			if node != nodeName(startingNode) {
				t.Log("Result is:", node)
				t.Error("Must be:", startingNode)
			}
			return
		case <-taskResult:
			t.Error("There is should not to be a result here")
		case err := <-errCh:
			a := errorData{urlError, nodeName(startingNode)}
			if err != a {
				t.Log("Result is:", err)
				t.Error("Must be:", a)
			}
			return
		}
	}

}
