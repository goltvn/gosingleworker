package main

import (
	"fmt"
	"net/http"
	"time"

	"database/sql"
	"log"
	"math/rand"
	"sync"
	//"time"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

var db, errDB = sql.Open("sqlite3", "jobdatabase.sqlite")
var waitGroup sync.WaitGroup
var data chan string

const (
	setupSQL = `
	CREATE TABLE IF NOT EXISTS jobs ( JobNumber	TEXT PRIMARY KEY,Status	INTEGER,Message	TEXT);		

	`
)

// update a jobid in database with a new status / message

func update(db *sql.DB, jobnumber string, status string, message string) int {
	result, err := db.Exec(`UPDATE jobs set Status = (?) where jobs.JobNumber = (?);`, status, jobnumber)
	if err != nil {
		fmt.Printf("user insert. Exec error=%s\r\n", err)
		return -1
	}
	_, err = result.LastInsertId()
	if err != nil {
		fmt.Printf("user updater. LastInsertId error=%s\r\n", err)
		return -1
	}
	fmt.Printf("update : %s\r\n", status)
	return 1
}

// write a new entry in the database

func write(db *sql.DB, jobnumber string, status string, message string) int {
	result, err := db.Exec(`INSERT INTO jobs (jobnumber,status) VALUES (?, ?);`, jobnumber, status)
	if err != nil {
		fmt.Printf("user insert. Exec error=%s\r\n", err)
		return -1
	}
	_, err = result.LastInsertId()
	if err != nil {
		fmt.Printf("user writer. LastInsertId error=%s\r\n", err)
		return -1
	}
	fmt.Printf("+")
	return 1
}

// read the status of a jobnumber in the database

func read(db *sql.DB, jobnumber string) int {
	var status string
	var message string
	sqlStatement := `SELECT status,message FROM jobs WHERE jobnumber=$1`
	row := db.QueryRow(sqlStatement, jobnumber)
	err := row.Scan(&status, &message)

	if err != nil {
		if err == sql.ErrNoRows {
			//			fmt.Println("Zero rows found") // just not found - not added in the queue
			return -1
		}
	}
	i1, err := strconv.Atoi(status)
	if err != nil {
		fmt.Println("error itoa")
		return -1
	}
	return i1
}

// read the current active job ( status == 2)

func readActiveJob(db *sql.DB) string {
	var jobnumber string
	sqlStatement := `SELECT jobnumber FROM jobs WHERE status=2`
	row := db.QueryRow(sqlStatement)
	err := row.Scan(&jobnumber)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
			jobnumber = "NONE"
		}
	}
	return jobnumber
}

// count with variable "where " statement

func countquery(db *sql.DB, countcmd string) string {
	var searchresult string

	sqlStatement := `SELECT count(*) as searchresult FROM jobs WHERE  ` + countcmd
	row := db.QueryRow(sqlStatement)
	err := row.Scan(&searchresult)

	if err != nil {
		if err == sql.ErrNoRows {

			return "NONE"
		}
	}
	return searchresult
}

// APi: get finished jobs ( returns which job is status 3)

func getFinishedJobs(w http.ResponseWriter, r *http.Request) {
	reply := countquery(db, "status=3;")
	fmt.Fprintf(w, "FinishedJobs: %s\n", reply)
}

// APi: get waiting jobs ( returns which job is status 1)

func getWaitingJobs(w http.ResponseWriter, r *http.Request) {
	reply := countquery(db, "status=1;")
	fmt.Fprintf(w, "WaitingJobs: %s\n", reply)
}

// APi: get active job ( returns which job is status 2 - should only be one)

func getActiveJob(w http.ResponseWriter, r *http.Request) {
	reply := readActiveJob(db)
	fmt.Fprintf(w, "activejob: %s\n", reply)
}

// APi: get status of a jobnumber

func getJobStatus(w http.ResponseWriter, r *http.Request) {

	jobnumber := r.FormValue("id")
	err := read(db, jobnumber)
	if err == -1 {
		fmt.Fprintf(w, "No Such Job: %s\n", jobnumber)
		return
	}

	fmt.Fprintf(w, "Jobstatus: %d\n", err)
}

var counter = 0

// APi: Start a job with a jobnumber, it will first check database if the job exists,
// if the job already exists then it wont be added.
// A new job starts as wait and will get status 1 in the database

func startJob(w http.ResponseWriter, r *http.Request) {

	var status string
	var message string
	counter++
	jobnumber := r.FormValue("id")
	err := read(db, jobnumber)
	if err != -1 {
		fmt.Fprintf(w, "JobID already Exists: %s\n", jobnumber)

	} else {
		status = "1"
		message = "queue"
		err := write(db, jobnumber, status, message)
		if err == -1 {
			fmt.Fprintf(w, "Write Error: %s\n", time.Now())
		} else {
			fmt.Fprintf(w, "Job Created: %s\n", jobnumber)
		}
		data <- (jobnumber)
	}
	fmt.Printf("incomings:%d\r\n", counter)
}

// worker thread - processes each channel job (simulating it with a wait time of up to 6 seconds )
// when job is starting being processed, status will change to 2 ( job is processing )
// when job is done, status is being changed to 3 (finished job)

func worker() {
	var status string
	var message string

	fmt.Println("Worker Started")
	defer func() {
		fmt.Println("Destroy Worker")
		waitGroup.Done()
	}()
	for {
		value, ok := <-data
		if !ok {
			fmt.Println("Channel is closed!")
			break
		}
		status = "2"
		message = "processing"
		update(db, value, status, message)

		time.Sleep(time.Duration(rand.Intn(8000)) * time.Millisecond)

		status = "3"
		message = "done"
		update(db, value, status, message)

		fmt.Println("worker - ID done:", value)
	}
}

// main
func main() {
	data = make(chan string)
	go worker()

	fmt.Println("Starting server on port :9090")
	db.SetMaxOpenConns(1) // prevents locked database error 
	_, err := db.Exec(setupSQL)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/StartJob", startJob)
	http.HandleFunc("/GetJobStatus", getJobStatus)
	http.HandleFunc("/GetActiveJob", getActiveJob)
	http.HandleFunc("/GetFinishedJobs", getFinishedJobs)
	http.HandleFunc("/GetWaitingJobs", getWaitingJobs)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
