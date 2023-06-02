package dao

type Order struct {
	OrderId					string		`gorm:"column:order_id"`
	InvoiceNo				string		`gorm:"column:invoice_no"`
	InvoicePrefix			string		`gorm:"column:invoice_prefix"`
	InvoiceDataTime			string		`gorm:"column:receipt_date"`
	StatusId  				string		`gorm:"column:order_status_id"`
	UserId    				string		`json:"user_id"`
	Seller       			string		`json:"seller"`
	ReceiptTitle			string		`gorm:"column:receipt_title"`
	BuyerUniform            string		`json:"buyer_uniform" gorm:"column:receipt_uniform_no"`
	FirstName				string		`gorm:"column:firstname"`
	LastName				string		`gorm:"column:lastname"`
	BuyerName               string		`json:"buyer_name"`
	BuyerEmail				string		`json:"buyer_email" gorm:"column:email"`
	BuyerPhone              string		`json:"buyer_phone" gorm:"column:telephone"`
	PaymentAddress			string		`gorm:"column:payment_address_1"`
	PaymentCity				string		`gorm:"column:payment_city"`
	PaymentZone				string		`gorm:"column:payment_zone"`
	BuyerAddress			string		`json:"buyer_address"`
	CarrierType				string		`json:"carrier_type" gorm:"column:carrier_type"`
	CarrierId1				string		`json:"carrier_id1" gorm:"column:carrier_id1"`
	CarrierId2				string		`json:"carrier_id2" gorm:"column:carrier_id2"`
	IsDonate				string		`json:"is_donate"`
	DonateCode             	string		`json:"donate_code" gorm:"column:donate_no"`
	InvoiceTo				string		`json:"invoice_to"`
	IsExchange				string		`json:"is_exchange"`
	IsPrint					string		`json:"is_print"`
	RandomNum				string		`json:"random_num"`
	IsProfit                string
	IsExist					bool
	SalesAmount				float64		`json:"sales_amount"`
	TaxType					int			`json:"tax_type"`
	TaxRate					float64		`json:"tax_rate"`
	TaxAmount				float64		`json:"tax_amount"`
	TotalAmount             string		`json:"total_amount" gorm:"column:total"`
	Products				[]Product	`json:"details" gorm:"foreignKey:OrderId;references:OrderId"`
}

func (o *Order) GetFullInvoiceNo() string{
	return o.InvoicePrefix + o.InvoiceNo
}

type Product struct {
	OrderId				string			`gorm:"column:order_id"`
	ProductId			string			`json:"product_id" gorm:"column:product_id"`
	ProductName 		string			`json:"product_name" gorm:"column:name"`
	ProductQuantity 	string			`json:"quantity" gorm:"column:quantity"`
	ProductUnitPrice	string			`json:"unit_price" gorm:"column:price"`
	ProductTotalAmount  string			`json:"amount" gorm:"column:total"`
}