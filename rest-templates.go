package main

import (
	"fmt"
	"strings"
)

// generateRestAPIMain creates the main REST API file
func generateRestAPIMain(moduleName string, tables []Table) string {
	var serviceFields strings.Builder
	var serviceParams strings.Builder
	var serviceInit strings.Builder
	var routeRegistrations strings.Builder
	var interactorImport string

	// Only import interactor package if there are tables
	if len(tables) > 0 {
		interactorImport = fmt.Sprintf("\n\n\t\"%s/internal/interactor\"", moduleName)
	}

	// Generate service interfaces and initialization for each table
	for i, table := range tables {
		structName := toCamelCase(table.Name)
		if strings.HasSuffix(structName, "s") {
			structName = structName[:len(structName)-1]
		}

		interfaceName := fmt.Sprintf("interactor.I%sService", structName)
		fieldName := fmt.Sprintf("%sService", strings.ToLower(structName))
		entityName := strings.ToLower(structName)
		entityPlural := entityName + "s"

		serviceFields.WriteString(fmt.Sprintf("\t%s %s\n", fieldName, interfaceName))

		if i > 0 {
			serviceParams.WriteString(", ")
		}
		serviceParams.WriteString(fmt.Sprintf("%s %s", fieldName, interfaceName))

		serviceInit.WriteString(fmt.Sprintf("\t\t%s: %s,\n", fieldName, fieldName))

		// Generate route registrations using template
		routeVars := map[string]string{
			"struct_name":   structName,
			"entity_name":   entityName,
			"entity_plural": entityPlural,
			"plural_name":   structName + "s",
			"field_name":    fieldName,
		}

		routeResult, err := processTemplate("rest-routes", routeVars)
		if err != nil {
			panic(fmt.Sprintf("Error processing rest-routes template: %v", err))
		}
		routeRegistrations.WriteString("\n")
		routeRegistrations.WriteString(routeResult)
	}

	serviceParamsStr := ""
	if serviceParams.Len() > 0 {
		serviceParamsStr = ", " + serviceParams.String()
	}

	vars := map[string]string{
		"module_name":         moduleName,
		"service_fields":      serviceFields.String(),
		"service_params":      serviceParamsStr,
		"service_init":        serviceInit.String(),
		"route_registrations": routeRegistrations.String(),
		"interactor_import":   interactorImport,
	}

	result, err := processTemplate("rest-api-main", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-api-main template: %v", err))
	}
	return result
}

// generateRestHandler creates REST handler for individual table
func generateRestHandler(moduleName string, table Table) string {
	structName := toCamelCase(table.Name)
	if strings.HasSuffix(structName, "s") {
		structName = structName[:len(structName)-1]
	}

	entityName := strings.ToLower(structName)
	entityPlural := entityName + "s"
	singularName := structName
	pluralName := structName + "s"
	dtoName := structName
	entityVar := entityName

	vars := map[string]string{
		"module_name":     moduleName,
		"struct_name":     structName,
		"entity_singular": entityName,
		"entity_plural":   entityPlural,
		"entity_snake":    entityName,
		"singular_name":   singularName,
		"plural_name":     pluralName,
		"dto_name":        dtoName,
		"entity_var":      entityVar,
	}

	var handler strings.Builder

	// Package and imports
	headerResult, err := processTemplate("rest-handler-header", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-handler-header template: %v", err))
	}
	handler.WriteString(headerResult)

	// Individual handler methods
	getAllResult, err := processTemplate("rest-func-get-all", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-func-get-all template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(getAllResult)

	createResult, err := processTemplate("rest-func-create", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-func-create template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(createResult)

	getByIDResult, err := processTemplate("rest-func-get-by-id", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-func-get-by-id template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(getByIDResult)

	updateResult, err := processTemplate("rest-func-update", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-func-update template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(updateResult)

	deleteResult, err := processTemplate("rest-func-delete", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-func-delete template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(deleteResult)

	return handler.String()
}

// generateRestParameter creates the REST parameter file for filtering and sorting
func generateRestParameter(tables []Table) string {
	var allContent strings.Builder

	// Add package header and imports using template
	headerResult, err := processTemplate("rest-parameter-header", map[string]string{})
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-parameter-header template: %v", err))
	}
	allContent.WriteString(headerResult)

	// Generate filter and sorting variables for each table
	for i, table := range tables {
		structName := toCamelCase(table.Name)
		if strings.HasSuffix(structName, "s") {
			structName = structName[:len(structName)-1]
		}

		entitySnake := strings.ToLower(structName)

		// Generate filter fields
		var filterFields strings.Builder
		var sortingFields strings.Builder

		// Add ID field first
		filterFields.WriteString("\n\t\t{Omitempty: true, DBKey: \"id\", Kind: reflect.Int64, QueryKey: \"id\"},")
		sortingFields.WriteString("\n\t\t{DBKey: \"id\", QueryKey: \"id\", Kind: reflect.Int64},")

		// Add other fields
		for _, col := range table.Columns {
			if strings.ToLower(col.Name) == "id" ||
				strings.ToLower(col.Name) == "created_at" ||
				strings.ToLower(col.Name) == "updated_at" ||
				strings.ToLower(col.Name) == "deleted_at" ||
				strings.ToLower(col.Name) == "is_deleted" {
				continue
			}

			reflectType := "reflect.String"
			if col.GoType == "int64" {
				reflectType = "reflect.Int64"
			} else if col.GoType == "float64" {
				reflectType = "reflect.Float64"
			} else if col.GoType == "bool" {
				reflectType = "reflect.Bool"
			}

			filterFields.WriteString(fmt.Sprintf("\n\t\t{Omitempty: true, DBKey: \"%s\", Kind: %s, QueryKey: \"%s\"},", col.Name, reflectType, col.Name))
			sortingFields.WriteString(fmt.Sprintf("\n\t\t{DBKey: \"%s\", QueryKey: \"%s\", Kind: %s},", col.Name, col.Name, reflectType))
		}

		// Process template for this table
		variables := map[string]string{
			"entity_snake":   entitySnake,
			"filter_fields":  filterFields.String(),
			"sorting_fields": sortingFields.String(),
		}

		result, err := processTemplate("rest-parameter", variables)
		if err != nil {
			panic(fmt.Sprintf("Error processing rest-parameter template: %v", err))
		}

		allContent.WriteString(result)

		// Add spacing between tables (except for the last one)
		if i < len(tables)-1 {
			allContent.WriteString("\n")
		}
	}

	// Close the var block
	allContent.WriteString("\n)\n")

	return allContent.String()
}
