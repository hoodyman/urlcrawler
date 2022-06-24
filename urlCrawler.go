package UrlCrawler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

const (
	DefaultStartingNode                 = "github.com"
	DefaultMaxThreads                   = 100
	DefaultMaxCrawLevel                 = 1000
	DefaultProcessingTimeoutSeconds     = 60
	DefaultHostRequestTimeoutSeconds    = 5
	DefaultOutputErrorsToConsole        = true
	DefaultOutputFileName               = "output.dot"
	DefaultOutputNewEdgeToConsole       = true
	DefaultOutputNewEdgeToFile          = false
	DefaultOutputWithErrorFileName      = "outputWithError.txt"
	DefaultOutputNodeWithErrorToFile    = false
	DefaultOutputNodeWithErrorToConsole = true
	DefaultMaxQueuedNames               = 100000
	DefaultMaxProcessedNames            = 100000
	DefaultMaxNodeEdges                 = 20000
	DefaultNameCacheSize                = 100
	DefaultNodeBodyReadBufferMaxBytes   = 1024 * 1024 * 100
	DefaultRegexpUrlCompiled            = `https?://(\w+\.)+\w+`
	DefaultRegexpHttpReplCompiled       = `https?://`
	DefaultStatOutputToConsole          = true
)

var (
	ErrMaximumProcessedNamesReached = errors.New("maximum processed urls reached")
	ErrMaximumQueuedNamesReached    = errors.New("maximum queued urls reached")
	ErrProcessTimedOut              = errors.New("process timed out")
	ErrMaximumCrawLevelReached      = errors.New("maximum craw level reached")
	ErrMaximumNodeEdgesReached      = errors.New("maximum node edges reached")
	ErrNodeNameNotFound             = errors.New("no such node idx")
)

type Config struct {
	StartingNode                 string
	MaxThreads                   int
	MaxCrawLevel                 int
	ProcessingTimeoutSeconds     int
	HostRequestTimeoutSeconds    int
	OutputErrorsToConsole        bool
	OutputFileName               string
	OutputNewEdgeToConsole       bool
	OutputNewEdgeToFile          bool
	OutputWithErrorFileName      string
	OutputNodeWithErrorToFile    bool
	OutputNodeWithErrorToConsole bool
	MaxQueuedNames               int
	MaxProcessedNames            int
	MaxNodeEdges                 int
	NodeNameCacheSize            int
	NodeBodyReadBufferMaxBytes   int
	RegexpUrl                    string
	RegexpHttpRepl               string
	StatOutputToConsole          bool
}

func NewConfig() Config {
	return Config{
		StartingNode:                 DefaultStartingNode,
		MaxThreads:                   DefaultMaxThreads,
		MaxCrawLevel:                 DefaultMaxCrawLevel,
		ProcessingTimeoutSeconds:     DefaultProcessingTimeoutSeconds,
		HostRequestTimeoutSeconds:    DefaultHostRequestTimeoutSeconds,
		OutputErrorsToConsole:        DefaultOutputErrorsToConsole,
		OutputFileName:               DefaultOutputFileName,
		OutputNewEdgeToConsole:       DefaultOutputNewEdgeToConsole,
		OutputNewEdgeToFile:          DefaultOutputNewEdgeToFile,
		OutputWithErrorFileName:      DefaultOutputWithErrorFileName,
		OutputNodeWithErrorToFile:    DefaultOutputNodeWithErrorToFile,
		OutputNodeWithErrorToConsole: DefaultOutputNodeWithErrorToConsole,
		MaxQueuedNames:               DefaultMaxQueuedNames,
		MaxProcessedNames:            DefaultMaxProcessedNames,
		MaxNodeEdges:                 DefaultMaxNodeEdges,
		NodeNameCacheSize:            DefaultNameCacheSize,
		NodeBodyReadBufferMaxBytes:   DefaultNodeBodyReadBufferMaxBytes,
		RegexpUrl:                    DefaultRegexpUrlCompiled,
		RegexpHttpRepl:               DefaultRegexpHttpReplCompiled,
		StatOutputToConsole:          DefaultStatOutputToConsole,
	}
}

type UrlCrawler struct {
	name_mapper            *nameMapper
	NodeEdgeArray          nodeEdgeMap
	regexpUrlCompiled      *regexp.Regexp
	regexpHttpReplCompiled *regexp.Regexp
	Config                 *Config
}

func NewUrlCrawler(config *Config) UrlCrawler {
	uc := UrlCrawler{}
	if config != nil {
		uc.Config = &Config{}
		*uc.Config = *config
	} else {
		newConfig := NewConfig()
		uc.Config = &newConfig
	}
	uc.regexpUrlCompiled = regexp.MustCompile(config.RegexpUrl)
	uc.regexpHttpReplCompiled = regexp.MustCompile(config.RegexpHttpRepl)
	return uc
}

func (y *UrlCrawler) ProcessDeep() error {

	y.name_mapper = newNameMapper(y.Config.NodeNameCacheSize)
	y.NodeEdgeArray = nodeEdgeMap{}

	inQueue := newNodeMap(y.name_mapper)
	inProcess := newNodeMap(y.name_mapper)
	toProcess := newNodeMap(y.name_mapper)
	processed := newNodeMap(y.name_mapper)

	inQueue.append(nodeName(y.Config.StartingNode))

	var outputFileHandle *os.File
	var outputWithErrorFileHandle *os.File

	if y.Config.OutputNewEdgeToFile {
		var err error
		outputFileHandle, err = os.Create(y.Config.OutputFileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	if y.Config.OutputNodeWithErrorToFile {
		var err error
		outputWithErrorFileHandle, err = os.Create(y.Config.OutputWithErrorFileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*time.Duration(y.Config.ProcessingTimeoutSeconds))
	defer cancelFunc()

	taskData := make(chan nodeName)
	taskResult := make(chan [2]nodeName)
	taskDone := make(chan nodeName)
	errCh := make(chan errorData)

	for i := 0; i < y.Config.MaxThreads; i++ {
		go taskWorker(ctx, taskData, taskDone, taskResult, errCh, y)
	}

	crawLevel := 0

	if y.Config.StatOutputToConsole {
		go func() {
			var lastNodesVal int
			var lastEdgesVal int
			prn := func() {
				curNodesVal := y.name_mapper.len()
				curEdgesVal := len(y.NodeEdgeArray)
				if curNodesVal != lastNodesVal || curEdgesVal != lastEdgesVal {
					fmt.Printf("%v nodes, %v edges\n", curNodesVal, curEdgesVal)
					lastNodesVal = curNodesVal
					lastEdgesVal = curEdgesVal
				}
			}
			tick := func() {
				ticker := time.NewTicker(time.Second)
				for {
					select {
					case <-ctx.Done():
						prn()
						return
					case <-ticker.C:
						prn()
					}
				}
			}
			tick()
		}()
	}

	for {
		select {
		case <-ctx.Done():
			log.Println(ErrProcessTimedOut)
			return ErrProcessTimedOut
		case name := <-taskResult:
			if !inProcess.isExist(name[1]) && !processed.isExist(name[1]) && !toProcess.isExist(name[1]) {
				toProcess.append(name[1])
				if len(toProcess.m) == y.Config.MaxQueuedNames {
					if y.Config.OutputErrorsToConsole {
						log.Println(ErrMaximumQueuedNamesReached)
					}
					return ErrMaximumQueuedNamesReached
				}
				isIn := y.NodeEdgeArray.append(newNodeEdge(name[0], name[1], y.name_mapper))
				if !isIn {
					if y.Config.OutputNewEdgeToConsole {
						fmt.Println(name[0], "--", name[1])
					}
					if y.Config.OutputNewEdgeToFile {
						_, err := outputFileHandle.WriteString(fmt.Sprintf("\t\"%v\" %v \"%v\";\n", name[0], "--", name[1]))
						if err != nil {
							log.Fatal(err)
						}
					}
				}
				if y.NodeEdgeArray.len() == y.Config.MaxNodeEdges {
					log.Println(ErrMaximumNodeEdgesReached.Error())
					return ErrMaximumNodeEdgesReached
				}
			}
		case name := <-taskDone:
			processed.append(name)
			inProcess.delete(name)
			if len(processed.m) == y.Config.MaxProcessedNames {
				if y.Config.OutputErrorsToConsole {
					log.Println(ErrMaximumProcessedNamesReached)
				}
				return ErrMaximumProcessedNamesReached
			}
		case err := <-errCh:
			processed.append(err.Name)
			inProcess.delete(err.Name)
			if len(processed.m) == y.Config.MaxProcessedNames {
				if y.Config.OutputErrorsToConsole {
					log.Println(ErrMaximumProcessedNamesReached)
				}
				return ErrMaximumProcessedNamesReached
			}
			if y.Config.OutputNodeWithErrorToConsole {
				log.Println(err)
			}
			if y.Config.OutputNodeWithErrorToFile {
				_, err2 := outputWithErrorFileHandle.WriteString(fmt.Sprintln(err.Name, ":", err.Error.Error()))
				if err2 != nil {
					log.Fatal(err2)
				}
			}
		default:
			if name, ok := inQueue.getFirst(); ok {
				select {
				case taskData <- name:
					inProcess.append(name)
					inQueue.delete(name)
				default:
				}
			} else if len(inProcess.m) > 0 {
			} else if len(toProcess.m) > 0 {
				crawLevel++
				if crawLevel == y.Config.MaxCrawLevel {
					if y.Config.OutputErrorsToConsole {
						log.Println(ErrMaximumCrawLevelReached)
					}
					return ErrMaximumCrawLevelReached
				}
				inQueue, toProcess = toProcess, inQueue
			} else {
				return nil
			}
		}
	}
}

func (y *UrlCrawler) GetNodeNameCacheSuccess() float64 {
	if y.name_mapper != nil {
		return y.name_mapper.getCacheSuccess()
	}
	return 0
}

// for k, v := range x.NodeEdgeArray {
// 	t, err := x.TranslateDataItem(k)
// 	if err != nil {
// 		log.Println("Error:", k)
// 	}
// 	fmt.Println(t, ":", v)
// }
func (y *UrlCrawler) TranslateDataItem(node_edge nodeEdge) ([2]string, error) {
	nA, ok := y.name_mapper.getName(node_edge.NodeA)
	if !ok {
		return [2]string{}, ErrNodeNameNotFound
	}
	nB, ok := y.name_mapper.getName(node_edge.NodeB)
	if !ok {
		return [2]string{}, ErrNodeNameNotFound
	}
	return [2]string{string(nA), string(nB)}, nil
}

func (y *UrlCrawler) OutputResultAsGraphVizDotFile(file_name string) error {

	f, err := os.Create(file_name)
	if err != nil {
		return err
	}

	f.WriteString("graph G {\n")
	defer f.WriteString("}\n")

	for k, v := range y.NodeEdgeArray {
		t, _ := y.TranslateDataItem(k)
		_, err = f.WriteString(fmt.Sprintf("\t\"%v\" -- \"%v\" [weight=%v]\n", t[0], t[1], v))
		if err != nil {
			return err
		}
	}

	return nil
}

func (y *UrlCrawler) OutputResultAsDLFile(file_name string) error {

	f, err := os.Create(file_name)
	if err != nil {
		return err
	}

	f.WriteString("dl\nformat = edgelist1\n")
	f.WriteString(fmt.Sprintf("n = %v\n", y.name_mapper.len()))
	f.WriteString("labels embedded:\ndata\n")

	for k, v := range y.NodeEdgeArray {
		t, _ := y.TranslateDataItem(k)
		_, err = f.WriteString(fmt.Sprintf("%v %v %.1f\n", t[0], t[1], float64(v)))
		if err != nil {
			return err
		}
	}

	return nil
}
