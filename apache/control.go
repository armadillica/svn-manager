// Package apache interfaces between Apache and SVNMan.
package apache

import (
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Control models the interface for controlling Apache.
type Control struct {
	mutex         sync.Mutex
	restartQueued bool
	timer         *time.Timer
	restartDelay  time.Duration
}

// Restarter describes the interface for delayed-restarting of Apache.
type Restarter interface {
	QueueRestart()
	Flush()
	PerformRestart()
}

func apachectl(subcmd string) (string, error) {
	deadline := time.Now().Add(10 * time.Second)
	ctx, cancelFunc := context.WithDeadline(context.Background(), deadline)
	defer cancelFunc()

	cmd := exec.CommandContext(ctx, "sudo", "--non-interactive", "apache2ctl", subcmd)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Check that we can run 'sudo apache2ctl configtest' successfully.
func testApachectl() {
	log.Info("testing Apache configuration")
	output, err := apachectl("configtest")
	if err != nil {
		out := strings.TrimSpace(output)
		log.WithField("output", out).WithError(err).Fatal("error running sudo apache2ctl configtest")
	}
}

// CreateControl creates a new Control object.
func CreateControl(restartDelay time.Duration) *Control {
	testApachectl()

	return &Control{
		restartDelay: restartDelay,
	}
}

// QueueRestart queues a graceful restart that'll take place in a few seconds.
// Any call to QueueRestart during that time is a no-op. This prevents Apache
// from being restarted too often.
func (am *Control) QueueRestart() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if am.restartQueued {
		log.Debug("Apache restart already queued")
		return
	}

	log.WithField("delay", am.restartDelay).Info("queueing Apache graceful restart")
	am.timer = time.AfterFunc(am.restartDelay, am.PerformRestart)
	am.restartQueued = true
}

// PerformRestart performs an immediate Apache restart
func (am *Control) PerformRestart() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	log.Info("performing graceful Apache restart")
	if am.timer != nil {
		am.timer.Stop()
		am.timer = nil
	}
	am.restartQueued = false

	output, err := apachectl("graceful")
	if err != nil {
		log.WithField("output", output).WithError(err).Error("error running apache2ctl")
		return
	}
	log.Info("Apache gracefully restarted")
}

// Flush performs a scheduled restart immediately and then returns.
func (am *Control) Flush() {
	am.mutex.Lock()

	if !am.restartQueued {
		am.mutex.Unlock()
		log.Debug("Apache restart not queued")
		return
	}

	log.Info("flushing Apache restart")
	am.mutex.Unlock()
	am.PerformRestart()
}
