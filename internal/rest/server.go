package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prismon/synthesis/internal/postgres"
)

// Server represents the REST API server
type Server struct {
	router       *mux.Router
	tenantRepo   *postgres.TenantRepository
	notebookRepo *postgres.NotebookRepository
	vectorRepo   *postgres.VectorRepository
	graphRepo    *postgres.GraphRepository
}

// NewServer creates a new REST API server
func NewServer(db *postgres.DB) *Server {
	s := &Server{
		router:       mux.NewRouter(),
		tenantRepo:   postgres.NewTenantRepository(db),
		notebookRepo: postgres.NewNotebookRepository(db),
		vectorRepo:   postgres.NewVectorRepository(db),
		graphRepo:    postgres.NewGraphRepository(db),
	}

	s.registerRoutes()

	return s
}

// registerRoutes registers all API routes
func (s *Server) registerRoutes() {
	// API v1 routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Tenant routes
	api.HandleFunc("/tenants", s.listTenants).Methods("GET")
	api.HandleFunc("/tenants", s.createTenant).Methods("POST")
	api.HandleFunc("/tenants/{id}", s.getTenant).Methods("GET")
	api.HandleFunc("/tenants/{id}", s.updateTenant).Methods("PUT")
	api.HandleFunc("/tenants/{id}", s.deleteTenant).Methods("DELETE")

	// Library routes
	api.HandleFunc("/libraries/by-tenant/{tenantId}", s.listLibraries).Methods("GET")
	api.HandleFunc("/libraries/by-tenant/{tenantId}", s.createLibrary).Methods("POST")
	api.HandleFunc("/libraries/{id}", s.getLibrary).Methods("GET")
	api.HandleFunc("/libraries/{id}", s.updateLibrary).Methods("PUT")
	api.HandleFunc("/libraries/{id}", s.deleteLibrary).Methods("DELETE")

	// Notebook routes
	api.HandleFunc("/notebooks/by-library/{libraryId}", s.listNotebooks).Methods("GET")
	api.HandleFunc("/notebooks/by-library/{libraryId}", s.createNotebook).Methods("POST")
	api.HandleFunc("/notebooks/{id}", s.getNotebook).Methods("GET")
	api.HandleFunc("/notebooks/{id}", s.updateNotebook).Methods("PUT")
	api.HandleFunc("/notebooks/{id}", s.deleteNotebook).Methods("DELETE")

	// Health check
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting REST API server on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}

// Health check handler
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "synthesis-api",
	})
}

// Tenant handlers

func (s *Server) listTenants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenants, err := s.tenantRepo.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenants)
}

func (s *Server) createTenant(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Create tenant not yet implemented",
	})
}

func (s *Server) getTenant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	tenant, err := s.tenantRepo.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

func (s *Server) updateTenant(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Update tenant not yet implemented",
	})
}

func (s *Server) deleteTenant(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.tenantRepo.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Library handlers (placeholders)

func (s *Server) listLibraries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

func (s *Server) createLibrary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Create library not yet implemented"})
}

func (s *Server) getLibrary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Get library not yet implemented"})
}

func (s *Server) updateLibrary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Update library not yet implemented"})
}

func (s *Server) deleteLibrary(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// Notebook handlers

func (s *Server) listNotebooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	libraryID := vars["libraryId"]

	notebooks, err := s.notebookRepo.ListByLibrary(ctx, libraryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notebooks)
}

func (s *Server) createNotebook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Create notebook not yet implemented"})
}

func (s *Server) getNotebook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	notebook, err := s.notebookRepo.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notebook)
}

func (s *Server) updateNotebook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Update notebook not yet implemented"})
}

func (s *Server) deleteNotebook(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.notebookRepo.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
