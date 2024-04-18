package dbApp

import(
	"fmt"
)

func (db *GormDB) CreateMachine(machineData map[string]interface{}) error {
	

	machineType, ok := machineData["machine_type"].(string)
	if !ok || machineType == "" {
	  return fmt.Errorf("missing or invalid machine type in machine data")
	}

	machineUUID, ok := machineData["machine_uuid"].(string)
	if !ok || machineUUID == "" {
		return fmt.Errorf("missing or invalid machine uuid in machine data")
	}

	userName, ok := machineData["owner_name"].(string)
	if !ok || userName == "" {
		return fmt.Errorf("missing or invalid machine user name in machine data")
	}

	// Get user
	var user User
	result := db.Where("username = ?", userName).Preload("Machines").First(&user)
	if result.Error != nil {
		return result.Error
	}

	//check if machine already exists:
	if db.MachineExists(machineUUID){
		return fmt.Errorf("a machine with a similar UUID already exists")
	}
    fmt.Print("The user found is:"); fmt.Println(user.Username)
	machine := Machine{
		MachineType: machineType,
		MachineUUID: machineUUID,
		User: user,
	 }
	 err := db.Create(&machine).Error
	 if err != nil {
	   return err
	 }
  
	return nil
}

func (db *GormDB) MachineExists(machineUUID string) bool {
	var count int64
	db.Model(&Machine{}).Where("machine_uuid = ?", machineUUID).Count(&count)
	return count > 0
}

func (db *GormDB) GetMachineByID(machineUUID string)(map[string]interface{}, error) {
	var machine Machine
	if !db.MachineExists(machineUUID){
		return nil, fmt.Errorf("machine with ID %d does not exist", machine.ID)
	}
	result := db.Preload("User").Where("machine_uuid = ?", machineUUID).First(&machine)
	if result.Error != nil {
	  return nil, result.Error
	}

	machineMap := map[string]interface{}{
	  "machine_type": machine.MachineType,
	  "machine_uuid": machine.MachineUUID,
	  "user":machine.User.Username,
	  "admin_cash":machine.AdminCash,
	  "vending_card_id":machine.VendingCardID,
	  "admin_card_id":machine.AdminCardID,
	  "service_card_id":machine.ServiceCardID,
	}
  
	return machineMap, nil
}

func (db *GormDB) UpdateMachinebyID(machineUUID string, fieldsToUpdate map[string]interface{}) error {
	if !db.MachineExists(machineUUID){
		return fmt.Errorf("Machine with given uuid does not exist")
	  }
	
	  // Get existing machine Data
	  var existingMachine Machine
	  result := db.DB.Where("machine_uuid = ?", machineUUID).First(&existingMachine)
	  if result.Error != nil {
		return result.Error
	  }
	
	  // Update user fields based on the provided map
	  for field, newValue := range fieldsToUpdate {
		switch field {
		case "new_machine_uuid": 
		  existingMachine.MachineUUID = newValue.(string)
		case "new_machine_type":
		  existingMachine.MachineType = newValue.(string)
		case "new_admin_cash":
		  existingMachine.AdminCash = newValue.(int)
		default:
		  // Handle updates for other supported fields (if any)
		}
	  }
	
	  result = db.DB.Save(&existingMachine)
	  if result.Error != nil {
		return result.Error
	  }
	  return nil
}

func (db *GormDB) DeleteMachineByID(machineUUID string) error {
	if !db.MachineExists(machineUUID) {
		return fmt.Errorf("machine with uuid given does not exist")
	  }
	
	  result := db.DB.Where("machine_uuid = ?", machineUUID).Delete(&Machine{}) 
	  if result.Error != nil {
		return result.Error
	  }
	
	  return nil
}

func (db *GormDB) GetAllMachines() ([]map[string]interface{}, error) {
	var machines []Machine
	result := db.Preload("User").Find(&machines)
	if result.Error != nil {
	  return nil, result.Error
	}
	// Convert user structs to maps
	machineData := make([]map[string]interface{}, len(machines))
	for i, machine := range machines {
	  machineData[i] = map[string]interface{}{
		"machine_type": machine.MachineType,
		"machineOwner": machine.User.Username,
		"machine_uid": machine.MachineUUID,
		"vending_card_id": machine.VendingCardID,
		"admin_card_id": machine.VendingCardID,
		"service_card_id": machine.ServiceCardID,
	  }
	}

	return machineData, nil
}