package authController

import (
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	loginUsecase "chirp/backend/services/auth/v1/internal/usecase/login"
	"log/slog"
	"net/http"
	"os"
)

func (a *AuthController) PostLogin(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := a.UseCases.Login.Execute(loginUsecase.LoginAccountInput{
		Email:    value_object.Email(r.FormValue("email")),
		Password: value_object.PasswordPlainText(r.FormValue("password")),
	})
	if err != nil {
		a.Logger.Error("failed to login", slog.String("error", err.Error()))
		http.Error(w, "failed to login", http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    jwtToken.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("AUTH_SERVICE_APP_ENV") == "production",
		MaxAge:   60 * 60 * 24, // 1day
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusNoContent)
}
