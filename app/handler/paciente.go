package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/katoozi/golang-mongodb-rest-api/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// results count per page
var limit int64 = 10

// CreatePaciente will handle the create person post request
func CreatePaciente(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	person := new(model.Paciente)
	err := json.NewDecoder(req.Body).Decode(person)
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "body json request have issues!!!", nil)
		return
	}
	result, err := db.Collection("pacientes").InsertOne(nil, person)
	if err != nil {
		switch err.(type) {
		case mongo.WriteException:
			ResponseWriter(res, http.StatusNotAcceptable, "username or email already exists in database.", nil)
		default:
			ResponseWriter(res, http.StatusInternalServerError, "Error while inserting data.", nil)
		}
		return
	}
	person.ID = result.InsertedID.(primitive.ObjectID)
	ResponseWriter(res, http.StatusCreated, "", person)
}

// GetPacientes will handle pacientes list get request
func GetPacientes(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var personList []model.Paciente
	pageString := req.FormValue("page")
	page, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		page = 0
	}
	page = page * limit
	findOptions := options.FindOptions{
		Skip:  &page,
		Limit: &limit,
		Sort: bson.M{
			"_id": -1, // -1 for descending and 1 for ascending
		},
	}
	curser, err := db.Collection("pacientes").Find(nil, bson.M{}, &findOptions)
	if err != nil {
		log.Printf("Error while quering collection: %v\n", err)
		ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	err = curser.All(context.Background(), &personList)
	if err != nil {
		log.Fatalf("Error in curser: %v", err)
		ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	ResponseWriter(res, http.StatusOK, "", personList)
}

// GetPaciente will give us person with special id
func GetPaciente(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var params = mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	var person model.Paciente
	err = db.Collection("pacientes").FindOne(nil, model.Paciente{ID: id}).Decode(&person)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			ResponseWriter(res, http.StatusNotFound, "person not found", nil)
		default:
			log.Printf("Error while decode to go struct:%v\n", err)
			ResponseWriter(res, http.StatusInternalServerError, "there is an error on server!!!", nil)
		}
		return
	}
	ResponseWriter(res, http.StatusOK, "", person)
}

// UpdatePaciente will handle the person update endpoint
func UpdatePaciente(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var updateData map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&updateData)
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "json body is incorrect", nil)
		return
	}
	// we dont handle the json decode return error because all our fields have the omitempty tag.
	var params = mux.Vars(req)
	oid, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	update := bson.M{
		"$set": updateData,
	}
	result, err := db.Collection("pacientes").UpdateOne(context.Background(), model.Paciente{ID: oid}, update)
	if err != nil {
		log.Printf("Error while updateing document: %v", err)
		ResponseWriter(res, http.StatusInternalServerError, "error in updating document!!!", nil)
		return
	}
	if result.MatchedCount == 1 {
		ResponseWriter(res, http.StatusAccepted, "", &updateData)
	} else {
		ResponseWriter(res, http.StatusNotFound, "person not found", nil)
	}
}
