package main

import (
	"fmt"
	"backend_go/pkg/utils"
	"backend_go/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	db := utils.ConnectDB()
	
	repo := repository.NewUserRepository(db)
	user, err := repo.FindByEmail("test@rme.com")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if user == nil {
		fmt.Println("User not found!")
		return
	}
	fmt.Printf("User: %+v\n", user)
	fmt.Printf("User.Password: %q\n", user.Password)
}
