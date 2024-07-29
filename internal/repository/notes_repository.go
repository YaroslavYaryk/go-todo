package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"simpleRestApi/internal/domain"
	"simpleRestApi/pkg/psql"
	"strings"
	"time"
)

type NotePostgres struct {
	db *sqlx.DB
}

func NewNotePostgres(db *sqlx.DB) *NotePostgres {
	return &NotePostgres{
		db: db,
	}
}

func (r *NotePostgres) Create(userId int, note domain.Note) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createNoteQuery := fmt.Sprintf("INSERT INTO %s (title, text, created_at, updated_at, tags, creator) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", psql.NoteTable)

	row := tx.QueryRow(createNoteQuery, note.Title, note.Text, time.Now(), time.Now(), note.Tags, userId)
	if err := row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createUserNoteQuery := fmt.Sprintf("INSERT INTO %s (note_id, user_id) VALUES ($1, $2)", psql.UserNoteTable)
	_, err = tx.Exec(createUserNoteQuery, id, userId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *NotePostgres) GetAll(userId int) ([]domain.Note, error) {
	var items []domain.Note

	query := fmt.Sprintf(`SELECT 
        ti.id, ti.title, ti.text, ti.created_at, ti.updated_at, 
        ti.tags , ti.creator
    FROM 
        %s ti
    JOIN 
        %s un ON un.note_id = ti.id
	WHERE
	    un.user_id = $1
	
	`, psql.NoteTable, psql.UserNoteTable)

	err := r.db.Select(&items, query, userId)

	return items, err

}

func (r *NotePostgres) GetById(userId int, noteId int) (domain.Note, error) {
	var note domain.Note

	query := fmt.Sprintf(`SELECT 
        ti.id, ti.title, ti.text, ti.created_at, ti.updated_at, 
        ti.tags, ti.creator
    FROM 
        %s ti
    JOIN 
        %s un ON un.note_id = ti.id
	WHERE
	    un.user_id = $1 and ti.id = $2
	
	`, psql.NoteTable, psql.UserNoteTable)
	err := r.db.Get(&note, query, userId, noteId)

	return note, err
}

func (r *NotePostgres) Update(noteId int, input domain.Note) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, input.Title)
		argId++
	}

	if input.Text != "" {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argId))
		args = append(args, input.Text)
		argId++
	}

	if input.Tags != "" {
		setValues = append(setValues, fmt.Sprintf("tags=$%d", argId))
		args = append(args, input.Tags)
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, time.Now())
	argId++

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, psql.NoteTable, setQuery, argId)
	args = append(args, noteId)

	_, err := r.db.Exec(query, args...)

	logrus.Debugf("updated todo list with id %d", noteId)

	return err
}

func (r *NotePostgres) Delete(noteId int, userId int) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var deletedNote sql.Result

	note, err := r.GetById(noteId, userId)
	if int(note.Creator) == userId {
		query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, psql.NoteTable)
		deletedNote, err = tx.Exec(query, noteId)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE note_id = $1`, psql.UserNoteTable)
	_, err = tx.Exec(query, userId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	affected_rows, err := deletedNote.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(affected_rows), nil
}

func (r *NotePostgres) ShareNote(userId int, noteId int, sharedUserId int) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	note, err := r.GetById(userId, noteId)
	if int(note.Creator) == userId {
		createUserNoteQuery := fmt.Sprintf("INSERT INTO %s (note_id, user_id) VALUES ($1, $2)", psql.UserNoteTable)
		_, err = tx.Exec(createUserNoteQuery, noteId, sharedUserId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	} else {
		_ = tx.Rollback()
		return errors.New("You dont have the permission to share this file")
	}

	return tx.Commit()
}
