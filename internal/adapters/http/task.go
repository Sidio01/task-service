package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/g6834/team26/task/internal/domain/errors"
	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/pkg/uuid"
)

func (s *Server) authHandlers() http.Handler {
	r := chi.NewRouter()
	r.Delete("/tasks/{taskID}", s.DeleteTaskHandler)
	r.Post("/tasks/{taskID}/approve/{approvalLogin}", s.ApproveTaskHandler)
	r.Post("/tasks/{taskID}/decline/{approvalLogin}", s.DeclineTaskHandler)
	r.Get("/tasks/", s.GetTasksListHandler)
	r.Post("/tasks/run", s.RunTaskHandler)
	return r
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: провалидировать куки
	// TODO: если кука валидна удалить задачу
	id := chi.URLParam(r, "taskID")
	err := s.task.DeleteTask(id)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("{\"status\": \"deleted\"}"))
}

func (s *Server) ApproveTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: провалидировать куки
	// TODO: если кука валидна выполнить задачу
	id := chi.URLParam(r, "taskID")
	login := chi.URLParam(r, "approvalLogin")
	err := s.task.ApproveTask(id, login)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("{\"status\": \"approved\"}"))
}

func (s *Server) DeclineTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: провалидировать куки
	// TODO: если кука валидна выполнить задачу
	id := chi.URLParam(r, "taskID")
	login := chi.URLParam(r, "approvalLogin")
	err := s.task.DeclineTask(id, login)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("{\"status\": \"declined\"}"))
}

func (s *Server) GetTasksListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: провалидировать куки
	// TODO: отдать только список задач пользователя
	t, err := s.task.ListTasks()
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
		http.Error(w, errors.ErrInvalidJsonBody.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: провалидировать куки
	// TODO: сверить логин из боди с логином из куки
	// runnedTask.InitiatorLogin == login
	// w.WriteHeader(http.StatusForbidden)
	// w.Write([]byte("{\"error\": \"invalid login\"}"))
	// return

	approvals := make([]*models.Approval, len(runnedTask.ApprovalLogins))
	for idx, al := range runnedTask.ApprovalLogins {
		approvals[idx] = &models.Approval{
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
