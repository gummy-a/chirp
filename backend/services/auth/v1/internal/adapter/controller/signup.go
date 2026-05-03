package authController

import (
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	signupUsecase "chirp/backend/services/auth/v1/internal/usecase/signup"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

func (a *AuthController) PostSignup(w http.ResponseWriter, r *http.Request) {
	var token value_object.TemporaryAccountID
	err := token.ParseString(r.FormValue("signup_token"))
	if err != nil {
		a.Logger.Error("invalid signup token", slog.String("error", err.Error()), slog.String("token", token.String()))
		http.Error(w, "invalid signup token", http.StatusBadRequest)
		return
	}

	numberCode, err := strconv.Atoi(r.FormValue("number_code"))
	if err != nil {
		a.Logger.Error("invalid number code", slog.String("error", err.Error()), slog.Int("number_code", numberCode))
		http.Error(w, "invalid number code", http.StatusBadRequest)
		return
	}

	jwtToken, err := a.UseCases.Signup.Execute(&signupUsecase.SignupAccountInput{
		SignupToken: token,
		NumberCode:  value_object.NumberCode(numberCode),
	})

	if err != nil {
		a.Logger.Error("Failed to execute account signup", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusBadRequest)
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
