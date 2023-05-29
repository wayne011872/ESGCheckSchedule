package dao

type InvoiceStatus struct {
	InvoiceNo  				string	`json:"invoice_no"`
	InvoiceStatus 			string	`json:"STATUS"`
	InvoiceCategoryType  	string	`json:"CATEGORY_TYPE"`
	InvoiceMessageType		string	`json:"MESSAGE_TYPE"`
	InvoiceDateTime  		string  `json:"MESSAGE_DTS"`
}