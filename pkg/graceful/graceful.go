// Package graceful implements general graceful shutdown functions.
package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	goroutinelock "github.com/teamsorghum/go-common/pkg/goroutine_lock"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/teamsorghum/go-common/pkg/util"
)

var (
	preShutdownHooks     []func()
	postShutdownHooks    []func()
	hooksMutex           sync.RWMutex
	registerShutdownOnce sync.Once
)

// RegisterPreShutdownHook registers a hook function that will be run before shutdown.
func RegisterPreShutdownHook(hook func()) {
	hooksMutex.Lock()
	defer hooksMutex.Unlock()
	preShutdownHooks = append(preShutdownHooks, hook)
}

// RegisterPostShutdownHook register a hook function that will be run after shutdown.
func RegisterPostShutdownHook(hook func()) {
	hooksMutex.Lock()
	defer hooksMutex.Unlock()
	postShutdownHooks = append(postShutdownHooks, hook)
}

// RegisterShutdown registers a function that will run when the process receives a kill signal. There is also a timeout
// time to control the maximum running time of the function. If this time is exceeded, execution will be forced to be
// interrupted. One common way of using RegisterShutdown is to register a function that stops your http server.
func RegisterShutdown(timeout time.Duration, shutdown func()) {
	registerShutdownOnce.Do(func() {
		go func() {
			l := log.GetDefault()
			signalCtx, signalCancel := signal.NotifyContext(context.Background(),
				syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
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
				hooksMutex.RLock()
				for _, hook := range preShutdownHooks {
					hook()
				}
				hooksMutex.RUnlock()
				shutdown()
				hooksMutex.RLock()
				for _, hook := range postShutdownHooks {
					hook()
				}
				hooksMutex.RUnlock()
			}()
			// Wait for shutdown function.
			select {
			case <-shutdownCtx.Done():
			case <-timeoutCtx.Done():
				l.Error("Shutdown times out.", "cost_time", util.ToStr(time.Since(startTime)))
				os.Exit(1)
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
				os.Exit(1)
			}
		}()
	})
}
