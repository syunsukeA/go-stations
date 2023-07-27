package handler

import (
	"context"
	"net/http"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		todo_request := model.CreateTODORequest{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&todo_request)
		if len(todo_request.Subject) == 0 {
			w.WriteHeader(http.StatusBadRequest)
    		return
		}
		create_todo_response, err := h.Create(r.Context(), &todo_request)
		jsonData, err := json.Marshal(create_todo_response)
		if err != nil {
			panic(err) // 何かしらwに書き込むべきかも？
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	} else if r.Method == "PUT" {
		todo_request := model.UpdateTODORequest{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&todo_request)
		if todo_request.ID == 0 || len(todo_request.Subject) == 0 {
			w.WriteHeader(http.StatusBadRequest)
    		return
		}
		create_todo_response, err := h.Update(r.Context(), &todo_request)
		if reflect.TypeOf(err) == reflect.TypeOf(model.ErrNotFound{}) {
			w.WriteHeader(http.StatusNotFound)
    		return
		}
		jsonData, err := json.Marshal(create_todo_response)
		if err != nil {
			panic(err) // 何かしらwに書き込むべきかも？
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	} else if r.Method == "GET" {
		todo_request := model.ReadTODORequest{}
		query := r.URL.Query()
		prev_id, _ := strconv.ParseInt(query.Get("prev_id"), 10, 64)
		size, _ := strconv.ParseInt(query.Get("size"), 10, 64)
		todo_request.PrevID = prev_id
		todo_request.Size = size
		read_todo_response, err := h.Read(r.Context(), &todo_request)
		if err != nil {
			panic(err)
		}
		jsonData, err := json.Marshal(read_todo_response)
		if err != nil {
			panic(err) // 何かしらwに書き込むべきかも？
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	} else if r.Method == "DELETE" {
		todo_request := model.DeleteTODORequest{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&todo_request)
		if len(todo_request.IDs) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo_response, err := h.Delete(r.Context(), &todo_request)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		jsonData, err := json.Marshal(todo_response)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, _ := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	return &model.ReadTODOResponse{TODOs: todos}, err
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *todo}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	return &model.DeleteTODOResponse{}, err
}
