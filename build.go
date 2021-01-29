package docs

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const docsOutPath = "./internal/dist/openapi.yaml"

// FIXME: Split into 3 fn units.
func (o *OAS) BuildDocs(customOutPath string) error {
	err := o.initCallStackForRoutes()
	if err != nil {
		return fmt.Errorf("failed initiating call stack for registered routes: %w", err)
	}

	transformedOAS := o.transformToMap()

	yml, err := yaml.Marshal(transformedOAS)
	if err != nil {
		return fmt.Errorf("failed marshaling to yaml: %w", err)
	}

	outYAML, err := os.Create(docsOutPath)
	if err != nil {
		return fmt.Errorf("failed creating yaml output file: %w", err)
	}
	defer outYAML.Close()

	writer := bufio.NewWriter(outYAML)

	_, err = writer.Write(yml)
	if err != nil {
		return fmt.Errorf("failed writing to yaml output file: %w", err)
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed flishing output writer: %w", err)
	}

	return nil
}

type (
	pathsMap        map[string]methodsMap
	methodsMap      map[string]interface{}
	pathSecurityMap map[string][]string
)

// FIXME: Validations and refactoring needed.
// Ditch interface{} in place for concrete types.
func (o *OAS) transformToMap() map[string]interface{} {
	oasPrep := make(map[string]interface{})

	oasPrep["openapi"] = o.OASVersion
	oasPrep["info"] = o.Info
	oasPrep["externalDocs"] = o.ExternalDocs
	oasPrep["servers"] = o.Servers
	oasPrep["tags"] = o.Tags

	// FIXME: All will need validations, e.g. if doesn't exist in the struct, then don't register the key in map...
	allPaths := make(pathsMap, len(o.Paths))
	for _, path := range o.Paths { //nolint:gocritic //consider indexing?
		if allPaths[path.Route] == nil {
			allPaths[path.Route] = make(methodsMap)
		}

		reqBodyMap := make(map[string]interface{})
		reqBodyMap["description"] = path.RequestBody.Description
		reqBodyMap["content"] = makeContentSchemaMap(path.RequestBody.Content)

		responsesMap := make(map[uint]interface{}, len(path.Responses))

		for _, resp := range path.Responses {
			codeBodyMap := make(map[string]interface{})
			codeBodyMap["description"] = resp.Description
			codeBodyMap["content"] = makeContentSchemaMap(resp.Content)

			responsesMap[resp.Code] = codeBodyMap
		}

		var securityMaps []pathSecurityMap

		for _, sec := range path.Security {
			inner := make(pathSecurityMap)
			inner[sec.AuthName] = sec.PermTypes

			securityMaps = append(securityMaps, inner)
		}

		pathMap := make(map[string]interface{})
		pathMap["tags"] = path.Tags
		pathMap["summary"] = path.Summary
		pathMap["operationId"] = path.OperationID
		pathMap["security"] = securityMaps
		pathMap["requestBody"] = reqBodyMap
		pathMap["responses"] = responsesMap

		allPaths[path.Route][strings.ToLower(path.HTTPMethod)] = pathMap
	}

	oasPrep["paths"] = allPaths

	componentsMap := make(map[string]interface{}, len(o.Components))

	for _, cm := range o.Components {
		schemesMap := make(map[string]interface{}, len(cm.Schemas))

		for _, s := range cm.Schemas {
			scheme := make(map[string]interface{})
			scheme["type"] = s.Type
			scheme["properties"] = s.Properties
			scheme["$ref"] = s.Ref

			if s.XML.Name != "" {
				scheme["xml"] = s.XML
			}

			schemesMap[s.Name] = scheme
		}

		secSchemesMap := make(map[string]interface{}, len(cm.SecuritySchemes))

		for _, ss := range cm.SecuritySchemes {
			scheme := make(map[string]interface{})
			scheme["name"] = ss.Name
			scheme["type"] = ss.Type

			if ss.In != "" {
				scheme["in"] = ss.In
			}

			secSchemesMap[ss.Name] = scheme
		}

		componentsMap["schemas"] = schemesMap
		componentsMap["securitySchemes"] = secSchemesMap
	}

	oasPrep["components"] = componentsMap

	return oasPrep
}

func makeContentSchemaMap(content ContentTypes) map[string]interface{} {
	contentSchemaMap := make(map[string]interface{})

	for _, ct := range content {
		refMap := make(map[string]string)
		refMap["$ref"] = ct.Schema

		schemaMap := make(map[string]map[string]string)
		schemaMap["schema"] = refMap

		contentSchemaMap[ct.Name] = schemaMap
	}

	return contentSchemaMap
}
