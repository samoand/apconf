package apconf

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

func processYamlDirs(dirs []string, templateParams map[string]any) map[string]any {
	finalDict := make(map[string]any)

	for _, dirPath := range dirs {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			panic(err) // Handle error as needed
		}

		for _, entry := range entries {
			if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || filepath.Ext(entry.Name()) != ".yaml" {
				continue
			}

			content, err := os.ReadFile(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				panic(err) // Handle error as needed
			}
			processedContent := renderTemplate(content, templateParams)
			processYamlContent(processedContent, finalDict)
		}
	}
	return finalDict
}

func processYamlContent(content []byte, finalDict map[string]any) {
	// First try to unmarshal into a single document map
	var singleDoc map[string]any
	if err := yaml.Unmarshal(content, &singleDoc); err == nil {
		mergeYamlDocument(singleDoc, finalDict)
		return
	}

	// If unmarshalling into a single map failed, try a slice of maps
	var multiDocs []map[string]any
	if err := yaml.Unmarshal(content, &multiDocs); err == nil {
		for _, doc := range multiDocs {
			mergeYamlDocument(doc, finalDict)
		}
	} else {
		// Handle the error properly (log, panic, or return an error)
		panic(err)
	}
}

func mergeYamlDocument(doc map[string]any, finalDict map[string]any) {
	if kind, ok := doc["kind"].(string); ok {
		if metadata, ok := doc["metadata"].(map[string]any); ok {
			if name, ok := metadata["name"].(string); ok {
				if finalDict[kind] == nil {
					finalDict[kind] = make(map[string]any)
				}
				finalDict[kind].(map[string]any)[name] = doc
			}
		}
	}
}

func preprocessTemplateForGo(content []byte) []byte {
	// This regular expression matches '{{ var_name }}' and captures 'var_name'
	re := regexp.MustCompile(`{{\s*(\w+)\s*}}`)

	// Replace '{{ var_name }}' with '{{ .VarName }}'
	processedContent := re.ReplaceAllFunc(content, func(match []byte) []byte {
		// Extract the variable name
		varName := re.FindSubmatch(match)[1]
		// Convert to Go style '{{ .VarName }}'
		pascalVarName := toPascalCase(string(varName))
		if !strings.HasPrefix(pascalVarName, ".") {
			pascalVarName = "." + pascalVarName
		}
		return []byte("{{ " + pascalVarName + " }}")
	})

	return processedContent
}

// Helper function to convert a string from snake_case to PascalCase
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			// Convert the first character to uppercase
			runes := []rune(part)
			runes[0] = unicode.ToUpper(runes[0])
			// Convert the rest of the characters to lowercase (optional)
			parts[i] = string(runes)
		}
	}
	// Join all parts together to form PascalCase
	return strings.Join(parts, "")
}

func renderTemplate(content []byte, templateParams map[string]any) []byte {
	content = preprocessTemplateForGo(content)

	// Create a new template and parse the content into it.
	tmpl, err := template.New("configTemplate").Parse(string(content))
	if err != nil {
		panic(err)
	}

	// Use a buffer to capture the output of the template execution.
	var renderedContent bytes.Buffer
	err = tmpl.Execute(&renderedContent, templateParams)
	if err != nil {
		panic(err)
	}

	return renderedContent.Bytes()
}
