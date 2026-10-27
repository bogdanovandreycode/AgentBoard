package httpapi

import (
	"net/http"

	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/go-chi/chi/v5"
)

func (a *API) getGlobalSettings(w http.ResponseWriter, r *http.Request) {
	value, err := a.Service.GetGlobalSettings(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	write(w, 200, value)
}

func (a *API) putGlobalSettings(w http.ResponseWriter, r *http.Request) {
	var in service.GlobalSettings
	if !decode(w, r, &in) {
		return
	}
	value, err := a.Service.UpdateGlobalSettings(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	write(w, 200, value)
}

func (a *API) updateProject(w http.ResponseWriter, r *http.Request) {
	var in service.ProjectUpdate
	if !decode(w, r, &in) {
		return
	}
	value, err := a.Service.UpdateProject(r.Context(), chi.URLParam(r, "projectID"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	write(w, 200, value)
}

func (a *API) listDrafts(w http.ResponseWriter, r *http.Request) {
	value, err := a.Service.ListTaskDrafts(r.Context(), chi.URLParam(r, "projectID"))
	if err != nil {
		writeError(w, err)
		return
	}
	write(w, 200, value)
}

func (a *API) saveDraft(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	var in service.DraftSave
	if !decode(w, r, &in) {
		return
	}
	value, err := a.Service.SaveTaskDraft(r.Context(), chi.URLParam(r, "projectID"), chi.URLParam(r, "draftID"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	write(w, 200, value)
}

func (a *API) discardDraft(w http.ResponseWriter, r *http.Request) {
	if err := a.Service.DiscardTaskDraft(r.Context(), chi.URLParam(r, "projectID"), chi.URLParam(r, "draftID")); err != nil {
		writeError(w, err)
		return
	}
	write(w, 200, map[string]bool{"discarded": true})
}

func (a *API) publishDraft(w http.ResponseWriter, r *http.Request) {
	var in struct{ Revision int }
	if !decode(w, r, &in) {
		return
	}
	value, err := a.Service.PublishTaskDraft(r.Context(), chi.URLParam(r, "projectID"), chi.URLParam(r, "draftID"), in.Revision)
	if err != nil {
		writeError(w, err)
		return
	}
	write(w, 201, value)
}
