// config provides the functions and basic structs which are necessary to perform
// dynamic DNS updates using Froxlor.
package config

import (
	"fmt"
	"strings"
)

const credentialMask = "********"

// Config holds all of the necessary configuration to run the application appropriately.
//
// In case a DNS update of the embedded updates does not have a username or a password
// the config's credentials will be used as a fallback.
type Config struct {
	Target   string     `json:"target"`
	Interval int        `json:"interval"`
	Username string     `json:"username"`
	Password credential `json:"password"`
	Updates  DNSUpdates `json:"updates"`
}

// String representation of a configuration used for logging.
func (c Config) String() string {
	return fmt.Sprintf(
		"\nTarget: %s\nInterval (seconds): %v\ndefault username: %s\ndefault password: %s \nupdates: [ \n %v]",
		c.Target,
		c.Interval,
		c.Username,
		c.Password.mask(),
		c.Updates)
}

// StdzeUpdates returns a normalized representation of a DNS update by setting the
// credentials of the parent config element if missing on an update
func (c Config) StdzeUpdates() []DNSUpdate {
	var ups []DNSUpdate
	for _, u := range c.Updates {
		up := DNSUpdate{u.Domains, u.IPv4, u.IPv6, u.Username, u.Password}
		if u.Username == nil {
			up.Username = &c.Username
		}
		if u.Password == nil {
			up.Password = &c.Password
		}
		ups = append(ups, up)
	}
	return ups
}

// checks if all standardized updates are valid otherwise the config itself is not valid
func (c Config) valid() bool {
	for _, u := range c.StdzeUpdates() {
		if !u.valid() {
			return false
		}
	}
	return true
}

// DNSUpdate represents a single DNS Update which can be passed using a single
type DNSUpdate struct {
	Domains  []string    `json:"domains"`
	IPv4     *string     `json:"ipv4"`
	IPv6     *string     `json:"ipv6"`
	Username *string     `json:"username"`
	Password *credential `json:"password"`
}

// String representation of a DNSUpdate used for logging
func (du DNSUpdate) String() string {
	return fmt.Sprintf(
		"domains: %v, IPv4: %s, IPv6: %s, username: %s, password: %s",
		strings.Join(du.Domains, ","),
		toLogStr(du.IPv4),
		toLogStr(du.IPv6),
		toLogStr(du.Username),
		du.Password.mask())
}

// IPSet returns whether one of the IPs is set for a DNS update
func (du DNSUpdate) IPSet() bool {
	return du.IPv4 != nil || du.IPv6 != nil
}

func (du DNSUpdate) valid() bool {
	return len(du.Domains) > 0 && du.Username != nil && du.Password != nil
}

// Convenience struct to ease logging
type DNSUpdates []DNSUpdate

// String representation of multiple DNSUpdates used for logging.
func (ups DNSUpdates) String() string {
	str := ""
	for _, du := range ups {
		str += "\t" + du.String() + "; \n"
	}
	return str
}

type credential string

func (c *credential) mask() string {
	if c == nil {
		return "<nil>"
	}
	cred := *c
	if len(cred) <= 4 {
		return credentialMask
	}
	return string(cred[:2] + credentialMask + cred[len(cred)-2:])
}

func toLogStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
