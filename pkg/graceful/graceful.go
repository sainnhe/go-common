// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

// Package graceful implements general graceful shutdown functions.
//
// The idea of graceful shutdown is that when a kill signal like [syscall.SIGINT] is received, instead of exiting
// directly, the program will perform a custom cleanup process to release resources.
//
// This package provides 3 functions to complete this task:
//
//   - [RegisterShutdown]: Registers a custom shutdown function that will be executed when a kill signal is received.
//   - [RegisterPreShutdownHook]: Register a hook that will be run before shutdown.
//   - [RegisterPostShutdownHook]: Register a hook that will be run after shutdown.
//
// The registered hook functions will be executed in the order of registration.
package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/teamsorghum/go-common/pkg/glock"
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
	if hook == nil {
		return
	}
	hooksMutex.Lock()
	defer hooksMutex.Unlock()
	preShutdownHooks = append(preShutdownHooks, hook)
}

// RegisterPostShutdownHook register a hook function that will be run after shutdown.
func RegisterPostShutdownHook(hook func()) {
	if hook == nil {
		return
	}
	hooksMutex.Lock()
	defer hooksMutex.Unlock()
	postShutdownHooks = append(postShutdownHooks, hook)
}

// RegisterShutdown registers a function that will run when the process receives a kill signal. To be precise, these
// signals include [syscall.SIGINT], [syscall.SIGTERM] and [syscall.SIGQUIT].
//
// There is also a timeout time to control the maximum running time of the function. If this time is exceeded, execution
// will be forced to be interrupted.
//
// One common way of using RegisterShutdown is to register a function that stops your server.
//
// NOTE: The shutdown process will wait for goroutine locks implemented in [glock] to be released, and the waiting time
// respects the timeout argument.
func RegisterShutdown(timeout time.Duration, shutdown func()) {
	if shutdown == nil {
		return
	}
	registerShutdownOnce.Do(func() {
		go func() {
			l := log.Global()
			signalCtx, signalCancel := signal.NotifyContext(context.Background(),
				syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			defer signalCancel()

			// Wait for signals and start graceful shutdown.
			<-signalCtx.Done()
			l.Info("[shutdown] Graceful shutdown started.")
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
				l.Error("[shutdown] Shutdown times out.", "cost", util.ToStr(time.Since(startTime)))
				os.Exit(1)
			}

			// Wait for goroutine locks.
			glCtx, glCancel := context.WithCancel(context.Background())
			go func() {
				defer glCancel()
				glock.Wait()
			}()
			select {
			case <-glCtx.Done():
				l.Info("[shutdown] Graceful shutdown finish.", "cost", util.ToStr(time.Since(startTime)))
			case <-timeoutCtx.Done():
				l.Error("[shutdown] Wait for goroutine locks times out.", "cost", util.ToStr(time.Since(startTime)))
				os.Exit(1)
			}
		}()
	})
}
