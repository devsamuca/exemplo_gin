package models

type Sala struct {
	ID         string   `json:"id" binding:"required"`
	Nome       string   `json:"nome" binding:"required"`
	Capacidade int      `json:"capacidade" binding:"required,gt=0"`
	Recursos   []string `json:"recursos"`
	Ativa      bool     `json:"ativa"`
}

type Aluno struct {
	ID    string `json:"id" binding:"required"`
	Nome  string `json:"nome" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type Turma struct {
	ID         string    `json:"id" binding:"required"`
	Nome       string    `json:"nome" binding:"required"`
	Disciplina string    `json:"disciplina" binding:"required"`
	Professor  string    `json:"professor" binding:"required"`
	Ativa      bool      `json:"ativa"`
	AlunosIDs  []string  `json:"-"`
	Alocacao   *Alocacao `json:"-"`
}

type Alocacao struct {
	TurmaID    string `json:"turma_id"`
	SalaID     string `json:"sala_id"`
	DiaSemana  string `json:"dia_semana"`
	HoraInicio string `json:"horario_inicio"`
	HoraFim    string `json:"horario_fim"`
}

type MatriculaRequest struct {
	AlunoID string `json:"aluno_id" binding:"required"`
}

type AlocacaoRequest struct {
	SalaID     string `json:"sala_id" binding:"required"`
	DiaSemana  string `json:"dia_semana" binding:"required"`
	HoraInicio string `json:"horario_inicio" binding:"required"`
	HoraFim    string `json:"horario_fim" binding:"required"`
}

type TurmaResposta struct {
	ID               string    `json:"id"`
	Nome             string    `json:"nome"`
	Disciplina       string    `json:"disciplina"`
	Professor        string    `json:"professor"`
	Ativa            bool      `json:"ativa"`
	QuantidadeAlunos int       `json:"quantidade_alunos"`
	Alocada          bool      `json:"alocada"`
	Alocacao         *Alocacao `json:"alocacao"`
}
