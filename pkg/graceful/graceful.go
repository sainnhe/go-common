// Package graceful implements general graceful shutdown functions.
package graceful

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	goroutinelock "github.com/teamsorghum/go-common/pkg/goroutine_lock"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/teamsorghum/go-common/pkg/util"
)

var preShutdownHooks []func()
var postShutdownHooks []func()

// RegisterShutdown registers a function that will run when the process receives a kill signal. There is also a timeout
// time to control the maximum running time of the function. If this time is exceeded, execution will be forced to be
// interrupted. One common way of using RegisterShutdown is to register a function that stops your http server.
func RegisterShutdown(timeout time.Duration, shutdown func()) {
	go func() {
		l := log.GetDefault()
		signalCtx, signalCancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSEGV)
		defer signalCancel()
		// Wait for signals and start graceful shutdown.
		<-signalCtx.Done()
		l.Info("Graceful shutdown started.")
		startTime := time.Now()
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
		defer timeoutCancel()
		// Run shutdown function.
		shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
		go func() {
			// We must use defer to avoid panic when running shutdown() and hooks.
			defer shutdownCancel()
			defer util.Recover()
			// Run hooks and the shutdown function.
			for _, hook := range preShutdownHooks {
				hook()
			}
			shutdown()
			for _, hook := range postShutdownHooks {
				hook()
			}
		}()
		// Wait for shutdown function.
		select {
		case <-shutdownCtx.Done():
		case <-timeoutCtx.Done():
			l.Error("Shutdown times out.", "cost_time", util.ToStr(time.Since(startTime)))
			return
		}
		// Wait for goroutine locks.
		glCtx, glCancel := context.WithCancel(context.Background())
		go func() {
			defer glCancel()
			goroutinelock.Wait()
		}()
		select {
		case <-glCtx.Done():
			l.Info("Graceful shutdown finish.", "cost_time", util.ToStr(time.Since(startTime)))
		case <-timeoutCtx.Done():
			l.Error("Wait for goroutine locks times out.", "cost_time", util.ToStr(time.Since(startTime)))
		}
	}()
}

// RegisterPreShutdownHook registers a hook function that will be run before shutdown.
func RegisterPreShutdownHook(hook func()) {
	preShutdownHooks = append(preShutdownHooks, hook)
}

// RegisterPostShutdownHook register a hook function that will be run after shutdown.
func RegisterPostShutdownHook(hook func()) {
	postShutdownHooks = append(postShutdownHooks, hook)
}
