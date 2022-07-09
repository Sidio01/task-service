package postgres

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	e "gitlab.com/g6834/team26/task/internal/domain/errors"
	"gitlab.com/g6834/team26/task/internal/domain/models"
)

type PostgresDatabase struct {
	psqlClient *sql.DB
}

func New(ctx context.Context, pgconn string) (*PostgresDatabase, error) {
	// _, cancel := context.WithTimeout(ctx, 5 * time.Second)
	// defer cancel()
	db, err := sql.Open("postgres", pgconn+"?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &PostgresDatabase{psqlClient: db}, nil
}

func (pdb *PostgresDatabase) List(login string) ([]*models.Task, error) {
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

func (pdb *PostgresDatabase) Run(t *models.Task) error {
	ctx := context.Background()
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

func (pdb *PostgresDatabase) Delete(login, id string) error { // TODO: логин можно не использовать
	query := `DELETE FROM "tasks" WHERE "uuid" = $1 AND "login" = $2` // TODO: отправлять письма всем участникам об отмене операции
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

func (pdb *PostgresDatabase) Approve(login, id, approvalLogin string) error { // TODO: логин не используется
	query := `UPDATE "approvals" SET "approved" = $1 WHERE "task_uuid" = $2 AND "approval_login" = $3`
	result, err := pdb.psqlClient.Exec(query, true, id, approvalLogin)
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

func (pdb *PostgresDatabase) Decline(login, id, approvalLogin string) error { // TODO: логин не используется
	query := `UPDATE "approvals" SET "approved" = $1 WHERE "task_uuid" = $2 AND "approval_login" = $3`
	result, err := pdb.psqlClient.Exec(query, false, id, approvalLogin)
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
