package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/samgozman/go-bloggy/pkg/server"
	"net/http"
)

func (h *Handler) PostPosts(ctx echo.Context) error {
	var req server.PostRequest
	if err := ctx.Bind(&req); err != nil {
		var errorMessage string
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			errorMessage = fmt.Sprintf("%v", echoErr.Message)
		}

		return ctx.JSON(http.StatusBadRequest, server.RequestError{
			Code:    errRequestBodyBinding,
			Message: fmt.Sprintf("Error binding request body: %v", errorMessage),
		})
	}

	var externalUserID string
	if s := ctx.Get("externalUserID"); s != nil {
		externalUserID = s.(string)
	}

	if externalUserID == "" {
		return ctx.JSON(http.StatusUnauthorized, server.RequestError{
			Code:    errUnauthorized,
			Message: "Unauthorized",
		})
	}

	user, err := h.db.Models.Users.GetByExternalID(ctx.Request().Context(), externalUserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, server.RequestError{
			Code:    errGetUser,
			Message: "Post author is not found",
		})
	}

	var keywords string
	if req.Keywords != nil && len(*req.Keywords) > 0 {
		keywords = (*req.Keywords)[0]
		for i := 1; i < len(*req.Keywords); i++ {
			keywords += "," + (*req.Keywords)[i]
		}
	}

	post := models.Post{
		UserID:      user.ID,
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		Description: req.Description,
		Keywords:    keywords,
	}

	if err := h.db.Models.Posts.Create(ctx.Request().Context(), &post); err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicate):
			return ctx.JSON(http.StatusConflict, server.RequestError{
				Code:    errDuplicatePost,
				Message: "Post with this URL slug already exists",
			})
		case errors.Is(err, models.ErrValidationFailed):
			return ctx.JSON(http.StatusBadRequest, server.RequestError{
				Code:    errValidationFailed,
				Message: "Post validation failed",
			})

		default:
			return ctx.JSON(http.StatusInternalServerError, server.RequestError{
				Code:    errCreatePost,
				Message: "Error creating post",
			})
		}
	}

	return ctx.JSON(http.StatusCreated, server.PostResponse{
		Id:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Description: post.Description,
		Content:     post.Content,
		Keywords:    req.Keywords,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	})
}
