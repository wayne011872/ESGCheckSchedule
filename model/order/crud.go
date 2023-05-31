package order

import (

	"github.com/wayne011872/ESGCheckSchedule/dao"
	"github.com/wayne011872/goSterna/log"
	"gorm.io/gorm"
)

func NewCRUD(di *dao.Di) CRUD {
	return &mysqlCRUD{
		mslDB: di.MySQLConf.NewMySQLDB(),
		log:di.LoggerConf.NewLogger("order"),
	}

}

type CRUD interface {
	FindOrdersNotIssue(orderId []string) ([]*dao.Order,error)
	FindPaid(paidStatusCode []uint16) ([]*dao.Order,error)
	UpdateIssuedOrder(orderId string ,invoicePre string, invoiceNo string)
	UpdateIsCheck(isCheckNo []string)
}

type mysqlCRUD struct {
	mslDB *gorm.DB
	log log.Logger
}

func (m *mysqlCRUD) FindOrdersNotIssue(orderId []string) ([]*dao.Order,error){
	var orders []*dao.Order
	selectColumn := []string{"order.order_id","order.receipt_title","order.receipt_uniform_no","order.firstname","order.lastname","order.email","order.telephone","order.carrier_type","order.carrier_id1","order.carrier_id2","order.donate_no","order.total"}
	result := m.mslDB.Table("order").Select(selectColumn).Where("order_id IN ?",orderId).Find(&orders)
	if result.Error != nil {
		return nil,result.Error
	}
	for _,o := range orders {
		o.Products = m.findProduct(o.OrderId)
	}
	return orders,nil
}

func (m *mysqlCRUD) FindPaid(unpaidStatusCode []uint16) ([]*dao.Order,error){
	var orders []*dao.Order
	selectColumn := []string{"order_id","invoice_no","invoice_prefix","order_status_id","receipt_date"}
	result := m.mslDB.Table("order").Select(selectColumn).Where("order_status_id NOT IN ? AND epayment_amount is not null AND isCheck = 0",unpaidStatusCode).Find(&orders)
	return orders,result.Error
}

func (m *mysqlCRUD) findProduct(orderId string)[]dao.Product{
	var product []dao.Product
	selectColumn := []string{"order_id","product_id","name","quantity","price","total"}
	m.mslDB.Table("order_product").Select(selectColumn).Where("order_id = ?",orderId).Find(&product)
	return product
}

func (m *mysqlCRUD) UpdateIssuedOrder(orderId string ,invoicePre string, invoiceNo string) {
	m.mslDB.Table("order").Where("order_id = ?",orderId).Update("invoice_prefix",invoicePre)
	m.mslDB.Table("order").Where("order_id = ?",orderId).Update("invoice_no",invoiceNo)
}

func (m *mysqlCRUD) UpdateIsCheck(checkedNo []string) {
	m.mslDB.Table("order").Where("CONCAT(`invoice_prefix`,`invoice_no`) IN ?",checkedNo).Update("isCheck",1)
}