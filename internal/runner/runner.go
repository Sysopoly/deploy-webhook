package runner

import (
	"log"
	"os/exec"
	"sync"
	"time"
)

// Runner executes apply.sh with dirty-flag coalescing.
// If a trigger arrives while running, it sets a pending flag.
// After the current run finishes, if pending is set, it runs again.
// Multiple triggers during one run collapse into a single re-run.
type Runner struct {
	script  string
	mu      sync.Mutex
	running bool
	pending bool
}

type Status struct {
	Running bool
	Pending bool
}

func New(script string) *Runner {
	return &Runner{script: script}
}

func (r *Runner) Trigger() {
	r.mu.Lock()
	if r.running {
		r.pending = true
		r.mu.Unlock()
		log.Printf("Apply already in progress — queued for re-run")
		return
	}
	r.running = true
	r.mu.Unlock()

	go r.loop()
}

func (r *Runner) Status() Status {
	r.mu.Lock()
	defer r.mu.Unlock()
	return Status{Running: r.running, Pending: r.pending}
}

func (r *Runner) loop() {
	for {
		r.run()

		r.mu.Lock()
		if !r.pending {
			r.running = false
			r.mu.Unlock()
			return
		}
		r.pending = false
		r.mu.Unlock()
		log.Printf("Pending apply detected — running again")
	}
}

func (r *Runner) run() {
	start := time.Now()
	log.Printf("Running: bash %s", r.script)

	cmd := exec.Command("bash", r.script)
	cmd.Dir = "/opt/sysopoly-infrastructure"
	output, err := cmd.CombinedOutput()

	elapsed := time.Since(start).Round(time.Millisecond)
	if err != nil {
		log.Printf("Apply FAILED after %s: %v\n%s", elapsed, err, string(output))
	} else {
		log.Printf("Apply completed in %s", elapsed)
	}
}
