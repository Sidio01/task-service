package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/g6834/team26/task/handlers/task"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/task/v1", func(r chi.Router) {
		r.Get("/{id}", task.Welcome)
		r.Delete("/tasks/{taskID}", task.DeleteTask)
		r.Post("/tasks/{taskID}/approve/{approvalLogin}", task.ApproveTask)
		r.Post("/tasks/{taskID}/decline/{approvalLogin}", task.DeclineTask)
		r.Get("/tasks/", task.GetTasksList)
		r.Post("/tasks/run", task.RunTask)
	})
	http.ListenAndServe(":3000", r)
}
