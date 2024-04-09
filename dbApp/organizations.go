package dbApp

import(
	"fmt"

)

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

// update user details: 
func (db *GormDB) UpdateOrgByName(orgname string, fieldsToUpdate map[string]interface{}) error {

	if !db.OrganizationExists(orgname) {
	  return fmt.Errorf("org with name '%s' does not exist", orgname)
	}
  
	// Get existing user data
	var existingOrg Organization 
	result := db.DB.Where("organization_name = ?", orgname).First(&existingOrg)
	if result.Error != nil {
	  return result.Error
	}
  
	// Update org fields based on the provided map
	for field, newValue := range fieldsToUpdate {
	  switch field {
	  case "new_orgname": 
		existingOrg.OrganizationName = newValue.(string)
	  default:
		// Handle updates for other supported fields (if any)
	  }
	}
  
	result = db.DB.Save(&existingOrg)
	if result.Error != nil {
	  return result.Error
	}
	return nil
  }

func (db *GormDB) DeleteOrgByName(orgname string) error {
	if !db.OrganizationExists(orgname) {
	  return fmt.Errorf("organization with name '%s' does not exist", orgname)
	}
  
	result := db.DB.Where("organization_name = ?", orgname).Delete(&Organization{}) 
	if result.Error != nil {
	  return result.Error
	}
  
	return nil
}