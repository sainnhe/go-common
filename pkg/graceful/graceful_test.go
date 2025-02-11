// nolint:paralleltest
package graceful_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/teamsorghum/go-common/pkg/graceful"
)

func TestRegisterShutdown(t *testing.T) {
	var preShutdownHookCalled bool
	var shutdownFuncCalled bool
	var postShutdownHookCalled bool

	// Register pre-shutdown hook
	graceful.RegisterPreShutdownHook(func() {
		preShutdownHookCalled = true
	})

	// Register post-shutdown hook
	graceful.RegisterPostShutdownHook(func() {
		postShutdownHookCalled = true
	})

	// Register shutdown function
	graceful.RegisterShutdown(5*time.Second, func() {
		shutdownFuncCalled = true
	})

	// Wait to ensure that the signal handling goroutine is set up
	time.Sleep(100 * time.Millisecond)

	// Send a SIGTERM signal to the process
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		t.Fatalf("Failed to send SIGTERM: %v", err)
	}

	// Wait for the shutdown function and hooks to be called
	time.Sleep(1 * time.Second)

	if !preShutdownHookCalled {
		t.Errorf("Pre-shutdown hook was not called")
	}

	if !shutdownFuncCalled {
		t.Errorf("Shutdown function was not called")
	}

	if !postShutdownHookCalled {
		t.Errorf("Post-shutdown hook was not called")
	}
}

func TestRegisterShutdownOrder(t *testing.T) {
	var callOrder []string

	// Register pre-shutdown hook
	graceful.RegisterPreShutdownHook(func() {
		callOrder = append(callOrder, "preShutdownHook")
	})

	// Register shutdown function
	graceful.RegisterShutdown(5*time.Second, func() {
		callOrder = append(callOrder, "shutdownFunc")
	})

	// Register post-shutdown hook
	graceful.RegisterPostShutdownHook(func() {
		callOrder = append(callOrder, "postShutdownHook")
	})

	// Wait to ensure that the signal handling goroutine is set up
	time.Sleep(100 * time.Millisecond)

	// Send a SIGTERM signal to the process
	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	if err != nil {
		t.Fatalf("Failed to send SIGTERM: %v", err)
	}

	// Wait for the shutdown function and hooks to be called
	time.Sleep(1 * time.Second)

	// Expected order
	expectedOrder := []string{"preShutdownHook", "shutdownFunc", "postShutdownHook"}

	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Number of calls does not match expected: got %v, expected %v", len(callOrder), len(expectedOrder))
	}

	for i, call := range callOrder {
		if call != expectedOrder[i] {
			t.Errorf("Call order mismatch at position %d: got %s, expected %s", i, call, expectedOrder[i])
		}
	}
}

func TestRegisterMultiplePreShutdownHooks(t *testing.T) {
	var callOrder []string

	// Register multiple pre-shutdown hooks
	graceful.RegisterPreShutdownHook(func() {
		callOrder = append(callOrder, "preShutdownHook1")
	})

	graceful.RegisterPreShutdownHook(func() {
		callOrder = append(callOrder, "preShutdownHook2")
	})

	// Register shutdown function
	graceful.RegisterShutdown(5*time.Second, func() {
		callOrder = append(callOrder, "shutdownFunc")
	})

	// Register post-shutdown hook
	graceful.RegisterPostShutdownHook(func() {
		callOrder = append(callOrder, "postShutdownHook")
	})

	// Wait to ensure that the signal handling goroutine is set up
	time.Sleep(100 * time.Millisecond)

	// Send a SIGTERM signal to the process
	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	if err != nil {
		t.Fatalf("Failed to send SIGTERM: %v", err)
	}

	// Wait for the shutdown function and hooks to be called
	time.Sleep(1 * time.Second)

	// Expected order
	expectedOrder := []string{"preShutdownHook1", "preShutdownHook2", "shutdownFunc", "postShutdownHook"}

	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Number of calls does not match expected: got %v, expected %v", len(callOrder), len(expectedOrder))
	}

	for i, call := range callOrder {
		if call != expectedOrder[i] {
			t.Errorf("Call order mismatch at position %d: got %s, expected %s", i, call, expectedOrder[i])
		}
	}
}

func TestRegisterShutdownTimeout(t *testing.T) {
	var preShutdownHookCalled bool
	var shutdownFuncCalled bool
	var postShutdownHookCalled bool

	// Register pre-shutdown hook
	graceful.RegisterPreShutdownHook(func() {
		preShutdownHookCalled = true
	})

	// Register post-shutdown hook
	graceful.RegisterPostShutdownHook(func() {
		postShutdownHookCalled = true
	})

	// Register shutdown function that takes longer than timeout
	graceful.RegisterShutdown(1*time.Second, func() {
		shutdownFuncCalled = true
		// Simulate a long-running shutdown function
		time.Sleep(2 * time.Second)
	})

	// Wait to ensure that the signal handling goroutine is set up
	time.Sleep(100 * time.Millisecond)

	// Send a SIGTERM signal to the process
	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	if err != nil {
		t.Fatalf("Failed to send SIGTERM: %v", err)
	}

	// Wait enough time for the timeout to occur
	time.Sleep(3 * time.Second)

	if !preShutdownHookCalled {
		t.Errorf("Pre-shutdown hook was not called")
	}

	if !shutdownFuncCalled {
		t.Errorf("Shutdown function was not called")
	}

	if !postShutdownHookCalled {
		t.Errorf("Post-shutdown hook was not called")
	}
}
