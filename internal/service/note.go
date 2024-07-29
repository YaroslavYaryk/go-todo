package service

import (
	"simpleRestApi/internal/domain"
	"simpleRestApi/internal/repository"
)

type NoteService struct {
	repo repository.Note
}

func NewNoteService(repo repository.Note) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) Create(userId int, note domain.Note) (int, error) {
	return s.repo.Create(userId, note)
}

func (s *NoteService) GetAll(userId int) ([]domain.Note, error) {
	return s.repo.GetAll(userId)
}

func (s *NoteService) GetById(userId int, NoteId int) (domain.Note, error) {
	return s.repo.GetById(userId, NoteId)
}

func (s *NoteService) Update(noteId int, input domain.Note) error {
	return s.repo.Update(noteId, input)
}

func (s *NoteService) Delete(noteId int, userId int) (int, error) {
	return s.repo.Delete(noteId, userId)
}

func (s *NoteService) ShareNote(userId int, noteId int, sharedUserId int) error {
	return s.repo.ShareNote(userId, noteId, sharedUserId)
}
