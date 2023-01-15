package main

import (
	"encoding/json"
	"fmt"
	"os/user"
	"strings"
	errorHandler "github.com/AdityaSubrahmanyaBhat/golang/dashDB/Error"
	dbFunctions "github.com/AdityaSubrahmanyaBhat/golang/dashDB/Functions"
	address "github.com/AdityaSubrahmanyaBhat/golang/dashDB/models/Address"
	u "github.com/AdityaSubrahmanyaBhat/golang/dashDB/models/User"
)

func main() {
	dir := "./"
	driver, err := dbFunctions.CreateNewDB(dir,nil)
	errorHandler.HandleError(err)
	employees := []u.User{
		{Name:"Aditya",Age:  "21",Company:  "jpmc",Address:  address.Address{City:"Mysuru",State:  "Karnataka",Country:  "India",Pincode:  "570002"}},
		{Name:"pranav", Age:"21", Company:"siemens",Address: address.Address{City:"Mysuru",State:  "Karnataka",Country:  "India",Pincode:  "570002"}},
		{Name:"prashanth",Age:  "21",Company:  "halodoc",Address:  address.Address{City:"Bangalore",State:  "Karnataka",Country:  "India",Pincode:  "560002"}},
		{Name:"Gajanana",Age:  "21",Company:  "halodoc",Address:  address.Address{City:"Mysuru",State:  "Karnataka",Country:  "India",Pincode:  "570002"}},
	}

	for _, record := range employees {
		dbFunctions.Write("users", strings.ToLower(record.Name), u.User{
			Name:    record.Name,
			Age:     record.Age,
			Company: record.Company,
			Address: record.Address,
		}, driver)
	}

	records, err := dbFunctions.ReadAll("users", driver)
	errorHandler.HandleError(err)
	fmt.Println(records)

	allUsers := []user.User{}

	for _, record := range records {
		employeeFound := user.User{}
		err := json.Unmarshal([]byte(record), &employeeFound)
		errorHandler.HandleError(err)
		allUsers = append(allUsers, employeeFound)
	}

	fmt.Println(allUsers)

	// if err := dbFunctions.Delete("users", "Prashanth", driver); err != nil{
	// 	errorHandler.HandleError(err)
	// }
	
	// if err := dbFunctions.Delete("users", "", driver); err != nil{
	// 	errorHandler.HandleError(err)
	// }
}
