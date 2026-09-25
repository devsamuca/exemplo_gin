package handlers

import (
	"net/http"
	"sort"

	"api-gin/database"
	"api-gin/models"

	"github.com/gin-gonic/gin"
)

type AlunoHandler struct {
	Database *database.Database
}

func (h *AlunoHandler) CreateAluno(c *gin.Context) {
	var aluno models.Aluno

	if err := c.ShouldBindJSON(&aluno); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados invalidos: " + err.Error()})
		return
	}

	h.Database.Lock()
	defer h.Database.Unlock()

	if _, existe := h.Database.Alunos[aluno.ID]; existe {
		c.JSON(http.StatusConflict, gin.H{"erro": "ja existe um aluno com essa matricula"})
		return
	}

	h.Database.Alunos[aluno.ID] = &aluno

	c.JSON(http.StatusCreated, aluno)
}

func (h *AlunoHandler) ListarAlunos(c *gin.Context) {
	h.Database.Lock()
	defer h.Database.Unlock()

	alunos := []models.Aluno{}
	for _, a := range h.Database.Alunos {
		alunos = append(alunos, *a)
	}

	sort.Slice(alunos, func(i, j int) bool {
		return alunos[i].ID < alunos[j].ID
	})

	c.JSON(http.StatusOK, alunos)
}

func (h *AlunoHandler) BuscarAluno(c *gin.Context) {
	id := c.Param("id")

	h.Database.Lock()
	defer h.Database.Unlock()

	aluno, existe := h.Database.Alunos[id]
	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"erro": "aluno nao encontrado"})
		return
	}

	c.JSON(http.StatusOK, aluno)
}
