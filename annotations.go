package docs

// MapAnnotationsInPath scanIn is relevant from initiator calling it.
func (o *OAS) MapAnnotationsInPath(scanIn string) error {
	filesInPath, err := scanForChangesInPath(scanIn)
	if err != nil {
		return err
	}

	for _, file := range filesInPath {
		err = o.MapDocAnnotations(file)
		if err != nil {
			return err
		}
	}

	return nil
}
