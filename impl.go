package docs

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

func New() OAS {
	initRoutes := RegRoutes{}

	return OAS{
		registeredRoutes: initRoutes,
	}
}

const (
	goFileExt         = ".go"
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

	registeredRoutes RegRoutes
}

type RegRoutes map[string]interface{}

func (o *OAS) AttachRoutes(fns []interface{}) {
	for _, fn := range fns {
		// TODO: Benchmark performance of this function with 1-1.5k routes.
		fnDeclaration := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		fields := strings.SplitAfter(fnDeclaration, ".")
		fnName := fields[len(fields)-1]

		o.registeredRoutes[fnName] = fn
	}
}

func (o *OAS) GetRegisteredRoutes() RegRoutes {
	return o.registeredRoutes
}

func (o *OAS) GetPathByIndex(index int) *Path {
	return &o.Paths[index]
}

// scanForChangesInPath scans for annotations changes on handlers in passed path, which is relative to the caller's point of view.
func scanForChangesInPath(handlersPath string) (files []string, err error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return files, fmt.Errorf("failed getting current working directory: %w", err)
	}

	files, err = walkFilepath(filepath.Join(currentPath, handlersPath))
	if err != nil {
		return files, fmt.Errorf("failed walking tree of the given path: %w", err)
	}

	return files, nil
}

func walkFilepath(pathToTraverse string) ([]string, error) {
	var files []string

	err := filepath.Walk(pathToTraverse, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != goFileExt {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return files, err //nolint: wrapcheck // it will be wrapped by consumer.
	}

	return files, nil
}

func (o *OAS) mapDocAnnotations(path string) error {
	if o == nil {
		return errors.New("pointer to OASHandlers can not be nil") // fixme: migrate to validator!
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file in path %s :%w", path, err)
	}
	defer f.Close() // FIXME: Consume this error.

	scanner := bufio.NewScanner(f)

	line := 1
	for scanner.Scan() {
		mapIfLineContainsOASTag(scanner.Text(), o)
		line++
	}

	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("scanner failure :%w", err)
	}

	return nil
}

func mapIfLineContainsOASTag(lineText string, o *OAS) {
	if strings.Contains(lineText, OASAnnotationInit) {
		// TODO: Can this be more performance cautious?
		fields := strings.Fields(lineText)

		// TODO: Implement getters for these fields?
		var newRoute Path
		newRoute.handlerFuncName = fields[2]
		newRoute.Route = fields[3]
		newRoute.HTTPMethod = fields[4]

		o.Paths = append(o.Paths, newRoute)
	}
}
