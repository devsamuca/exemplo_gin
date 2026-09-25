package handlers

import (
	"net/http"
	"slices"
	"sort"

	"api-gin/database"
	"api-gin/models"

	"github.com/gin-gonic/gin"
)

type SalaHandler struct {
	Database *database.Database
}

func (h *SalaHandler) CreateSala(c *gin.Context) {
	var sala models.Sala

	if err := c.ShouldBindJSON(&sala); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados invalidos: " + err.Error()})
		return
	}

	h.Database.Lock()
	defer h.Database.Unlock()

	if _, existe := h.Database.Salas[sala.ID]; existe {
		c.JSON(http.StatusConflict, gin.H{"erro": "ja existe uma sala com esse id"})
		return
	}

	if sala.Recursos == nil {
		sala.Recursos = []string{}
	}
	sala.Ativa = true

	h.Database.Salas[sala.ID] = &sala

	c.JSON(http.StatusCreated, sala)
}

func (h *SalaHandler) ListarSalas(c *gin.Context) {
	h.Database.Lock()
	defer h.Database.Unlock()

	salas := []models.Sala{}
	for _, s := range h.Database.Salas {
		salas = append(salas, *s)
	}

	sort.Slice(salas, func(i, j int) bool {
		return salas[i].ID < salas[j].ID
	})

	c.JSON(http.StatusOK, salas)
}

func (h *SalaHandler) GradeSala(c *gin.Context) {
	id := c.Param("id")
	dia := c.Query("dia")

	h.Database.Lock()
	defer h.Database.Unlock()

	sala, existe := h.Database.Salas[id]
	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"erro": "sala nao encontrada"})
		return
	}

	diasValidos := []string{"segunda", "terca", "quarta", "quinta", "sexta", "sabado", "domingo"}
	if dia != "" && !slices.Contains(diasValidos, dia) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dia da semana invalido"})
		return
	}

	grade := []models.Alocacao{}
	for _, turma := range h.Database.Turmas {
		if turma.Alocacao == nil || turma.Alocacao.SalaID != sala.ID {
			continue
		}
		if dia != "" && turma.Alocacao.DiaSemana != dia {
			continue
		}
		grade = append(grade, *turma.Alocacao)
	}

	sort.Slice(grade, func(i, j int) bool {
		if grade[i].DiaSemana != grade[j].DiaSemana {
			return grade[i].DiaSemana < grade[j].DiaSemana
		}
		return grade[i].HoraInicio < grade[j].HoraInicio
	})

	c.JSON(http.StatusOK, gin.H{
		"sala":  sala,
		"grade": grade,
	})
}
