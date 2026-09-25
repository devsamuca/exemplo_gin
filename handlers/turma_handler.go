package handlers

import (
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"api-gin/database"
	"api-gin/models"

	"github.com/gin-gonic/gin"
)

type TurmaHandler struct {
	Database *database.Database
}

func (h *TurmaHandler) CreateTurma(c *gin.Context) {
	var turma models.Turma

	if err := c.ShouldBindJSON(&turma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados invalidos: " + err.Error()})
		return
	}

	h.Database.Lock()
	defer h.Database.Unlock()

	if _, existe := h.Database.Turmas[turma.ID]; existe {
		c.JSON(http.StatusConflict, gin.H{"erro": "ja existe uma turma com esse id"})
		return
	}

	turma.Ativa = true
	turma.AlunosIDs = []string{}
	turma.Alocacao = nil

	h.Database.Turmas[turma.ID] = &turma

	c.JSON(http.StatusCreated, montarResposta(&turma))
}

func (h *TurmaHandler) ListarTurmas(c *gin.Context) {
	h.Database.Lock()
	defer h.Database.Unlock()

	turmas := []models.TurmaResposta{}
	for _, t := range h.Database.Turmas {
		turmas = append(turmas, montarResposta(t))
	}

	sort.Slice(turmas, func(i, j int) bool {
		return turmas[i].ID < turmas[j].ID
	})

	c.JSON(http.StatusOK, turmas)
}

func (h *TurmaHandler) MatricularAluno(c *gin.Context) {
	turmaID := c.Param("id")

	var req models.MatriculaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados invalidos: " + err.Error()})
		return
	}

	h.Database.Lock()
	defer h.Database.Unlock()

	turma, existe := h.Database.Turmas[turmaID]
	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"erro": "turma nao encontrada"})
		return
	}

	aluno, existe := h.Database.Alunos[req.AlunoID]
	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"erro": "aluno nao encontrado"})
		return
	}

	if slices.Contains(turma.AlunosIDs, aluno.ID) {
		c.JSON(http.StatusConflict, gin.H{"erro": "aluno ja esta matriculado nessa turma"})
		return
	}

	if turma.Alocacao != nil {
		sala := h.Database.Salas[turma.Alocacao.SalaID]
		if len(turma.AlunosIDs)+1 > sala.Capacidade {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": "capacidade da sala insuficiente para matricular mais um aluno"})
			return
		}

		for _, outra := range h.Database.Turmas {
			if outra.ID == turma.ID || outra.Alocacao == nil {
				continue
			}

			mesmoDia := outra.Alocacao.DiaSemana == turma.Alocacao.DiaSemana
			sobrepoe := turma.Alocacao.HoraInicio < outra.Alocacao.HoraFim && turma.Alocacao.HoraFim > outra.Alocacao.HoraInicio

			if slices.Contains(outra.AlunosIDs, aluno.ID) && mesmoDia && sobrepoe {
				c.JSON(http.StatusConflict, gin.H{"erro": "conflito de agenda: aluno ja esta na turma " + outra.ID + " nesse horario"})
				return
			}
		}
	}

	turma.AlunosIDs = append(turma.AlunosIDs, aluno.ID)

	c.JSON(http.StatusCreated, gin.H{
		"mensagem":          "aluno matriculado com sucesso",
		"turma_id":          turma.ID,
		"aluno_id":          aluno.ID,
		"quantidade_alunos": len(turma.AlunosIDs),
	})
}

func (h *TurmaHandler) ListarAlunosDaTurma(c *gin.Context) {
	turmaID := c.Param("id")

	h.Database.Lock()
	defer h.Database.Unlock()

	turma, existe := h.Database.Turmas[turmaID]
	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"erro": "turma nao encontrada"})
		return
	}

	alunos := []models.Aluno{}
	for _, id := range turma.AlunosIDs {
		alunos = append(alunos, *h.Database.Alunos[id])
	}

	c.JSON(http.StatusOK, alunos)
}

func (h *TurmaHandler) AlocarSala(c *gin.Context) {
	turmaID := c.Param("id")

	var req models.AlocacaoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados invalidos: " + err.Error()})
		return
	}

	req.DiaSemana = strings.ToLower(req.DiaSemana)
	diasValidos := []string{"segunda", "terca", "quarta", "quinta", "sexta", "sabado", "domingo"}
	if !slices.Contains(diasValidos, req.DiaSemana) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dia da semana invalido, use: segunda, terca, quarta, quinta, sexta, sabado ou domingo"})
		return
	}

	inicio, err := time.Parse("15:04", req.HoraInicio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "horario_inicio invalido, use o formato HH:MM"})
		return
	}
	fim, err := time.Parse("15:04", req.HoraFim)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "horario_fim invalido, use o formato HH:MM"})
		return
	}
	if !inicio.Before(fim) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "horario_inicio deve ser menor que horario_fim"})
		return
	}

	h.Database.Lock()
	defer h.Database.Unlock()

	turma, existe := h.Database.Turmas[turmaID]
	if !existe || !turma.Ativa {
		c.JSON(http.StatusNotFound, gin.H{"erro": "turma nao encontrada"})
		return
	}

	sala, existe := h.Database.Salas[req.SalaID]
	if !existe || !sala.Ativa {
		c.JSON(http.StatusNotFound, gin.H{"erro": "sala nao encontrada"})
		return
	}

	if len(turma.AlunosIDs) > sala.Capacidade {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": "capacidade da sala menor que a quantidade de alunos da turma"})
		return
	}

	nova := models.Alocacao{
		TurmaID:    turma.ID,
		SalaID:     sala.ID,
		DiaSemana:  req.DiaSemana,
		HoraInicio: inicio.Format("15:04"),
		HoraFim:    fim.Format("15:04"),
	}

	for _, outra := range h.Database.Turmas {
		if outra.ID == turma.ID || outra.Alocacao == nil {
			continue
		}

		mesmoDia := outra.Alocacao.DiaSemana == nova.DiaSemana
		sobrepoe := nova.HoraInicio < outra.Alocacao.HoraFim && nova.HoraFim > outra.Alocacao.HoraInicio

		if !mesmoDia || !sobrepoe {
			continue
		}

		if outra.Alocacao.SalaID == sala.ID {
			c.JSON(http.StatusConflict, gin.H{"erro": "conflito de agenda: sala ja ocupada pela turma " + outra.ID + " nesse horario"})
			return
		}

		for _, alunoID := range turma.AlunosIDs {
			if slices.Contains(outra.AlunosIDs, alunoID) {
				c.JSON(http.StatusConflict, gin.H{"erro": "conflito de agenda: aluno " + alunoID + " ja esta na turma " + outra.ID + " nesse horario"})
				return
			}
		}
	}

	turma.Alocacao = &nova

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "turma alocada com sucesso",
		"alocacao": nova,
	})
}

func montarResposta(t *models.Turma) models.TurmaResposta {
	return models.TurmaResposta{
		ID:               t.ID,
		Nome:             t.Nome,
		Disciplina:       t.Disciplina,
		Professor:        t.Professor,
		Ativa:            t.Ativa,
		QuantidadeAlunos: len(t.AlunosIDs),
		Alocada:          t.Alocacao != nil,
		Alocacao:         t.Alocacao,
	}
}
