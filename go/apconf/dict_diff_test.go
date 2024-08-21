package apconf

import (
	"reflect"
	"testing"
)

// nolint: funlen
func TestConfigDiff(t *testing.T) {
	t.Run("no_difference", func(t *testing.T) {
		configNew := msa{"a": 1, "b": 2}
		configOld := msa{"a": 1, "b": 2}
		expectedChanged := msa{}
		expectedAdded := msa{}
		expectedRemoved := msa{}
		configDiffResult := ConfigDiff(configNew, configOld)
		if !reflect.DeepEqual(configDiffResult.Changed, expectedChanged) {
			t.Errorf("Expected changed: %v, got: %v", expectedChanged, configDiffResult.Changed)
		}
		if !reflect.DeepEqual(configDiffResult.Added, expectedAdded) {
			t.Errorf("Expected added: %v, got: %v", expectedAdded, configDiffResult.Added)
		}
		if !reflect.DeepEqual(configDiffResult.Removed, expectedRemoved) {
			t.Errorf("Expected removed: %v, got: %v", expectedRemoved, configDiffResult.Removed)
		}
	})

	t.Run("simple_difference", func(t *testing.T) {
		configNew := msa{"a": 1, "b": 2}
		configOld := msa{"a": 1, "b": 3}
		expectedChanged := msa{"b": 2}
		expectedAdded := msa{}
		expectedRemoved := msa{}
		configDiffResult := ConfigDiff(configNew, configOld)
		if !reflect.DeepEqual(configDiffResult.Changed, expectedChanged) {
			t.Errorf("Expected changed: %v, got: %v", expectedChanged, configDiffResult.Changed)
		}
		if !reflect.DeepEqual(configDiffResult.Added, expectedAdded) {
			t.Errorf("Expected added: %v, got: %v", expectedAdded, configDiffResult.Added)
		}
		if !reflect.DeepEqual(configDiffResult.Removed, expectedRemoved) {
			t.Errorf("Expected removed: %v, got: %v", expectedRemoved, configDiffResult.Removed)
		}
	})

	t.Run("nested_difference", func(t *testing.T) {
		configNew := msa{"a": 1, "b": msa{"c": 2, "d": 3}}
		configOld := msa{"a": 1, "b": msa{"c": 2, "d": 4}}
		expectedChanged := msa{"b": msa{"d": 3}}
		expectedAdded := msa{}
		expectedRemoved := msa{}
		configDiffResult := ConfigDiff(configNew, configOld)
		if !reflect.DeepEqual(configDiffResult.Changed, expectedChanged) {
			t.Errorf("Expected changed: %v, got: %v", expectedChanged, configDiffResult.Changed)
		}
		if !reflect.DeepEqual(configDiffResult.Added, expectedAdded) {
			t.Errorf("Expected added: %v, got: %v", expectedAdded, configDiffResult.Added)
		}
		if !reflect.DeepEqual(configDiffResult.Removed, expectedRemoved) {
			t.Errorf("Expected removed: %v, got: %v", expectedRemoved, configDiffResult.Removed)
		}
	})

	t.Run("extra_keys_in_config_new", func(t *testing.T) {
		configNew := msa{"a": 1, "b": 2, "c": 3}
		configOld := msa{"a": 1, "b": 2}
		expectedChanged := msa{}
		expectedAdded := msa{"c": 3}
		expectedRemoved := msa{}
		configDiffResult := ConfigDiff(configNew, configOld)
		if !reflect.DeepEqual(configDiffResult.Changed, expectedChanged) {
			t.Errorf("Expected changed: %v, got: %v", expectedChanged, configDiffResult.Changed)
		}
		if !reflect.DeepEqual(configDiffResult.Added, expectedAdded) {
			t.Errorf("Expected added: %v, got: %v", expectedAdded, configDiffResult.Added)
		}
		if !reflect.DeepEqual(configDiffResult.Removed, expectedRemoved) {
			t.Errorf("Expected removed: %v, got: %v", expectedRemoved, configDiffResult.Removed)
		}
	})

	t.Run("extra_keys_in_config_old", func(t *testing.T) {
		configNew := msa{"a": 1}
		configOld := msa{"a": 1, "b": 2}
		expectedChanged := msa{}
		expectedAdded := msa{}
		expectedRemoved := msa{"b": 2}
		configDiffResult := ConfigDiff(configNew, configOld)
		if !reflect.DeepEqual(configDiffResult.Changed, expectedChanged) {
			t.Errorf("Expected changed: %v, got: %v", expectedChanged, configDiffResult.Changed)
		}
		if !reflect.DeepEqual(configDiffResult.Added, expectedAdded) {
			t.Errorf("Expected added: %v, got: %v", expectedAdded, configDiffResult.Added)
		}
		if !reflect.DeepEqual(configDiffResult.Removed, expectedRemoved) {
			t.Errorf("Expected removed: %v, got: %v", expectedRemoved, configDiffResult.Removed)
		}
	})

	t.Run("complex_nested_difference", func(t *testing.T) {
		configNew := msa{"a": 1, "b": msa{"c": 2, "d": msa{"e": 3, "f": 4}}}
		configOld := msa{"a": 1, "b": msa{"c": 2, "d": msa{"e": 3, "f": 5}, "g": 6}}
		expectedChanged := msa{"b": msa{"d": msa{"f": 4}}}
		expectedAdded := msa{}
		expectedRemoved := msa{"b": msa{"g": 6}}
		configDiffResult := ConfigDiff(configNew, configOld)
		if !reflect.DeepEqual(configDiffResult.Changed, expectedChanged) {
			t.Errorf("Expected changed: %v, got: %v", expectedChanged, configDiffResult.Changed)
		}
		if !reflect.DeepEqual(configDiffResult.Added, expectedAdded) {
			t.Errorf("Expected added: %v, got: %v", expectedAdded, configDiffResult.Added)
		}
		if !reflect.DeepEqual(configDiffResult.Removed, expectedRemoved) {
			t.Errorf("Expected removed: %v, got: %v", expectedRemoved, configDiffResult.Removed)
		}
	})
}
