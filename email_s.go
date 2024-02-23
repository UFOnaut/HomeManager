package main

import (
	"home_manager/entities"
	"home_manager/utils"
)

func main() {
	utils.SendVerificationEmail("illia.mondok@gmail.com", entities.VerificationToken{
		UserId: 1,
		Token:  "test_token",
	})
}
