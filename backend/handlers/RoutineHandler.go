package handlers

import (
	"AppFitness/dto"
	"AppFitness/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoutineHandler struct {
	RoutineService services.RoutineService
}

func NewRoutineHandler(routineService services.RoutineService) *RoutineHandler {
	return &RoutineHandler{
		RoutineService: routineService,
	}
}

func (h *RoutineHandler) PostRoutine(c *gin.Context) {
	idUser, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var routine dto.RoutineRegisterDTO
	if err := c.ShouldBindJSON(&routine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	routine.CreatorUserID = idUser.(string)

	result, err := h.RoutineService.PostRoutine(&routine)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no puede estar vacío"),
			strings.Contains(msg, "no puede estar vacía"):
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) //400
			return

		case strings.Contains(msg, "dicho nombre de rutina ya existe"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) //409
			return

		case strings.Contains(msg, "no se pudo verificar si existe una rutina"),
			strings.Contains(msg, "error al crear la rutina"),
			strings.Contains(msg, "no se pudo obtener el ObjectID insertado"),
			strings.Contains(msg, "error al obtener la rutina creada"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al registrar la rutina"}) //500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) GetRoutines(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //401
		return
	}

	result, err := h.RoutineService.GetRoutines()
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no existen rutinas registradas"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) //404
			return

		case strings.Contains(msg, "error al obtener rutinas"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener rutinas"}) //500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) GetRoutineByID(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //401
		return
	}

	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "se requiere un ID de rutina para hacer la busqueda"})
		return
	}

	result, err := h.RoutineService.GetRoutineByID(id)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no existe ninguna rutina"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) //404
			return

		case strings.Contains(msg, "error al obtener la rutina por ID"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener la rutina"}) //500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) PutRoutine(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	idRoutine := c.Param("id")

	var routineModify dto.RoutineModifyDTO
	if err := c.ShouldBindJSON(&routineModify); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	routineModify.IDRoutine = idRoutine
	result, err := h.RoutineService.PutRoutine(routineModify)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no puede estar vacío"),
			strings.Contains(msg, "no puede ser igual al anterior"):
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) //400
			return

		case strings.Contains(msg, "no se modificó ninguna rutina"):
			c.JSON(http.StatusNotModified, gin.H{"error": msg}) //304
			return

		case strings.Contains(msg, "error al obtener la rutina a modificar"),
			strings.Contains(msg, "error al modificar la rutina"),
			strings.Contains(msg, "error al obtener la rutina modificada"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al modificar la rutina"}) //500
			return

		case strings.Contains(msg, "no existe ninguna rutina con ese ID"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) //404
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}
