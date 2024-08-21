package apconf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

// Utility function to check if a path exists in a nested map.
func pathExistsInMap(data map[string]any, path []string) bool {
	current := data
	for i, key := range path {
		val, exists := current[key]
		if !exists {
			return false
		}
		if i == len(path)-1 {
			return true
		}
		next, ok := val.(map[string]any)
		if !ok {
			return false
		}
		current = next
	}

	return false
}

type ConfigDiffResult struct {
	Changed map[string]any
	Added   map[string]any
	Removed map[string]any
}

func (entity ConfigDiffResult) Contains(path []string) bool {
	return pathExistsInMap(entity.Changed, path) ||
		pathExistsInMap(entity.Added, path) ||
		pathExistsInMap(entity.Removed, path)
}

func findParentDir(startPath string, matchers []func(string) bool) (string, error) {
	currentPath := startPath
	for currentPath != filepath.Dir(currentPath) {
		for _, matcher := range matchers {
			if matcher(currentPath) {
				return currentPath, nil
			}
		}
		currentPath = filepath.Dir(currentPath)
	}
	return "", errors.New("parent directory not found")
}

func findGitRoot(startPath string) (string, error) {
	return findParentDir(startPath, []func(string) bool{
		func(path string) bool {
			_, err := os.Stat(filepath.Join(path, ".git"))
			return !os.IsNotExist(err)
		},
	})
}

func ConfigDiff(configNew, configOld map[string]any) ConfigDiffResult {
	isMap := func(value any) bool {
		_, ok := value.(map[string]any)
		return ok
	}
	changed := make(map[string]any)
	added := make(map[string]any)
	removed := make(map[string]any)

	allKeys := make(map[string]struct{})
	for key := range configNew {
		allKeys[key] = struct{}{}
	}
	for key := range configOld {
		allKeys[key] = struct{}{}
	}

	for key := range allKeys {
		newValue, newExists := configNew[key]
		oldValue, oldExists := configOld[key]

		switch {
		case newExists && oldExists:
			// Both exist, compare them
			switch {
			case reflect.TypeOf(newValue) == reflect.TypeOf(oldValue):
				switch {
				case reflect.DeepEqual(newValue, oldValue):
					// If both values are deeply equal, continue to the next iteration
					continue
				case isMap(newValue) && isMap(oldValue):
					// If both values are maps, perform a nested diff
					newMap := newValue.(map[string]any)
					oldMap := oldValue.(map[string]any)
					nestedResult := ConfigDiff(newMap, oldMap)

					if len(nestedResult.Changed) > 0 {
						changed[key] = nestedResult.Changed
					}
					if len(nestedResult.Added) > 0 {
						added[key] = nestedResult.Added
					}
					if len(nestedResult.Removed) > 0 {
						removed[key] = nestedResult.Removed
					}
				default:
					// Different values of the same type
					changed[key] = newValue
				}
			default:
				// Different types
				changed[key] = newValue
			}
		case newExists:
			// Only the new value exists
			added[key] = newValue
		case oldExists:
			// Only the old value exists
			removed[key] = oldValue
		}
	}

	return ConfigDiffResult{
		Changed: changed,
		Added:   added,
		Removed: removed,
	}
}

func ToInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		parsedValue, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string to int: %v", err)
		}
		return parsedValue, nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

// deepClone performs a deep copy of a map[string]any
func deepClone(src map[string]any) map[string]any {
	clone := make(map[string]any)
	for k, v := range src {
		switch v := v.(type) {
		case map[string]any:
			clone[k] = deepClone(v) // Recursively clone nested maps
		case []any:
			cloneArr := make([]any, len(v))
			for i, item := range v {
				switch item := item.(type) {
				case map[string]any:
					cloneArr[i] = deepClone(item)
				default:
					cloneArr[i] = item
				}
			}
			clone[k] = cloneArr
		default:
			clone[k] = v
		}
	}
	return clone
}
