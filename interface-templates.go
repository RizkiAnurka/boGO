package main

import (
	"fmt"
	"strings"
)

// generateUnifiedApplicationInterfaces creates a single adapter.go file with all application interfaces
func generateUnifiedApplicationInterfaces(moduleName string, tables []Table) string {
	var interfaces strings.Builder

	for _, table := range tables {
		structName := toCamelCase(table.Name)
		if strings.HasSuffix(structName, "s") {
			structName = structName[:len(structName)-1]
		}

		entityName := strings.ToLower(structName)

		// Use template for interface generation
		interfaceVars := map[string]string{
			"entity_name": entityName,
			"struct_name": structName,
		}

		interfaceResult, err := processTemplate("application-interface-content", interfaceVars)
		if err != nil {
			panic(fmt.Sprintf("Error processing application-interface-content template: %v", err))
		}
		interfaces.WriteString(interfaceResult)
		interfaces.WriteString("\n")
	}

	variables := map[string]string{
		"module_name": moduleName,
		"interfaces":  interfaces.String(),
	}

	result, err := processTemplate("application-interfaces", variables)
	if err != nil {
		panic(fmt.Sprintf("Error processing application-interfaces template: %v", err))
	}
	return result
}

// generateUnifiedInteractorInterfaces creates a single adapter.go file with all interactor interfaces
func generateUnifiedInteractorInterfaces(moduleName string, tables []Table) string {
	var interfaces strings.Builder

	for _, table := range tables {
		structName := toCamelCase(table.Name)
		if strings.HasSuffix(structName, "s") {
			structName = structName[:len(structName)-1]
		}

		entityName := strings.ToLower(structName)

		// Use template for interface generation
		interfaceVars := map[string]string{
			"entity_name": entityName,
			"struct_name": structName,
		}

		interfaceResult, err := processTemplate("interactor-interface-content", interfaceVars)
		if err != nil {
			panic(fmt.Sprintf("Error processing interactor-interface-content template: %v", err))
		}
		interfaces.WriteString(interfaceResult)
		interfaces.WriteString("\n")
	}

	variables := map[string]string{
		"module_name": moduleName,
		"interfaces":  interfaces.String(),
	}

	content, err := processTemplate("interactor-interfaces", variables)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate unified interactor interfaces: %v", err))
	}
	return content
}
