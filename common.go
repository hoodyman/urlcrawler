package UrlCrawler

type nodeIndex int32
type weightType int32
type nodeName string

type errorData struct {
	Error error
	Name  nodeName
}
