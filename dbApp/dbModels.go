package dbApp

import (
	"encoding/json"
	// "github.com/gin-gonic/gin"
	"github.com/google/uuid"     // Import for UUID generation
	"gorm.io/gorm"
)

type GormDB struct {
	*gorm.DB //gorm.DB object
  }

type User struct {
	gorm.Model
	Username       string       `json:"username"`
	UserPhone      string       `json:"user_phone_number"`       
	MachinesOwned  []Machine    `gorm:"foreignKey:OwnerID"`  //one 2 many
	Role           string       `json:"role"`                      // Permission level 
	OrganizationID uint         `json:"organization_id"`           // Foreign key to Organizations
	Organization   Organization  // Belongs to Organization
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
	MachineID             string          `json:"machine_id"`
	OwnerID               uint            `json:"owner_id"`     			// foreign key to Users 
	User                  User            `gorm:"foreignKey:OwnerID"` 		// belongs to User 
	AdminCardID           string          `json:"admin_card_id"`
	VendingCardID         string          `json:"vending_card_id"`
	ServiceCardID         string          `json:"service_card_id"`
	TapEnableStatus       int             `json:"tap_enable_status"`
	AdminCash             int             `json:"admin_cash"`
	Transactions          []Transactions  `gorm:"foreignKey:MachineID"` 	// One-to-Many relationship with Transactions
}

type Transactions struct {
	gorm.Model
	MachineID         uint      `json:"machine_id"`           			// Foreign key to Machine
	Machine           Machine   `gorm:"foreignKey:MachineID"` 			// Belongs to Machine
	TransactionType   int       `json:"transaction_type"`				//type of transaction
	Amount            int       `json:"amount"`							// transaction amount
	TransactionUUID   uuid.UUID `gorm:"type:uuid"`						// UUID of transaction
	TransactionCode   string    `json:"transaction_code"`            	// Optional field
	TransactionStatus int       `json:"transaction_status"`				// 1- pending, 0 - fulfilled, -1 -failed, -2, cancelled
}

type Events struct{
	gorm.Model
	MachineID			uint 	`json:"machine_id"` 					//Foreign key to Machine
	Machine				Machine `gorm:"foreignKey:MachineID"` 			//Belongs to Machine
	Rss                 int             `json:"rss"`                  	// Signal strength
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
  
  }
