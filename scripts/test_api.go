package main

import (
	"fmt"
	"backend_go/pkg/utils"
	"github.com/joho/godotenv"
	"net/http"
	"io/ioutil"
)

func main() {
	godotenv.Load(".env")
	
	// Generate token for admin
	token, _ := utils.GenerateToken("admin-id-here", "Admin")
	
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/patients", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))
}
