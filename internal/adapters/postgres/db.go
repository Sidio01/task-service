package postgres

import (
	"context"
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
	e "gitlab.com/g6834/team26/task/internal/domain/errors"
	"gitlab.com/g6834/team26/task/internal/domain/models"
)

type PostgresDatabase struct {
	psqlClient *sql.DB
}

func New(ctx context.Context, pgconn string) (*PostgresDatabase, error) {
	// _, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	db, err := sql.Open("postgres", pgconn+"?sslmode=disable")

	if err != nil {
		return nil, err
	}
	return &PostgresDatabase{psqlClient: db}, nil
}

func (pdb *PostgresDatabase) Stop(ctx context.Context) error {
	err := pdb.psqlClient.Close()
	if err != nil {
		return err
	}
	return nil
}

func (pdb *PostgresDatabase) List(ctx context.Context, login string) ([]*models.Task, error) {
	// _, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	var result []*models.Task

	taskQuery := `SELECT "uuid", "name", "text", "login", "status" FROM "tasks" WHERE "login" = $1`
	taskRows, err := pdb.psqlClient.Query(taskQuery, login)
	if err != nil {
		return nil, err
		// return nil, fmt.Errorf("no user with such login")
	}
	defer taskRows.Close()

	for taskRows.Next() {
		var task models.Task
		err := taskRows.Scan(&task.UUID, &task.Name, &task.Text, &task.InitiatorLogin, &task.Status)
		if err != nil {
			return nil, err
		}
		task.Name = strings.TrimSpace(task.Name)
		task.Text = strings.TrimSpace(task.Text)
		task.InitiatorLogin = strings.TrimSpace(task.InitiatorLogin)
		task.Status = strings.TrimSpace(task.Status)

		approvalQuery := `SELECT "approval_login", "approved", "sent", "n" FROM "approvals" WHERE "task_uuid" = $1`
		approvalRows, err := pdb.psqlClient.Query(approvalQuery, task.UUID)
		if err != nil {
			return nil, err
		}
		defer approvalRows.Close()

		for approvalRows.Next() {
			var approval models.Approval
			err := approvalRows.Scan(&approval.ApprovalLogin, &approval.Approved, &approval.Sent, &approval.N)
			if err != nil {
				return nil, err
			}
			approval.ApprovalLogin = strings.TrimSpace(approval.ApprovalLogin)
			task.Approvals = append(task.Approvals, &approval)
		}
		result = append(result, &task)
	}

	return result, nil
}

func (pdb *PostgresDatabase) Run(ctx context.Context, t *models.Task) error {
	// _, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	tx, err := pdb.psqlClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	taskQuery := `INSERT INTO "tasks" ("uuid", "name", "text", "login", "status") VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, taskQuery, t.UUID, t.Name, t.Text, t.InitiatorLogin, t.Status)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, approval := range t.Approvals {
		approvalsQuery := `INSERT INTO "approvals" ("task_uuid", "approval_login", "n") VALUES ($1, $2, $3)`
		_, err = tx.ExecContext(ctx, approvalsQuery, t.UUID, approval.ApprovalLogin, approval.N)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (pdb *PostgresDatabase) Update(ctx context.Context, id, login, name, text string) error {
	// _, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	var (
		query  string
		err    error
		result sql.Result
	)

	switch {
	case text == "" && name == "":
		log.Println("text == \"\" && name == \"\"")
		return e.ErrNothingToChange
	case text == "":
		log.Println("text == \"\"")
		query = `UPDATE "tasks" SET "name" = $1 WHERE "uuid" = $2 AND "login" = $3`
		result, err = pdb.psqlClient.Exec(query, name, id, login)
	case name == "":
		log.Println("name == \"\"")
		query = `UPDATE "tasks" SET "text" = $1 WHERE "uuid" = $2 AND "login" = $3`
		result, err = pdb.psqlClient.Exec(query, text, id, login)
	default:
		log.Println("default")
		query = `UPDATE "tasks" SET "name" = $1, "text" = $2 WHERE "uuid" = $3 AND "login" = $4`
		result, err = pdb.psqlClient.Exec(query, name, text, id, login)
	}
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return e.ErrNotFound
	}
	return nil
}

func (pdb *PostgresDatabase) Delete(ctx context.Context, login, id string) error {
	// _, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	query := `DELETE FROM "tasks" WHERE "uuid" = $1 AND "login" = $2`
	result, err := pdb.psqlClient.Exec(query, id, login)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return e.ErrNotFound
	}

	return nil
}

func (pdb *PostgresDatabase) Approve(ctx context.Context, login, id, approvalLogin string) error {
	// _, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	tx, err := pdb.psqlClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = pdb.checkTaskStatus(id)
	if err != nil {
		return err
	}

	err = pdb.checkApproval(id, approvalLogin)
	if err != nil {
		return err
	}

	query := `UPDATE "approvals" SET "approved" = $1 WHERE "task_uuid" = $2 AND "approval_login" = $3`
	result, err := tx.Exec(query, true, id, approvalLogin)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return e.ErrNotFound
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (pdb *PostgresDatabase) Decline(ctx context.Context, login, id, approvalLogin string) error {
	// _, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	tx, err := pdb.psqlClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = pdb.checkTaskStatus(id)
	if err != nil {
		return err
	}

	err = pdb.checkApproval(id, approvalLogin)
	if err != nil {
		return err
	}

	query := `UPDATE "approvals" SET "approved" = $1 WHERE "task_uuid" = $2 AND "approval_login" = $3`
	result, err := tx.Exec(query, false, id, approvalLogin)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return e.ErrNotFound
	}

	err = pdb.changeTaskStatus(tx, id, "declined")
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (pdb *PostgresDatabase) checkApproval(id, approvalLogin string) error {
	var approval models.Approval
	approvalQuery := `SELECT "approved" FROM "approvals" WHERE "task_uuid" = $1 AND "approval_login" = $2`
	approvalRow := pdb.psqlClient.QueryRow(approvalQuery, id, approvalLogin)
	approvalRow.Scan(&approval.Approved)
	if approval.Approved.Valid {
		return e.ErrApprovalHasBeenDone
	}
	return nil
}

func (pdb *PostgresDatabase) checkTaskStatus(id string) error {
	var task models.Task
	taskQuery := `SELECT "status" FROM "tasks" WHERE "uuid" = $1`
	taskRow := pdb.psqlClient.QueryRow(taskQuery, id)
	taskRow.Scan(&task.Status)
	if task.Status != "created" {
		return e.ErrTaskNotAvailableForApproval
	}
	return nil
}

func (pdb *PostgresDatabase) changeTaskStatus(tx *sql.Tx, id, status string) error {
	query := `UPDATE "tasks" SET "status" = $1 WHERE "uuid" = $2`
	_, err := tx.Exec(query, status, id)
	if err != nil {
		return err
	}
	return nil
}
