package helper

import (
	"crud-app/middleware"
	"fmt"
)

func main() {
	hash, _ := middleware.HashPassword("123456")
	fmt.Println("Hash baru untuk 123456:", hash)
}
