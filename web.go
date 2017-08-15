package main

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"fmt"

	"github.com/crazyprograms/sud/client"
	"github.com/crazyprograms/sud/core"

	"log"

	"github.com/crazyprograms/sud/httpserver"
	"github.com/crazyprograms/sud/storage"
	_ "github.com/crazyprograms/sud/test"
)

type d interface{}
type Query struct {
	Action string               `json:",omitempty"`
	Token  string               `json:",omitempty"`
	Int    map[string]int64     `json:",omitempty"`
	Float  map[string]float64   `json:",omitempty"`
	String map[string]string    `json:",omitempty"`
	Bytes  map[string][]byte    `json:",omitempty"`
	Time   map[string]time.Time `json:",omitempty"`
	V      d                    `json:",omitempty"`
}
type A struct {
	Value1 string
}
type B struct {
	Value2 string
}

func checkTest(c *core.Core) {
	var err error
	var tid string
	if tid, err = c.BeginTransaction(); err != nil {
		fmt.Println(err)
	}
	defer c.RollbackTransaction(tid)
	//c.CreateDatabase()
	/*fmt.Println(c.CheckConfiguration(tid, "Configuration"))
	fmt.Println(c.CheckConfiguration(tid, "Record"))
	fmt.Println(c.CheckConfiguration(tid, "Storage"))*/
	c.CommitTransaction(tid)
}
func storageNode() {
	var err error
	var client client.IClient
	if client, err = httpserver.NewClient("http://localhost:8080", "Storage", "Test", "Storage.Default"); err != nil {
		return
	}
	//client := c.NewClient("Test", "Test", "Storage.Default")
	storage := storage.StartStorage("Default", client, "D:/SUDStorage")
	fmt.Println(storage)

}
func loadConfigurations(c *core.Core, configDir string) error {
	configDir = path.Clean(configDir)

	var err error
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(configDir); err != nil {
		return err
	}
	for _, file := range files {
		n := file.Name()
		nf := path.Join(configDir, n)
		if file.IsDir() {
			loadConfigurations(c, nf)
		} else {
			var data []byte
			if data, err = ioutil.ReadFile(nf); err != nil {
				return err
			}
			confName := n[0 : len(n)-len(path.Ext(n))]
			conf := core.NewConfiguration([]string{})
			if err = conf.LoadJson(data); err != nil {
				return err
			}
			c.AddBaseConfiguration(confName, conf)
		}
	}
	return nil
}
func StartServer(end chan error) {
	var err error
	var c *core.Core
	if c, err = core.NewCore("test", "user=suduser dbname=test password=Pa$$w0rd sslmode=disable"); err != nil {
		fmt.Println(err)
	}
	loadConfigurations(c, "./configuration/")
	core.InitStdModule(c)
	storage.InitModule(c)
	checkTest(c)
	server := httpserver.NewServer(c, ":8080")
	end <- server.Start()
}
func StartClient() {
	var err error
	var tid string
	var client client.IClient
	if client, err = httpserver.NewClient("http://localhost:8080", "Test", "Test", "Storage"); err != nil {
		return
	}
	if tid, err = client.BeginTransaction(); err != nil {
		fmt.Println(err)
	}
	defer client.RollbackTransaction(tid)
	Result, err := client.Call(
		"Storage.SetStream",
		map[string]interface{}{
			"Storage":        "Default",
			"Stream":         ([]byte)("Stream1"),
			"TransactionUID": tid},
		time.Second*1000, "")
	//Result, err := c.Call("Storage", "Storage.SetStream", map[string]interface{}{"Storage": "Default", "Stream": ([]byte)("Stream1"), "TransactionUID": tid}, time.Second*1000)
	fmt.Println("Set Stream", Result, err)
	client.CommitTransaction(tid)
}
func main() {
	stop := make(chan error)
	//var err error
	go StartServer(stop)
	go storageNode()
	go StartClient()
	err := <-stop
	if err != nil {
		log.Fatalln(err)
	}
}
