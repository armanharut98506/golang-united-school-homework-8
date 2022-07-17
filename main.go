package main

import (
	"io"
	"os"
	"flag"
	"fmt"
	"errors"
	"strconv"
	"io/ioutil"
 "encoding/json"
)

type Arguments map[string]string

type User struct {
	Id string
	Email string 
	Age int 
}

func checkFile(filename string) error {
	_, err := os.Stat(filename)
			if os.IsNotExist(err) {
					_, err := os.Create(filename)
							if err != nil {
									return err
							}
			}
			return nil
}

func InitialCheck(operation, fileName string) error {
	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}
	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}

	operations := []string{"list", "add", "remove", "findById"}
	for _, op := range operations {
		if op == operation {
			return nil
		}
	}
	return fmt.Errorf("Operation %s not allowed!", operation)
}

func IdCheck(id string) error {
	if id == "" {
		return errors.New("-id flag has to be specified")
	}
	return nil
}

func ItemCheck(item string) error {
	if item == "" {
		return errors.New("-item flag is missing")
	}
	return nil
}

func List(fileName string, writer io.Writer) {
	err := checkFile(fileName)
	if err != nil {
		panic(err)
	}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(writer, string(data))
}

func Add(item string, fileName string) {
	var user User
	json.Unmarshal([]byte(item), &user)

	err = checkFile(fileName)
	if err != nil {
		panic(err)
	}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var users []User
	json.Unmarshal(data, &users)
	users = append(users, user)

	data, err = json.Marshal(users)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		panic(err)
	}
}

func RemoveById(id string, fileName string, writer io.Writer) {
	err := checkFile(fileName)
	if err != nil {
		panic(err)
	}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var users []User
	json.Unmarshal(data, &users)

	for index, user := range users {
		if user.Id == id {
			users = append(users[:index], users[index+1:]...)
			data, err := json.Marshal(users)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(fileName, data, 0644)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(writer, "Item with id %d not found", i)
}

func FindById(id string, fileName string, writer io.Writer) {
	err := checkFile(fileName)
	if err != nil {
		panic(err)
	}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var users []User
	json.Unmarshal(data, &users)

	for _, user := range users {
		if user.Id == id {
			data, err := json.Marshal(user)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(writer, string(data))
			return 
		}
	}
	fmt.Fprintln(writer, "")
}

func Perform(args Arguments, writer io.Writer) error {
	err := InitialCheck(args["operation"], args["fileName"])
	if err != nil {
		return err
	}
	switch args["operation"] {
	case "list":
		List(args["fileName"], writer)
	case "add":
		err := ItemCheck(args["item"])
		if err != nil {
			return err
		}
		Add(args["item"], args["fileName"])
	case "remove":
		err := IdCheck(args["id"])
		if err != nil {
			return err
		}
		RemoveById(args["id"], args["fileName"], writer)
	case "findById":
		err := IdCheck(args["id"])
		if err != nil {
			return err
		}
		FindById(args["id"], args["fileName"], writer)
	}
	return nil
}

func main() {
	args := Arguments{}
	operation := flag.String("operation", "", "")
	fileName := flag.String("fileName", "", "")
	item := flag.String("item", "", "")
	id := flag.String("id", "", "")

	flag.Parse()
	
	args["operation"] = *operation
	args["fileName"] = *fileName
	args["item"] = *item
	args["id"] = *id

	err := Perform(args, os.Stdout)
	if err != nil {
		panic(err)
	}
}
