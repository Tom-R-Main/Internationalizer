package formats

import (
	"fmt"
	"sort"
)

// UnitKind describes the semantic role of one independently translatable unit.
// Adapters may add format-specific structure while the translation pipeline
// continues to address units by stable ID.
type UnitKind string

const (
	UnitMessage   UnitKind = "message"
	UnitDocument  UnitKind = "document"
	UnitAttribute UnitKind = "attribute"
	UnitTerm      UnitKind = "term"
)

// Unit is the format-neutral boundary between source adapters and translation.
// Context is translator-facing information; Structure is a deterministic
// adapter-owned signature used to detect semantic message changes.
type Unit struct {
	ID        string   `json:"id"`
	Value     string   `json:"value"`
	Kind      UnitKind `json:"kind"`
	Context   string   `json:"context,omitempty"`
	Structure string   `json:"structure,omitempty"`
}

// UnitFormat is implemented by formats with richer semantics than a flat
// key/value catalog. Existing formats automatically receive a compatibility
// adapter through ParseUnits and SerializeUnits.
type UnitFormat interface {
	ParseUnits(data []byte) ([]Unit, error)
	SerializeUnits(units []Unit, original []byte) ([]byte, error)
}

// ParseUnits parses a document into stable semantic translation units.
func ParseUnits(format Format, data []byte) ([]Unit, error) {
	if adapter, ok := format.(UnitFormat); ok {
		units, err := adapter.ParseUnits(data)
		if err != nil {
			return nil, err
		}
		if err := ValidateUnits(units); err != nil {
			return nil, err
		}
		return units, nil
	}
	entries, err := format.Parse(data)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	units := make([]Unit, 0, len(keys))
	for _, key := range keys {
		units = append(units, Unit{ID: key, Value: entries[key], Kind: UnitMessage})
	}
	return units, nil
}

// SerializeUnits writes semantic units through a rich adapter or the legacy
// key/value serializer when the format has no additional structure.
func SerializeUnits(format Format, units []Unit, original []byte) ([]byte, error) {
	if err := ValidateUnits(units); err != nil {
		return nil, err
	}
	if adapter, ok := format.(UnitFormat); ok {
		return adapter.SerializeUnits(units, original)
	}
	return format.Serialize(UnitValues(units), original)
}

// UnitValues returns the compatibility key/value view of semantic units.
func UnitValues(units []Unit) map[string]string {
	values := make(map[string]string, len(units))
	for _, unit := range units {
		values[unit.ID] = unit.Value
	}
	return values
}

// MergeUnitValues preserves baseline unit metadata while replacing values and
// appending source units absent from a new target document.
func MergeUnitValues(baseline, source []Unit, values map[string]string) []Unit {
	merged := make([]Unit, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, unit := range baseline {
		value, ok := values[unit.ID]
		if !ok {
			continue
		}
		unit.Value = value
		merged = append(merged, unit)
		seen[unit.ID] = struct{}{}
	}
	for _, sourceUnit := range source {
		if _, ok := seen[sourceUnit.ID]; ok {
			continue
		}
		value, ok := values[sourceUnit.ID]
		if !ok {
			continue
		}
		sourceUnit.Value = value
		merged = append(merged, sourceUnit)
		seen[sourceUnit.ID] = struct{}{}
	}
	missing := make([]string, 0)
	for id := range values {
		if _, ok := seen[id]; !ok {
			missing = append(missing, id)
		}
	}
	sort.Strings(missing)
	for _, id := range missing {
		unit := Unit{ID: id, Kind: UnitMessage}
		unit.Value = values[id]
		merged = append(merged, unit)
	}
	return merged
}

// ValidateUnits enforces stable non-empty unique unit identities.
func ValidateUnits(units []Unit) error {
	seen := make(map[string]struct{}, len(units))
	for _, unit := range units {
		if unit.ID == "" {
			return fmt.Errorf("translation unit ID must not be empty")
		}
		if _, duplicate := seen[unit.ID]; duplicate {
			return fmt.Errorf("duplicate translation unit ID %q", unit.ID)
		}
		seen[unit.ID] = struct{}{}
	}
	return nil
}
