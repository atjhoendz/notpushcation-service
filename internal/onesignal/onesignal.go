package onesignal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/atjhoendz/notpushcation-service/internal/model"

	"github.com/atjhoendz/notpushcation-service/internal/config"
	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
)

type client struct {
	httpClient *http.Client
	appID      string
}

var hostURL = "https://onesignal.com/api/v1/"

// NewClient :nodoc:
func NewClient(httpClient *http.Client, appID string) model.OnesignalClient {
	return &client{
		httpClient: httpClient,
		appID:      appID,
	}
}

func (c client) Deliver(ctx context.Context, message *model.OnesignalPayload) error {
	message.AppID = c.appID

	logger := log.WithFields(log.Fields{
		"ctx":     utils.Dump(ctx),
		"message": utils.Dump(message),
	})

	url, err := utils.JoinURL(hostURL, "notifications")
	if err != nil {
		logger.Error(err)
		return err
	}

	msgJSON, err := json.Marshal(message)
	if err != nil {
		logger.Error(err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msgJSON))
	if err != nil {
		logger.Error(err)
		return err
	}

	req.Header.Set("Authorization", "Basic "+config.OnesignalAPIKey())
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.WithField("req", fmt.Sprintf("%+v", req)).Error(err)
		if e, ok := err.(net.Error); ok && e.Timeout() {
			logger.Error(ErrOnesignalTimeout)
			return ErrOnesignalTimeout
		}
		return err
	}

	defer func() {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		logger.WithField("response", string(responseBody)).Info("onesignal response body")
	}()

	if resp.StatusCode >= http.StatusInternalServerError {
		logger.Error(ErrOnesignalServerSide)
		return ErrOnesignalServerSide
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error(ErrOnesignalUnexpectedResponse)
		return ErrOnesignalUnexpectedResponse
	}

	return nil
}
