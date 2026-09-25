package main

import (
	"net/http"
	"time"

	"api-gin/database"
	"api-gin/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	db := database.NewDatabase()

	salaHandler := &handlers.SalaHandler{Database: db}
	alunoHandler := &handlers.AlunoHandler{Database: db}
	turmaHandler := &handlers.TurmaHandler{Database: db}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		})

		v1.POST("/salas", salaHandler.CreateSala)
		v1.GET("/salas", salaHandler.ListarSalas)
		v1.GET("/salas/:id/grade", salaHandler.GradeSala)

		v1.POST("/alunos", alunoHandler.CreateAluno)
		v1.GET("/alunos", alunoHandler.ListarAlunos)
		v1.GET("/alunos/:id", alunoHandler.BuscarAluno)

		v1.POST("/turmas", turmaHandler.CreateTurma)
		v1.GET("/turmas", turmaHandler.ListarTurmas)
		v1.POST("/turmas/:id/alunos", turmaHandler.MatricularAluno)
		v1.GET("/turmas/:id/alunos", turmaHandler.ListarAlunosDaTurma)
		v1.POST("/turmas/:id/alocar", turmaHandler.AlocarSala)
	}

	r.Run(":8080")
}
