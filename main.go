package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run *.go <module-name> <sql-schema-file>")
		fmt.Println("Example: go run *.go user-service schema.sql")
		os.Exit(1)
	}

	moduleName := os.Args[1]
	sqlSchemaFile := os.Args[2]

	fmt.Printf("Creating hexagonal architecture for module: %s\n", moduleName)
	fmt.Printf("Using SQL schema from: %s\n", sqlSchemaFile)

	err := createHexagonalArchitecture(moduleName, sqlSchemaFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s!\n", moduleName)
}
