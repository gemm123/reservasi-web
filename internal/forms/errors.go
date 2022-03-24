package forms

type errors map[string][]string

//menambahkan pesan error untuk form
func (err errors) Add(field, message string) {
	err[field] = append(err[field], message)
}

//mengambil pesan error pertama
func (err errors) Get(field string) string {
	errString := err[field]
	if len(errString) == 0 {
		return ""
	}
	return errString[0]
}
