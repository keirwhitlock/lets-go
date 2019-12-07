package forms

// errors type to hold error messages for forms
type errors map[string][]string

// add and get methods on the errors type.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
