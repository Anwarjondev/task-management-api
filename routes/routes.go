package routes

import (
	"net/http"

	"github.com/Anwarjondev/task-management-api/handlers"
	"github.com/Anwarjondev/task-management-api/middleware"
)


func SetUpRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /regitser", handlers.Register)
	mux.HandleFunc("POST /login", handlers.Login)


	protected := http.NewServeMux()
	protected.HandleFunc("POST /createproject", handlers.CreateTask)
	protected.HandleFunc("GET /getproject", handlers.GetProject)
	protected.HandleFunc("PUT /updateproject", handlers.UpdateProject)
	protected.HandleFunc("DELETE /deleteproject", handlers.DeleteProject)
	protected.HandleFunc("POST /projects/{id}/members", handlers.AddProjectMember)
	protected.HandleFunc("POST /createtask", handlers.CreateTask)
	protected.HandleFunc("GET /gettask", handlers.GetTask)
	protected.HandleFunc("PUT /updatetask", handlers.Updatetask)
	protected.HandleFunc("DELETE /deletetask", handlers.DeleteTask)
	protected.HandleFunc("PUT /updateuser", handlers.UpdateUser)

	admiMux := http.NewServeMux()
	admiMux.HandleFunc("GET /users", handlers.GetUsers)
	admiMux.HandleFunc("DELETE /deleteusers", handlers.DeleteUser)

	mux.Handle("/", middleware.AuthMiddleware(protected))
	mux.Handle("/admin/", middleware.AdminMiddleware(middleware.AuthMiddleware(admiMux)))
	return mux
}