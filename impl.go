package docs

import (
	"bufio"
	"errors"
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
func scanForChangesInPath(handlersPath string) ([]string, error) {
	var files []string

	currentPath, err := os.Getwd()
	if err != nil {
		return files, err
	}

	pathToScan := filepath.Join(currentPath, handlersPath)

	err = filepath.Walk(pathToScan, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files, nil
}

const OASAnnotationInit = "// @OAS "

func (o *OAS) mapDocAnnotations(path string) error {
	if o == nil {
		return errors.New("ptr to OASHandlers can not be nil")
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close() // FIXME: Consume this error.

	scanner := bufio.NewScanner(f)

	line := 1

	for scanner.Scan() {
		lineText := scanner.Text()
		if strings.Contains(lineText, OASAnnotationInit) {
			// TODO: Can this be more performance cautious?
			fields := strings.Fields(lineText)

			// TODO: Implement getters for these fields?
			var newRoute Path
			newRoute.handlerFuncName = fields[2]
			newRoute.Route = fields[3]
			newRoute.HTTPMethod = fields[4]

			o.Paths = append(o.Paths, newRoute)

			return nil
		}

		line++
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	return nil
}
