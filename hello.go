package hello

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("templates/*"))

var store = sessions.NewCookieStore([]byte("magical-secret-store-pickles"))
var passcode = "7852"

type JSONResponse map[string]interface{}

func (r JSONResponse) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

type QuestionResponse struct {
	Question Question
	Passcode string
}

type Question struct {
	Text    string
	Answers []Answer
}

type Answer struct {
	Text    string
	Correct bool
}

var questions = []Question{}

func (q *Question) isCorrect(index int) bool {
	return q.Answers[index].Correct
}

func init() {
	qbytes, err := ioutil.ReadFile("questions.json")
	if err != nil {
		log.Print("Find your files:", err)
		return
	}

	err = json.Unmarshal(qbytes, &questions)
	if err != nil {
		log.Print("Fix your JSON:", err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/", PageHandler).
		Methods("GET")
	r.HandleFunc("/", QuestionHandler).
		Methods("POST")
	http.Handle("/", r)
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "pickles")
	if err != nil {
		log.Printf("Error getting session: %v", err)
	}

	if session.Values["count"] == nil {
		session.Values["count"] = 0
	}

	index := session.Values["count"].(int)
	if index == len(questions) {
		err = templates.ExecuteTemplate(w, "main", QuestionResponse{Passcode: passcode})
	} else {
		question := questions[index]
		session.Save(r, w)
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Error saving index: %v", err)
		}
		err = templates.ExecuteTemplate(w, "main", QuestionResponse{Question: question})
	}
	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func QuestionHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "pickles")
	if err != nil {
		log.Printf("Error getting session: %v", err)
	}

	answer, err := strconv.ParseInt(r.FormValue("answer"), 10, 0)
	if err != nil {
		log.Printf("Malformed answer, let's call it 0: %v", err)
		answer = 0
	}

	index := session.Values["count"].(int)
	correct := questions[index].isCorrect(int(answer))

	if correct {
		index = index + 1
	} else {
		index = 0
	}

	buffer := new(bytes.Buffer)
	if index == len(questions) {
		session.Values["count"] = index
		err = session.Save(r, w)
		templates.ExecuteTemplate(buffer, "passcode", QuestionResponse{Passcode: passcode})
		fmt.Fprint(w, JSONResponse{"correct": true, "html": buffer.String()})
	} else {
		question := questions[index]
		session.Values["count"] = index
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Error saving index: %v", err)
		}
		templates.ExecuteTemplate(buffer, "question", QuestionResponse{Question: question})
		fmt.Fprint(w, JSONResponse{"correct": correct, "html": buffer.String()})
	}
}
