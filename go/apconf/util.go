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

func (self ConfigDiffResult) Contains(path []string) bool {
	return pathExistsInMap(self.Changed, path) ||
		pathExistsInMap(self.Added, path) ||
		pathExistsInMap(self.Removed, path)
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

		if newExists && oldExists {
			// Both exist, compare them
			if reflect.TypeOf(newValue) == reflect.TypeOf(oldValue) {
				if reflect.DeepEqual(newValue, oldValue) {
					continue
				}

				if newMap, ok := newValue.(map[string]any); ok {
					if oldMap, ok := oldValue.(map[string]any); ok {
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
					}
				} else {
					changed[key] = newValue
				}
			} else {
				changed[key] = newValue
			}
		} else if newExists {
			added[key] = newValue
		} else if oldExists {
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
