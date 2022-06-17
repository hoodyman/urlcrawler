urlCrawler

```
func main() {
	config := UrlCrawler.NewConfig()
	config.StartingNode = "github.com"
	config.OutputNewEdgeToConsole = false
	config.OutputErrorsToConsole = true
	config.OutputNodeWithErrorToConsole = false

	x := UrlCrawler.NewUrlCrawler(&config)
	x.ProcessDeep()
	x.OutputResultAsGraphVizDotFile("output.dot")
}
```
