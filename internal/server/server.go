package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/bblanton/gspro-control/internal/actions"
	"github.com/bblanton/gspro-control/internal/input"
)

type commandResponse struct {
	Action string `json:"action"`
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// test seam: allows tests to stub key execution without OS side effects
var executeComboFunc = input.ExecuteCombo

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		registry := actions.GetRegistry()
		list := registry.List()
		_ = json.NewEncoder(w).Encode(list)
	})

	mux.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		actionName := r.URL.Query().Get("action")
		if actionName == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(commandResponse{Action: "", Result: "error", Error: "missing action param"})
			return
		}

		registry := actions.GetRegistry()
		combo, ok := registry.Lookup(actionName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(commandResponse{Action: actionName, Result: "error", Error: "unknown action"})
			return
		}

		comboDisplay := strings.Join([]string(combo), "+")
		log.Printf("command start action=%s combo=%s", actionName, comboDisplay)

		if err := executeComboFunc([]string(combo)); err != nil {
			log.Printf("command result=error action=%s combo=%s error=%v", actionName, comboDisplay, err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(commandResponse{Action: actionName, Result: "error", Error: err.Error()})
			return
		}

		log.Printf("command result=success action=%s combo=%s", actionName, comboDisplay)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(commandResponse{Action: actionName, Result: "success"})
	})
}
