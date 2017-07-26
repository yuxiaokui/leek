package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Result struct {
	Url    string
	Status string
}

type Task struct {
	Id   bson.ObjectId "_id"
	Url  string
	Time string
}

func httpStatus(url string, ch chan<- string) {
	c := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := c.Head(url)
	if err != nil {
		fmt.Println("Connect Error!")
		ch <- fmt.Sprintf(url + "#" + "error")
		return
	}
	ch <- fmt.Sprintf(url+"#%d", resp.StatusCode)
	return
}

func getTask() string {
	session, err := mgo.Dial("mongodb.t2.daoapp.io:61257")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("task").C("tasks")
	result := Task{}
	err = c.Find(bson.M{}).One(&result)
	if err != nil {
		fmt.Println("Find error!")
	}

	err = c.Remove(bson.M{"url": result.Url})
	if err != nil {
		fmt.Println("Remove error!")
	}
	return result.Url

}

func run() {
	ch := make(chan string)
	tgtUrl := getTask()
	if tgtUrl == "" {
		time.Sleep(5 * time.Second)
		return
	}
	tgtUrl = strings.TrimSpace(tgtUrl)
	tgtUrl = "http://" + tgtUrl
	go httpStatus(tgtUrl, ch)
	result := fmt.Sprintf(<-ch)
	fmt.Println(result)
	res := strings.Split(result, "#")

	session, err := mgo.Dial("mongodb.t2.daoapp.io:61257")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("task").C("result")
	err = c.Insert(&Result{res[0], res[1]})
	if err != nil {
		panic(err)
	}

}

func main() {
	for {
		run()
	}
}
