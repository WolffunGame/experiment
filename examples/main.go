// to run the CPU profiling: go build -ldflags "-X main.RunCPUProfile=true" main.go && ./main
// to run the Mem profiling: go build -ldflags "-X main.RunMemProfile=true" main.go && ./main

package main

import (
	"fmt"
	"time"

	"github.com/WolffunService/experiment"
	"github.com/WolffunService/experiment/pkg/client"
	"github.com/WolffunService/experiment/pkg/config"
	"github.com/WolffunService/experiment/pkg/event"
	"github.com/WolffunService/experiment/pkg/logging"
)

func main() {
	sdkKey := "4SLpaJA1r1pgE6T2CoMs9q"
	logging.SetLogLevel(logging.LogLevelDebug)

	user := optimizely.UserContext(
		"mike ng",
		map[string]interface{}{
			"country":      "Unknown",
			"likes_donuts": true,
		},
	)

	/************* Bad SDK Key  ********************/

	optimizelyClient, err := optimizely.Client("some_key")
	enabled, err := optimizelyClient.IsFeatureEnabled("mutext_feat", user)
	if err == config.Err403Forbidden {
		fmt.Println("A Valid 403 error received:", config.Err403Forbidden)
	}

	/************* Simple usage ********************/

	optimizelyClient, err = optimizely.Client(sdkKey)
	enabled, _ = optimizelyClient.IsFeatureEnabled("mutext_feat", user)

	fmt.Printf("Is feature enabled? %v\n", enabled)

	/************* StaticClient ********************/

	optimizelyFactory := &client.OptimizelyFactory{
		Datafile: []byte("DATAFILE_JSON_STRING_HERE"),
		SDKKey:   sdkKey,
	}

	optimizelyClient, err = optimizelyFactory.StaticClient()

	if err != nil {
		fmt.Printf("Error instantiating client: %s", err)
		return
	}

	enabled, _ = optimizelyClient.IsFeatureEnabled("mutext_feat", user)
	fmt.Printf("Is feature enabled? %v\n", enabled)

	fmt.Println()
	optimizelyClient.Close() //  user can close dispatcher
	fmt.Println()

	/************* Client ********************/

	optimizelyFactory = &client.OptimizelyFactory{
		SDKKey: sdkKey,
	}

	optimizelyClient, err = optimizelyFactory.Client()

	if err != nil {
		fmt.Printf("Error instantiating client: %s", err)
		return
	}

	enabled, _ = optimizelyClient.IsFeatureEnabled("mutext_feat", user)
	fmt.Printf("Is feature enabled? %v\n", enabled)
	optimizelyClient.Close() //  user can close dispatcher

	/************* Setting Polling Interval ********************/

	optimizelyClient, _ = optimizelyFactory.Client(
		client.WithPollingConfigManager(time.Second, nil),
		client.WithBatchEventProcessor(event.DefaultBatchSize, event.DefaultEventQueueSize, event.DefaultEventFlushInterval),
	)

	optimizelyClient.Close()
}
