package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

// TemplateData defines the data injected into the script template.
type TemplateData struct {
	ProjectName   string
	CreatedAt     string
	UserName      string
	ContainerName string
	Network       string
	DNS           string
	PUID          string
	PGID          string
	TZ            string
	Port          string
	Volume        string
	HostIP        string
}

// GenerateScript loads the appropriate template for the given OS and returns the rendered script.
func GenerateScript(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	osType := vars["osType"] // Expected values: "linux", "windows", "osx", etc.

	// Construct the template file path.
	templateFile := filepath.Join("templates", "jelly-bash-" + osType + ".sh")
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		http.Error(w, "Template for OS type not found", http.StatusNotFound)
		return
	}

	// Set default values.
	data := TemplateData{
		ProjectName:   "Nexus Creator Vault",
		CreatedAt:     time.Now().Format(time.RFC1123),
		UserName:      "Administrator",
		ContainerName: "nexus-creator-vault",
		Network:       "Inner-Athena",
		DNS:           "10.20.0.20",
		PUID:          "1050",
		PGID:          "1050",
		TZ:            "America/Colorado",
		Port:          "1050",
		Volume:        "creator-vault000",
		HostIP:        "10.20.0.1",
	}

	// If Content-Type is JSON, override defaults with posted values.
	if r.Header.Get("Content-Type") == "application/json" {
		var input TemplateData
		if err := json.NewDecoder(r.Body).Decode(&input); err == nil {
			if input.ProjectName != "" {
				data.ProjectName = input.ProjectName
			}
			if input.UserName != "" {
				data.UserName = input.UserName
			}
			if input.ContainerName != "" {
				data.ContainerName = input.ContainerName
			}
			if input.Network != "" {
				data.Network = input.Network
			}
			if input.DNS != "" {
				data.DNS = input.DNS
			}
			if input.PUID != "" {
				data.PUID = input.PUID
			}
			if input.PGID != "" {
				data.PGID = input.PGID
			}
			if input.TZ != "" {
				data.TZ = input.TZ
			}
			if input.Port != "" {
				data.Port = input.Port
			}
			if input.Volume != "" {
				data.Volume = input.Volume
			}
			if input.HostIP != "" {
				data.HostIP = input.HostIP
			}
		}
	}

	// Parse and execute the template.
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Error parsing template: " + err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template: " + err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	router := mux.NewRouter()

	// Endpoint to generate scripts for a given OS.
	// Example: GET /scripts/linux
	router.HandleFunc("/scripts/{osType}", GenerateScript).Methods("GET", "POST")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
