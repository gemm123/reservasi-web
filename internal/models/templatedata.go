package models

import "github.com/gemm123/reservasi-web/internal/forms"

//menyimpan data untuk mengirim dari handler ke templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Success         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}
