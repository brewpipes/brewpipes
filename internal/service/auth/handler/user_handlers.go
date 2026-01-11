package handler

import (
	"context"
	"net/http"

	"github.com/brewpipes/brewpipesproto/internal/service"
	"github.com/brewpipes/brewpipesproto/internal/service/auth/handler/dto"
	"github.com/brewpipes/brewpipesproto/internal/service/auth/storage"
)

type UserLister interface {
	ListUsers(ctx context.Context) ([]storage.User, error)
}

// HandleListUsers handles [GET /users]
func HandleListUsers(db UserLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := db.ListUsers(r.Context())
		if err != nil {
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewUsersResponse(users))
	}
}
