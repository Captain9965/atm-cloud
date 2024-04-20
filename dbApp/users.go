package dbApp

import (
	"fmt"
	"golang.org/x/crypto/bcrypt" // Import for password hashing
)

// check if a user exists: 
func (db *GormDB) UserExists(username string) bool {

	var user User
	result := db.DB.Where("username = ?", username).First(&user)
	return result.Error == nil && result.RowsAffected > 0
}

func (db *GormDB)AuthenticateUserByName(username string, password string) (bool,error){
	User, err := db.GetUserByName(username)

	if err != nil{
		return false, fmt.Errorf("User with the username %s does not exist", username)
	}

	currentPassword, ok := User["password"].(string)
	if !ok || currentPassword == "" {
	  return false, fmt.Errorf("missing or invalid password in db user object")
	}

	err = bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(password))
	if err!= nil {
		return false, fmt.Errorf("wrong password provided")
	}
	return true, nil
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
	"password": user.Password,
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
	  case "new_username": 
		existingUser.Username = newValue.(string)
	  case "new_password":
		hashedPassword, err := HashPassword(newValue.(string))
		if err != nil {
		  return err
		}
		existingUser.Password = hashedPassword
	  case "new_role":
		existingUser.Role = newValue.(string) // Ensure newValue is cast to string
	  case "new_phonenumber": // Assuming a single phone number field in User struct
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
  
	role, ok := userData["role"].(string)
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
	result := db.DB.Preload("Organization").Preload("Machines").Find(&users)
	if result.Error != nil {
	  return nil, result.Error
	}
  
	// Convert user structs to maps
	userData := make([]map[string]interface{}, len(users))
	for i, user := range users {
		machineUUIDs := make([]string, len(user.Machines))
		for j, machine := range user.Machines {
		machineUUIDs[j] = machine.MachineUUID
		}
		userData[i] = map[string]interface{}{
			"username": user.Username,
			"phone_number": user.UserPhone,
			"role":     user.Role,
			"machines": machineUUIDs,
			"org": user.Organization,
			"org_id": user.OrganizationID,
		}
	}
  
	return userData, nil
}