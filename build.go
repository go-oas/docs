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
	pathsMap         map[string]methodsMap
	methodsMap       map[string]interface{}
	pathSecurityMap  map[string][]string
	pathSecurityMaps []pathSecurityMap
)

// TODO: Validations and refactoring needed.
func (o *OAS) transformToMap() map[string]interface{} {
	oasPrep := make(map[string]interface{})

	oasPrep["openapi"] = o.OASVersion
	oasPrep["info"] = o.Info
	oasPrep["externalDocs"] = o.ExternalDocs
	oasPrep["servers"] = o.Servers
	oasPrep["tags"] = o.Tags

	oasPrep["paths"] = makeAllPathsMap(&o.Paths)
	oasPrep["components"] = makeComponentsMap(&o.Components)

	return oasPrep
}

func makeAllPathsMap(paths *Paths) pathsMap {
	allPaths := make(pathsMap, len(*paths))
	for _, path := range *paths { //nolint:gocritic //consider indexing?
		if allPaths[path.Route] == nil {
			allPaths[path.Route] = make(methodsMap)
		}

		pathMap := make(map[string]interface{})
		pathMap["tags"] = path.Tags
		pathMap["summary"] = path.Summary
		pathMap["operationId"] = path.OperationID
		pathMap["security"] = makeSecurityMap(&path.Security)
		pathMap["requestBody"] = makeRequestBodyMap(&path.RequestBody)
		pathMap["responses"] = makeResponsesMap(&path.Responses)

		allPaths[path.Route][strings.ToLower(path.HTTPMethod)] = pathMap
	}

	return allPaths
}

func makeRequestBodyMap(reqBody *RequestBody) map[string]interface{} {
	reqBodyMap := make(map[string]interface{})

	reqBodyMap["description"] = reqBody.Description
	reqBodyMap["content"] = makeContentSchemaMap(reqBody.Content)

	return reqBodyMap
}

func makeResponsesMap(responses *Responses) map[uint]interface{} {
	responsesMap := make(map[uint]interface{}, len(*responses))

	for _, resp := range *responses {
		codeBodyMap := make(map[string]interface{})
		codeBodyMap["description"] = resp.Description
		codeBodyMap["content"] = makeContentSchemaMap(resp.Content)

		responsesMap[resp.Code] = codeBodyMap
	}

	return responsesMap
}

func makeSecurityMap(se *SecurityEntities) pathSecurityMaps {
	securityMaps := make(pathSecurityMaps, 0, len(*se))

	for _, sec := range *se {
		securityMap := make(pathSecurityMap)
		securityMap[sec.AuthName] = sec.PermTypes

		securityMaps = append(securityMaps, securityMap)
	}

	return securityMaps
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

func makeComponentsMap(components *Components) map[string]interface{} {
	componentsMap := make(map[string]interface{}, len(*components))

	for _, cm := range *components {
		componentsMap["schemas"] = makeComponentSchemasMap(&cm.Schemas)
		componentsMap["securitySchemes"] = makeComponentSecuritySchemesMap(&cm.SecuritySchemes)
	}

	return componentsMap
}

func makeComponentSchemasMap(schemas *Schemas) map[string]interface{} {
	schemesMap := make(map[string]interface{}, len(*schemas))

	for _, s := range *schemas {
		scheme := make(map[string]interface{})
		scheme["type"] = s.Type
		scheme["properties"] = s.Properties
		scheme["$ref"] = s.Ref

		if s.XML.Name != "" {
			scheme["xml"] = s.XML
		}

		schemesMap[s.Name] = scheme
	}

	return schemesMap
}

func makeComponentSecuritySchemesMap(secSchemes *SecuritySchemes) map[string]interface{} {
	secSchemesMap := make(map[string]interface{}, len(*secSchemes))

	for _, ss := range *secSchemes {
		scheme := make(map[string]interface{})
		scheme["name"] = ss.Name
		scheme["type"] = ss.Type

		if ss.In != "" {
			scheme["in"] = ss.In
		}

		secSchemesMap[ss.Name] = scheme
	}

	return secSchemesMap
}
