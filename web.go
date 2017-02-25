package main

import (
	"reflect"
	"time"

	"fmt"

	"github.com/crazyprograms/sud/core"
	/* "github.com/crazyprograms/sud/httpserver" */
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
	fmt.Println(c.CheckConfiguration(tid, "Configuration"))
	fmt.Println(c.CheckConfiguration(tid, "Document"))
	fmt.Println(c.CheckConfiguration(tid, "Storage"))
	c.CommitTransaction(tid)
}
func storageNode(c *core.Core) {
	client := c.NewClient("Test", "Test", "Storage.Default")
	storage := storage.StartStorage("Default", client, "D:/SUDStorage")
	fmt.Println(storage)
}
func sq() string {
	return "ABC"
}
func sh(a core.Object) {
	s1, ok1 := a.(string)
	s2, ok2 := a.(*string)
	fmt.Println(a, s1, ok1, s2, ok2)

}

type i1 interface {
	Check(value interface{})
}
type st1 struct {
}

func (s *st1) Check(value interface{}) {
	s1, ok1 := value.(string)
	fmt.Println(value, s1, ok1)
	t := reflect.TypeOf(value)
	fmt.Println(t.Name())

}
func main() {
	s1 := sq()
	st := st1{}
	Poles := map[string]interface{}{"A": s1}
	for pole, value := range Poles {
		fmt.Println(pole)
		st.Check(value)
	}

	var err error
	var c *core.Core
	if c, err = core.NewCore("test", "user=suduser dbname=test password=Pa$$w0rd sslmode=disable"); err != nil {
		fmt.Println(err)
	}
	/*server := httpserver.NewServer(c, ":8080")
	server.Start()
	fmt.Println("end")
	return*/

	checkTest(c)
	storageNode(c)
	var tid string
	if tid, err = c.BeginTransaction(); err != nil {
		fmt.Println(err)
	}
	defer c.RollbackTransaction(tid)

	Result, err := c.Call("Storage", "Storage.SetStream", map[string]interface{}{"Storage": "Default", "Stream": ([]byte)("Stream1"), "TransactionUID": tid}, time.Second*1000)
	fmt.Println("Set Stream", Result, err)
	c.CommitTransaction(tid)

	time.Sleep(time.Second * 2)
	/*go func() {
		defer fmt.Println("end")
		for {
			_, Result, err := c.Listen("Test", "TestAsync", time.Second)
			if err != nil {
				return
			}
			Result <- "Test async Ok"
		}
	}()

	fmt.Println(c.Call("Test", "TestStd", map[string]interface{}{}, time.Second))
	fmt.Println(c.Call("Test", "TestAsync", map[string]interface{}{}, time.Second))
	time.Sleep(time.Second * 2)
	/*doc, err := server.NewDocument(tid, "Document", "Document")
	fmt.Println(doc.SetPoleValue("Document.DocumentType", "Q1"))
	fmt.Println(server.SaveDocument(tid, doc))
	//server.GetDocuments(tid, "Document", "Test", []string{"Document.*"}, []storage.IDocumentWhere{})
	server.CommitTransaction(tid)
	fmt.Println("end")
	/*c, _ := storage.Connect("test", "user=suduser dbname=test password=Pa$$w0rd sslmode=disable")
	defer c.Close()
	config := &storage.Configuration{}
	config.LoadConfiguration("Configuration")
	fmt.Println(c.CheckConfig(config))
	/*
		c, err := storage.Connect("Data Source=192.168.1.102;Initial Catalog=TestDB;User ID=sa;Password=Pa$$w0rd")
		if err != nil {
			log.Fatalln(err)
			return
		}
		defer c.Close()
		config := &storage.Configuration{}
		config.LoadConfiguration("Configuration")
		c.CreateDatabase()
		fmt.Println(c.CheckConfig(config))
		docs, err := c.GetDocuments(config, "Configuration", []string{"Configuration.DocumentType.Date1"}, []storage.IDocumentWhere{})
		fmt.Println(err)
		var doc storage.IDocument
		/*var i int
		for i, doc = range docs {
			poles := doc.GetPoleNames()
			for _, p := range poles {
				//fmt.Println(i,doc.GetPole(p).)
			}
		}
		/*fmt.Println(c.DatabaseExists())
		q1 := Query{
			Action: "Test1",
			Link:   map[string]uuid.UUID{"DocumentUID": uuid.NewV4()},
			Time:   map[string]time.Time{"Time1": time.Now()},
			String: map[string]string{"Name": "Name1"},
			V:      B{Value2: "qwe"},
		}
		q2 := Query{}
		str, err := json.Marshal(q1)
		fmt.Println(string(str), err)
		err2 := json.Unmarshal(str, &q2)
		fmt.Println(q1)
		fmt.Println(q2, err2)
		//http.HandleFunc("/json/", viewHandler)
		//http.ListenAndServe(":80", nil)*/
}
