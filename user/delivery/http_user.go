package delivery

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/virph/sc-project/models"
	"github.com/virph/sc-project/user"
	"github.com/virph/sc-project/visitorCount"
)

type userHandler struct {
	userUsecase         user.UserUsecase
	visitorCountUsecase visitorCount.VisitorCountUsecase
}

type UserTemplateData struct {
	SearchTerm   string
	Users        []models.User
	VisitorCount int
}

func (h *userHandler) handleFindUser(w http.ResponseWriter, r *http.Request) {
	h.visitorCountUsecase.PublishIncrease()

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
		SearchTerm:   searchTerm,
		Users:        users,
		VisitorCount: h.visitorCountUsecase.Get(),
	})
}

func NewUserHandler(u *user.UserUsecase, v *visitorCount.VisitorCountUsecase) {
	handler := userHandler{
		userUsecase:         *u,
		visitorCountUsecase: *v,
	}

	http.HandleFunc("/user", handler.handleFindUser)
}
