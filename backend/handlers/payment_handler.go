package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type CheckoutRequest struct {
	TotalAmount	int64	`json:"total_amount"`
	FirstName	string	`json:"first_name"`
	LastName	string	`json:"last_name"`
	Email		string	`json:"email"`
	Phone		string	`json:"phone"`
	Note		string	`json:"note"`
}

func HandleCheckout(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request payload"})
		return
	}

	midtrans.Environment = midtrans.Sandbox

	midtrans.ServerKey = "Mid-server-ZM2VH9KSHXRCMBsC_cdP_Xk7"

	orderID := "ORDER-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	grossAmount := int64(req.TotalAmount)

	fName := req.FirstName
	if fName == "" {
		fName = "Customer"
	}
	lName := req.LastName
	if lName == "" {
		lName = "OpenPeo"
	}
	email := req.Email
	if email == "" {
		email = "test@openpeo.com"
	}
	phone := req.Phone
	if phone == "" {
		phone = "08123456789"
	}

	reqSnap := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:	orderID,
			GrossAmt:	grossAmount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:	fName,
			LName:	lName,
			Email:	email,
			Phone:	phone,
		},
		CustomField1:	req.Note,
	}

	snapResp, midErr := snap.CreateTransaction(reqSnap)
	if midErr != nil {
		log.Printf("Midtrans Error: %v", midErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":	false,
			"message":	"Failed to create transaction",
			"details":	midErr.GetMessage(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"order_id":	orderID,
		"snap_token":	snapResp.Token,
	})
}
