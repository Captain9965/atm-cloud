package dbApp

import (
	"encoding/json"
	// "time"
	// "github.com/gin-gonic/gin"
	"github.com/google/uuid" // Import for UUID generation
	"gorm.io/gorm"
)

type GormDB struct {
	*gorm.DB //gorm.DB object
  }

type User struct {
	gorm.Model
	Username       string       `json:"username"`
	UserPhone      string       `json:"user_phone_number"`       
	Machines  	   []Machine    `gorm:"foreignKey:UserID"`  //one 2 many
	Role           string       `json:"role"`                      // Permission level 
	OrganizationID uint         `json:"organization_id"`           // Foreign key to Organizations
	Organization   Organization 
	Password string       		`json:"-"`                         // Excluded from JSON marshaling
}

type Organization struct {
	gorm.Model
	OrganizationName 	string  `json:"organization_name"`         //name of organisation
	Users          		[]User `gorm:"foreignKey:OrganizationID"` // One-to-Many relationship with Users
}

type Machine struct {
	gorm.Model
	MachineType           string          `json:"machine_type"`
	MachineSerial		  uint 			  `json:"machine_serial"`
	MachineUUID			  string          `json:"machine_uuid"`
	UserID                uint            `json:"user_id"`     			// foreign key to Users 
	User                  User            
	AdminCardID           string          `json:"admin_card_id"`
	VendingCardID         string          `json:"vending_card_id"`
	ServiceCardID         string          `json:"service_card_id"`
	TapEnableStatus       int             `json:"tap_enable_status"`
	AdminCash             int             `json:"admin_cash"`
	Transactions          []Transactions  `gorm:"foreignKey:MachineID"` 	// One-to-Many relationship with Transactions
	Events				  []Events        `gorm:"foreignKey:MachineID"`		// One-to-Many relationship with Events
}

type Transactions struct {
	gorm.Model
	Machine           Machine   
	MachineID         uint      `json:"machine_id"`
	TransactionType   int       `json:"transaction_type"`				//type of transaction
	Amount            int       `json:"amount"`							// transaction amount
	TransactionUUID   uuid.UUID `gorm:"type:uuid"`						// UUID of transaction
	TransactionCode   string    `json:"transaction_code"`            	// Optional field
	TransactionStatus int       `json:"transaction_status"`				// 1- pending, 0 - fulfilled, -1 -failed, -2, cancelled
}

type Events struct{
	gorm.Model
	MachineID			uint 	`json:"machine_id"` 					//Foreign key to Machine
	Machine				Machine 
	Rss                 int             `json:"rss"`         			// Signal strength
	Location            json.RawMessage `json:"location"`            	// Cell tower location as JSON
	EventType			string `json:"event_type"`						//EventType eg. error, heartbeat, pay
	Data 				json.RawMessage `json:"data"`					//Data in the event
}

type Database interface {
	// Db operations
	Connect() error  // Connects to the database

	//users
	CreateUser(userData map[string]interface{}) error  // Creates a new user
  	GetUserByName(username string) (map[string]interface{}, error)    // Retrieves user details by name (returns a map)
  	UserExists(username string) bool //method to check user existence
  	UpdateUserByName(username string, fieldsToUpdate map[string]interface{}) error //method to update by name
  	DeleteUserByName(username string) error                        // Deletes a user by name
	GetAllUsers() ([]map[string]interface{}, error)
	AuthenticateUserByName(username string, password string) (bool,error)

	//organization
	CreateOrganization(orgData map[string]interface{}) error  // Creates a new organization
	GetOrganizationByName(orgname string) (map[string]interface{}, error)    // Retrieves org details by name (returns a map)
	OrganizationExists(orgname string) bool //method to check organization's existence
	UpdateOrgByName(orgname string, fieldsToUpdate map[string]interface{}) error // update org by name
	DeleteOrgByName(orgname string) error
	GetAllOrganizations()([]map[string]interface{}, error)
	
	//machine
	CreateMachine(machineData map[string]interface{})error // creates a new machine
	GetMachineByID(machineUUID string)(map[string]interface{}, error) // Retrieves machine by its unique id...usually uuid of mcu or modem IMEI
	MachineExists(machineUUID string)bool // Check whether a machine exists
	UpdateMachinebyID(machineUUID string, fieldsToUpdate map[string]interface{}) error //update machine fields by ID
	DeleteMachineByID(machineUUID string)error
	GetAllMachines()([]map[string]interface{}, error)

	// //transactions
	// CreateTransaction(transactionData map[string]interface{})error // creates a new transactions from payment api or mqtt
	// GetTransactionByUUID(uid uuid.UUID)(map[string]interface{}, error) // Retrieve transaction by UID
	// TransactionExists(uid uuid.UUID)bool
	// UpdateTransactionByUUID(uid uuid.UUID, fieldsToUpdate map[string]interface{}) error
	// DeleteTransactionByUUID(uid uuid.UUID)error
	// GetTransactions(from time.Time, to time.Time)([]map[string]interface{}, error)

	// //events
	// CreateEvent(eventData map[string]interface{})error
	// GetEventByID(eventID uint)(map[string]interface{}, error)
	// EventExists(eventID uint)bool
	// GetEvents(from time.Time, to time.Time)([]map[string]interface{}, error)
  }
