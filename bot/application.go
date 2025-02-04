package bot

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/johnnyeven/bot-sdk-go/bot/model"
)

type Application struct {
	AppId              string
	DisableCertificate bool
	DisableVerifyJson  bool
	Handler            func(rawRequest string) string
}

// ServeHTTP 创建一个HTTP服务
func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	//心跳请求
	if r.Method == "HEAD" {
		// 返回204
		w.WriteHeader(http.StatusNoContent)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Fatalf("[DuerOS][Application] ServeHTTP: request read failed: %s", err.Error())
		HTTPError(w, "request read failed", "Server Error", http.StatusInternalServerError)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "requestBody", body))

	if !a.Verify(w, r) {
		return
	}

	ret := a.Handler(string(body))
	_, err = w.Write([]byte(ret))
	if err != nil {
		logrus.Errorf("[DuerOS][Application] ServeHTTP: response write failed: %s", err.Error())
		return
	}
}

// Start 启动HTTP服务
func (a *Application) Start(host string) {
	err := http.ListenAndServe(host, a)

	if err != nil {
		logrus.Fatalf("[DuerOS][Application] Start failed: %s", err.Error())
	}
}

// Verify 验证请求是否合法
func (a *Application) Verify(w http.ResponseWriter, r *http.Request) bool {
	if !a.DisableVerifyJson && !verifyJSON(w, r, a.AppId) {
		return false
	}

	if !a.DisableCertificate && !validateRequest(w, r) {
		return false
	}
	return true
}

func HTTPError(w http.ResponseWriter, logMsg string, err string, errCode int) {
	if logMsg != "" {
		logrus.Errorf("[DuerOS][Application] HTTPError: %s", logMsg)
	}

	http.Error(w, err, errCode)
}

// Decode the JSON request and verify it.
func verifyJSON(w http.ResponseWriter, r *http.Request, appId string) bool {
	req := model.Request{}

	body := r.Context().Value("requestBody").([]byte)

	if err := json.Unmarshal(body, &req.Common); err != nil {
		logrus.Errorf("[DuerOS][Application] verifyJSON: %s", err.Error())
		HTTPError(w, err.Error(), "Bad Request", http.StatusBadRequest)
		return false
	}

	// Check the timestamp
	if !req.VerifyTimestamp() && r.URL.Query().Get("_dev") == "" {
		logrus.Errorf("[DuerOS][Application] verifyJSON: Request too old to continue (>180s).")
		HTTPError(w, "Request too old to continue (>180s).", "Bad Request", http.StatusBadRequest)
		return false
	}

	// Check the app id
	if !req.VerifyBotID(appId) {
		logrus.Errorf(
			"[DuerOS][Application] verifyJSON: DuerOS BotID mismatch! request botId: %s, app botId: %s",
			req.GetBotId(),
			appId,
		)
		HTTPError(w, "DuerOS BotID mismatch!", "Bad Request", http.StatusBadRequest)
		return false
	}
	return true
}

// Run all mandatory DuerOS security checks on the request.
func validateRequest(w http.ResponseWriter, r *http.Request) bool {
	// Check for debug bypass flag
	devFlag := r.URL.Query().Get("_dev")

	isDev := devFlag != ""

	if !isDev {
		isRequestValid := IsValidRequest(w, r)
		if !isRequestValid {
			return false
		}
	}
	return true
}

// IsValidRequest handles all the necessary steps to validate that an incoming http.Request has actually come from
// the DuerOS service. If an error occurs during the validation process, an http.Error will be written to the provided http.ResponseWriter.
// The required steps for request validation can be found on this page:
// https://dueros.baidu.com/didp/doc/dueros-bot-platform/dbp-deploy/authentication_markdown
func IsValidRequest(w http.ResponseWriter, r *http.Request) bool {
	certURL := r.Header.Get("SignatureCertUrl")

	// Verify certificate URL
	if !verifyCertURL(certURL) {
		logrus.Errorf("[DuerOS][Application] IsValidRequest: invalid cert URL: %s", certURL)
		HTTPError(w, "Invalid cert URL: "+certURL, "Not Authorized", http.StatusUnauthorized)
		return false
	}

	// Fetch certificate data
	certContents, err := readCert(certURL)
	if err != nil {
		logrus.Errorf("[DuerOS][Application] IsValidRequest: failed to read cert file: %s", err.Error())
		HTTPError(w, err.Error(), "Not Authorized", http.StatusUnauthorized)
		return false
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil {
		logrus.Errorf("[DuerOS][Application] IsValidRequest: failed to parse certificate PEM.")
		HTTPError(w, "Failed to parse certificate PEM.", "Not Authorized", http.StatusUnauthorized)
		return false
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		logrus.Errorf("[DuerOS][Application] IsValidRequest: failed to parse certificate: %s", err.Error())
		HTTPError(w, err.Error(), "Not Authorized", http.StatusUnauthorized)
		return false
	}

	// Check the certificate date
	if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
		logrus.Errorf("[DuerOS][Application] IsValidRequest: DuerOS certificate expired.")
		HTTPError(w, "DuerOS certificate expired.", "Not Authorized", http.StatusUnauthorized)
		return false
	}

	// Check the certificate alternate names
	foundName := false
	for _, altName := range cert.Subject.Names {
		if altName.Value == "dueros-api.baidu.com" {
			foundName = true
		}
	}

	if !foundName {
		logrus.Errorf("[DuerOS][Application] IsValidRequest: DuerOS certificate invalid.")
		HTTPError(w, "DuerOS certificate invalid.", "Not Authorized", http.StatusUnauthorized)
		return false
	}

	// Verify the key
	publicKey := cert.PublicKey
	encryptedSig, _ := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))

	// Make the request body SHA1 and verify the request with the public key
	//var bodyBuf bytes.Buffer
	hash := sha1.New()
	_, err = io.Copy(hash, bytes.NewReader(r.Context().Value("requestBody").([]byte)))
	if err != nil {
		logrus.Errorf("[DuerOS][Application] IsValidRequest: failed to hash request body: %s", err.Error())
		HTTPError(w, err.Error(), "Internal Error", http.StatusInternalServerError)
		return false
	}

	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSig)
	if err != nil {
		logrus.Errorf("[DuerOS][Application] IsValidRequest: signature match failed: %s", err.Error())
		HTTPError(w, "Signature match failed.", "Not Authorized", http.StatusUnauthorized)
		return false
	}

	return true
}

func readCert(certURL string) ([]byte, error) {
	cert, err := http.Get(certURL)
	if err != nil {
		return nil, errors.New("could not download DuerOS cert file")
	}
	defer cert.Body.Close()
	certContents, err := io.ReadAll(cert.Body)
	if err != nil {
		return nil, errors.New("could not read DuerOS cert file")
	}

	return certContents, nil
}

func verifyCertURL(path string) bool {
	link, _ := url.Parse(path)

	if link.Scheme != "https" {
		return false
	}

	if link.Host != "duer.bdstatic.com" {
		return false
	}

	return true
}
