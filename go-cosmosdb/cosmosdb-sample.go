package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

func main() {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{"jiren-acs.documents.azure.com:10255"},
		Database: "",
		Username: "jiren-acs",
		Password: "C6nb062RuVzZkjnqkipzokcxmv2YLNlLCpO2rDgo30XgRMZeHdJBrVX360CSIIQQyYzsLtessvrl4lZyBHq7cw==",
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		},
		Timeout: time.Second * 10,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Cla"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}
