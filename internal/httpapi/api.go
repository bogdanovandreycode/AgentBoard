package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/go-chi/chi/v5"
)

type API struct{ Service *service.Service }

func New(s *service.Service) http.Handler {
	a := &API{Service: s}
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/projects", a.listProjects)
		r.Post("/projects", a.createProject)
		r.Get("/projects/{projectID}/board", a.board)
		r.Post("/projects/{projectID}/tasks", a.createTask)
		r.Get("/tasks/{taskID}", a.getTask)
		r.Patch("/tasks/{taskID}", a.updateTask)
		r.Delete("/tasks/{taskID}", a.deleteTask)
		r.Post("/tasks/{taskID}/move", a.moveTask)
		r.Post("/tasks/{taskID}/assign", a.assignTask)
		r.Get("/tasks/{taskID}/history", a.history)
		r.Post("/tasks/{taskID}/comments", a.comment)
		r.Get("/tasks/{taskID}/tests", a.tests)
		r.Post("/tasks/{taskID}/tests", a.addTest)
		r.Get("/tasks/{taskID}/artifacts", a.artifacts)
		r.Post("/tasks/{taskID}/artifacts", a.addArtifact)
		r.Get("/tasks/{taskID}/usage", a.usage)
		r.Get("/projects/{projectID}/workers", a.workers)
		r.Post("/projects/{projectID}/workers", a.createWorker)
		r.Get("/workers/{workerID}", a.worker)
		r.Patch("/workers/{workerID}", a.updateWorker)
		r.Delete("/workers/{workerID}", a.deleteWorker)
		r.Get("/projects/{projectID}/properties", a.properties)
		r.Post("/projects/{projectID}/properties", a.createProperty)
		r.Patch("/properties/{propertyID}", a.updateProperty)
		r.Delete("/properties/{propertyID}", a.deleteProperty)
	})
	return r
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, core.ErrInvalidInput)
		return false
	}
	return true
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, core.ErrNotFound) {
		status = http.StatusNotFound
	} else if errors.Is(err, core.ErrInvalidInput) || errors.Is(err, core.ErrInvalidTransition) || errors.Is(err, core.ErrHumanTestPending) {
		status = http.StatusBadRequest
	} else if errors.Is(err, core.ErrTaskNotAccessible) {
		status = http.StatusForbidden
	}
	write(w, status, map[string]string{"error": err.Error()})
}

func (a *API) listProjects(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.Store.ListProjects(r.Context())
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) createProject(w http.ResponseWriter, r *http.Request) {
	var in struct{ Name, Path string }
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.Store.EnsureProject(r.Context(), in.Name, in.Path)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 201, v)
}
func (a *API) board(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.HumanBoard(r.Context(), chi.URLParam(r, "projectID"))
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) createTask(w http.ResponseWriter, r *http.Request) {
	var in core.TaskInput
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.Store.CreateTask(r.Context(), chi.URLParam(r, "projectID"), in)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 201, v)
}
func (a *API) getTask(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.HumanTask(r.Context(), chi.URLParam(r, "taskID"))
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) updateTask(w http.ResponseWriter, r *http.Request) {
	var in core.TaskUpdate
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.HumanUpdateTask(r.Context(), chi.URLParam(r, "taskID"), in)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) deleteTask(w http.ResponseWriter, r *http.Request) {
	if e := a.Service.HumanDeleteTask(r.Context(), chi.URLParam(r, "taskID")); e != nil {
		writeError(w, e)
		return
	}
	w.WriteHeader(204)
}
func (a *API) moveTask(w http.ResponseWriter, r *http.Request) {
	var in struct{ State string }
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.HumanMove(r.Context(), chi.URLParam(r, "taskID"), in.State)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) assignTask(w http.ResponseWriter, r *http.Request) {
	var in core.Assignee
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.HumanAssign(r.Context(), chi.URLParam(r, "taskID"), in)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) history(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.ListHistory(r.Context(), chi.URLParam(r, "taskID"))
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) comment(w http.ResponseWriter, r *http.Request) {
	var in struct{ Content string }
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.HumanComment(r.Context(), chi.URLParam(r, "taskID"), in.Content)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 201, v)
}
func (a *API) tests(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.ListTests(r.Context(), chi.URLParam(r, "taskID"))
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) addTest(w http.ResponseWriter, r *http.Request) {
	var in core.TestRunInput
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.HumanRecordTest(r.Context(), chi.URLParam(r, "taskID"), in)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 201, v)
}
func (a *API) artifacts(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.ListArtifacts(r.Context(), chi.URLParam(r, "taskID"))
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) addArtifact(w http.ResponseWriter, r *http.Request) {
	var in core.ArtifactInput
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.HumanAddArtifact(r.Context(), chi.URLParam(r, "taskID"), in)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 201, v)
}
func (a *API) usage(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.ListUsage(r.Context(), chi.URLParam(r, "taskID"))
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) workers(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.Store.ListWorkers(r.Context(), chi.URLParam(r, "projectID"), false)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) createWorker(w http.ResponseWriter, r *http.Request) {
	var in core.WorkerInput
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.Store.CreateWorker(r.Context(), chi.URLParam(r, "projectID"), in)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 201, v)
}
func (a *API) worker(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.Store.GetWorker(r.Context(), chi.URLParam(r, "workerID"))
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) updateWorker(w http.ResponseWriter, r *http.Request) {
	var in core.WorkerInput
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.Store.UpdateWorker(r.Context(), chi.URLParam(r, "workerID"), in)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) deleteWorker(w http.ResponseWriter, r *http.Request) {
	if e := a.Service.Store.ArchiveWorker(r.Context(), chi.URLParam(r, "workerID")); e != nil {
		writeError(w, e)
		return
	}
	w.WriteHeader(204)
}
func (a *API) properties(w http.ResponseWriter, r *http.Request) {
	v, e := a.Service.ListPropertyDefinitions(r.Context(), chi.URLParam(r, "projectID"))
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) createProperty(w http.ResponseWriter, r *http.Request) {
	var in core.PropertyDefinitionInput
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.CreatePropertyDefinition(r.Context(), chi.URLParam(r, "projectID"), in)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 201, v)
}
func (a *API) updateProperty(w http.ResponseWriter, r *http.Request) {
	var in core.PropertyDefinitionInput
	if !decode(w, r, &in) {
		return
	}
	v, e := a.Service.UpdatePropertyDefinition(r.Context(), chi.URLParam(r, "propertyID"), in)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func (a *API) deleteProperty(w http.ResponseWriter, r *http.Request) {
	if e := a.Service.DeletePropertyDefinition(r.Context(), chi.URLParam(r, "propertyID")); e != nil {
		writeError(w, e)
		return
	}
	w.WriteHeader(204)
}
