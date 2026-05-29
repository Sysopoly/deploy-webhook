package runner

import (
	"os"
	"testing"
	"time"
)

func TestRunner_CoalescesMultipleTriggers(t *testing.T) {
	if os.Getenv("CI") != "" && os.Getenv("GOOS") == "windows" {
		t.Skip("requires bash")
	}

	tmp := t.TempDir() + "/apply.sh"
	os.WriteFile(tmp, []byte("#!/bin/bash\nsleep 0.2\n"), 0755)

	r := New(tmp)

	// Trigger first run
	r.Trigger()
	time.Sleep(50 * time.Millisecond) // let it start

	// Trigger 3 more while running — should coalesce into 1 re-run
	r.Trigger()
	r.Trigger()
	r.Trigger()

	// Wait for both runs to complete (initial + 1 coalesced re-run)
	time.Sleep(800 * time.Millisecond)

	status := r.Status()
	if status.Running {
		t.Error("expected not running after completion")
	}
	if status.Pending {
		t.Error("expected no pending after completion")
	}
}

func TestRunner_StatusWhileIdle(t *testing.T) {
	r := New("/bin/true")
	status := r.Status()
	if status.Running || status.Pending {
		t.Error("expected idle status")
	}
}
