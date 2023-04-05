package docs

// WARNING:
// Most structures in here are an representation of what is defined in default
//		Open API Specification documentation, v3.0.3.
//
// [More about it can be found on this link](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.3.md)

// New returns a new instance of OAS structure.
func New() OAS {
	initRoutes := RegRoutes{}

	return OAS{
		RegisteredRoutes: initRoutes,
	}
}

const (
	oasAnnotationInit = "// @OAS "
)

// OAS - represents Open API Specification structure, in its approximated Go form.
type OAS struct {
	OASVersion       OASVersion   `yaml:"openapi"`
	Info             Info         `yaml:"info"`
	ExternalDocs     ExternalDocs `yaml:"externalDocs"`
	Servers          Servers      `yaml:"servers"`
	Tags             Tags         `yaml:"tags"`
	Paths            Paths        `yaml:"paths"`
	Components       Components   `yaml:"components"`
	RegisteredRoutes RegRoutes    `yaml:"-"`
}

type (
	// Version represents a SemVer version.
	Version string

	// URL represents and URL which is casted from string.
	URL string

	// OASVersion represents the OpenAPISpecification version which will be used.
	OASVersion Version
)

// Info represents OAS info object.
type Info struct {
	Title          string  `yaml:"title"`
	Description    string  `yaml:"description"`
	TermsOfService URL     `yaml:"termsOfService"`
	Contact        Contact `yaml:"contact"`
	License        License `yaml:"license"`
	Version        Version `yaml:"version"`
}

// Contact represents OAS contact object, used by Info.
type Contact struct {
	Email string `yaml:"email"`
}

// License represents OAS license object, used by Info.
type License struct {
	Name string `yaml:"name"`
	URL  URL    `yaml:"url"`
}

// ExternalDocs represents OAS externalDocs object.
//
// Aside from base OAS structure, this is also used by Tag object.
type ExternalDocs struct {
	Description string `yaml:"description"`
	URL         URL    `yaml:"url"`
}

// Servers is a slice of Server objects.
type Servers []Server

// Server represents OAS server object.
type Server struct {
	URL URL `yaml:"url"`
}

// Tags is a slice of Tag objects.
type Tags []Tag

// Tag represents OAS tag object.
type Tag struct {
	Name         string       `yaml:"name"`
	Description  string       `yaml:"description"`
	ExternalDocs ExternalDocs `yaml:"externalDocs"`
}

// Paths is a slice of Path objects.
type Paths []Path

// Path represents OAS path object.
type Path struct {
	Route           string           `yaml:"route"`
	HTTPMethod      string           `yaml:"httpMethod"`
	Tags            []string         `yaml:"tags"`
	Summary         string           `yaml:"summary"`
	Description     string           `yaml:"description"`
	OperationID     string           `yaml:"operationId"`
	RequestBody     RequestBody      `yaml:"requestBody"`
	Responses       Responses        `yaml:"responses"`
	Security        SecurityEntities `yaml:"security,omitempty"`
	Parameters      Parameters       `yaml:"parameters,omitempty"`
	HandlerFuncName string           `yaml:"-"`
}

// RequestBody represents OAS requestBody object, used by Path.
type RequestBody struct {
	Description string       `yaml:"description"`
	Content     ContentTypes `yaml:"content"`
	Required    bool         `yaml:"required"`
}

// ContentTypes is a slice of ContentType objects.
type ContentTypes []ContentType

// ContentType represents OAS content type object, used by RequestBody and Response.
type ContentType struct {
	Name   string `yaml:"ct-name"`   // e.g. application/json
	Schema string `yaml:"ct-schema"` // e.g. $ref: '#/components/schemas/Pet'
}

// Responses is a slice of Response objects.
type Responses []Response

// Response represents OAS response object, used by Path.
type Response struct {
	Code        uint         `yaml:"code"`
	Description string       `yaml:"description"`
	Content     ContentTypes `yaml:"content"`
}

// SecurityEntities is a slice of Security objects.
type SecurityEntities []Security

// Security represents OAS security object.
type Security struct {
	AuthName  string
	PermTypes []string // write:pets , read:pets etc.
}

// Components is a slice of Component objects.
type Components []Component

// Component represents OAS component object.
type Component struct {
	Schemas         Schemas         `yaml:"schemas"`
	SecuritySchemes SecuritySchemes `yaml:"securitySchemes"`
}

// Schemas is a slice of Schema objects.
type Schemas []Schema

// Schema represents OAS schema object, used by Component.
type Schema struct {
	Name       string
	Type       string
	Properties SchemaProperties
	XML        XMLEntry    `yaml:"xml, omitempty"`
	Ref        string      `yaml:"$ref,omitempty"` // $ref: '#/components/schemas/Pet' // TODO: Should this be omitted if empty?
	Items      *ArrayItems `yaml:"items,omitempty"`
}

// XMLEntry represents name of XML entry in Schema object.
type XMLEntry struct {
	Name string
}

// SchemaProperties is a slice of SchemaProperty objects.
type SchemaProperties []SchemaProperty

// SchemaProperty represents OAS schema object, used by Schema.
type SchemaProperty struct {
	Name        string            `yaml:"-"`
	Type        string            // OAS3.0 data types - e.g. integer, boolean, string
	Format      string            `yaml:"format,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Ref         string            `yaml:"$ref,omitempty"`
	Enum        []string          `yaml:"enum,omitempty"`
	Default     interface{}       `yaml:"default,omitempty"`
	Properties  *SchemaProperties `yaml:"properties,omitempty"`
	Items       *ArrayItems       `yaml:"items,omitempty"`
}

type ArrayItems struct {
	Properties *SchemaProperties `yaml:"properties,omitempty"`
	Ref        string            `yaml:"$ref,omitempty"`
}

// SecuritySchemes is a slice of SecuritySchemes objects.
type SecuritySchemes []SecurityScheme

// SecurityScheme represents OAS security object, used by Component.
type SecurityScheme struct {
	Name  string        `yaml:"name,omitempty"`
	Type  string        `yaml:"type,omitempty"`
	In    string        `yaml:"in,omitempty"`
	Flows SecurityFlows `yaml:"flows,omitempty"`
}

// SecurityFlows is a slice of SecurityFlow objects.
type SecurityFlows []SecurityFlow

// SecurityFlow represents OAS Flows object, used by SecurityScheme.
type SecurityFlow struct {
	Type    string         `yaml:"type,omitempty"`
	AuthURL URL            `yaml:"authorizationUrl,omitempty"`
	Scopes  SecurityScopes `yaml:"scopes,omitempty"`
}

// SecurityScopes is a slice of SecurityScope objects.
type SecurityScopes []SecurityScope

// SecurityScope represents OAS SecurityScope object, used by  SecurityFlow.
type SecurityScope struct {
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// Parameters is a slice of Parameter objects.
type Parameters []Parameter

// Parameter represents OAS parameter object.
type Parameter struct {
	// If in is "path", the name field MUST correspond to a template expression occurring within
	// the path field in the Paths Object. See Path Templating for further information.
	// If in is "header" and the name field is "Accept", "Content-Type" or "Authorization",
	//  the parameter definition SHALL be ignored.
	// For all other cases, the name corresponds to the parameter name used by the in property.
	Name            string `yaml:"name,omitempty"`
	In              string `yaml:"in,omitempty"` // "query", "header", "path" or "cookie".
	Description     string `yaml:"description,omitempty"`
	Required        bool   `yaml:"required,omitempty"`
	Deprecated      bool   `yaml:"deprecated,omitempty"`
	AllowEmptyValue bool   `yaml:"allowEmptyValue,omitempty"`
	Schema          Schema
}

// isEmpty checks if *ExternalDocs struct is empty.
func (ed *ExternalDocs) isEmpty() bool {
	if ed == nil {
		return true
	}

	if ed.Description == emptyStr && ed.URL == emptyStr {
		return true
	}

	return false
}
