package handler

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"net/http"
	"regexp"
)

func (h *Handler) PostSubscribers(ctx echo.Context) error {
	var req api.CreateSubscriberRequest
	if err := ctx.Bind(&req); err != nil {
		var errorMessage string
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			errorMessage = fmt.Sprintf("%v", echoErr.Message)
		}

		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errRequestBodyBinding,
			Message: fmt.Sprintf("Error binding request body: %v", errorMessage),
		})
	}

	if hr := h.hcaptchaService.VerifyToken(req.Captcha); !hr.Success {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errValidationCaptcha,
			Message: "Invalid captcha",
		})
	}

	// validate email
	if !isValidEmail(req.Email) {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errValidationEmail,
			Message: fmt.Sprintf("Invalid email: %v", req.Email),
		})
	}

	subscription := models.Subscriber{
		Email: req.Email,
	}

	if err := h.db.Models().Subscribers().Create(ctx.Request().Context(), &subscription); err != nil {
		// Note: we shouldn't tell duplicate error to the user for security reasons
		if !errors.Is(err, models.ErrDuplicate) {
			return ctx.JSON(http.StatusInternalServerError, api.RequestError{
				Code:    errCreateSubscription,
				Message: "Error creating subscription",
			})
		}
	}

	// Note: for confirmation code can be used internal ID of the subscription just for simplicity
	err := h.mailerService.SendConfirmationEmail(req.Email, subscription.ID.String())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.RequestError{
			Code:    errSendConfirmationEmail,
			Message: "Error sending confirmation email",
		})
	}

	return ctx.NoContent(http.StatusCreated)
}

func (h *Handler) DeleteSubscribers(ctx echo.Context) error {
	var req api.UnsubscribeRequest
	if err := ctx.Bind(&req); err != nil {
		var errorMessage string
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			errorMessage = fmt.Sprintf("%v", echoErr.Message)
		}

		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errRequestBodyBinding,
			Message: fmt.Sprintf("Error binding request body: %v", errorMessage),
		})
	}

	subscriptionID, err := h.db.Models().Subscribers().GetByID(ctx.Request().Context(), req.SubscriptionId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errGetSubscription,
			Message: "Subscriber is not found or error getting subscription by ID",
		})
	}

	if err := h.db.Models().Subscribers().Delete(ctx.Request().Context(), subscriptionID.ID.String()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.RequestError{
			Code:    errDeleteSubscription,
			Message: "Error deleting subscription",
		})
	}

	if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("reason", req.Reason)
			scope.SetExtra("subscription_id", req.SubscriptionId)
			scope.SetLevel(sentry.LevelWarning)
			hub.CaptureMessage(fmt.Sprintf("User unsubscribed: %v", req.SubscriptionId))
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *Handler) PostSubscribersConfirm(ctx echo.Context) error {
	var req api.ConfirmSubscriberRequest
	if err := ctx.Bind(&req); err != nil {
		var errorMessage string
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			errorMessage = fmt.Sprintf("%v", echoErr.Message)
		}

		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errRequestBodyBinding,
			Message: fmt.Sprintf("Error binding request body: %v", errorMessage),
		})
	}

	if hr := h.hcaptchaService.VerifyToken(req.Captcha); !hr.Success {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errValidationCaptcha,
			Message: "Invalid captcha",
		})
	}

	// Note: Token is used as subscription ID for simplicity
	subscription, err := h.db.Models().Subscribers().GetByID(ctx.Request().Context(), req.Token)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.RequestError{
			Code:    errGetSubscription,
			Message: "Subscriber is not found or error getting subscription by ID",
		})
	}

	if subscription.IsConfirmed {
		return ctx.NoContent(http.StatusOK)
	}

	subscription.IsConfirmed = true
	if err := h.db.Models().Subscribers().Update(ctx.Request().Context(), subscription); err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.RequestError{
			Code:    errUpdateSubscription,
			Message: "Error updating subscription",
		})
	}

	return ctx.NoContent(http.StatusOK)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,24}$`)
	return re.MatchString(email)
}
