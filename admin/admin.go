package main

import (
	"fmt"

	"github.com/osoderholm/eletab-lite/app/bundles/accountsbundle"
)

/*
	This part of the eletab-lite system is used for creating a SuperAdmin-account.

	SuperAdmins are the highest admins and have the ability to create new "normal" admins.
	Super admins also have an account with a balance, even though it is not recommended to
	use them for purchases. Or actually, who cares, do what ever you want... :D
*/

func main() {
	fmt.Println("Create a SuperAdmin-account")
	fmt.Println("To cancel, use the normal escape char ^C")
	var name string
	var username string
	var pass string
	var balance int64

	readString("Name: ", &name)
	readString("Username: ", &username)
	readString("Password: ", &pass)
	fmt.Print("Balance: ")
	fmt.Scan(&balance)
	if balance < 0 {
		balance = 0
	}

	account := accountsbundle.AddAccount(name, username, pass, balance, accountsbundle.LevelSuperAdmin)
	if account == nil {
		fmt.Println("SuperAdmin account was not created!")
	}
	fmt.Println("SuperAdmin account created with following information:")
	fmt.Printf("%v", *account)

}

// Read string input into variable until string is not empty
func readString(prompt string, v *string) {
	for len(*v) == 0 {
		fmt.Print(prompt)
		fmt.Scan(v)
	}
}
