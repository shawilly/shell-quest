package world

import "github.com/shanewilliams/shell-quest/internal/shell"

// MissionRunner tracks progress through a mission's objectives.
type MissionRunner struct {
	mission    Mission
	currentIdx int
}

// NewMissionRunner creates a new runner for the given mission.
func NewMissionRunner(m Mission) *MissionRunner {
	return &MissionRunner{mission: m}
}

// IsComplete returns true when all objectives have been met.
func (r *MissionRunner) IsComplete() bool {
	return r.currentIdx >= len(r.mission.Objectives)
}

// CurrentObjectiveIndex returns the index of the current objective.
func (r *MissionRunner) CurrentObjectiveIndex() int {
	return r.currentIdx
}

// CurrentObjective returns the current objective, or nil if complete.
func (r *MissionRunner) CurrentObjective() *Objective {
	if r.IsComplete() {
		return nil
	}
	return &r.mission.Objectives[r.currentIdx]
}

// HandleEvent checks if the event matches the current objective and advances if so.
// Returns true if the objective was advanced.
func (r *MissionRunner) HandleEvent(e *shell.Event) bool {
	if e == nil || r.IsComplete() {
		return false
	}
	obj := r.mission.Objectives[r.currentIdx]
	if e.Type == obj.Command && e.Path == obj.Path {
		r.currentIdx++
		return true
	}
	return false
}

// CurrentHint returns the hint for the objective that just completed.
// Call after HandleEvent returns true.
func (r *MissionRunner) CurrentHint() string {
	idx := r.currentIdx - 1
	hints := r.mission.ObjectiveHints
	if idx >= 0 && idx < len(hints) {
		return hints[idx]
	}
	return ""
}

// Mission returns the mission this runner is tracking.
func (r *MissionRunner) Mission() Mission {
	return r.mission
}
