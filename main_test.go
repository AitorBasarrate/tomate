package main

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/timer"
)

func TestNextPhase(t *testing.T) {
	focusDur := 10 * time.Second
	restDur := 5 * time.Second
	m := model{
		phase:     focusPhase,
		totalReps: 2,
		focusDur:  focusDur,
		restDur:   restDur,
		timer:     timer.New(focusDur),
	}

	// 1. Focus -> Rest
	newModel, _ := m.nextPhase()
	m = newModel.(model)
	if m.phase != restPhase {
		t.Errorf("Expected restPhase, got %v", m.phase)
	}
	if m.timer.Timeout != restDur {
		t.Errorf("Expected timer timeout %v, got %v", restDur, m.timer.Timeout)
	}

	// 2. Rest -> Focus (2nd repetition)
	newModel, _ = m.nextPhase()
	m = newModel.(model)
	if m.phase != focusPhase {
		t.Errorf("Expected focusPhase, got %v", m.phase)
	}
	if m.repetition != 1 {
		t.Errorf("Expected repetition 1, got %d", m.repetition)
	}

	// 3. Focus -> Rest (2nd repetition)
	newModel, _ = m.nextPhase()
	m = newModel.(model)
	if m.phase != restPhase {
		t.Errorf("Expected restPhase, got %v", m.phase)
	}

	// 4. Rest -> Finished
	newModel, _ = m.nextPhase()
	m = newModel.(model)
	if m.phase != finishedPhase {
		t.Errorf("Expected finishedPhase, got %v", m.phase)
	}
}

func TestNoRest(t *testing.T) {
	focusDur := 10 * time.Second
	m := model{
		phase:     focusPhase,
		totalReps: 1,
		focusDur:  focusDur,
		restDur:   0,
		timer:     timer.New(focusDur),
	}

	// 1. Focus -> Finished (since restDur is 0)
	newModel, _ := m.nextPhase()
	m = newModel.(model)
	if m.phase != finishedPhase {
		t.Errorf("Expected finishedPhase when restDur is 0, got %v", m.phase)
	}
}
