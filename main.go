package companies

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const (
	// remove this once GAE moves to Go 1.7
	StatusUnprocessableEntity = 422
)

type Company struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

func (c *Company) buildID() string {
	return fmt.Sprintf("%v:%v", c.Exchange, c.Symbol)
}

func init() {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/companies", companiesIndexHandler).Methods("GET")
	api.HandleFunc("/companies", companiesCreateHandler).Methods("POST")

	http.Handle("/", router)
}

func ParseCompany(body io.Reader) (*Company, error) {
	var company Company

	dec := json.NewDecoder(body)
	err := dec.Decode(&company)

	if err != nil {
		return nil, err
	}

	if company.ID == "" {
		company.ID = company.buildID()
	}

	return &company, err
}

func companiesIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// create a new App Engine context from the HTTP request.
	ctx := appengine.NewContext(r)

	var companies []Company

	// create a new query on the kind Company
	q := datastore.NewQuery("Company")

	_, err := q.GetAll(ctx, &companies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)

	// use the encoder to encode companies, which could fail.
	err = enc.Encode(companies)

	// if it failed, log the error and stop execution.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func companiesCreateHandler(w http.ResponseWriter, r *http.Request) {
	company, err := ParseCompany(r.Body)

	if err != nil {
		http.Error(w, err.Error(), StatusUnprocessableEntity)
		return
	}

	// create a new App Engine context from the HTTP request.
	ctx := appengine.NewContext(r)

	// create a new complete key of kind Company, using Company ID as the key.
	key := datastore.NewKey(ctx, "Company", company.ID, 0, nil)

	// put company in the datastore.
	key, err = datastore.Put(ctx, key, company)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	enc := json.NewEncoder(w)

	// use the encoder to encode companies, which could fail.
	if err = enc.Encode(company); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
