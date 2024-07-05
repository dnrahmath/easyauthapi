package docs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/swaggo/swag/gen"
	"gopkg.in/yaml.v2"
)

/*
go get -u github.com/swaggo/swag/gen
go get -u gopkg.in/yaml.v2

//jika tetap error , coba kembali ke cadangan sebelumnya
//dan lakukan ulang
go mod tidy
*/

// createDocs is the main function to generate swagger documentation
func CreateDocs() {
	// Path to your docs directory
	// docsDir := "D:\\file-kodingan\\easyauthapi\\docs"
	docsDir := "./"

	// Ensure the docs directory exists
	err := os.MkdirAll(docsDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create docs directory: %v", err)
	}

	// Generate docs.go
	docsGoPath := filepath.Join(docsDir, "docs.go")
	err = GenerateDocsGo(docsGoPath)
	if err != nil {
		log.Fatalf("Failed to generate docs.go: %v", err)
	}
	log.Printf("Generated docs.go at %s", docsGoPath)

	// Generate swagger.json
	swaggerJSONPath := filepath.Join(docsDir, "swagger.json")
	err = GenerateSwaggerFile(swaggerJSONPath, "json")
	if err != nil {
		log.Fatalf("Failed to generate swagger.json: %v", err)
	}
	log.Printf("Generated swagger.json at %s", swaggerJSONPath)

	// Generate swagger.yaml
	swaggerYAMLPath := filepath.Join(docsDir, "swagger.yaml")
	err = GenerateSwaggerFile(swaggerYAMLPath, "yaml")
	if err != nil {
		log.Fatalf("Failed to generate swagger.yaml: %v", err)
	}
	log.Printf("Generated swagger.yaml at %s", swaggerYAMLPath)
}

// GenerateDocsGo generates the docs.go file
func GenerateDocsGo(outputPath string) error {
	conf := &gen.Config{
		SearchDir:           "./",
		Excludes:            "",
		OutputDir:           filepath.Dir(outputPath),
		MainAPIFile:         "main.go",
		PropNamingStrategy:  "",
		MarkdownFilesDir:    "",
		CodeExampleFilesDir: "",
		ParseVendor:         false,
	}

	// Generate swagger documentation
	return gen.New().Build(conf)
}

// GenerateSwaggerFile generates the swagger file in the specified format (json or yaml)
func GenerateSwaggerFile(outputPath, format string) error {
	// Read the generated swagger
	swaggerFile := "D:\\file-kodingan\\easyauthapi\\docs\\swagger.json"
	swagger, err := ioutil.ReadFile(swaggerFile)
	if err != nil {
		return err
	}

	// Marshal the swagger into the specified format
	var data []byte
	if format == "json" {
		data = swagger
	} else if format == "yaml" {
		var yamlData map[string]interface{}
		err = json.Unmarshal(swagger, &yamlData)
		if err != nil {
			return err
		}
		data, err = yaml.Marshal(yamlData)
		if err != nil {
			return err
		}
	} else {
		return errors.New("unsupported format")
	}

	// Write the swagger file
	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// func main() {
// 	createDocs()
// }
