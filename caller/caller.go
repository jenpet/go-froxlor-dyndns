package caller

import (
	"fmt"
	cfg "github.com/jenpet/froxlor-dyndns/config"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type callResult struct {
	err    error
	status *int
	body   *string
}

func (cr callResult) String() string {
	return fmt.Sprintf("Status: %v, body: %v, error: %v", *cr.status, *cr.body, cr.err)
}

func (c callResult) successfull() bool {
	return c.err != nil && *c.status == http.StatusOK
}

type httpCaller struct {
	http.Client
	Target string
}

func (c httpCaller) call(update cfg.DNSUpdate, out chan<- callResult) {
	result := callResult{}
	req := c.createRequest(update)
	log.Printf("Calling target %v", req.URL)
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Printf("Failed calling target for domain(s) %v. Error: %v\n", update, err)
		result.err = err
	} else {
		b, _ := ioutil.ReadAll(resp.Body)
		body := string(b)
		log.Printf("Called target for domain(s) '%v' and received status '%v' and body: '%v'\n", update, resp.Status, body)
		result.status = &resp.StatusCode
		result.body = &body
	}
	out <- result
}

func (c httpCaller) createRequest(update cfg.DNSUpdate) *http.Request {
	req, err := http.NewRequest("GET", c.createURL(update), nil)
	if err != nil {
		log.Panicf("Failed to create request for domain(s) %v. Error: %v\n", update, err)
	}
	req.SetBasicAuth(*update.Username, string(*update.Password))
	return req
}

func (c httpCaller) createURL(update cfg.DNSUpdate) string {
	url := fmt.Sprintf("%s/froxlor/updateip.php?domain=%s", c.Target, strings.Join(update.Domains, ","))

	// if no IPs are set use auto ip detection
	if !update.IPSet() {
		return url + "&detect"
	}

	if update.IPv4 != nil {
		url = url + "&ipv4=" + *update.IPv4
	}

	if update.IPv6 != nil {
		url = url + "&ipv6" + *update.IPv6
	}
	return url
}

func defaultCaller(target string) *httpCaller {
	return &httpCaller{
		http.Client{},
		target}
}
