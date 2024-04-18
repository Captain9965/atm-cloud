package restapi

import(
	"github.com/gin-gonic/gin"
	dbapp "cloud/dbApp"
	"net/http"
	"fmt"
)

func RegisterMachinesApi(r *gin.Engine){
	// create a machine
	r.POST("/machines/create", func(c *gin.Context) {
		createMachine(c)
	})
	// udpate machine
	r.POST("/machines/update", func(c *gin.Context) {
		updateMachine(c)
	})
	  // get all machines
	r.GET("/machines/all", func(c *gin.Context) {
		getAllMachines(c)
	})
	//delete machine
	r.DELETE("/machines/delete", func(c *gin.Context) {
		deleteMachine(c)
	})
}

func createMachine(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	// Bind JSON from request body
	var newMachine struct {
		MachineType string  `json:"machine_type" binding:"required"`
		MachineUUID string 	`json:"machine_uuid" binding:"required"`
		OwnerName string	`json:"owner_name" binding:"required"`
	}

	if err := c.Bind(&newMachine); err != nil {
		fmt.Println(newMachine)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := Db.CreateMachine(map[string]interface{}{
		"machine_type": newMachine.MachineType,
		"owner_name": newMachine.OwnerName,
		"machine_uuid": newMachine.MachineUUID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create machine", "Description":err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Machine created successfully"})

}

func updateMachine(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)

	// Bind JSON from request body
	var updatedMachine struct {
		MachineUUID string  		`json:"machine_uuid" binding:"required"`
		NewMachineUUID string  		`json:"new_machine_uuid"`

	}

	if err := c.Bind(&updatedMachine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	updatedData := map[string]interface{}{"machine_uuid": updatedMachine.MachineUUID}
	if updatedMachine.NewMachineUUID != ""{
		updatedData["new_machine_uuid"] = updatedMachine.NewMachineUUID
	}

	err := Db.UpdateMachinebyID(updatedMachine.MachineUUID, updatedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update machine", "Description":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Machine updated successfully"})
}

func getAllMachines(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	
	// Get all machines
	machines, err := Db.GetAllMachines()

	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve machines", "Description":err.Error()})
	  return
	
	}
	c.JSON(http.StatusOK, machines)
}

func deleteMachine(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)

	var machine struct{
		MachineUUID string `json:"machine_uuid" binding:"required"`
	}
	if err := c.Bind(&machine); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
	}
	err := Db.DeleteMachineByID(machine.MachineUUID)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete machine", "Description":err.Error()})
	  return
	}
	c.JSON(http.StatusOK, gin.H{"Machine deleted": machine.MachineUUID})
}