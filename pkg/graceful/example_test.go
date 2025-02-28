// nolint:gosec
package graceful_test

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/graceful"
	"github.com/teamsorghum/go-common/pkg/log"
)

// This example demonstrates how to implement graceful shutdown in a web server using this package.
func Example_gracefulShutdown() {
	// Get logger
	logger := log.GetDefault()

	// Register pre-shutdown hooks that will be executed before shutdown. These hook functions will be executed in the
	// order of registration.
	graceful.RegisterPreShutdownHook(func() {
		logger.Info("Pre-shutdown hook1")
	})
	graceful.RegisterPreShutdownHook(func() {
		logger.Info("Pre-shutdown hook2")
	})

	// Register post-shutdown hooks that will be executed after shutdown. These hook functions will be executed in the
	// order of registration.
	graceful.RegisterPostShutdownHook(func() {
		logger.Info("Post-shutdown hook1")
	})
	graceful.RegisterPostShutdownHook(func() {
		logger.Info("Post-shutdown hook2")
	})

	// Create a web server.
	server := &http.Server{
		Addr:    "localhost:7788",
		Handler: nil, // Your handler here. We use nil as a temporary placeholder because we don't need to process any
		// requests in our example.
	}

	// Register shutdown function that will be executed when the process receives a kill signal.
	graceful.RegisterShutdown(time.Second, func() {
		server.Close()

		// And you can do more cleanup here.
		logger.Info("Cleaning up...")
	})

	// Before starting the server, let's launch a goroutine that will send a kill signal after 500ms. We need to launch
	// it before starting the server because server.ListenAndServe() will block further operations until the server is
	// closed.
	go func() {
		time.Sleep(time.Duration(500) * time.Millisecond)
		if err := syscall.Kill(os.Getpid(), syscall.SIGINT); err != nil {
			logger.Error("Send kill signal failed.", constant.LogAttrError, err.Error())
		}
	}()

	// Start the server.
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Listen and serve failed.", constant.LogAttrError, err.Error())
	}

	// This message will be printed if the kill signal is successfully captured.
	fmt.Println("Shutdown completed.")

	// Output: Shutdown completed.
}
