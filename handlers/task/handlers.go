package task

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("Welcome, %v", id)))
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "taskID")
	w.Write([]byte(fmt.Sprintf("Welcome, %v", id)))
}

func ApproveTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "taskID")
	login := chi.URLParam(r, "approvalLogin")
	w.Write([]byte(fmt.Sprintf("Approve, %v, %v", id, login)))
}

func DeclineTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "taskID")
	login := chi.URLParam(r, "approvalLogin")
	w.Write([]byte(fmt.Sprintf("Decline, %v, %v", id, login)))
}

func GetTasksList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get tasks"))
}

func RunTask(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	w.Write([]byte(fmt.Sprintf("Run task - %v", string(data))))
}
