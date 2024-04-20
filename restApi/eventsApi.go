package restapi

import (
	dbapp "cloud/dbApp"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

func RegisterEventsApi(r *gin.Engine){
	// create an event
	r.POST("/events/create", func(c *gin.Context) {
		createEvent(c)
	})
	// udpate event
	r.POST("/events/update", func(c *gin.Context) {
		updateEvent(c)
	})
	  // get a slice of events by type
	r.GET("/events/get/type", func(c *gin.Context) {
		getEventsbyType(c)
	})
	  // get a slice of events by machine
	r.GET("/events/get/machine", func(c *gin.Context) {
		getEventsbyMachine(c)
	})
	//delete event
	r.DELETE("/events/delete", func(c *gin.Context) {
		deleteEvent(c)
	})
}

func createEvent(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	// Bind JSON from request body
	var newEvent struct {
		EventType string `json:"event_type" binding:"required"`
		MachineUUID string 	`json:"machine_uuid" binding:"required"`
		Data 		json.RawMessage	`json:"data" binding:"required"`
	}

	if err := c.Bind(&newEvent); err != nil {
		fmt.Println(newEvent)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := Db.CreateEvent(map[string]interface{}{
		"event_type": newEvent.EventType,
		"machine_uuid": newEvent.MachineUUID,
		"data": newEvent.Data,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event", "Description":err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created successfully"})

}

func updateEvent(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)

	// Bind JSON from request body
	var updatedEvent struct {
		ID uint  	`json:"id" binding:"required"`
		NewData json.RawMessage 		`json:"new_data"`
		NewType string `json:"new_type"`
	}

	if err := c.Bind(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedData := map[string]interface{}{
		"new_data": updatedEvent.NewData,
		"new_type": updatedEvent.NewType,

	}
	
	err := Db.UpdateEventByID(updatedEvent.ID, updatedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event", "Description":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event updated successfully"})
}

func getEventsbyType(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	
	// Get from and to time parameters from the request
	fromStr := c.Query("from")
	toStr := c.Query("to")
	eventType := c.Query("type")
  
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
	// Get events in the time range
	events, err := Db.GetEventsByType(eventType, from, to)

	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events", "Description":err.Error()})
	  return
	
	}
	c.JSON(http.StatusOK, events)
}

func getEventsbyMachine(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	
	// Get from and to time parameters from the request
	fromStr := c.Query("from")
	toStr := c.Query("to")
	machineUUID := c.Query("uuid")
  
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
	// Get events in the time range
	events, err := Db.GetEventsByMachine(machineUUID, from, to)

	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events", "Description":err.Error()})
	  return
	
	}
	c.JSON(http.StatusOK, events)
}
func deleteEvent(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)

	var event struct{
		ID uint `json:"id" binding:"required"`
	}
	if err := c.Bind(&event); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
	}
	
	deleted, err := Db.DeleteEventByID(event.ID)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete event", "Description":err.Error()})
	  return
	}
	if !deleted{
		c.JSON(http.StatusNotFound, gin.H{"error": "Event does not exist"})
	  return
	}
	c.JSON(http.StatusOK, gin.H{"event deleted": event.ID})
}