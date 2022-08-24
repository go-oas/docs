package docs

// SetOASVersion sets the OAS version, by casting string to OASVersion type.
func (oas *OAS) SetOASVersion(ver string) {
	oas.OASVersion = OASVersion(ver)
}

// GetInfo returns pointer to the Info struct.
func (oas *OAS) GetInfo() *Info {
	return &oas.Info
}

// SetContact setts the contact on the Info struct.
func (i *Info) SetContact(email string) {
	i.Contact = Contact{Email: email}
}

// SetLicense sets the license on the Info struct.
func (i *Info) SetLicense(licType, url string) {
	i.License = License{
		Name: licType,
		URL:  URL(url),
	}
}

// SetTag is used to define a new tag based on input params, and append it to the slice of tags its being called from.
func (tt *Tags) SetTag(name, tagDescription string, extDocs ExternalDocs) {
	var tag Tag

	if !isStrEmpty(name) {
		tag.Name = name
	}

	if !isStrEmpty(tagDescription) {
		tag.Description = tagDescription
	}

	if !extDocs.isEmpty() {
		tag.ExternalDocs = extDocs
	}

	tt.AppendTag(&tag)
}

// AppendTag is used to append an Tag to the slice of Tags its being called from.
func (tt *Tags) AppendTag(tag *Tag) {
	*tt = append(*tt, *tag)
}
