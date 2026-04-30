package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Estrutura de contacto
type Contact struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// Estrutura de lista de contactos (map)
type ContactService struct {
	Contacts map[int]Contact
}

// Handler do POST
func handleCreateContact(w http.ResponseWriter, r *http.Request, service *ContactService) {
	service.Create(w, r)
}

// Função para criar contacto - POST
func (c *ContactService) Create(w http.ResponseWriter, r *http.Request) {
	var contact Contact // Variável de contacto

	err := json.NewDecoder(r.Body).Decode(&contact) // Variável de erro

	// Se existir erro
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Caso contrário, atribui um id
	id := len(c.Contacts) + 1;
	contact.Id = id;

	c.Contacts[id] = contact // Adiciona o contacto à lista de contactos

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contact)
	w.WriteHeader(http.StatusCreated)
}

// Handler do GET (GET ALL e GET ONE)
func handleGetContacts(w http.ResponseWriter, r *http.Request, service *ContactService) {
	q := r.URL.Query() // Query params

	// Se os qp não estiverem vazios, o usuário quer um contacto especifico, caso contrário, listar tudo
	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Get(w, r, id)
	} else {
		service.List(w, r)
	}
}

// Função listar todos os contactos - GET ALL
func (c *ContactService) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contacts []Contact // Lista de contactos

	// Adicionar todos os contactos existenstes à lista
	for _, c := range c.Contacts {
		contacts = append(contacts, c)
	}

	json.NewEncoder(w).Encode(contacts) // Retornar a lista em json
}

// Função para listar apenas um contacto - GET ONE
func (c* ContactService) Get(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	// Verifica no map se o contacto existe com base na chave (id)
	if val, ok := c.Contacts[id]; ok {
		json.NewEncoder(w).Encode(val)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

// Handler do DELETE
func handleDeleteContact(w http.ResponseWriter, r *http.Request, service *ContactService) {
	q := r.URL.Query()
	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Delete(w, r, id)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

// Função para eliminar um contacto - DELETE
func (c *ContactService) Delete(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	// Verifica no map se o contacto existe com base na chave (id)
	if _, ok := c.Contacts[id]; ok {
		delete(c.Contacts, id)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

// Handler do PUT
func handleUpdateContact(w http.ResponseWriter, r *http.Request, service *ContactService) {
	q := r.URL.Query()
	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Update(w, r, id)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

// Função para atualizar um contacto - PUT
func (c* ContactService) Update(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	var contact Contact
	err := json.NewDecoder(r.Body).Decode(&contact)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := c.Contacts[id]; ok {
		contact.Id = id
		c.Contacts[id] = contact
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}

}


func main() {
	service := &ContactService{Contacts: make(map[int]Contact)}
	mux := http.NewServeMux()

	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		// Tratamento dos métodos HTTP
		switch r.Method {
		// GET (ALL OR ONE)
		case http.MethodGet:
			handleGetContacts(w, r, service)
		// POST
		case http.MethodPost:
			handleCreateContact(w, r, service)
		// DELETE
		case http.MethodDelete:
			handleDeleteContact(w, r, service)
		// PUT
		case http.MethodPut:
			handleUpdateContact(w, r, service)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}