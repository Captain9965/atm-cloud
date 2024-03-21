package ussdApp

import(
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"cloud/stkApp"

)

func UssdCallback(c *gin.Context)  {
	
	err := c.Request.ParseForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
			return
		}

		// Extract values
		sessionId := c.PostForm("sessionId")
		phoneNumber := c.PostForm("phoneNumber")
		networkCode := c.PostForm("networkCode")
		serviceCode := c.PostForm("serviceCode")
		text := c.PostForm("text")
		fmt.Printf("%s, %s, %s, %s, %s \n", sessionId, phoneNumber, networkCode, serviceCode, text)

		if len(text) == 0 {
			// we are at the beginning of a session:
			resp := "CON Welcome to Yu n Mi beverage machines"
			resp += "\n Type in the amount you wish to pay"
			c.Data(200, "text/plain; charset=utf-8", []byte(resp))
		} else {

			num, err := strconv.Atoi(text)
			if err != nil {
				fmt.Printf("Error occured %s, %s\n", err, text)
				c.Data(200, "text/plain; charset=utf-8", []byte("END Invalid input \n Try again"))
			}else {
				fmt.Printf("Sending stk push to %s, for amount %d\n", phoneNumber, num)
				stkapp.SendStkPush(phoneNumber, num)
				c.Data(200, "text/plain; charset=utf-8", []byte("END Input your pin \n Thank you!"))
			}
		}
}