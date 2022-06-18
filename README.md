urlCrawler

```
func main() {
	config := UrlCrawler.NewConfig()
	config.StartingNode = "github.com"
	config.OutputNewEdgeToConsole = false
	config.OutputErrorsToConsole = true
	config.OutputNodeWithErrorToConsole = false
	config.MaxNodeEdges = 10

	x := UrlCrawler.NewUrlCrawler(&config)
	x.ProcessDeep()
	x.OutputResultAsGraphVizDotFile("output.dot")
}
```

output:

```
graph G {
	"github.com" -- "www.w3.org" [weight=1]
	"github.com" -- "avatars.githubusercontent.com" [weight=1]
	"github.com" -- "api.github.com" [weight=1]
	"github.com" -- "skills.github.com" [weight=1]
	"github.com" -- "stars.github.com" [weight=1]
	"github.com" -- "education.github.com" [weight=1]
	"github.com" -- "github.githubassets.com" [weight=1]
	"github.com" -- "collector.github.com" [weight=1]
	"github.com" -- "opensource.guide" [weight=1]
	"github.com" -- "github.community" [weight=1]
}
```
