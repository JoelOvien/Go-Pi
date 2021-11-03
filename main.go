package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//create a struct and slice for our dummy database
type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

///our db holding various objects
//the db is a slice so we could create and delete items in it freely
var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

//homelink or home page when the try to access ourr server from the "/"
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome homeeee!")
}


//Create Event

//When creating an event, we get data from the user’s end. The user enters data which is in the form of http Request data.
func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event

	//When creating an event, we get data from the user’s end. The user enters data which is in the form of http Request data.
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

// Once it is converted to a slice, we unmarshal it to fit into our event struct. Once it is successful,
// we append the new event struct to the events slice and also display
// the new event as an http Responsewith 201 Created Status Code.
	json.Unmarshal(reqBody, &newEvent)

//  Once it is successful, we append the new event struct to the events slice
	events = append(events, newEvent)
	
// display the new event as an http Responsewith 201 Created Status Code.
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

//Get one event
func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}


//Get all events
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}


///Update an event
func updateEvent(w http.ResponseWriter, r *http.Request) {
	//get the value from the request endpoint and pass it to eventId variable
	eventID := mux.Vars(r)["id"]

	//create a variable updated event which is of type event
	var updatedEvent event

	//use ioutil to read the value coming fron the request body and store it to err and reqBody
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	//unmarshal the data and from our request body and updated event
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}


//Delete an event
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func main(){
	// initEvents()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event/create", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/getOne/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/update/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/delete/{id}", deleteEvent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

