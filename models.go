package docs

func New() OAS {
	initRoutes := RegRoutes{}

	return OAS{
		registeredRoutes: initRoutes,
	}
}

const (
	OASAnnotationInit = "// @OAS "
)

// OAS - represents Open API Specification structure, in its approximated Go form.
type OAS struct {
	OASVersion   OASVersion   `yaml:"openapi"`
	Info         Info         `yaml:"info"`
	ExternalDocs ExternalDocs `yaml:"externalDocs"`
	Servers      Servers      `yaml:"servers"`
	Tags         Tags         `yaml:"tags"`
	Paths        Paths        `yaml:"paths"`
	Components   Components   `yaml:"components"`

	registeredRoutes RegRoutes `yaml:"-"`
}

// Version is represented in SemVer format.
type (
	Version    string
	URL        string
	OASVersion Version
)

type Info struct {
	Title          string  `yaml:"title"`
	Description    string  `yaml:"description"`
	TermsOfService string  `yaml:"termsOfService"`
	Contact        Contact `yaml:"contact"`
	License        License `yaml:"license"`
	Version        Version `yaml:"version"`
}

type Contact struct {
	Email string `yaml:"email"`
}

type License struct {
	Name string `yaml:"name"`
	URL  URL    `yaml:"url"`
}

type ExternalDocs struct {
	Description string `yaml:"description"`
	URL         URL    `yaml:"url"`
}

type Servers []Server

type Server struct {
	URL URL `yaml:"url"`
}

type Tags []Tag

type Tag struct {
	Name         string       `yaml:"name"`
	Description  string       `yaml:"description"`
	ExternalDocs ExternalDocs `yaml:"externalDocs"`
}

type Paths []Path

type Path struct {
	Route       string           `yaml:"route"`
	HTTPMethod  string           `yaml:"httpMethod"`
	Tags        []string         `yaml:"tags"`
	Summary     string           `yaml:"summary"`
	OperationID string           `yaml:"operationId"`
	RequestBody RequestBody      `yaml:"requestBody"`
	Responses   Responses        `yaml:"responses"`
	Security    SecurityEntities `yaml:"security,omitempty"`

	handlerFuncName string `yaml:"-"`
}

type RequestBody struct {
	Description string       `yaml:"description"`
	Content     ContentTypes `yaml:"content"`
	Required    bool         `yaml:"required"`
	// TODO: Further develop/research.
}

type ContentTypes []ContentType

type ContentType struct {
	Name   string `yaml:"ct-name"`   // e.g. application/json
	Schema string `yaml:"ct-schema"` // e.g. $ref: '#/components/schemas/Pet'
}

type Responses []Response

type Response struct {
	Code        uint         `yaml:"code"`
	Description string       `yaml:"description"`
	Content     ContentTypes `yaml:"content"`
}

type SecurityEntities []Security

type Security struct {
	AuthName  string
	PermTypes []string // write:pets , read:pets etc.
}

type Components []Component

type Component struct {
	Schemas         Schemas         `yaml:"schemas"`
	SecuritySchemes SecuritySchemes `yaml:"securitySchemes"`
}

type Schemas []Schema

type Schema struct {
	Name       string
	Type       string
	Properties SchemaProperties
	XML        XMLEntry `yaml:"xml, omitempty"`
	Ref        string   // $ref: '#/components/schemas/Pet' // TODO: Should this be omitted if empty?
}

type XMLEntry struct {
	Name string
}
type SchemaProperties []SchemaProperty

type SchemaProperty struct {
	Type        string      // OAS3.0 data types - e.g. integer, boolean, string
	Format      string      `yaml:"format,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Enum        []string    `yaml:"enum,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
}

type SecuritySchemes []SecurityScheme

type SecurityScheme struct { // TODO: Lots of variants, yet to be researched.
	Name  string
	Type  string
	In    string
	Flows interface{}
}
