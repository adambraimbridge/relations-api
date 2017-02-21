package relations

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type HttpHandlers struct {
	relationsDriver    Driver
	cacheControlHeader string
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func NewHttpHandlers(relationsDriver Driver, cacheControlHeader string) HttpHandlers {
	return HttpHandlers{relationsDriver, cacheControlHeader}
}

func (hh *HttpHandlers) HealthCheck(neoURL string) v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to Relations API requests",
		Name:             "Check connectivity to Neo4j",
		PanicGuide:       "https://dewey.ft.com/upp-relations-api.html",
		Severity:         1,
		TechnicalSummary: fmt.Sprintf(`Cannot connect to Neo4j (%v). Check that Neo4j instance is up and running`, neoURL),
		Checker:          hh.Checker,
	}
}

func (hh *HttpHandlers) Checker() (string, error) {
	err := hh.relationsDriver.checkConnectivity()
	if err == nil {
		return "Connectivity to Neo4j is ok", err
	}
	return "Error connecting to Neo4j", err
}

func (hh *HttpHandlers) GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := hh.Checker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}

}

func (hh *HttpHandlers) GetRelations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	uuid := vars["uuid"]
	err := validateUuid(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg, _ := json.Marshal(ErrorMessage{fmt.Sprintf("The given uuid is not valid, err=%v", err)})
		w.Write([]byte(msg))
		return
	}
	relations, found, err := hh.relationsDriver.read(uuid)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		msg, _ := json.Marshal(ErrorMessage{fmt.Sprintf("Error retrieving relations for %s, err=%v", uuid, err)})
		w.Write([]byte(msg))
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		msg, _ := json.Marshal(ErrorMessage{fmt.Sprintf("No relations found for content with uuid %s", uuid)})
		w.Write([]byte(msg))
		return
	}

	w.Header().Set("Cache-Control", hh.cacheControlHeader)
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(relations); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorMessage{fmt.Sprintf("Error parsing result for content with uuid %s, err=%v", uuid, err)})
		w.Write([]byte(msg))
	}
}

func validateUuid(contentUUID string) error {
	parsedUUID, err := uuid.FromString(contentUUID)
	if err != nil {
		return err
	}
	if contentUUID != parsedUUID.String() {
		return fmt.Errorf("Parsed UUID (%v) is different than the given uuid (%v).", parsedUUID, contentUUID)
	}
	return nil
}
