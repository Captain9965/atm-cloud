package dbApp

import (
	"encoding/json"
	"fmt"
	// "github.com/gin-gonic/gin"
	"github.com/google/uuid"     // Import for UUID generation
	"github.com/joho/godotenv"   // Import for loading environment variables
	"golang.org/x/crypto/bcrypt" // Import for password hashing
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"os"
	"time"
)

type GormDB struct {
	*gorm.DB //gorm.DB object
  }

type User struct {
	gorm.Model
	Username       string       `json:"username"`
	UserPhone      string       `json:"user_phone_number"`         // Optional field
	MachinesOwned  []Machine    `gorm:"foreignKey:OwnerID"`  //one 2 many
	Role           int          `json:"role"`                      // Permission level (1 - Superuser, 2 - Admin, 3 - Owner, 4 - Operator)
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
	OwnerID               uint            `json:"owner_id"`     // foreign key to Users 
	User                  User            `gorm:"foreignKey:OwnerID"` // belongs to User 
	AdminCardID           string          `json:"admin_card_id"`
	VendingCardID         string          `json:"vending_card_id"`
	ServiceCardID         string          `json:"service_card_id"`
	MachineLastUpdateTime time.Time       `json:"machine_last_update_time"`
	TapEnableStatus       int             `json:"tap_enable_status"`
	AdminCash             int             `json:"admin_cash"`
	Rss                   int             `json:"rss"`                  // Signal strength
	Locations             json.RawMessage `json:"locations"`            // Cell tower location as JSON
	Transactions          []Transactions  `gorm:"foreignKey:MachineID"` // One-to-Many relationship with Transactions
}

type Transactions struct {
	gorm.Model
	MachineID         uint      `json:"machine_id"`           // Foreign key to Machine
	Machine           Machine   `gorm:"foreignKey:MachineID"` // Belongs to Machine
	TransactionType   int       `json:"transaction_type"`
	Amount            int       `json:"amount"`
	TransactionUUID   uuid.UUID `gorm:"type:uuid"`
	TransactionCode   string    `json:"transaction_code"`                     // Optional field
	TransactionStatus int       `json:"transaction_status"`
}

type Database interface {
	// Connection management
	Connect() error  // Connects to the database
	//users
	CreateUser(userData map[string]interface{}) error  // Creates a new user
  	GetUserByName(username string) (map[string]interface{}, error)    // Retrieves user details by name (returns a map)
  	UserExists(username string) bool //method to check user existence
  	UpdateUserByName(username string, fieldsToUpdate map[string]interface{}) error //method to update by name
  	DeleteUserByName(username string) error                        // Deletes a user by name
	GetAllUsers() ([]map[string]interface{}, error)

	//organization
	CreateOrganization(orgData map[string]interface{}) error  // Creates a new organization
	GetOrganizationByName(orgname string) (map[string]interface{}, error)    // Retrieves org details by name (returns a map)
	OrganizationExists(orgname string) bool //method to check organization's existence
	GetAllOrganizations()([]map[string]interface{}, error)
  
  }

  func (db *GormDB) Connect() error {
	err := godotenv.Load() // Load environment variables from .env file (optional)
	if err != nil {
	  return fmt.Errorf("error loading environment variables: %w", err)
	}
  
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
  
	// Build the connection string using environment variables
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUsername, dbPassword, dbName, dbPort)
  
	// Configure connection options (optional)
	config := &gorm.Config{}
	// Add other configuration options that can be made: 
	// config.Logger = logger.Default // Set a custom logger
  
	dbClient, err := gorm.Open(postgres.Open(connectionString), config)
	if err != nil {
	  return fmt.Errorf("failed to connect to database: %w", err)
	}
  
	// Set automigrations mode (optional)
	err = dbClient.AutoMigrate(&User{}, &Organization{},&Machine{},&Transactions{})
	if err != nil{
		fmt.Println("Error occured during migration: ",err)
	}

	db.DB = dbClient
	return nil
}

// check if a user exists: 
func (db *GormDB) UserExists(username string) bool {

	var user User
	result := db.DB.Where("username = ?", username).First(&user)
	return result.Error == nil && result.RowsAffected > 0
}

func (db *GormDB)GetUserByName(username string) (map[string]interface{}, error){

	if !db.UserExists(username) {
		return nil, fmt.Errorf("user with username '%s' does not exist", username)
	  }
	var user User
	result := db.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
	  return nil, result.Error
	}

	userData := map[string]interface{}{
	"username": user.Username,
	"phone_number": user.UserPhone,
	"role":     user.Role,
	}

	return userData, nil

}

// update user details: 
func (db *GormDB) UpdateUserByName(username string, fieldsToUpdate map[string]interface{}) error {

	if !db.UserExists(username) {
	  return fmt.Errorf("user with username '%s' does not exist", username)
	}
  
	// Get existing user data
	var existingUser User 
	result := db.DB.Where("username = ?", username).First(&existingUser)
	if result.Error != nil {
	  return result.Error
	}
  
	// Update user fields based on the provided map
	for field, newValue := range fieldsToUpdate {
	  switch field {
	  case "username": 
		existingUser.Username = newValue.(string)
	  case "password":
		hashedPassword, err := HashPassword(newValue.(string))
		if err != nil {
		  return err
		}
		existingUser.Password = hashedPassword
	  case "role":
		existingUser.Role = newValue.(int) // Ensure newValue is cast to int
	  case "phone_number": // Assuming a single phone number field in User struct
		existingUser.UserPhone = newValue.(string)
	  default:
		// Handle updates for other supported fields (if any)
	  }
	}
  
	result = db.DB.Save(&existingUser)
	if result.Error != nil {
	  return result.Error
	}
	return nil
  }

//create new user	
func (db *GormDB) CreateUser(userData map[string]interface{}) error {

	// Validate and extract data from the map
	username, ok := userData["username"].(string)
	if !ok || username == "" {
	  return fmt.Errorf("missing or invalid username in user data")
	}
  
	password, ok := userData["password"].(string)
	if !ok || password == "" {
	  return fmt.Errorf("missing or invalid password in user data")
	}
  
	role, ok := userData["role"].(int)
	if !ok {
	  return fmt.Errorf("missing or invalid role in user data")
	}

	orgname, ok := userData["org"].(string)
	if !ok{
		return fmt.Errorf("missing or invalid organization name in user data")
	}

	phonenumber, ok := userData["phonenumber"].(string)
	// validate phone number maybe
	if !ok{
		return fmt.Errorf("missing or invalid phone number")
	}

	if !db.OrganizationExists(orgname){
		return fmt.Errorf("organization does not exist")
	}
  
	hashedPassword, err := HashPassword(password)
	if err != nil {
	  return err
	}
	// Check if user with the same username already exists
	if db.UserExists(username) {
		return fmt.Errorf("User with username '%s' already exists", username)
	}
	// Get organization to attach user to:
	var org Organization
	result := db.Where("organization_name = ?", orgname).Preload("Users").First(&org)
	if result.Error != nil {
		return result.Error
	}
		// User doesn't exist, create directly using model
	newUser := User{Username: username, Password: hashedPassword, Role: role, Organization: org, UserPhone: phonenumber}
	result = db.DB.Create(&newUser)
	if result.Error != nil {
	  return result.Error
	}
	return nil
}

func (db *GormDB) DeleteUserByName(username string) error {
	if !db.UserExists(username) {
	  return fmt.Errorf("user with username '%s' does not exist", username)
	}
  
	result := db.DB.Where("username = ?", username).Delete(&User{}) 
	if result.Error != nil {
	  return result.Error
	}
  
	return nil
}

func (db *GormDB) GetAllUsers() ([]map[string]interface{}, error) {
	var users []User 
	result := db.DB.Preload("Organization").Find(&users)
	if result.Error != nil {
	  return nil, result.Error
	}
  
	// Convert user structs to maps
	userData := make([]map[string]interface{}, len(users))
	for i, user := range users {
	  userData[i] = map[string]interface{}{
		"username": user.Username,
		"phone_number": user.UserPhone,
		"role":     user.Role,
		"machines": user.MachinesOwned,
		"org": user.Organization,
		"org_id": user.OrganizationID,
	  }
	}
  
	return userData, nil
}

  // check if an organization exists: 
func (db *GormDB) OrganizationExists(orgname string) bool {

	var org Organization
	result := db.DB.Where("organization_name = ?", orgname).First(&org)
	return result.Error == nil && result.RowsAffected > 0
}

func (db *GormDB) CreateOrganization(orgData map[string]interface{}) error {

	// Validate and extract data from the map
	orgname, ok := orgData["orgname"].(string)
	if !ok || orgname == "" {
	  return fmt.Errorf("missing or invalid orgname in org data")
	}
  
	// Check if organization with the same name already exists
	if db.OrganizationExists(orgname) {
		return fmt.Errorf("organization with name '%s' already exists", orgname)
	}
	
	// User doesn't exist, create directly using model
	newOrg := Organization{OrganizationName: orgname}
	result := db.DB.Create(&newOrg)
	if result.Error != nil {
	  return result.Error
	}
	return nil
}

func (db *GormDB)GetOrganizationByName(orgname string) (map[string]interface{}, error){

	if !db.OrganizationExists(orgname) {
		return nil, fmt.Errorf("Organization with name '%s' does not exist", orgname)
	  }
	var org Organization
	result := db.DB.Where("organization_name = ?", orgname).First(&org)
	if result.Error != nil {
	  return nil, result.Error
	}

	orgData := map[string]interface{}{
	"orgname": org.OrganizationName,
	"role":   org.Users,
	}

	return orgData, nil
}

func (db *GormDB) GetAllOrganizations() ([]map[string]interface{}, error) {
	var orgs []Organization
	result := db.DB.Preload("Users").Find(&orgs)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert organization structs to maps, extracting only usernames
	orgData := make([]map[string]interface{}, len(orgs))
	for i, org := range orgs {
		usernames := make([]string, len(org.Users))
		for j, user := range org.Users {
		usernames[j] = user.Username
		}

		orgData[i] = map[string]interface{}{
		"orgname": org.OrganizationName,
		"users":   usernames,  // Include the filtered usernames
		"id":      org.ID,
		}
	}
  
	return orgData, nil
}

// HashPassword function to hash passwords using bcrypt
func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
