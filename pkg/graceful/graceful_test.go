package graceful // nolint:testpackage

import "testing"

func TestGraceful_nilHooks(t *testing.T) {
	t.Parallel()

	RegisterPreShutdownHook(nil)
	RegisterPostShutdownHook(nil)

	if len(preShutdownHooks)+len(postShutdownHooks) != 0 {
		t.Fatalf("Expect len(preShutdownHooks) + len(postShutdownHooks) == 0, got %d",
			len(preShutdownHooks)+len(postShutdownHooks))
	}
}
