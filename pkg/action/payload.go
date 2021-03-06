/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package action

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"math"
	"strconv"
)

// Payload is an actin payload.
type Payload map[string]interface{}

// Creates a payload with an action name and fields.
func CreatePayload(actionName string, fields map[string]interface{}) Payload {
	payload := Payload{
		"action": actionName,
	}

	for k, v := range fields {
		payload[k] = v
	}

	return payload
}

// Extracts a GroupVersionKind from a payload.
func (p Payload) GroupVersionKind() (schema.GroupVersionKind, error) {
	group, err := p.String("group")
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	version, err := p.String("version")
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	kind, err := p.String("kind")
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	return schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}, nil
}

// Returns a uint16 from the  payload.
func (p Payload) Uint16(key string) (uint16, error) {
	i, found, err := unstructured.NestedFloat64(p, key)
	if err != nil {
		return 0, err
	}

	if !found {
		return 0, errors.Errorf("payload does not contain %q", key)
	}

	if i > math.MaxUint16 || i < 0 {
		return 0, errors.Errorf("value %v is not a valid uint16", i)
	}

	return uint16(i), nil
}

// Returns a string from the payload.
func (p Payload) String(key string) (string, error) {
	s, ok := p[key].(string)
	if !ok {
		return "", errors.Errorf("payload does not contain %q", key)
	}

	return s, nil
}

// Returns a string from the payload. If the string does not exist, it returns an empty string.
func (p Payload) OptionalString(key string) (string, error) {
	s, _, err := unstructured.NestedString(p, key)
	if err != nil {
		return "", err
	}

	return s, nil
}

// Returns a string slice from the payload.
func (p Payload) StringSlice(key string) ([]string, error) {
	sli, ok := p[key].([]interface{})
	if !ok {
		return nil, errors.Errorf("payload does not contain %q", key)
	}

	var list []string
	for i := range sli {
		s, ok := sli[i].(string)
		if !ok {
			return nil, errors.New("could not convert slice entry to string")
		}
		list = append(list, s)
	}

	return list, nil
}

// Returns a float64 from the payload.
func (p Payload) Float64(key string) (float64, error) {
	switch v := p[key].(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	default:
		return 0, errors.Errorf("unable to handle type %T for %q; got %#v", p[key], key, v)
	}
}
