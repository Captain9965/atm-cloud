package restapi

import (
	dbapp "cloud/dbApp"
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"time"
)

func RegisterTransactionsApi(r *gin.Engine){
	// create a transaction
	r.POST("/transactions/create", func(c *gin.Context) {
		createTransaction(c)
	})
	// udpate transaction
	r.POST("/transactions/update", func(c *gin.Context) {
		updateTransaction(c)
	})
	  // get a slice of transactions
	r.GET("/transactions/get", func(c *gin.Context) {
		getTransactions(c)
	})
	//delete transaction
	r.DELETE("/transactions/delete", func(c *gin.Context) {
		deleteTransaction(c)
	})
}

func createTransaction(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	// Bind JSON from request body
	var newTransaction struct {
		TransactionType string `json:"transaction_type" binding:"required"`
		TransactionCode string  `json:"transaction_code" binding:"required"`
		TransactionStatus int `json:"transaction_status" binding:"required"`
		MachineUUID string 	`json:"machine_uuid" binding:"required"`
		Amount 		int	`json:"amount" binding:"required"`
	}

	if err := c.Bind(&newTransaction); err != nil {
		fmt.Println(newTransaction)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := Db.CreateTransaction(map[string]interface{}{
		"transaction_type": newTransaction.TransactionType,
		"transaction_code": newTransaction.TransactionCode,
		"transaction_status": newTransaction.TransactionStatus,
		"machine_uuid": newTransaction.MachineUUID,
		"amount": newTransaction.Amount,
		"transaction_uuid": uuid.New(),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction", "Description":err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully"})

}

func updateTransaction(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)

	// Bind JSON from request body
	var updatedTransaction struct {
		TransactionUUID string  	`json:"transaction_uuid" binding:"required"`
		NewAmount int  		`json:"new_amount"`
		NewTransactionStatus int `json:"new_status"`
	}

	if err := c.Bind(&updatedTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uuid, err := uuid.Parse(updatedTransaction.TransactionUUID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Transaction UUID format"})
		return
	}
	
	updatedData := map[string]interface{}{
		"new_amount": updatedTransaction.NewAmount,
		"new_transaction_status": updatedTransaction.NewTransactionStatus,

	}
	
	err = Db.UpdateTransactionByUUID(uuid, updatedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction", "Description":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}

func getTransactions(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	
	// Get from and to time parameters from the request
	fromStr := c.Query("from")
	toStr := c.Query("to")
  
	// Parse time strings into time.Time format
	from, err := time.Parse(time.RFC3339, fromStr)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid from time format", "Description":err.Error()})
	  	return
	}
	to, err := time.Parse(time.RFC3339, toStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid to time format", "Description":err.Error()})
	  	return
	}
  
	if to.Before(from) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "To must come after from", "Description": "Nil"})
		return
		
	}
	// Get transactions in the time range
	transactions, err := Db.GetTransactions(from, to)

	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions", "Description":err.Error()})
	  return
	
	}
	c.JSON(http.StatusOK, transactions)
}

func deleteTransaction(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)

	var transaction struct{
		TransactionUUID string `json:"Transaction_uuid" binding:"required"`
	}
	if err := c.Bind(&transaction); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
	}
	uuid, err := uuid.Parse(transaction.TransactionUUID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Transaction UUID format"})
		return
	}
	err = Db.DeleteTransactionByUUID(uuid)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete transaction", "Description":err.Error()})
	  return
	}
	c.JSON(http.StatusOK, gin.H{"transaction deleted": transaction.TransactionUUID})
}