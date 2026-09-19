package httpapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/oapi-codegen/runtime/types"

	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth"
	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/httpapi/generated"
)

type AuthService interface {
	Register(context.Context, auth.RegisterInput) (auth.Authenticated, error)
	Login(context.Context, auth.LoginInput) (auth.Authenticated, error)
	CurrentUser(context.Context, string) (auth.User, error)
}

type AuthHandler struct {
	service      AuthService
	logger       *slog.Logger
	cookieSecure bool
}

func NewAuthHandler(service AuthService, logger *slog.Logger, cookieSecure bool) *AuthHandler {
	return &AuthHandler{service: service, logger: logger, cookieSecure: cookieSecure}
}

func (handler *AuthHandler) Register(ctx context.Context, request generated.RegisterRequestObject) (generated.RegisterResponseObject, error) {
	if request.Body == nil || request.Body.Password == nil {
		return registerError(auth.ErrInvalidPassword), nil
	}
	result, err := handler.service.Register(ctx, auth.RegisterInput{
		Email:     string(request.Body.Email),
		FirstName: request.Body.FirstName,
		LastName:  request.Body.LastName,
		Password:  *request.Body.Password,
	})
	if err != nil {
		if response := registerError(err); response != nil {
			return response, nil
		}
		handler.logger.ErrorContext(ctx, "register user", "error", err)
		return generated.Register500JSONResponse{InternalErrorJSONResponse: internalError()}, nil
	}
	cookie := handler.sessionCookie(result)
	return generated.Register201JSONResponse{
		Body:    apiUser(result.User),
		Headers: generated.Register201ResponseHeaders{SetCookie: &cookie},
	}, nil
}

func (handler *AuthHandler) Login(ctx context.Context, request generated.LoginRequestObject) (generated.LoginResponseObject, error) {
	if request.Body == nil || request.Body.Password == nil {
		return generated.Login400JSONResponse(errorResponse("invalid_password", "Укажите пароль.")), nil
	}
	result, err := handler.service.Login(ctx, auth.LoginInput{
		Email:    string(request.Body.Email),
		Password: *request.Body.Password,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return generated.Login401JSONResponse(errorResponse("invalid_credentials", "Неверный email или пароль.")), nil
		}
		handler.logger.ErrorContext(ctx, "login user", "error", err)
		return generated.Login500JSONResponse{InternalErrorJSONResponse: internalError()}, nil
	}
	cookie := handler.sessionCookie(result)
	return generated.Login200JSONResponse{
		Body:    apiUser(result.User),
		Headers: generated.Login200ResponseHeaders{SetCookie: &cookie},
	}, nil
}

func (handler *AuthHandler) GetCurrentUser(ctx context.Context, _ generated.GetCurrentUserRequestObject) (generated.GetCurrentUserResponseObject, error) {
	token, _ := ctx.Value(sessionTokenKey{}).(string)
	if token == "" {
		return unauthenticated(), nil
	}
	user, err := handler.service.CurrentUser(ctx, token)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthenticated) {
			return unauthenticated(), nil
		}
		handler.logger.ErrorContext(ctx, "get current user", "error", err)
		return generated.GetCurrentUser500JSONResponse{InternalErrorJSONResponse: internalError()}, nil
	}
	return generated.GetCurrentUser200JSONResponse(apiUser(user)), nil
}

func (handler *AuthHandler) sessionCookie(result auth.Authenticated) string {
	return (&http.Cookie{
		Name:     "session_id",
		Value:    result.SessionToken,
		Path:     "/",
		Expires:  result.ExpiresAt,
		HttpOnly: true,
		Secure:   handler.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	}).String()
}

func registerError(err error) generated.RegisterResponseObject {
	switch {
	case errors.Is(err, auth.ErrInvalidEmail):
		return generated.Register400JSONResponse(errorResponse("invalid_email", "Укажите корректный email."))
	case errors.Is(err, auth.ErrInvalidFirstName):
		return generated.Register400JSONResponse(errorResponse("invalid_first_name", "Имя должно содержать от 1 до 100 символов."))
	case errors.Is(err, auth.ErrInvalidLastName):
		return generated.Register400JSONResponse(errorResponse("invalid_last_name", "Фамилия должна содержать от 1 до 100 символов."))
	case errors.Is(err, auth.ErrInvalidPassword):
		return generated.Register400JSONResponse(errorResponse("invalid_password", "Пароль должен содержать от 12 до 128 символов."))
	case errors.Is(err, auth.ErrEmailAlreadyExists):
		return generated.Register409JSONResponse(errorResponse("email_already_registered", "Пользователь с таким email уже зарегистрирован."))
	default:
		return nil
	}
}

func unauthenticated() generated.GetCurrentUser401JSONResponse {
	return generated.GetCurrentUser401JSONResponse(errorResponse("unauthenticated", "Требуется вход в систему."))
}

func apiUser(user auth.User) generated.User {
	email := types.Email(user.Email)
	firstName := user.FirstName
	lastName := user.LastName
	role := generated.UserRole(user.Role)
	return generated.User{Id: user.ID, Email: &email, FirstName: &firstName, LastName: &lastName, Role: &role}
}

func errorResponse(code, message string) generated.ErrorResponse {
	return generated.ErrorResponse{Code: code, Message: message}
}

func internalError() generated.InternalErrorJSONResponse {
	return generated.InternalErrorJSONResponse(errorResponse("internal_error", "Внутренняя ошибка сервера."))
}

var _ generated.StrictServerInterface = (*AuthHandler)(nil)
