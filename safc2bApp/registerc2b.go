package safc2bapp

import(
	"io"
	"fmt"
	"net/http"
	stkapp "cloud/stkApp"
	"bytes"
)

func RegisterURL() {
	url := "https://sandbox.safaricom.co.ke/mpesa/c2b/v1/registerurl"
	method := "POST"
	client := &http.Client{}
	payload := []byte(`{   
		"ShortCode": "600999",
		"ResponseType":"Completed",
		"ConfirmationURL":"https://google.com/confirmation",
		"ValidationURL":"https://google.com/validation"
	 }`)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+stkapp.Token.AccessToken)

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Error ->", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error ->", err)
		return
	}
	fmt.Println(string(body))
}