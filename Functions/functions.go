package functions

import (
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"sync"
	"strings"
	"path/filepath"
	"github.com/jcelliott/lumber"
	errorHandler "github.com/AdityaSubrahmanyaBhat/golang/dashDB/Error"
	o "github.com/AdityaSubrahmanyaBhat/golang/dashDB/models/Options"
	d "github.com/AdityaSubrahmanyaBhat/golang/dashDB/models/Driver"
)

type Functions interface{
	CreateNewDB(dir string)
	Write(collection string, recordItem string, v interface{}, driver *d.Driver)
	Read(collection string, recordItem string, v interface{}, driver *d.Driver)
	ReadAll(collection string, driver *d.Driver)
	GetOrCreateMutex(collection string, driver *d.Driver)
	Delete(collection string, recordItem string, driver *d.Driver)
	stat(path string)
}

func CreateNewDB(dir string, options *o.Options) (dr *d.Driver, err error){
	dir = filepath.Clean(dir)
	opts := o.Options{}
	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := d.Driver{
		Dir:     dir,
		Mutexes: make(map[string]*sync.Mutex),
		Log:     opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func Write(collection string, recordItem string, v interface{}, driver *d.Driver) error{
	if collection == ""{
		return fmt.Errorf("missing collection name")
	}
	if recordItem == ""{
		return fmt.Errorf("missing record item")
	}

	mutex := GetOrCreateMutex(collection, driver)
	mutex.Lock()
	defer mutex.Unlock()
	dir := filepath.Join(driver.Dir, collection)
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

func Read(collection string, recordItem string, v interface{}, driver *d.Driver) error{
	if collection == ""{
		return fmt.Errorf("missing collection")
	}
	if recordItem == ""{
		return fmt.Errorf("missing record item")
	}
	record := filepath.Join(driver.Dir, collection, recordItem)

	if _, err := stat(record); err!= nil {
		return err
	}

	b, err := ioutil.ReadFile(record + ".json")
	if err != nil{
		return err
	}
	return json.Unmarshal(b, &v)
}

func ReadAll(collection string, driver *d.Driver) ([]string, error){
	 	if collection == ""{
				return nil, fmt.Errorf("missing collection")
			}
		
			dir := filepath.Join(driver.Dir, collection)
		
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

func GetOrCreateMutex(collection string, driver *d.Driver) *sync.Mutex{
	driver.Mutex.Lock()
	defer driver.Mutex.Unlock()
	m, ok := driver.Mutexes[collection]

	if !ok{
		m = &sync.Mutex{}
		driver.Mutexes[collection] = m
	}
	return m
}

func Delete(collection string, recordItem string, driver *d.Driver) error{
	path := filepath.Join(collection, strings.ToLower(recordItem))
	mutex :=GetOrCreateMutex(collection, driver)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(driver.Dir, path)

	switch fi, err :=stat(dir); {
	case fi == nil, err != nil : return fmt.Errorf("unable to find file")
	case fi.Mode().IsDir() : return os.RemoveAll(dir)
	case fi.Mode().IsRegular() : return os.RemoveAll(dir + ".json")
	}
	return nil
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return 
}