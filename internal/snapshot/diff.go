package snapshot

// ChangeKind describes the type of change between two snapshots.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// Change represents a single key-level difference between two snapshots.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Diff compares a previous snapshot against a next snapshot and returns
// the list of changes. Secret values are included so callers can decide
// whether to redact them before logging.
func Diff(prev, next *Snapshot) []Change {
	var changes []Change

	prevMap := prev.ToMap()
	nextMap := next.ToMap()

	// Detect changed and removed keys.
	for k, oldVal := range prevMap {
		newVal, exists := nextMap[k]
		if !exists {
			changes = append(changes, Change{
				Key:      k,
				Kind:     Removed,
				OldValue: oldVal,
			})
		} else if oldVal != newVal {
			changes = append(changes, Change{
				Key:      k,
				Kind:     Changed,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	// Detect added keys.
	for k, newVal := range nextMap {
		if _, exists := prevMap[k]; !exists {
			changes = append(changes, Change{
				Key:      k,
				Kind:     Added,
				NewValue: newVal,
			})
		}
	}

	return changes
}
