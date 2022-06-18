package UrlCrawler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

// test changes
func taskWorkerDefaultUrlFunction(node_name nodeName, timeout int) (*io.ReadCloser, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	resp, err := client.Get("https://" + string(node_name))
	if err != nil {
		return nil, err
	}
	return &resp.Body, err
}

var taskWorkerUrlFunction func(node_name nodeName, timeout int) (*io.ReadCloser, error) = taskWorkerDefaultUrlFunction

func taskWorker(ctx context.Context, taskData, taskDone chan nodeName, taskResult chan [2]nodeName, errCh chan errorData, uc *UrlCrawler) {
	var bufferOveflow = false
	for {
		select {
		case <-ctx.Done():
			return
		case node_name := <-taskData:
			data, err := taskWorkerUrlFunction(node_name, uc.Config.HostRequestTimeoutSeconds)
			if err != nil {
				errCh <- errorData{err, node_name}
			} else {
				buffer := bytes.Buffer{}
				for {
					b := make([]byte, 1024)
					n, err := (*data).Read(b)
					if err != nil || n == 0 {
						if !errors.Is(err, io.EOF) {
							errCh <- errorData{err, node_name}
						}
						break
					}
					if uc.Config.OutputErrorsToConsole {
						if !bufferOveflow {
							if buffer.Len()+n > uc.Config.NodeBodyReadBufferMaxBytes {
								bufferOveflow = true
								log.Println("Node body read buffer overflow at", node_name)
							}
						}
					}
					buffer.Write(b)
				}
				str := buffer.String()
				res := uc.regexpUrlCompiled.FindAllString(str, -1)
				for _, v1 := range res {
					v2 := uc.regexpHttpReplCompiled.ReplaceAllString(v1, "")
					taskResult <- [2]nodeName{node_name, nodeName(v2)}
				}

			}
			taskDone <- node_name
		}
	}
}
