package UrlCrawler

import (
	"io"
	"strings"
	"testing"
)

func TestUrlCrawler(t *testing.T) {

	config := NewConfig()

	if config.HostRequestTimeoutSeconds != DefaultHostRequestTimeoutSeconds {
		t.Log("HostRequestTimeoutSeconds is:", config.HostRequestTimeoutSeconds)
		t.Error("Must be:", DefaultHostRequestTimeoutSeconds)
	}
	if config.MaxCrawLevel != DefaultMaxCrawLevel {
		t.Log("MaxCrawLevel is:", config.MaxCrawLevel)
		t.Error("Must be:", DefaultMaxCrawLevel)
	}
	if config.MaxProcessedNames != DefaultMaxProcessedNames {
		t.Log("MaxProcessedDNames is:", config.MaxProcessedNames)
		t.Error("Must be:", DefaultMaxProcessedNames)
	}
	if config.MaxQueuedNames != DefaultMaxQueuedNames {
		t.Log("MaxQueuedDNames is:", config.MaxQueuedNames)
		t.Error("Must be:", DefaultMaxQueuedNames)
	}
	if config.MaxThreads != DefaultMaxThreads {
		t.Log("MaxThreads is:", config.MaxThreads)
		t.Error("Must be:", DefaultMaxThreads)
	}
	if config.NodeNameCacheSize != DefaultNameCacheSize {
		t.Log("NodeNameCacheSize is:", config.NodeNameCacheSize)
		t.Error("Must be:", DefaultNameCacheSize)
	}
	if config.OutputFileName != DefaultOutputFileName {
		t.Log("OutputFileName is:", config.OutputFileName)
		t.Error("Must be:", DefaultOutputFileName)
	}
	if config.OutputNewEdgeToConsole != DefaultOutputNewEdgeToConsole {
		t.Log("StartingNode is:", config.OutputNewEdgeToConsole)
		t.Error("Must be:", DefaultOutputNewEdgeToConsole)
	}
	if config.OutputNewEdgeToFile != DefaultOutputNewEdgeToFile {
		t.Log("StartingNode is:", config.OutputNewEdgeToFile)
		t.Error("Must be:", DefaultOutputNewEdgeToFile)
	}
	if config.OutputWithErrorFileName != DefaultOutputWithErrorFileName {
		t.Log("OutputWithErrorFileName is:", config.OutputWithErrorFileName)
		t.Error("Must be:", DefaultOutputWithErrorFileName)
	}
	if config.OutputNodeWithErrorToFile != DefaultOutputNodeWithErrorToFile {
		t.Log("StartingNode is:", config.OutputNodeWithErrorToFile)
		t.Error("Must be:", DefaultOutputNodeWithErrorToFile)
	}
	if config.ProcessingTimeoutSeconds != DefaultProcessingTimeoutSeconds {
		t.Log("ProcessingTimeoutSeconds is:", config.ProcessingTimeoutSeconds)
		t.Error("Must be:", DefaultProcessingTimeoutSeconds)
	}
	if config.StartingNode != DefaultStartingNode {
		t.Log("StartingNode is:", config.StartingNode)
		t.Error("Must be:", DefaultStartingNode)
	}

	config.StartingNode = "example.com"
	config.MaxCrawLevel = 1
	config.OutputNewEdgeToFile = false
	config.OutputNewEdgeToConsole = false
	config.OutputNodeWithErrorToFile = false
	config.OutputErrorsToConsole = false
	config.ProcessingTimeoutSeconds = 3
	config.StatOutputToConsole = false

	t_data := "www.example.com"
	fakeUrlFunc := func(node_name nodeName, timeout int) (*io.ReadCloser, error) {
		data := io.NopCloser(strings.NewReader("https://" + t_data))
		return &data, nil
	}
	taskWorkerUrlFunction = fakeUrlFunc

	x := NewUrlCrawler(&config)

	if x.GetNodeNameCacheSuccess() != 0 {
		t.Error("Cache success must be zero here")
	}

	err := x.ProcessDeep()
	if err != ErrMaximumCrawLevelReached {
		t.Error("ProcessDeep return:", err)
	}
}
