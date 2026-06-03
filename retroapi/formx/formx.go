package formx

import "github.com/go-playground/form/v4"

func NewDecoder() (d *form.Decoder) {
	d = form.NewDecoder()
	d.SetTagName("json")
	d.SetNamespacePrefix("[")
	d.SetNamespaceSuffix("]")
	return d
}

func NewEncoder() (e *form.Encoder) {
	e = form.NewEncoder()
	e.SetTagName("json")
	e.SetNamespacePrefix("[")
	e.SetNamespaceSuffix("]")

	return e
}
