package search

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Course representa la estructura de un curso buscado en Solr
type Course struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Instructor   string `json:"instructor"`
	Duration     int    `json:"duration"`
	Level        string `json:"level"`
	Availability bool   `json:"availability"`
}

// SearchCourses maneja la búsqueda de cursos utilizando Solr
// SearchCourses maneja la búsqueda de cursos utilizando Solr
func SearchCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "El parámetro 'q' es requerido", http.StatusBadRequest)
		return
	}

	solrQuery := "*" + query + "*"
	solrURL := "http://localhost:8983/solr/courses/select?q=title:" + url.QueryEscape(solrQuery)

	resp, err := http.Get(solrURL)
	if err != nil {
		log.Println("Error al conectar con Solr:", err)
		http.Error(w, "Error al conectar con el motor de búsqueda", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error al leer la respuesta de Solr:", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
