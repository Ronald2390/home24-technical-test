package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"home24-technical-test/internal/user"
	userAdapter "home24-technical-test/internal/user/adapter"
	userPublic "home24-technical-test/internal/user/public"
	"home24-technical-test/pkg/appcontext"
	"home24-technical-test/pkg/data"
	"home24-technical-test/pkg/http/response"
)

// UserController represents the user controller
type UserController struct {
	getLoginSessionAdapter userAdapter.GetLoginSessionAdapter
	loginAdapter           userAdapter.LoginAdapter
	logoutAdapter          userAdapter.LogoutAdapter
	changePasswordAdapter  userAdapter.ChangePasswordAdapter
	dataManager            *data.Manager
}

// Login POST /login
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var params userPublic.LoginParams
	err := decoder.Decode(&params)
	if err != nil {
		response.Error(w, "Bad Request", http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	var errLogin error
	var sess *userPublic.LoginResponse
	err = uc.dataManager.RunInTransaction(ctx, func(tctx context.Context) error {
		sess, errLogin = uc.loginAdapter.Execute(r.Context(), &params)
		return errLogin
	})
	if err != nil {
		if err == user.ErrWrongPassword || err == data.ErrNotFound {
			response.Error(w, "Email or password is wrong", http.StatusBadRequest, err)
		} else {
			fmt.Printf("Error: %v", err)
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "sessionId",
		Value: sess.SessionID,
	})

	response.JSON(w, http.StatusOK, sess)
}

// GetLoginSession GET /session
func (uc *UserController) GetLoginSession(w http.ResponseWriter, r *http.Request) {
	var token = appcontext.SessionID(r.Context())

	sess, err := uc.getLoginSessionAdapter.Execute(r.Context(), token)
	if err != nil {
		if err == user.ErrWrongPassword || err == data.ErrNotFound {
			response.Error(w, "Email or password is wrong", http.StatusBadRequest, err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, err)
		}
		return
	}

	response.JSON(w, http.StatusOK, sess)
}

// ChangePassword PUT /users/password
func (uc *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var params userPublic.ChangePasswordParams
	err := decoder.Decode(&params)
	if err != nil {
		response.Error(w, "Bad Request", http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	userID := appcontext.UserID(ctx)
	err = uc.dataManager.RunInTransaction(ctx, func(tctx context.Context) error {
		err = uc.changePasswordAdapter.Execute(ctx, userID, params.OldPassword, params.NewPassword)
		return err
	})
	if err != nil {
		if err == user.ErrWrongPassword {
			response.Error(w, "Wrong old password", http.StatusBadRequest, err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, err)
		}
		return
	}

	response.JSON(w, http.StatusNoContent, "")
}

// Logout POST /v1/logout
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// get token from the context
	// log it out!
	loginToken, ok := ctx.Value(appcontext.KeySessionID).(string)
	if !ok {
		err := errors.New("failed to get user id from request context")
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, err)
		return
	}

	err := uc.dataManager.RunInTransaction(ctx, func(tctx context.Context) error {
		err := uc.logoutAdapter.Execute(ctx, loginToken)
		return err
	})
	if err != nil {
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, "")
}

// NewUserController creates a new user controller
func NewUserController(
	getLoginSessionAdapter userAdapter.GetLoginSessionAdapter,
	loginAdapter userAdapter.LoginAdapter,
	logoutAdapter userAdapter.LogoutAdapter,
	changePasswordAdapter userAdapter.ChangePasswordAdapter,
	dataManager *data.Manager,
) *UserController {
	return &UserController{
		getLoginSessionAdapter: getLoginSessionAdapter,
		loginAdapter:           loginAdapter,
		logoutAdapter:          logoutAdapter,
		changePasswordAdapter:  changePasswordAdapter,
		dataManager:            dataManager,
	}
}
