/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package action

import "time"

const (
	// DefaultAlertExpiration is the default expiration for alerts.
	DefaultAlertExpiration = 10 * time.Second
)

// AlertType is the type of alert.
type AlertType string

const (
	AlertTypeError   AlertType = "ERROR"
	AlertTypeWarning AlertType = "WARNING"
	AlertTypeInfo    AlertType = "INFO"
)

// Alert is an alert message.
type Alert struct {
	Type       AlertType  `json:"type"`
	Message    string     `json:"message"`
	Expiration *time.Time `json:"expiration,omitempty"`
}

// CreateAlert creates an alert with optional expiration. If the expireAt is < 1 That Expiration will be nil.
func CreateAlert(alertType AlertType, message string, expireAt time.Duration) Alert {
	alert := Alert{
		Type:    alertType,
		Message: message,
	}

	if expireAt > 0 {
		t := time.Now().Add(expireAt)
		alert.Expiration = &t
	}

	return alert
}

type Alerter interface {
	SendAlert(alert Alert)
}
