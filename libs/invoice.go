package libs

import (
	"fmt"
	"math"
	"regexp"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/wayne011872/ESGCheckSchedule/dao"
	"github.com/wayne011872/ESGCheckSchedule/mail"
)

func CheckInvoiceError(paidInvoice []*dao.InvoiceStatus) {
	invoiceErrStr := ""
	invoiceOver12Str := ""
	for _,invoice := range paidInvoice{
		var timeLayoutStr = "2006-01-02 15:04:05"
		dateTime := TransferStrToTimeStr(invoice.InvoiceDateTime[0:14])
		st, _ := time.Parse(timeLayoutStr, dateTime)
		nt := time.Now()
		sub := nt.Sub(st)
		if invoice.InvoiceStatus == "E" {
			mailContent := fmt.Sprintf("<h3><strong>--------傳承發票開立錯誤--------</strong></h3></br><p>錯誤發票號碼 :%s</p></br><p>錯誤發票狀態 :%s</p></br><p>錯誤發票類型 :%s</p></br><p>錯誤發票日期時間 :%s</p>", invoice.InvoiceNo,invoice.InvoiceStatus,invoice.InvoiceMessageType,dateTime)
			mail.SendMail("傳承訂單發票異常通知",mailContent)
			invoiceErrStr += invoice.InvoiceNo + "\n"
		}else if (invoice.InvoiceMessageType != "A0101") && (invoice.InvoiceStatus == "G") && (sub.Hours() >= 12) {
			mailContent := fmt.Sprintf("<h3><strong>--傳承發票開立超過12小時未上傳--</strong></h3></br><p>錯誤發票號碼 :%s</p></br><p>錯誤發票狀態 :%s</p></br><p>錯誤發票類型 :%s<p></br><p>錯誤發票日期時間 :%s</p></br>", invoice.InvoiceNo,invoice.InvoiceStatus,invoice.InvoiceMessageType,dateTime)
			mail.SendMail("傳承訂單發票異常通知",mailContent)
			invoiceOver12Str += invoice.InvoiceNo + "\n"
		}
	}
	if invoiceErrStr == "" && invoiceOver12Str == ""{
		fmt.Printf("[%s] 無檢測到發票異常\n",time.Now().Format("2006-01-02 15:04:05"))
	}else{
		if invoiceErrStr != "" {
			fmt.Printf("[%s]檢測到發票異常 發票號碼: %s",time.Now().Format("2006-01-02 15:04:05"),invoiceErrStr)
		}
		if invoiceOver12Str != "" {
			fmt.Printf("[%s]檢測到發票超過12小時未上傳 發票號碼: %s",time.Now().Format("2006-01-02 15:04:05"),invoiceOver12Str)
		}
	}
}

func TransferStrToTimeStr(dateTime string)string {
	timeStr := dateTime[0:4] + "-" + dateTime[4:6] + "-" + dateTime[6:8] + " " + dateTime[8:10] + ":" + dateTime[10:12] + ":" + dateTime[12:14]
	return timeStr
}

func GetBuyerAddress(order *dao.Order) {
	pattern := "\\d+"
    isCityContainNum,_ := regexp.MatchString(pattern,order.PaymentCity)
	isZoneContainNum,_ := regexp.MatchString(pattern,order.PaymentZone)
	order.BuyerAddress = order.PaymentAddress
	if !strings.Contains(order.PaymentAddress, order.PaymentCity) && !isCityContainNum{
		order.BuyerAddress = order.PaymentCity + order.BuyerAddress
	}
	if !strings.Contains(order.PaymentAddress, order.PaymentZone) && !isZoneContainNum{
		order.BuyerAddress = order.PaymentZone + order.BuyerAddress
	}
}

func TransferToPostInvoice(orders []*dao.Order) {
	for _,o := range orders {
		o.UserId = "37"
		o.Seller = "28703305"
		o.IsExchange = "0"
		o.RandomNum = createCode()
		o.TaxRate = 0.05
		if o.DonateCode == "" {
			o.IsDonate = "N"
		}else{
			o.IsDonate = "Y"
		}
		o.IsPrint = "N"
		if o.ReceiptTitle == "" {
			o.BuyerName = o.LastName + o.FirstName
		} else {
			o.BuyerName = o.ReceiptTitle
		}
		
		o.InvoiceTo = "C"
		o.TaxType = 1
		totalAmount,_ := strconv.ParseFloat(o.TotalAmount,32)
		o.SalesAmount = float64(math.Round(totalAmount / 1.05))
		o.TaxAmount = totalAmount - o.SalesAmount
		GetBuyerAddress(o)
	}
}

func createCode() string {
	return fmt.Sprintf("%04v",rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}