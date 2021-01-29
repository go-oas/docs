package docs

import "fmt"

// MapAnnotationsInPath scanIn is relevant from initiator calling it.
//
// It accepts the path in which to scan for annotations within Go files.
func (o *OAS) MapAnnotationsInPath(scanIn string) error {
	filesInPath, err := scanForChangesInPath(scanIn)
	if err != nil {
		return fmt.Errorf(" :%w", err)
	}

	for _, file := range filesInPath {
		err = o.mapDocAnnotations(file)
		if err != nil {
			return fmt.Errorf(" :%w", err)
		}
	}

	return nil
}
