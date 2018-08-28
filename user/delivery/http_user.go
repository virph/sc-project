package delivery

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/virph/sc-project/models"
	"github.com/virph/sc-project/user"
)

type userHandler struct {
	userUsecase user.UserUsecase
}

type UserTemplateData struct {
	SearchTerm string
	Users      []models.User
}

func (h *userHandler) handleFindUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
	}
	searchTerm := r.FormValue("user-name")

	users := h.userUsecase.Find(searchTerm)

	tmpl, err := template.ParseFiles("user/delivery/template_user.html")
	if err != nil {
		log.Println(err)
	}
	tmpl.Execute(w, UserTemplateData{
		SearchTerm: searchTerm,
		Users:      users,
	})
}

func NewUserHandler(userUsecase *user.UserUsecase) {
	handler := userHandler{
		userUsecase: *userUsecase,
	}

	http.HandleFunc("/", handler.handleFindUser)
}
