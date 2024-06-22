package server

type InvoiceService interface {
	FetchData() []map[string]interface{}
	WriteBack(resp string)
}
