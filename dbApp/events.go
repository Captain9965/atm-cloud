package dbApp

import (
	"encoding/json"
	"fmt"
	"time"
)

var Allowed_events = map[string]struct{}{
	"pay":   {},
	"check": {},
	"loc":   {},
	"hello": {},
}

func (db *GormDB) CreateEvent(eventData map[string]interface{})error {
	

	eventType, ok := eventData["event_type"].(string)
	if !ok || eventType == "" {
	  return fmt.Errorf("missing or invalid event type in event data")
	}

	machineUUID, ok := eventData["machine_uuid"].(string)
	if !ok || machineUUID == "" {
		return fmt.Errorf("missing or invalid machine uuid in event data")
	}

	jsonData, ok := eventData["data"].(json.RawMessage)
	if !ok || string(jsonData) == "" {
		return fmt.Errorf("missing or invalid data in event Data")
	}


	// Preload machine that event belongs to:
	var machine Machine
	result := db.Where("machine_uuid = ?", machineUUID).Preload("Events").First(&machine)
	if result.Error != nil {
		return result.Error
	}

	//validate event type:
	
	if _, ok := Allowed_events[eventType]; !ok {
		return fmt.Errorf("unsupported event")
	} 
	event := Events{
		EventType: eventType,
		Data: jsonData,
		Machine: machine,
	 }
	 err := db.Create(&event).Error
	 if err != nil {
	   return err
	 }
  
	return nil
}


func (db *GormDB) GetEventsByType(eventType string, from time.Time, to time.Time)([]map[string]interface{}, error) {
	//check if event type is valid:
	if _, ok := Allowed_events[eventType]; !ok {
		return nil, fmt.Errorf("unsupported event")
	}

	var Events []Events
	result := db.Preload("Machine").Where("created_at BETWEEN ? AND ?", from, to).Where("event_type = ?", eventType).Find(&Events)
	if result.Error != nil {
	  return nil, result.Error
	}
	//Convert event structs to maps
	eventData := make([]map[string]interface{}, len(Events))
	for i, event := range Events {
	   eventData[i] = map[string]interface{}{
		 "event_type": event.EventType,
		 "data" : event.Data,
		 "machine_uuid": event.Machine.MachineUUID,
		 "id":event.ID,
	  }
    }

	return eventData, nil
}

func (db *GormDB) GetEventsByMachine(machineUUID string,from time.Time, to time.Time)([]map[string]interface{}, error) {
	//check if event type is valid:
	if !db.MachineExists(machineUUID){
		return nil, fmt.Errorf("machine does not exist")
	}
	
	var Events []Events
	result := db.Preload("Machine").
    Joins("JOIN machines ON machines.id = events.machine_id").
    Where("machines.machine_uuid = ? AND events.created_at BETWEEN ? AND ?", machineUUID, from, to).
    Find(&Events)

	if result.Error != nil {
		return nil, result.Error
	}

	//Convert event structs to maps
	eventData := make([]map[string]interface{}, len(Events))
	for i, event := range Events {
	   eventData[i] = map[string]interface{}{
		 "event_type": event.EventType,
		 "data" : event.Data,
		 "machine_uuid": event.Machine.MachineUUID,
		 "id": event.ID,
	  }
    }

	return eventData, nil
}

func (db *GormDB) UpdateEventByID(eventID uint, fieldsToUpdate map[string]interface{})error {
		
	// Get existing event:
	var existingEvent Events
	result := db.Preload("Machine"). Where("id = ?", eventID).First(&existingEvent)
	if result.Error != nil {
	return result.Error
	}

	// Update user fields based on the provided map
	for field, newValue := range fieldsToUpdate {
		switch field {
		case "new_type":
			event := newValue.(string)
			if _, ok := Allowed_events[event]; !ok {
				return fmt.Errorf("unsupported event")
			}
			existingEvent.EventType = event
		case "new_data":
			data, ok := newValue.(json.RawMessage)
			if !ok || string(data) == "" {
				return fmt.Errorf("missing or invalid data")
			}
			existingEvent.Data = data
		default:
			// Handle updates for other supported fields (if any)
		}
	}

	result = db.DB.Save(&existingEvent)
	if result.Error != nil {
	return result.Error
	}
	return nil
}

func (db *GormDB) DeleteEventByID(eventID uint)(bool, error) {

	result := db.DB.Where("id = ?", eventID).Delete(&Events{}) 
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

