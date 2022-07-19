package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	e "gitlab.com/g6834/team26/task/internal/domain/errors"
	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/pkg/api"
	"gitlab.com/g6834/team26/task/pkg/uuid"
)

// @title Сервис создания и согласования задач
// @version 1.0
// @description Сервис для создания и согласования задач и последующей отправкой писем последовательно всем участвующим лицам.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /task/v1

// @securityDefinitions.apikey access_token
// @in header
// @name access_token
// @securityDefinitions.apikey refresh_token
// @in header
// @name refresh_token

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

func (s *Server) setCookie(w http.ResponseWriter, c models.Cookie) {
	cookie := http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     "/",
		HttpOnly: true,
		Expires:  c.Expiration,
	}

	http.SetCookie(w, &cookie)
}

func (s *Server) updateCookies(w http.ResponseWriter, g *api.AuthResponse) {
	s.setCookie(w, models.Cookie{
		// Name:       s.config.Server.AccessCookie,
		Name:  "access_token",
		Value: g.AccessToken.GetValue(),
		// Expiration: time.Now().Add(time.Minute),
		Expiration: time.Unix(g.AccessToken.GetExpires(), 0),
	})
	s.setCookie(w, models.Cookie{
		// Name:       s.config.Server.RefreshCookie,
		Name:  "refresh_token",
		Value: g.RefreshToken.GetValue(),
		// Expiration: time.Now().Add(time.Hour),
		Expiration: time.Unix(g.RefreshToken.GetExpires(), 0),
	})
}

func (s *Server) getValidationResult(ctx context.Context, w http.ResponseWriter, r *http.Request) (string, error) {
	refreshToken, accessToken := s.getCookies(r)
	grpcResponse, err := s.task.Validate(ctx, refreshToken, accessToken)
	// log.Println("grpcResponse, err -", grpcResponse, err)
	if err != nil {
		return "", err
		// return err
	}
	// log.Printf("grpc result: %v, grpc login: %v", authResponseResult, authResponseLogin)
	// log.Println(grpcResponse)
	if !grpcResponse.Result {
		return "", e.ErrAuthFailed
	}

	if grpcResponse.RefreshToken != nil && grpcResponse.AccessToken != nil {
		s.updateCookies(w, grpcResponse)
	}

	return grpcResponse.Login, nil
}

// Run Task
// @ID RunTask
// @Security access_token
// @Security refresh_token
// @tags Работа с сервисом создания и согласования задач
// @Summary Создание задачи согласования
// @Description Создание задачи согласования с последующей отправкой
// @Param RunTask body models.RunTask true "Run Task"
// @Success 200 {object} models.Task
// @Failure 400 {object} e.ErrApiBadRequest
// @Failure 403 {object} e.ErrApiAuthFailed
// @Failure 500 {object} e.ErrApiInternalServerError
// @Router /tasks/run [post]
func (s *Server) RunTaskHandler(w http.ResponseWriter, r *http.Request) { // TODO: добавить получение из боди названия и текста задачи
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	runnedTask := models.RunTask{}
	err = json.Unmarshal(data, &runnedTask)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.ErrInvalidJsonBody.Error(), http.StatusBadRequest)
		return
	}

	login, err := s.getValidationResult(ctx, w, r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusInternalServerError)
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

	err = s.task.RunTask(ctx, &createdTask)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(createdTask)
}

// Get Tasks List
// @ID GetTasksList
// @Security access_token
// @Security refresh_token
// @tags Работа с сервисом создания и согласования задач
// @Summary Получение списка задач
// @Description Получения списка задач пользователя
// @Success 200 {object} models.Task
// @Failure 400 {object} e.ErrApiBadRequest
// @Failure 403 {object} e.ErrApiAuthFailed
// @Failure 500 {object} e.ErrApiInternalServerError
// @Router /tasks/ [get]
func (s *Server) GetTasksListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	login, err := s.getValidationResult(ctx, w, r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusInternalServerError)
		return
	}

	t, err := s.task.ListTasks(ctx, login)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(t)
}

// TODO: добавить эндпойнт на изменение задачи

// Approve Task
// @ID ApproveTask
// @Security access_token
// @Security refresh_token
// @tags Работа с сервисом создания и согласования задач
// @Summary Согласование задачи
// @Description Согласование задачи. В результате очередь согласования перейдет к следующему в списке согласующих, либо, в случае последнего этапа согласования, задача будет считаться выполненной.
// @Param taskID path string required "Task ID" Format(uuid)
// @Param approvalLogin path string required "Approval Login"
// @Success 200 {object} models.StatusApproved
// @Failure 400 {object} e.ErrApiBadRequest
// @Failure 403 {object} e.ErrApiAuthFailed
// @Failure 404 {object} e.ErrApiNotFound
// @Failure 500 {object} e.ErrApiInternalServerError
// @Router /tasks/{taskID}/approve/{approvalLogin} [post]
func (s *Server) ApproveTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	login, err := s.getValidationResult(ctx, w, r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "taskID")
	approvalLogin := chi.URLParam(r, "approvalLogin")
	err = s.task.ApproveTask(ctx, login, id, approvalLogin)
	if errors.Is(err, e.ErrNotFound) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusBadRequest)
		return
	}
	// w.Write([]byte("{\"status\": \"approved\"}"))
	json.NewEncoder(w).Encode(models.StatusApproved{Status: "approved"})
}

// Decline Task
// @ID DeclineTask
// @Security access_token
// @Security refresh_token
// @tags Работа с сервисом создания и согласования задач
// @Summary Отклонение задачи
// @Description Отклонение согласования задачи. В этом случае всем участникам поступит письмо с завершением задачи.
// @Param taskID path string required "Task ID" Format(uuid)
// @Param approvalLogin path string required "Approval Login"
// @Success 200 {object} models.StatusDeclined
// @Failure 400 {object} e.ErrApiBadRequest
// @Failure 403 {object} e.ErrApiAuthFailed
// @Failure 404 {object} e.ErrApiNotFound
// @Failure 500 {object} e.ErrApiInternalServerError
// @Router /tasks/{taskID}/decline/{approvalLogin} [post]
func (s *Server) DeclineTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	login, err := s.getValidationResult(ctx, w, r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "taskID")
	approvalLogin := chi.URLParam(r, "approvalLogin")
	err = s.task.DeclineTask(ctx, login, id, approvalLogin)
	if errors.Is(err, e.ErrNotFound) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusBadRequest)
		return
	}
	// w.Write([]byte("{\"status\": \"declined\"}"))
	json.NewEncoder(w).Encode(models.StatusDeclined{Status: "declined"})
}

// Delete Task
// @ID DeleteTask
// @Security access_token
// @Security refresh_token
// @tags Работа с сервисом создания и согласования задач
// @Summary Удаление созданной задачи
// @Description Удаление созданной задачи (доступно для автора задачи)
// @Param taskID path string required "Task ID" Format(uuid)
// @Success 200 {object} models.StatusDeleted
// @Failure 400 {object} e.ErrApiBadRequest
// @Failure 403 {object} e.ErrApiAuthFailed
// @Failure 404 {object} e.ErrApiNotFound
// @Failure 500 {object} e.ErrApiInternalServerError
// @Router /tasks/{taskID} [delete]
func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	login, err := s.getValidationResult(ctx, w, r)
	if errors.Is(err, e.ErrAuthFailed) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusForbidden)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusInternalServerError)
		return
	}

	// login := "test123"
	id := chi.URLParam(r, "taskID")
	err = s.task.DeleteTask(ctx, login, id)
	if errors.Is(err, e.ErrNotFound) {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		s.logger.Error().Msg(err.Error())
		http.Error(w, e.JsonErrWrapper{E: err.Error()}.Error(), http.StatusBadRequest)
		return
	}
	// w.Write([]byte("{\"status\": \"deleted\"}"))
	json.NewEncoder(w).Encode(models.StatusDeleted{Status: "deleted"})
}
