package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"main/db"
	"main/models"
	"net/http"
	"strconv"
)

type Meta struct {
	Action string `json:"action"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func PageHome(w http.ResponseWriter, request *http.Request) {
	sendJSONResponse(w, http.StatusOK, Response{
		Meta: Meta{Action: "home"},
		Data: map[string]interface {
		}{},
	})
}

//	func GetThing(w http.ResponseWriter, r *http.Request) {
//		var things = []models.Thing{}
//		if err := db.DB.Db.Find(&things).Error; err != nil {
//			sendJSONResponse(w, http.StatusInternalServerError, Response{
//				Meta: Meta{Action: "error"},
//				Data: map[string]interface{}{"error": err.Error()},
//			})
//			return
//		}
//
//		if things == nil {
//			things = make([]models.Thing, 0)
//		}
//		responseData := map[string]interface{}{
//			"items": things,
//			"count": len(things),
//		}
//
//		sendJSONResponse(w, http.StatusOK, Response{
//			Meta: Meta{Action: "things"},
//			Data: responseData,
//		})
//	}
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
func GetThing(w http.ResponseWriter, r *http.Request) {
	things := []models.Thing{}

	//resultDB := db.DB.Db.Find(&things)
	//if resultDB.Error != nil {
	//	http.Error(w, resultDB.Error.Error(), http.StatusBadRequest)
	//	return
	//}

	resultDB := db.DB.Db.Model(&models.Thing{}).Select("id, name").Find(&things)
	if resultDB.Error != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": resultDB.Error.Error(),
		})
		return
	}

	resp, err := json.Marshal(things)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func PostThing(w http.ResponseWriter, r *http.Request) {
	thing := new(models.Thing)
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &thing); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(thing)

	resultDB := db.DB.Db.Create(&thing)
	if resultDB.Error != nil {
		http.Error(w, resultDB.Error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func GetThingByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thing := new(models.Thing)

	//resultDB := db.DB.Db.First(&thing, id)
	//if resultDB.Error != nil {
	//	http.Error(w, resultDB.Error.Error(), http.StatusBadRequest)
	//	return
	//}

	resultDB := db.DB.Db.Model(&models.Thing{}).Select("id, name").First(&thing, id)
	if resultDB.Error != nil {
		sendJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": "Thing not found",
		})
		return
	}

	resp, err := json.Marshal(thing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func DeleteThing(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultDB := db.DB.Db.Delete(&models.Thing{}, id)
	if resultDB.Error != nil {
		http.Error(w, resultDB.Error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
