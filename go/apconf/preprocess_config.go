package apconf

import "reflect"

func processList(
	lst []any,
	keyFilter func(string) bool,
	keyTransformer func(string) string,
	valueTransformer func(any) any,
	processedIDs map[uintptr]struct{},
	keepFiltKeys bool,
) {
	for _, item := range lst {
		switch v := item.(type) {
		case map[string]any:
			processDict(v, keyFilter, keyTransformer, valueTransformer, processedIDs, keepFiltKeys)
		case []any:
			processList(v, keyFilter, keyTransformer, valueTransformer, processedIDs, keepFiltKeys)
		}
	}
}

func processDict(
	d map[string]any,
	keyFilter func(string) bool,
	keyTransformer func(string) string,
	valueTransformer func(any) any,
	processedIDs map[uintptr]struct{},
	keepFiltKeys bool,
) {
	keysToAdd := make(map[string]any)
	var filtKeys []string

	for key, value := range d {
		if keyFilter(key) {
			newKey := keyTransformer(key)
			newValue := valueTransformer(value)

			id := reflect.ValueOf(&newValue).Pointer()
			if _, processed := processedIDs[id]; !processed {
				processedIDs[id] = struct{}{}
				switch nv := newValue.(type) {
				case map[string]any:
					processDict(nv, keyFilter, keyTransformer, valueTransformer, processedIDs, keepFiltKeys)
				case []any:
					processList(nv, keyFilter, keyTransformer, valueTransformer, processedIDs, keepFiltKeys)
				}
				keysToAdd[newKey] = newValue
				filtKeys = append(filtKeys, key)
			}
		} else {
			switch v := value.(type) {
			case map[string]any:
				processDict(v, keyFilter, keyTransformer, valueTransformer, processedIDs, keepFiltKeys)
			case []any:
				processList(v, keyFilter, keyTransformer, valueTransformer, processedIDs, keepFiltKeys)
			}
		}
	}

	for newKey, newValue := range keysToAdd {
		d[newKey] = newValue
	}
	if !keepFiltKeys {
		for _, key := range filtKeys {
			delete(d, key)
		}
	}
}

func Preprocessor(
	keyFilter func(string) bool,
	keyTransformer func(string) string,
	valueTransformer func(any) any,
	keepFiltKeys bool,
) func(map[string]any) {
	return func(config map[string]any) {
		processDict(
			config,
			keyFilter,
			keyTransformer,
			valueTransformer,
			make(map[uintptr]struct{}),
			keepFiltKeys)
	}
}
