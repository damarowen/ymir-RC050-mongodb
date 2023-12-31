// Package rest is port handler.
package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	pkgRest "github.com/kubuskotak/asgard/rest"
	pkgTracer "github.com/kubuskotak/asgard/tracer"
	"github.com/kubuskotak/ymir-test/pkg/entity"
	"github.com/kubuskotak/ymir-test/pkg/usecase/users"
)

// MongorestOption is a struct holding the handler options.
type MongorestOption func(Mongorest *Mongorest)

// Mongorest handler instance data.
type Mongorest struct {
	UsersUsecase users.T
}

// NewMongorest creates a new Mongorest handler instance.
//
//	var MongorestHandler = rest.NewMongorest()
//
//	You can pass optional configuration options by passing a Config struct:
//
//	var adaptor = &adapters.Adapter{}
//	var MongorestHandler = rest.NewMongorest(rest.WithMongorestAdapter(adaptor))
func NewMongorest(opts ...MongorestOption) *Mongorest {
	// Create a new handler.
	var handler = &Mongorest{}

	for index := range opts {
		// Get the current configuration function from the opts slice.
		// index represent  the function argument
		applyConfig := opts[index]

		// We're applying the config function to the handler.
		// This step will do something to the handler, depending on what the function is designed to do.
		// It could set a variable on the handler, initialize some data, etc.

		// in this case applyConfig wil behave same like WithUsersUsecase, bcuz we only passing one function (index 0)
		applyConfig(handler)
	}

	// Return handler.
	return handler
}

// Register is endpoint group for handler.
func (h *Mongorest) Register(router chi.Router) {
	router.Get("/users", pkgRest.HandlerAdapter[GetListUsersRequest](h.GetAll).JSON)
	router.Post("/user", pkgRest.HandlerAdapter[UpsertUserRequest](h.Create).JSON)
	router.Get("/user/{UserId}", pkgRest.HandlerAdapter[GetRequestParam](h.GetByID).JSON)
	router.Put("/user/{UserId}", pkgRest.HandlerAdapter[UpsertUserRequest](h.UpdateByID).JSON)
	router.Delete("/user/{UserId}", pkgRest.HandlerAdapter[GetRequestParam](h.DeleteByID).JSON)
}

// GetAll user.
func (h *Mongorest) GetAll(w http.ResponseWriter, r *http.Request) (GetListUsersResponse, error) {
	ctx, span, l := pkgTracer.StartSpanLogTrace(r.Context(), "GetAll")
	defer span.End()
	var (
		request GetListUsersRequest
	)

	request, err := pkgRest.GetBind[GetListUsersRequest](r)
	if err != nil {
		l.Info().Msg(err.Error())
		return GetListUsersResponse{}, pkgRest.ErrBadRequest(w, r, err)
	}

	payload := entity.RequestGetUsers{
		Pagination: entity.Pagination{Limit: request.Limit, Page: request.Page},
	}

	documents, err := h.UsersUsecase.GetAll(ctx, payload)
	if err != nil {
		l.Info().Msg(err.Error())
		return GetListUsersResponse{}, pkgRest.ErrBadRequest(w, r, err)
	}

	pkgRest.Paging(r, pkgRest.Pagination{
		Page:  documents.Page,
		Limit: documents.Limit,
	})

	l.Info().Msg("GetAll")
	return GetListUsersResponse{Data: documents.Users}, nil
}

// Create user.
func (h *Mongorest) Create(w http.ResponseWriter, r *http.Request) (GetUserResponse, error) {
	ctx, span, l := pkgTracer.StartSpanLogTrace(r.Context(), "CreateUser")
	defer span.End()

	request, err := pkgRest.GetBind[UpsertUserRequest](r)
	if err != nil {
		l.Info().Msg(err.Error())
		return GetUserResponse{}, pkgRest.ErrBadRequest(w, r, err)
	}

	payload := entity.User{
		Name:  request.Name,
		Email: request.Email,
		Age:   request.Age,
	}

	documents, err := h.UsersUsecase.Create(ctx, payload)
	if err != nil {
		return GetUserResponse{}, pkgRest.ErrBadRequest(w, r, err)
	}

	l.Info().Msg("CreateUser")
	return GetUserResponse{User: documents}, nil
}

// GetByID user.
func (h *Mongorest) GetByID(w http.ResponseWriter, r *http.Request) (GetUserResponse, error) {
	ctx, span, l := pkgTracer.StartSpanLogTrace(r.Context(), "GetByID")
	defer span.End()

	request, err := pkgRest.GetBind[GetRequestParam](r)
	if err != nil {
		l.Info().Msg(err.Error())
		return GetUserResponse{}, pkgRest.ErrBadRequest(w, r, err)
	}

	doc, err := h.UsersUsecase.GetByID(ctx, request.UserID)
	if err != nil {
		l.Info().Msg(err.Error())
		return GetUserResponse{}, pkgRest.ErrBadRequest(w, r, err)
	}

	l.Info().Msg("GetByID")
	return GetUserResponse{User: doc}, nil
}

// UpdateByID user.
func (h *Mongorest) UpdateByID(w http.ResponseWriter, r *http.Request) (GetUserResponse, error) {
	ctx, span, l := pkgTracer.StartSpanLogTrace(r.Context(), "UpdateByID")
	defer span.End()

	request, err := pkgRest.GetBind[UpsertUserRequest](r)
	if err != nil {
		l.Info().Msg(err.Error())
		return GetUserResponse{}, pkgRest.ErrBadRequest(w, r, err)
	}

	payload := entity.User{
		ID:    request.UserID,
		Name:  request.Name,
		Email: request.Email,
		Age:   request.Age,
	}

	doc, err := h.UsersUsecase.UpdateByID(ctx, payload)
	if err != nil {
		l.Info().Msg(err.Error())
		return GetUserResponse{}, pkgRest.ErrBadRequest(w, r, err)
	}

	l.Info().Msg("UpdateByID")
	return GetUserResponse{User: doc}, nil
}

// DeleteByID user.
func (h *Mongorest) DeleteByID(w http.ResponseWriter, r *http.Request) (ResponseMessage, error) {
	ctx, span, l := pkgTracer.StartSpanLogTrace(r.Context(), "DeleteByID")
	defer span.End()

	request, err := pkgRest.GetBind[GetRequestParam](r)
	if err != nil {
		l.Info().Msg(err.Error())
		return ResponseMessage{}, pkgRest.ErrBadRequest(w, r, err)
	}

	err = h.UsersUsecase.DeleteByID(ctx, request.UserID)
	if err != nil {
		l.Info().Msg(err.Error())
		return ResponseMessage{}, pkgRest.ErrBadRequest(w, r, err)
	}

	l.Info().Msg("DeleteByID")
	return ResponseMessage{Message: fmt.Sprintf("success delete %v", request.UserID)}, nil
}

// WithUsersUsecase allows setting the UsersUsecase during initialisation.
func WithUsersUsecase(uc users.T) MongorestOption {
	return func(m *Mongorest) {
		m.UsersUsecase = uc
	}
}
