package docs

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	goFileExt = ".go"
)

type configAnnotation struct {
	getWD getWorkingDirFn
}

type (
	getWorkingDirFn func() (dir string, err error)
	pathWalkerFn    func(path string, walker walkerFn) (files []string, err error)
	walkerFn        func(root string, walkFn filepath.WalkFunc) error
)

// MapAnnotationsInPath scanIn is relevant from initiator calling it.
//
// It accepts the path in which to scan for annotations within Go files.
func (oas *OAS) MapAnnotationsInPath(scanIn string, conf ...configAnnotation) error {
	filesInPath, err := scanForChangesInPath(scanIn, getWDFn(conf), walkFilepath)
	if err != nil {
		return fmt.Errorf(" :%w", err)
	}

	for _, file := range filesInPath {
		err = oas.mapDocAnnotations(file)
		if err != nil {
			return fmt.Errorf(" :%w", err)
		}
	}

	return nil
}

// scanForChangesInPath scans for annotations changes on handlers in passed path,
// which is relative to the caller's point of view.
func scanForChangesInPath(handlersPath string, getWD getWorkingDirFn, walker pathWalkerFn) (files []string, err error) {
	currentPath, err := getWD()
	if err != nil {
		return files, fmt.Errorf("failed getting current working directory: %w", err)
	}

	files, err = walker(filepath.Join(currentPath, handlersPath), filepath.Walk)
	if err != nil {
		return files, fmt.Errorf("failed walking tree of the given path: %w", err)
	}

	return files, nil
}

func (ca configAnnotation) getCurrentDirFetcher() getWorkingDirFn {
	if ca.getWD != nil {
		return ca.getWD
	}

	return os.Getwd
}

func getWDFn(configs []configAnnotation) getWorkingDirFn {
	if len(configs) != 0 {
		return configs[0].getCurrentDirFetcher()
	}

	return os.Getwd
}

func walkFilepath(pathToTraverse string, walker walkerFn) ([]string, error) {
	var files []string

	walkFn := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != goFileExt {
			return nil
		}

		files = append(files, path)

		return nil
	}

	err := walker(pathToTraverse, walkFn)
	if err != nil {
		return files, err
	}

	return files, nil
}

func (oas *OAS) mapDocAnnotations(path string) error {
	if oas == nil {
		return errors.New("pointer to OASHandlers can not be nil")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file in path %s :%w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	line := 1

	for scanner.Scan() {
		mapIfLineContainsOASTag(scanner.Text(), oas)
		line++
	}

	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("scanner failure :%w", err)
	}

	return nil
}

func mapIfLineContainsOASTag(lineText string, o *OAS) {
	if strings.Contains(lineText, oasAnnotationInit) {
		annotations := oasAnnotations(strings.Fields(lineText))

		var newRoute Path
		newRoute.HandlerFuncName = annotations.getHandlerFuncName()
		newRoute.Route = annotations.getRoute()
		newRoute.HTTPMethod = annotations.getHTTPMethod()

		o.Paths = append(o.Paths, newRoute)
	}
}

type oasAnnotations []string

func (oa oasAnnotations) getHandlerFuncName() string {
	return oa[2]
}

func (oa oasAnnotations) getRoute() string {
	return oa[3]
}

func (oa oasAnnotations) getHTTPMethod() string {
	return oa[4]
}
