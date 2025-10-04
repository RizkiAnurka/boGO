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

		// Generate route registrations
		routeRegistrations.WriteString(fmt.Sprintf("\n\t// %s routes\n", structName))
		routeRegistrations.WriteString(fmt.Sprintf("\t%sHandler := New%sHandler(r.ctx, r.%s)\n", entityName, structName, fieldName))
		routeRegistrations.WriteString(fmt.Sprintf("\trouter.GET(\"/%s\", %sHandler.GetAll%s)\n", entityPlural, entityName, structName+"s"))
		routeRegistrations.WriteString(fmt.Sprintf("\trouter.POST(\"/%s\", %sHandler.Create%s)\n", entityPlural, entityName, structName))
		routeRegistrations.WriteString(fmt.Sprintf("\trouter.GET(\"/%s/:id\", %sHandler.Get%sByID)\n", entityPlural, entityName, structName))
		routeRegistrations.WriteString(fmt.Sprintf("\trouter.PUT(\"/%s/:id\", %sHandler.Update%s)\n", entityPlural, entityName, structName))
		routeRegistrations.WriteString(fmt.Sprintf("\trouter.DELETE(\"/%s/:id\", %sHandler.Delete%s)\n", entityPlural, entityName, structName))
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
	getAllResult, err := processTemplate("rest-get-all", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-get-all template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(getAllResult)

	createResult, err := processTemplate("rest-create", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-create template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(createResult)

	getByIDResult, err := processTemplate("rest-get-by-id", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-get-by-id template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(getByIDResult)

	updateResult, err := processTemplate("rest-update", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-update template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(updateResult)

	deleteResult, err := processTemplate("rest-delete", vars)
	if err != nil {
		panic(fmt.Sprintf("Error processing rest-delete template: %v", err))
	}
	handler.WriteString("\n")
	handler.WriteString(deleteResult)

	return handler.String()
}
