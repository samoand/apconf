package apconf

import (
	"reflect"
	"testing"
)

func TestConfigDiff(t *testing.T) {
	t.Run("no_difference", func(t *testing.T) {
		configNew := map[string]any{"a": 1, "b": 2}
		configOld := map[string]any{"a": 1, "b": 2}
		expectedChanged := map[string]any{}
		expectedAdded := map[string]any{}
		expectedRemoved := map[string]any{}
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
		configNew := map[string]any{"a": 1, "b": 2}
		configOld := map[string]any{"a": 1, "b": 3}
		expectedChanged := map[string]any{"b": 2}
		expectedAdded := map[string]any{}
		expectedRemoved := map[string]any{}
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
		configNew := map[string]any{"a": 1, "b": map[string]any{"c": 2, "d": 3}}
		configOld := map[string]any{"a": 1, "b": map[string]any{"c": 2, "d": 4}}
		expectedChanged := map[string]any{"b": map[string]any{"d": 3}}
		expectedAdded := map[string]any{}
		expectedRemoved := map[string]any{}
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
		configNew := map[string]any{"a": 1, "b": 2, "c": 3}
		configOld := map[string]any{"a": 1, "b": 2}
		expectedChanged := map[string]any{}
		expectedAdded := map[string]any{"c": 3}
		expectedRemoved := map[string]any{}
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
		configNew := map[string]any{"a": 1}
		configOld := map[string]any{"a": 1, "b": 2}
		expectedChanged := map[string]any{}
		expectedAdded := map[string]any{}
		expectedRemoved := map[string]any{"b": 2}
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
		configNew := map[string]any{"a": 1, "b": map[string]any{"c": 2, "d": map[string]any{"e": 3, "f": 4}}}
		configOld := map[string]any{"a": 1, "b": map[string]any{"c": 2, "d": map[string]any{"e": 3, "f": 5}, "g": 6}}
		expectedChanged := map[string]any{"b": map[string]any{"d": map[string]any{"f": 4}}}
		expectedAdded := map[string]any{}
		expectedRemoved := map[string]any{"b": map[string]any{"g": 6}}
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
