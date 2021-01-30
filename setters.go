package docs

// SetOASVersion sets the OAS version, by casting string to OASVersion type.
func (o *OAS) SetOASVersion(ver string) {
	o.OASVersion = OASVersion(ver)
}

// GetInfo returns pointer to the Info struct.
func (o *OAS) GetInfo() *Info {
	return &o.Info
}

// SetContact setts the contact on the Info struct.
func (i *Info) SetContact(email string) {
	i.Contact = Contact{Email: email}
}

// SetLicense sets the license on the Info struct
func (i *Info) SetLicense(licType, url string) {
	i.License = License{
		Name: licType,
		URL:  URL(url),
	}
}
