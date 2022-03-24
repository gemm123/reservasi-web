package forms

import "net/url"

//membuat custom form
type Form struct {
	url.Values
	Errors errors
}
