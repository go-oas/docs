package docs

func New() OAS {
	initRoutes := RegRoutes{}

	return OAS{
		RegisteredRoutes: initRoutes,
	}
}

const (
	OASAnnotationInit = "// @OAS "
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

// Version is represented in SemVer format.
type (
	Version    string
	URL        string
	OASVersion Version
)

type Info struct {
	Title          string  `yaml:"title"`
	Description    string  `yaml:"description"`
	TermsOfService URL     `yaml:"termsOfService"`
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
	Route           string           `yaml:"route"`
	HTTPMethod      string           `yaml:"httpMethod"`
	Tags            []string         `yaml:"tags"`
	Summary         string           `yaml:"summary"`
	OperationID     string           `yaml:"operationId"`
	RequestBody     RequestBody      `yaml:"requestBody"`
	Responses       Responses        `yaml:"responses"`
	Security        SecurityEntities `yaml:"security,omitempty"`
	HandlerFuncName string           `yaml:"-"`
}

type RequestBody struct {
	Description string       `yaml:"description"`
	Content     ContentTypes `yaml:"content"`
	Required    bool         `yaml:"required"`
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
	Name        string      `yaml:"-"`
	Type        string      // OAS3.0 data types - e.g. integer, boolean, string
	Format      string      `yaml:"format,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Enum        []string    `yaml:"enum,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
}

type SecuritySchemes []SecurityScheme

type SecurityScheme struct {
	Name  string        `yaml:"name,omitempty"`
	Type  string        `yaml:"type,omitempty"`
	In    string        `yaml:"in,omitempty"`
	Flows SecurityFlows `yaml:"flows,omitempty"`
}

type SecurityFlows []SecurityFlow

type SecurityFlow struct {
	Type    string         `yaml:"type,omitempty"`
	AuthURL URL            `yaml:"authorizationUrl,omitempty"`
	Scopes  SecurityScopes `yaml:"scopes,omitempty"`
}

type SecurityScopes []SecurityScope

type SecurityScope struct {
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

func (ed *ExternalDocs) isEmpty() bool {
	if ed == nil {
		return true
	}

	if ed.Description == emptyStr && ed.URL == emptyStr {
		return true
	}

	return false
}
