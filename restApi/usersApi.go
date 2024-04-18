package restapi

import(
	"github.com/gin-gonic/gin"
	dbapp "cloud/dbApp"
	"net/http"
	"fmt"
)

func RegisterUsersApi(r *gin.Engine){
	// create a user
	r.POST("/users/create", func(c *gin.Context) {
		createUser(c)
	})
	// udpate user
	r.POST("/users/update", func(c *gin.Context) {
		updateUser(c)
	})
	  // get all users
	r.GET("/users/all", func(c *gin.Context) {
		getAllUsers(c)
	})
	// login user
	r.POST("/users/login", func(c *gin.Context) {
		loginUser(c)	
	})
	//delete user
	r.DELETE("/users/delete", func(c *gin.Context) {
		deleteUser(c)
	})
}


func createUser(c *gin.Context){

	Db := c.MustGet("db").(dbapp.Database)
	// Bind JSON from request body
	var newUser struct {
		Username string  `json:"username" binding:"required"`
		Password string  `json:"password" binding:"required"`
		Role string		   `json:"role" binding:"required"`
		Org string	   `json:"org" binding:"required"`
		PhoneNumber string `json:"phonenumber" binding:"required"`
	}

	if err := c.Bind(&newUser); err != nil {
		fmt.Println(newUser)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	err := Db.CreateUser(map[string]interface{}{
		"username": newUser.Username,
		"password": newUser.Password,
		"role"	: newUser.Role,
		"org"		: newUser.Org,
		"phonenumber" : newUser.PhoneNumber,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "Description":err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

}

func updateUser(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)

	// Bind JSON from request body
	var updatedUser struct {
		Username string  			`json:"username" binding:"required"`
		NewUsername string  		`json:"new_username"`
		NewPassword string  		`json:"new_password"`
		NewRole string			`json:"new_role"`
		NewOrg string	   			`json:"new_org"`
		NewPhoneNumber string 	`json:"new_phonenumber"`
	}

	if err := c.Bind(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	updatedData := map[string]interface{}{"username": updatedUser.Username}
	if updatedUser.NewUsername != ""{
		updatedData["new_username"] = updatedUser.NewUsername
	}

	if updatedUser.NewPassword != ""{
		updatedData["new_password"] = updatedUser.NewPassword
	}

	if updatedUser.NewRole != ""{
		updatedData["new_role"] = updatedUser.NewRole
	}

	if updatedUser.NewPhoneNumber != ""{
		updatedData["new_phonenumber"] = updatedUser.NewPhoneNumber
	}

	err := Db.UpdateUserByName(updatedUser.Username, updatedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "Description":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func getAllUsers(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
	
	// Get all users
	users, err := Db.GetAllUsers()

	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "Description":err.Error()})
	  return
	
	}
	c.JSON(http.StatusOK, users)
}

func loginUser(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)
		
	var user struct{
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.Bind(&user); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
	}

	success, err := Db.AuthenticateUserByName(user.Username, user.Password)

	if !success{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed", "Description":err.Error()})
	} else{
		c.JSON(http.StatusOK, gin.H{"User Login success": user.Username})
	}
}

func deleteUser(c *gin.Context){
	Db := c.MustGet("db").(dbapp.Database)

	var user struct{
		Username string `json:"username" binding:"required"`
	}
	if err := c.Bind(&user); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
	}
	err := Db.DeleteUserByName(user.Username)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete user", "Description":err.Error()})
	  return
	}
	c.JSON(http.StatusOK, gin.H{"User deleted": user.Username})
}