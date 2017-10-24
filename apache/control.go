// Package apache interfaces between Apache and SVNMan.
package apache

import log "github.com/sirupsen/logrus"

// Control models the interface for controlling Apache.
type Control struct{}

// Restarter describes the interface for delayed-restarting of Apache.
type Restarter interface {
	QueueRestart()
}

// QueueRestart queues a graceful restart that'll take place in a few seconds.
// Any call to QueueRestart during that time is a no-op. This prevents Apache
// from being restarted too often.
func (am *Control) QueueRestart() {
	log.Error("QueueRestart() not implemented")
}
