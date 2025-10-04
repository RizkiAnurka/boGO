package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// loadTemplate reads a template file and returns its content
func loadTemplate(templateName string) (string, error) {
	templatePath := filepath.Join("templates", templateName+".template")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templateName, err)
	}
	return string(content), nil
}

// replaceTemplateVariables replaces template variables with actual values
func replaceTemplateVariables(template string, variables map[string]string) string {
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("<%s>", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// processTemplate loads a template and replaces variables
func processTemplate(templateName string, variables map[string]string) (string, error) {
	template, err := loadTemplate(templateName)
	if err != nil {
		return "", err
	}

	return replaceTemplateVariables(template, variables), nil
}
