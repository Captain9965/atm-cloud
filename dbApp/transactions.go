package dbApp

import (
	"fmt"
	"time"
	"github.com/google/uuid"
)

func (db *GormDB) CreateTransaction(transactionData map[string]interface{})error {
	

	transactionType, ok := transactionData["transaction_type"].(string)
	if !ok || transactionType == "" {
	  return fmt.Errorf("missing or invalid transaction type in transaction data")
	}

	amount, ok := transactionData["amount"].(int)
	if !ok {
		return fmt.Errorf("missing or invalid transaction amount")
	}

	transactionCode, ok := transactionData["transaction_code"].(string)
	if !ok || transactionCode == "" {
		return fmt.Errorf("missing or invalid transaction code")
	}

	transactionStatus, ok := transactionData["transaction_status"].(int)
	if !ok {
		return fmt.Errorf("missing or invalid transaction amount")
	}

	transactionUUID, ok := transactionData["transaction_uuid"].(uuid.UUID)
	if !ok {
		return fmt.Errorf("missing or invalid transaction uuid")
	}

	machineUUID, ok := transactionData["machine_uuid"].(string)
	if !ok || machineUUID == "" {
		return fmt.Errorf("missing or invalid machine uuid in transaction data")
	}

	//check if transaction already exists:
	exists, err := db.TransactionExists(&transactionUUID, &transactionCode)
	if err != nil{
		return err
	}
	if exists{
		return fmt.Errorf("duplicate transaction")
	}
	// Preload machine that transaction belongs to:
	var machine Machine
	result := db.Where("machine_uuid = ?", machineUUID).Preload("Transactions").First(&machine)
	if result.Error != nil {
		return result.Error
	}

	transaction := Transactions{
		TransactionType: transactionType,
		Amount: amount,
		TransactionCode: transactionCode,
		TransactionStatus: transactionStatus,
		TransactionUUID: transactionUUID,
		Machine: machine,
	 }
	 err = db.Create(&transaction).Error
	 if err != nil {
	   return err
	 }
  
	return nil
}

func (db *GormDB) TransactionExists(uid *uuid.UUID, code *string)(bool, error) {
	var count int64
	query := db.Model(&Transactions{})

	if uid != nil {
		query = query.Where("transaction_uuid = ?", *uid)
	  }
	
	  if code != nil {
		if uid != nil {
		  query = query.Or("transaction_code = ?", *code)
		} else {
		  query = query.Where("transaction_code = ?", *code)
		}

	  }
	  
	  if code == nil && uid == nil{
		return false, fmt.Errorf("either transaction_uuid or transaction_code must be provided")
	  }

	result := query.Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
    fmt.Printf("Count is %d", count)
	return count > 0, nil
}

func (db *GormDB) GetTransactionByUUID(uid uuid.UUID)(map[string]interface{}, error) {
	var transaction Transactions
	exists, err:= db.TransactionExists(&uid, nil)
	if err != nil {
		return nil, err
	}
    if !exists{
		return nil, fmt.Errorf("transaction does not exist")
	}
	result := db.Preload("Machine").Where("transaction_uuid = ?", uid).First(&transaction)
	if result.Error != nil {
	  return nil, result.Error
	}

	transactionMap := map[string]interface{}{
	  "transaction_type": transaction.TransactionType,
	  "transaction_uuid": transaction.TransactionUUID,
	  "Machine": transaction.Machine,
	  "transaction_status": transaction.TransactionStatus,
	  "transaction_code": transaction.TransactionCode,
	}
  
	return transactionMap, nil
}

func (db *GormDB) UpdateTransactionByUUID(uid uuid.UUID, fieldsToUpdate map[string]interface{}) error {
	exists, err:= db.TransactionExists(&uid, nil)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("transaction does not exist")
	}
		
	// Get existing transaction:
	var existingTransaction Transactions
	result := db.DB.Where("transaction_uuid = ?", uid).First(&existingTransaction)
	if result.Error != nil {
	return result.Error
	}

	// Update user fields based on the provided map
	for field, newValue := range fieldsToUpdate {
		switch field {
		case "new_transaction_code":
			existingTransaction.TransactionCode = newValue.(string)
		case "new_transaction_status":
			existingTransaction.TransactionStatus = newValue.(int)
		case "new_amount":
			existingTransaction.Amount = newValue.(int)
		default:
			// Handle updates for other supported fields (if any)
		}
	}

	result = db.DB.Save(&existingTransaction)
	if result.Error != nil {
	return result.Error
	}
	return nil
}

func (db *GormDB) DeleteTransactionByUUID(uid uuid.UUID)error {
	exists, err:= db.TransactionExists(&uid, nil)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("transaction does not exist")
	}
	
	result := db.DB.Where("transaction_uuid = ?", uid).Delete(&Transactions{}) 
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *GormDB) GetTransactions(from time.Time, to time.Time)([]map[string]interface{}, error) {
	var transactions []Transactions
	result := db.Preload("Machine").Where("created_at BETWEEN ? AND ?", from, to).Find(&transactions)
	if result.Error != nil {
	  return nil, result.Error
	}
	//Convert user structs to maps
	transactionData := make([]map[string]interface{}, len(transactions))
	for i, transaction := range transactions {
	   transactionData[i] = map[string]interface{}{
		 "transaction_type": transaction.TransactionType,
		 "transaction_status": transaction.TransactionStatus,
		 "amount": transaction.Amount,
		 "transaction_status_code": transaction.TransactionCode,
		 "transaction_uuid": transaction.TransactionUUID.String(),
		 "machine_uuid": transaction.Machine.MachineUUID,
	  }
    }

	return transactionData, nil
}