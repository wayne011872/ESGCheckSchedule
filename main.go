package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/wayne011872/ESGCheckSchedule/api"
	"github.com/wayne011872/ESGCheckSchedule/dao"
	"github.com/wayne011872/ESGCheckSchedule/libs"
	"github.com/wayne011872/ESGCheckSchedule/mail"
	"github.com/wayne011872/ESGCheckSchedule/model/order"
	"github.com/wayne011872/goSterna"
)

var (
	service = flag.String("s", "cli", "service (auto, cli)")
)

func main() {
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	switch *service {
	case "auto":
		detectInv := os.Getenv(("DETECT_INTERVAL"))
		detectInvInt, _ := strconv.Atoi(detectInv)
		fmt.Printf("[%s] 傳承發票檢查排程啟動\n", time.Now().Format("2006-01-02 15:04:05"))
		for {
			err := runAutoCheck()
			if err != nil {
				panic(err)
			}
			fmt.Printf("---------------------------------每三小時檢驗一次(下一次檢驗時間 :%s)-----------------------\n", time.Now().Add(3*time.Hour).Format("2006-01-02 15:04:05"))
			time.Sleep(time.Duration(detectInvInt) * time.Hour)
		}
	case "cli":
		var serviceType int
		fmt.Println("環球傳承發票排程發票檢驗服務 請輸入要執行的服務:")
		fmt.Println("(0) 檢驗傳承發票")
		fmt.Scan(&serviceType)
		switch serviceType {
		case 0:
			err := runAutoCheck()
			if err != nil {
				mailContent := fmt.Sprintf("<h3><strong>--------傳承發票排程錯誤通知--------</strong></h3></br><p>以下為錯誤訊息 :%s</p></br>", err)
				mail.SendMail("傳承發票排程錯誤通知", mailContent)
				panic(err)
			}
		}
	}
}

func runAutoCheck() error {
	di := &dao.Di{}
	goSterna.InitDefaultConf(".", di)
	unpaidStatusCode := []uint16{1, 4, 7, 8, 10, 14, 16}
	crud := order.NewCRUD(di)
	paidOrders, _ := crud.FindPaid(unpaidStatusCode)
	issuedInvoiceNo := []string{}
	notIssuedOrderId := []string{}
	if len(paidOrders) != 0 {
		for _, order := range paidOrders {
			if len(order.InvoiceNo) != 0 && order.InvoiceNo != "0" {
				issuedInvoiceNo = append(issuedInvoiceNo, order.GetFullInvoiceNo())
			} else {
				notIssuedOrderId = append(notIssuedOrderId, order.OrderId)
			}
		}
		if len(issuedInvoiceNo) > 0 {
			issuedInvoice, err := api.RequestPostInvoiceStatus(issuedInvoiceNo)
			if err != nil && err.Error() != "Not Found" {
				return err
			}
			if len(issuedInvoice) > 0 {
				var invoiceCheckMessage string
				fmt.Printf("[%s] 共有%d筆訂單發票需要檢測\n", time.Now().Format("2006-01-02 15:04:05"),len(issuedInvoice))
				for _,i := range issuedInvoice {
					invoiceCheckMessage += fmt.Sprintf("發票號碼 :%s",i.InvoiceNo)
				}
				checkedInvoiceNo := libs.CheckInvoiceError(issuedInvoice)
				if len(checkedInvoiceNo) > 0 {
					crud.UpdateIsCheck(checkedInvoiceNo)
				}
			} else {
				fmt.Printf("[%s] 沒有發票需要檢測\n", time.Now().Format("2006-01-02 15:04:05"))
			}
		}
		if len(notIssuedOrderId) > 0 {
			var notIssuedMessage string
			for _, o := range notIssuedOrderId {
				notIssuedMessage += fmt.Sprintf("OrderId :%s\n", o)
			}
			fmt.Printf("[%s] 檢驗到未開立發票訂單 \n共有%d筆\n%s", time.Now().Format("2006-01-02 15:04:05"), len(notIssuedOrderId), notIssuedMessage)
			notIssuedOrders, err := crud.FindOrdersNotIssue(notIssuedOrderId)
			if err != nil {
				return err
			}
			uniformList, err := api.RequestPostCheckBan(notIssuedOrders)
			if err != nil {
				return err
			}
			libs.CheckBanProfit(uniformList, notIssuedOrders)
			libs.TransferToPostInvoice(notIssuedOrders)
			invoiceMailContent := ""
			for _, order := range notIssuedOrders {
				err, result := api.RequestPostInvoiceIssue(order)
				if err != nil {
					return err
				}
				var issuedInvoicePre,issuedInvoiceNo string
				if len(result) == 10 {
					issuedInvoicePre = result[0:2]
					issuedInvoiceNo = result[2:10]
				}else {
					return errors.New("issued invoice result length less than 10")
				}
				crud.UpdateIssuedOrder(order.OrderId,issuedInvoicePre,issuedInvoiceNo)
				invoiceMailContent += "<p>" + result + "</p></br>"
			}
			mailContent := fmt.Sprintf("<h3><strong>--------傳承發票補開立通知--------</strong></h3></br><p>補開立發票號碼 :%s</p></br>", invoiceMailContent)
			mail.SendMail("傳承訂單發票補開立通知", mailContent)
		}
	}
	return nil
}