package controller

import (
	"auth_backend/errors"
	"net/http"
)

func errorResponse(w http.ResponseWriter, err *errors.Error) {
	w.WriteHeader(err.GetHttpStatus())
	w.Write(err.ToJson())
}

func successResponse(w http.ResponseWriter, response []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func Router() http.Handler {

	mux := http.NewServeMux()

	// GET
	mux.Handle("/api/auth/basic", corsGet(http.HandlerFunc(authBasic)))
	mux.Handle("/api/auth/oauth42", corsGet(http.HandlerFunc(authOauth42)))
	mux.Handle("/api/auth/oauthFb", corsGet(http.HandlerFunc(authOauthFb)))
	mux.Handle("/api/auth/oauthVk", corsGet(http.HandlerFunc(authOauthVk)))
	mux.Handle("/api/auth/oauthGoogle", corsGet(http.HandlerFunc(authOauthGoogle)))
	mux.Handle("/api/info", corsGet(http.HandlerFunc(info)))
	mux.Handle("/api/profile/get", corsGet(authMW(http.HandlerFunc(profileGet))))

	mux.Handle("/api/passwd/repair/confirm", corsGet(http.HandlerFunc(passwdRepairConfirm)))
	mux.Handle("/api/email/confirm", corsGet(http.HandlerFunc(emailConfirm)))
	mux.Handle("/api/email/patch/confirm", corsGet(http.HandlerFunc(emailPatchConfirm)))

	// GET IMAGES
	mux.Handle("/api/image", corsGet(http.HandlerFunc(image)))

	// POST
	mux.Handle("/api/auth/check", corsPost(http.HandlerFunc(authCheck)))
	mux.Handle("/api/email/resend", corsPost(http.HandlerFunc(emailResend)))
	mux.Handle("/api/passwd/repair", corsPost(http.HandlerFunc(passwdRepair)))

	// PUT
	mux.Handle("/api/profile/create", corsPut(http.HandlerFunc(profileCreate)))

	// PATCH
	mux.Handle("/api/email/patch", corsPatch(authMW(http.HandlerFunc(emailPatch))))
	mux.Handle("/api/passwd/patch", corsPatch(authMW(http.HandlerFunc(passwdPatch))))
	mux.Handle("/api/profile/patch", corsPatch(authMW(http.HandlerFunc(profilePatch))))
	mux.Handle("/api/passwd/repair/patch", corsPatch(http.HandlerFunc(passwdRepairPatch)))

	// DELETE
	mux.Handle("/api/profile/delete", corsDelete(authMW(http.HandlerFunc(profileDelete))))

	serveMux := panicRecover(mux)

	return serveMux
}
