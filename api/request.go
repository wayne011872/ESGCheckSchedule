package api

import (
	"fmt"
	"time"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/wayne011872/ESGCheckSchedule/dao"
)

type invoiceRequestBody struct{
	Seller 		string    `json:"seller"`
	InvoiceNo 	[]string  `json:"invoice_no"`
}

type invoiceResponseBody struct{
	Status     string				`json:"status"`
	Msg        string   			`json:"msg"`
	Result     []*dao.InvoiceStatus	`json:"result"`
}
type invoiceIssueNo struct {
	InvoiceNo	string		`json:"invoice_no"`
}

type invoiceIssueResponseBody struct {
	Status     	string			`json:"status"`
	Msg        	string   		`json:"msg"`
	Result		*invoiceIssueNo	`json:"result"`
}

func RequestPostInvoiceStatus(i []string)([]*dao.InvoiceStatus,error){
	fmt.Printf("[%s] Send Post Turnkey Invoice Status To EinvoiceCenter\n",time.Now().Format("2006-01-02 15:04:05"))
	requestBody := &invoiceRequestBody{Seller: "28703305",InvoiceNo: i}
	rb,err := json.Marshal(requestBody)
	if err !=nil {
		return nil,err
	}
	requestURI := os.Getenv(("TURNKEY_URI"))
	client := &http.Client{}
	req,err := http.NewRequest("POST",requestURI,bytes.NewReader(rb))
	if err !=nil {
		return nil,err
	}
	req.Header.Add("Authorization","2RErbrOodU77ZOREF/2+2o80E/bHA8VKhQC42A+i78=z4+f")
	resp,err := client.Do(req)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()
	body,err := io.ReadAll(resp.Body)
	if err != nil{
		return nil,err
	}
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	invoiceResponse := &invoiceResponseBody{}
	err = json.Unmarshal(body,invoiceResponse)
	if err != nil {
		return nil,err
	}
	if invoiceResponse.Status != "success" {
		err := errors.New(invoiceResponse.Msg)
		return nil,err
	}
	return invoiceResponse.Result,nil
}

func RequestPostInvoiceIssue(o *dao.Order)(string,error){
	fmt.Printf("[%s] Send Issue Invoice Request To EinvoiceCenter\n",time.Now().Format("2006-01-02 15:04:05"))
	rb,err := json.Marshal(o)
	if err !=nil {
		return "",err
	}
	requestURI := os.Getenv(("INVOICE_URI"))
	client := &http.Client{}
	req,err := http.NewRequest("POST",requestURI,bytes.NewReader(rb))
	if err !=nil {
		return "",err
	}
	req.Header.Add("Authorization","2RErbrOodU77ZOREF/2+2o80E/bHA8VKhQC42A+i78=z4+f")
	resp,err := client.Do(req)
	if err != nil {
		return "",err
	}
	defer resp.Body.Close()
	body,err := io.ReadAll(resp.Body)
	if err != nil{
		return "",err
	}
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	invoiceResponse := &invoiceIssueResponseBody{}
	err = json.Unmarshal(body,invoiceResponse)
	if err != nil {
		return "",err
	}
	if invoiceResponse.Status != "success" {
		err := errors.New(invoiceResponse.Msg)
		return "",err
	}
	return invoiceResponse.Result.InvoiceNo,nil
}