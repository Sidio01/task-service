package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	e "gitlab.com/g6834/team26/task/internal/domain/errors"
	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/pkg/uuid"
)

func (s *Server) taskHandlers() http.Handler {
	r := chi.NewRouter()
	r.Delete("/tasks/{taskID}", s.DeleteTaskHandler)
	r.Post("/tasks/{taskID}/approve/{approvalLogin}", s.ApproveTaskHandler)
	r.Post("/tasks/{taskID}/decline/{approvalLogin}", s.DeclineTaskHandler)
	r.Get("/tasks/", s.GetTasksListHandler)
	r.Post("/tasks/run", s.RunTaskHandler)
	return r
}

func (s *Server) getCookies(r *http.Request) (refreshToken, accessToken string) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		switch cookie.Name {
		case "refresh_token":
			refreshToken = cookie.Value
		case "access_token":
			accessToken = cookie.Value
		}
	}
	// log.Println("refreshToken -", refreshToken)
	// log.Println("accessToken -", accessToken)
	return
}

func (s *Server) getValidationResult(r *http.Request) (string, error) {
	refreshToken, accessToken := s.getCookies(r)
	authResponseResult, authResponseLogin, err := s.task.Validate(refreshToken, accessToken)
	if err != nil {
		return "", err
		// return err
	}
	// log.Printf("grpc result: %v, grpc login: %v", authResponseResult, authResponseLogin)
	if !authResponseResult {
		return "", e.ErrAuthFailed
	}
	return authResponseLogin, nil
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	login, err := s.getValidationResult(r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "taskID")
	err = s.task.DeleteTask(login, id)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("{\"status\": \"deleted\"}"))
}

func (s *Server) ApproveTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	login, err := s.getValidationResult(r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "taskID")
	approvalLogin := chi.URLParam(r, "approvalLogin")
	err = s.task.ApproveTask(login, id, approvalLogin)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("{\"status\": \"approved\"}"))
}

func (s *Server) DeclineTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	login, err := s.getValidationResult(r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "taskID")
	approvalLogin := chi.URLParam(r, "approvalLogin")
	err = s.task.DeclineTask(login, id, approvalLogin)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("{\"status\": \"declined\"}"))
}

func (s *Server) GetTasksListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	login, err := s.getValidationResult(r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := s.task.ListTasks(login)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(t)
}

func (s *Server) RunTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	runnedTask := models.RunTask{}
	err = json.Unmarshal(data, &runnedTask)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.ErrInvalidJsonBody.Error(), http.StatusInternalServerError)
		return
	}

	login, err := s.getValidationResult(r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if login != runnedTask.InitiatorLogin {
		s.logger.Error().Msg(e.ErrTokenLoginNotEqualInitiatorLogin.Error())
		http.Error(w, e.ErrTokenLoginNotEqualInitiatorLogin.Error(), http.StatusForbidden)
		return
	}

	approvals := make([]*models.Approval, len(runnedTask.ApprovalLogins))
	for idx, al := range runnedTask.ApprovalLogins {
		approvals[idx] = &models.Approval{
			Approved:      sql.NullBool{Valid: false, Bool: false},
			Sent:          sql.NullBool{Valid: false, Bool: false},
			N:             idx + 1,
			ApprovalLogin: al,
		}
	}

	createdTask := models.Task{
		UUID:           uuid.GenUUID(),
		InitiatorLogin: runnedTask.InitiatorLogin,
		Status:         "created",
		Approvals:      approvals,
	}

	err = s.task.RunTask(&createdTask)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: отправляем письма по очереди approvalLogins

	json.NewEncoder(w).Encode(createdTask)
}
