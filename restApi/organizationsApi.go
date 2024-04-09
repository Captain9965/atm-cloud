package restapi

import(
	"github.com/gin-gonic/gin"
	dbapp "cloud/dbApp"
	"net/http"
)

func RegisterOrganizationApi(r *gin.Engine){
	// create organization
	r.POST("/createorg", func(c *gin.Context) {
		createOrg(c)
	})
	// update organization
	r.POST("/updateorg", func(c *gin.Context) {
		updateOrg(c)
	})
	// get all organizations
	r.GET("/orgs", func(c *gin.Context) {
		getAllOrganizations(c)
	})
	//delete organization
	r.DELETE("/deleteorg", func(c *gin.Context) {
		deleteOrganization(c)
	})
}

func createOrg(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	
	// Bind JSON from request body
	var newOrg struct {
	  Orgname string  `json:"orgname" binding:"required"`
	}

	if err := c.Bind(&newOrg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := Db.CreateOrganization(map[string]interface{}{
	  "orgname": newOrg.Orgname,
	})

	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization", "Description":err.Error()})
	  return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Organization created successfully"})
}

func updateOrg(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	
	// Bind JSON from request body
	var updatedOrg struct {
	  Orgname string  			`json:"orgname" binding:"required"`
	  NewOrgname string  		`json:"new_orgname"`
	}

	if err := c.Bind(&updatedOrg); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}
	
	updatedData := map[string]interface{}{"orgname": updatedOrg.Orgname}
	if updatedOrg.NewOrgname != ""{
		updatedData["new_orgname"] = updatedOrg.NewOrgname
	}
	err := Db.UpdateOrgByName(updatedOrg.Orgname, updatedData)
	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization", "Description":err.Error()})
	  return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization updated successfully"})
}

func getAllOrganizations(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	
	// Get all orgs
	orgs, err := Db.GetAllOrganizations()

	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organiizations", "Description":err.Error()})
	  return
	}

	c.JSON(http.StatusOK, orgs)
}

func deleteOrganization(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	var org struct{
		Orgname string `json:"orgname" binding:"required"`
	}
	if err := c.Bind(&org); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
	}
	err := Db.DeleteOrgByName(org.Orgname)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete organization", "Description":err.Error()})
	  return
	}
	c.JSON(http.StatusOK, gin.H{"Organization deleted": org.Orgname})
}