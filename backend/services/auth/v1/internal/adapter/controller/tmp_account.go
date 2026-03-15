package authController

import (
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	"log/slog"
	"net/http"
)

func (a *AuthController) GetTmpAccount(w http.ResponseWriter, r *http.Request) {
	account_id := r.PathValue("account_id")
	if account_id == "" {
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	var accountId value_object.TemporaryAccountID
	err := accountId.ParseString(account_id)
	if err != nil {
		a.Logger.Error("failed to parse account_id", slog.String("error", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	res, err := a.UseCases.TmpSignup.FindById(&accountId)
	if err != nil {
		a.Logger.Error("failed to find account", slog.String("error", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if res == nil {
		a.Logger.Error("account not found")
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	// success
	w.WriteHeader(http.StatusNoContent)
}
