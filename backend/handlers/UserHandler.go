package handlers

import (
	"AppFitness/dto"
	"AppFitness/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserInterface
}

func NewUserHandler(userService services.UserInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	collection, err := h.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"mensaje": "error al obtener usuarios"})
		return
	}

	c.JSON(http.StatusOK, collection)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	id := c.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no se encontro") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener cliente"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) PostUser(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //esto ns si esta bien aca
		return
	}

	var user dto.UserRegisterDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := h.userService.PostUser(&user) // dependiendo el error lanzamos un status
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "ya existe"): // email o username duplicados / eror 409
			c.JSON(http.StatusConflict, gin.H{"error": msg})
			return
		case strings.Contains(msg, "no se pudo verificar"),
			strings.Contains(msg, "error al insertar usuario"),
			strings.Contains(msg, "hashear contraseña"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al registrar usuario"}) //error 500
			return
		default:
			//  validaciones de negocio
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) // error 400
			return
		}
	}

	c.JSON(http.StatusCreated, resultado)
}
