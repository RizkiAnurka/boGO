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

		interfaces.WriteString(fmt.Sprintf(`// Adapter to %s repository
type i%s interface {
	Find(ctx context.Context, filter, sort map[string]any, limit, offset int) (res []model.%s, total int64, err error)
	Create(ctx context.Context, %s *model.%s) error
	Update(ctx context.Context, %s model.%s) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (model.%s, error)
}

`, entityName, structName, structName, entityName, structName, entityName, structName, structName))
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

		interfaces.WriteString(fmt.Sprintf(`// I%sService interface for %s business operations
type I%sService interface {
	Find(ctx context.Context, filter map[string]any, sort map[string]any, limit, offset int) (%ss dto.%ss, total int64, err error)
	Create(ctx context.Context, %s dto.%s) (int64, error)
	Update(ctx context.Context, %s dto.%s) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (dto.%s, error)
}

`, structName, entityName, structName, entityName, structName, entityName, structName, entityName, structName, structName))
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
