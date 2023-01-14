package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"strings"
	errorHandler "github.com/AdityaSubrahmanyaBhat/golang/dashDB/Error"
	// dbFunctions "github.com/AdityaSubrahmanyaBhat/golang/dashDB/Functions"
	address "github.com/AdityaSubrahmanyaBhat/golang/dashDB/models/Address"
	u "github.com/AdityaSubrahmanyaBhat/golang/dashDB/models/User"
	"github.com/jcelliott/lumber"
)

type (
	Logger interface {
		Fatal(string, ...interface{})
		Warn(string, ...interface{})
		Error(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}

	Options struct {
		Logger
	}
)

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return 
}

func (d *Driver) CreateNewDB(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection string, recordItem string, v interface{}) error {
	if collection == ""{
		return fmt.Errorf("missing collection name")
	}
	if recordItem == ""{
		return fmt.Errorf("missing record item")
	}

	mutex := d.GetOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	finalPath := filepath.Join(dir, recordItem + ".json")
	tmpPath := finalPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil{
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t"); 
	if err != nil{
		return err
	}
	
	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, finalPath)
}

func (d *Driver) Read(collection string, recordItem string, v interface{}) error{
	if collection == ""{
		return fmt.Errorf("missing collection")
	}
	if recordItem == ""{
		return fmt.Errorf("missing record item")
	}
	record := filepath.Join(d.dir, collection, recordItem)

	if _, err := stat(record); err!= nil {
		return err
	}

	b, err := ioutil.ReadFile(record + ".json")
	if err != nil{
		return err
	}
	return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string) ([]string, error){
	if collection == ""{
		return nil, fmt.Errorf("missing collection")
	}

	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(dir)
	errorHandler.HandleError(err)

	var records []string

	for _, file := range files{
		b, err := ioutil.ReadFile(filepath.Join(dir, strings.ToLower(file.Name())))
		if err != nil{
			return nil, err
		}
		records = append(records, string(b))
	}
	return records, nil
}

func (d *Driver) Delete(collection string, recordName string) error {
	path := filepath.Join(collection, strings.ToLower(recordName))
	mutex :=d.GetOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err :=stat(dir); {
	case fi == nil, err != nil : return fmt.Errorf("unable to find file")
	case fi.Mode().IsDir() : return os.RemoveAll(dir)
	case fi.Mode().IsRegular() : return os.RemoveAll(dir + ".json")
	}
	return nil
}

func (d *Driver) GetOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok{
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}


func main() {
	dir := "./"
	driver := Driver{}
	dashDB, err := driver.CreateNewDB(dir,nil)
	errorHandler.HandleError(err)

	employees := []u.User{
		{"Aditya", "21", "jpmc", address.Address{"Mysuru", "Karnataka", "India", "570002"}},
		{"pranav", "21", "siemens", address.Address{"Mysuru", "Karnataka", "India", "570002"}},
		{"prashanth", "21", "halodoc", address.Address{"Bangalore", "Karnataka", "India", "560002"}},
		{"Gajanana", "21", "halodoc", address.Address{"Mysuru", "Karnataka", "India", "570002"}},
	}

	for _, record := range employees {
		dashDB.Write("users", strings.ToLower(record.Name), u.User{
			Name:    record.Name,
			Age:     record.Age,
			Company: record.Company,
			Address: record.Address,
		})
	}

	records, err := dashDB.ReadAll("users")
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

	// if err := dashDB.Delete("users", "Prashanth"); err != nil{
	// 	errorHandler.HandleError(err)
	// }
	
	// if err := dashDB.Delete("users", ""); err != nil{
	// 	errorHandler.HandleError(err)
	// }
}
