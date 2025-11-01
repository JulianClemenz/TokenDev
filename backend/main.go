package main

import (
	"AppFitness/handlers"
	"AppFitness/middleware"
	"AppFitness/repositories"
	"AppFitness/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Iniciando AppFitness...")

	// 1. Conexión a la Base de Datos
	db := repositories.NewMongoDB()

	defer func() {
		log.Println("Cerrando conexion con MongoDB...")
		if err := db.Disconnect(); err != nil {
			log.Fatalf("Error al desconectar MongoDB: %v", err)
		}
	}()
	log.Println("Conectado a MongoDB exitosamente.")

	// 2. Dependencias

	// --- Repositorios ---
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	exerciseRepo := repositories.NewExcerciseRepository(db)
	routineRepo := repositories.NewRoutineRepository(db)
	workoutRepo := repositories.NewWorkoutRepository(db)

	// --- Servicios ---
	authService := services.NewAuthService(userRepo, sessionRepo)
	userService := services.NewUserService(userRepo)
	exerciseService := services.NewExcerciseService(exerciseRepo)
	routineService := services.NewRoutineService(routineRepo, exerciseRepo)
	workoutService := services.NewWorkoutService(workoutRepo, routineRepo, userRepo)

	// --- Handlers ---
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService)
	routineHandler := handlers.NewRoutineHandler(routineService)
	workoutHandler := handlers.NewWorkoutHadler(workoutService)

	router := gin.Default()

	// Configurar archivos státic y templates
	//router.Static("/static", "./static") // Para CSS, JS, imágenes
	//router.LoadHTMLGlob("templates/*")   // Para renderizado HTML del lado del servidor

	// Rutas Públicas (Autenticación y Registro)
	// no requieren el middleware de autenticación
	router.POST("/register", userHandler.PostUser)
	router.POST("/login", authHandler.PostLogin)
	router.POST("/logout", authHandler.PostLogout)
	router.POST("/refresh", authHandler.PostRefresh)

	// rutas que SÍ requieren autenticación
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	//  Rutas de Perfil de Usuario
	userRoutes := api.Group("/users")
	{
		// GET /api/users/:id - Ver perfil de usuario
		userRoutes.GET("/:id", userHandler.GetUserByID)
		// PUT /api/users/:id - Modificar perfil
		userRoutes.PUT("/:id", userHandler.PutUser)
		// POST /api/users/:id/password - Cambiar contraseña
		userRoutes.POST("/:id/password", userHandler.PasswordModify)
	}
	//Rutas de Ejercicios
	exerciseRoutes := api.Group("/exercises")
	{
		// Rutas públicas (GET)
		exerciseRoutes.GET("/", exerciseHandler.GetExcercises)
		exerciseRoutes.GET("/filter", exerciseHandler.GetByFilters) // Búsqueda y filtros
		exerciseRoutes.GET("/:id", exerciseHandler.GetExcerciseByID)

		// Rutas de Administrador (POST, PUT, DELETE)
		// Aplicamos CheckAdmin solo a estas rutas
		adminExercise := exerciseRoutes.Group("/")
		adminExercise.Use(middleware.CheckAdmin())
		{
			adminExercise.POST("/", exerciseHandler.PostExcercise)  // Alta
			adminExercise.PUT("/:id", exerciseHandler.PutExcercise) // Edición
			// TODO: Necesitas implementar el handler para DeleteExcercise
			// adminExercise.DELETE("/:id", exerciseHandler.DeleteExcercise)//////////////ATENCION
		}
	}

	// Rutas de Rutinas
	// Estas son para usuarios normales
	routineRoutes := api.Group("/routines")
	routineRoutes.Use(middleware.CheckUser()) // Aseguramos que solo los clientes puedan gestionar rutinas
	{
		routineRoutes.POST("/", routineHandler.PostRoutine)
		routineRoutes.GET("/", routineHandler.GetRoutines)
		routineRoutes.GET("/:id", routineHandler.GetRoutineByID)
		routineRoutes.PUT("/:id", routineHandler.PutRoutine)
		routineRoutes.DELETE("/:id", routineHandler.DeleteRoutine)

		// Manejo de ejercicios dentro de una rutina
		routineRoutes.POST("/:id/exercises", routineHandler.AddExcerciseToRoutine)
		routineRoutes.PUT("/:routine_id/exercises/:exercise_id", routineHandler.UpdateExerciseInRoutine)
		// handler espera un DTO en el body, así que no usamos params
		routineRoutes.DELETE("/exercises", routineHandler.RemoveExerciseFromRoutine)
	}

	// Rutas de Seguimiento (Workouts)
	workoutRoutes := api.Group("/workouts")
	workoutRoutes.Use(middleware.CheckUser()) // Solo para clientes
	{
		// GET /api/workouts/ - Ver historial de entrenamientos
		workoutRoutes.GET("/", workoutHandler.GetWorkouts)

		// POST /api/workouts/:id_routine - Registrar un entrenamiento completado
		workoutRoutes.POST("/:id_routine", workoutHandler.PostWorkout)

		// GET /api/workouts/stats - Ver estadísticas personales
		workoutRoutes.GET("/stats", workoutHandler.GetWorkoutStats)

		// TODO: Tu handler 'GetWorkoutByID' usa el ID del token en lugar del ID del workout.
		// Deberías corregir el handler para que use c.Param("id")
		workoutRoutes.GET("/:id", workoutHandler.GetWorkoutByID) // Ver un workout específico

		// DELETE /api/workouts/:id - Eliminar un registro de workout
		workoutRoutes.DELETE("/:id", workoutHandler.DeleteWorkout)
	}

	// --- Rutas del Panel de Administración ---
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middleware.CheckAdmin()) // Protegido solo para Admins
	{
		adminRoutes.GET("/users", userHandler.GetUsers) // Gestión de usuarios
		// TODO: Implementar handlers para estadísticas globales
		//adminRoutes.GET("/stats/users", ...)
		//adminRoutes.GET("/stats/exercises", ...) ////////////////////////////////////ATENCION
	}

	// 5. Iniciar Servidor
	log.Println("Servidor escuchando en http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor Gin: %v", err)
	}
}
