package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	operationFlag = "operation"
	itemFlag      = "item"
	fileNameFlag  = "fileName"
	idFlag        = "id"

	addCommand      = "add"
	listCommand     = "list"
	findByIdCommand = "findById"
	removeCommand   = "remove"
)

type Arguments map[string]string
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	operation, ok := args[operationFlag]
	if !ok || operation == "" {
		return buildMissingFlagError(operationFlag)
	}

	fileName, ok := args[fileNameFlag]
	if !ok || fileName == "" {
		return buildMissingFlagError(fileNameFlag)
	}

	switch operation {
	case listCommand:
		return listUsers(fileName, writer)
	case addCommand:
		item, ok := args[itemFlag]
		if !ok || item == "" {
			return buildMissingFlagError(itemFlag)
		}
		return addNewUser(item, fileName, writer)
	case removeCommand:
		id, ok := args[idFlag]
		if !ok || id == "" {
			return buildMissingFlagError(idFlag)
		}
		return removeUser(id, fileName, writer)
	case findByIdCommand:
		id, ok := args[idFlag]
		if !ok || id == "" {
			return buildMissingFlagError(idFlag)
		}
		return findUser(id, fileName, writer)
	}

	return fmt.Errorf("Operation %s not allowed!", operation)
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	result := make(map[string]string)
	operation := flag.String(operationFlag, "", "")
	item := flag.String(itemFlag, "", "")
	fileName := flag.String(fileNameFlag, "", "")
	id := flag.String(idFlag, "", "")
	flag.Parse()

	result[operationFlag] = *operation
	result[itemFlag] = *item
	result[fileNameFlag] = *fileName
	result[idFlag] = *id

	return result
}

func buildMissingFlagError(flagName string) error {
	return fmt.Errorf("-%s flag has to be specified", flagName)
}

func getUsersFromFile(fileName string) ([]User, error) {
	allUsers := []User{}
	fileData, err := ioutil.ReadFile(fileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &allUsers)
		if err != nil {
			return nil, err
		}
	}

	return allUsers, nil
}

func putUsersToFile(users []User, fileName string) error {
	usersJson, err := json.Marshal(&users)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, usersJson, os.ModePerm)
}

func addNewUser(item string, fileName string, writer io.Writer) error {
	user := User{}
	err := json.Unmarshal([]byte(item), &user)
	if err != nil {
		return err
	}

	allUsers, err := getUsersFromFile(fileName)
	if err != nil {
		return err
	}

	for _, v := range allUsers {
		if v.ID == user.ID {
			fmt.Fprint(writer, fmt.Sprintf("Item with id %s already exists", user.ID))
			return nil
		}
	}
	allUsers = append(allUsers, user)

	return putUsersToFile(allUsers, fileName)
}

func removeUser(id string, fileName string, writer io.Writer) error {
	allUsers, err := getUsersFromFile(fileName)
	if err != nil {
		return err
	}
	filteredUsers := []User{}
	for _, v := range allUsers {
		if v.ID != id {
			filteredUsers = append(filteredUsers, v)
		}
	}
	if len(allUsers) != len(filteredUsers) {
		return putUsersToFile(filteredUsers, fileName)
	} else {
		fmt.Fprint(writer, fmt.Sprintf("Item with id %s not found", id))
	}

	return nil
}

func findUser(id string, fileName string, writer io.Writer) error {
	allUsers, err := getUsersFromFile(fileName)
	if err != nil {
		return err
	}

	for _, v := range allUsers {
		if v.ID == id {
			stringUser, err := json.Marshal(&v)
			if err != nil {
				return err
			}
			fmt.Fprint(writer, string(stringUser))
			break
		}
	}
	return nil
}

func listUsers(fileName string, writer io.Writer) error {
	users, err := getUsersFromFile(fileName)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		jsonUsers, err := json.Marshal(&users)
		if err != nil {
			return err
		}
		fmt.Fprint(writer, string(jsonUsers))
	}
	return nil
}
