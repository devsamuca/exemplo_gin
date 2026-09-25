package database

import (
	"sync"

	"api-gin/models"
)

type Database struct {
	sync.Mutex

	Salas  map[string]*models.Sala
	Alunos map[string]*models.Aluno
	Turmas map[string]*models.Turma
}

func NewDatabase() *Database {
	return &Database{
		Salas:  make(map[string]*models.Sala),
		Alunos: make(map[string]*models.Aluno),
		Turmas: make(map[string]*models.Turma),
	}
}
