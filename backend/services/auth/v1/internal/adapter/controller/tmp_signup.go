package authController

import (
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	signupUsecase "chirp/backend/services/auth/v1/internal/usecase/signup"
	"encoding/json"
	"log/slog"
	"net/http"
)

type PostTmpSignupResponse struct {
	SignupToken string `json:"signup_token"`
}

func (a *AuthController) PostTmpSignup(w http.ResponseWriter, r *http.Request) {
	accountId, err := a.UseCases.TmpSignup.Execute(signupUsecase.SignupTemporaryAccountInput{
		Email:    value_object.Email(r.FormValue("email")),
		Password: value_object.PasswordHash(r.FormValue("password")),
	})

	if err != nil {
		a.Logger.Error("Failed to execute temporary account signup", slog.String("error", err.Error()))
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&PostTmpSignupResponse{
		SignupToken: accountId.String(),
	}); err != nil {
		a.Logger.Error("Failed to encode response", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
