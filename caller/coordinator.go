package caller

import (
	cfg "github.com/jenpet/froxlor-dyndns/config"
	"log"
	"sync"
	"time"
)

// Coordinator handles all of the DNS updates towards froxlor and performs them using a caller.
// Updates will be performed timer based depending on the value set in the config file
type Coordinator struct {
	config cfg.Config
	caller *httpCaller
	ticker *time.Ticker
}

// Init sets up the Coordinator with its required values
func (c *Coordinator) Init(cfg cfg.Config) {
	c.config = cfg
	log.Printf("Initialized coordinator with configuration: %v", cfg)
	c.caller = defaultCaller(cfg.Target)
	d := time.Duration(cfg.Interval) * time.Second
	c.ticker = time.NewTicker(d)
	log.Printf("Setting ticker for update interval to %v seconds.", d)
}

// Run triggers the coordinator to run endless and perform DNS updates after every interval
func (c *Coordinator) Run() {
	quit := make(chan bool)
	go func() {
		// perform it once and then use the ticker
		c.update()
		for {
			select {
			case <-c.ticker.C:
				c.update()
			}
		}
	}()
	<-quit
}

func (c *Coordinator) update() {
	updates := c.config.StdzeUpdates()
	var wg sync.WaitGroup
	wg.Add(len(updates))

	res := make(chan callResult, len(updates))

	for _, u := range updates {
		go func(update cfg.DNSUpdate) {
			defer wg.Done()
			c.caller.call(update, res)
		}(u)
	}
	wg.Wait()
}
