/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package component

import "encoding/json"

const (
	// WidthQuarter is a quarter width section.
	WidthQuarter int = 6
	// WidthHalf is a half width section.
	WidthHalf int = 12
	// WidthFull is a full width section.
	WidthFull int = 24
)

// FlexLayoutItem is an item in a flex layout.
type FlexLayoutItem struct {
	Width int       `json:"width,omitempty"`
	View  Component `json:"view,omitempty"`
}

func (fli *FlexLayoutItem) UnmarshalJSON(data []byte) error {
	x := struct {
		Width int
		View  TypedObject
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	fli.Width = x.Width
	var err error
	fli.View, err = x.View.ToComponent()
	if err != nil {
		return err
	}

	return nil
}

// FlexLayoutSection is a slice of items group together.
type FlexLayoutSection []FlexLayoutItem

// FlexLayoutConfig is configuration for the flex layout view.
type FlexLayoutConfig struct {
	Sections    []FlexLayoutSection `json:"sections,omitempty"`
	ButtonGroup *ButtonGroup        `json:"buttonGroup,omitempty"`
}

// FlexLayout is a flex layout view.
type FlexLayout struct {
	base
	Config FlexLayoutConfig `json:"config,omitempty"`
}

// NewFlexLayout creates an instance of FlexLayout.
func NewFlexLayout(title string) *FlexLayout {
	return &FlexLayout{
		base: newBase(typeFlexLayout, TitleFromString(title)),
		Config: FlexLayoutConfig{
			ButtonGroup: NewButtonGroup(),
		},
	}
}

// GetMetadata returns the metadata for the flex layout view.
func (fl *FlexLayout) GetMetadata() Metadata {
	return fl.Metadata
}

// AddSections adds one or more sections to the flex layout.
func (fl *FlexLayout) AddSections(sections ...FlexLayoutSection) {
	fl.Config.Sections = append(fl.Config.Sections, sections...)
}

type flexLayoutMarshal FlexLayout

// MarshalJSON marshals the flex layout to JSON.
func (fl *FlexLayout) MarshalJSON() ([]byte, error) {
	x := flexLayoutMarshal(*fl)
	x.Metadata.Type = typeFlexLayout
	return json.Marshal(&x)
}

func (fl *FlexLayout) SetButtonGroup(group *ButtonGroup) {
	fl.Config.ButtonGroup = group
}

// Tab represents a tab. A tab is a flex layout with a name.
type Tab struct {
	Name     string
	Contents FlexLayout
}

// NewTabWithContents creates a tab with contents.
func NewTabWithContents(flexLayout FlexLayout) *Tab {
	name, err := TitleFromTitleComponent(flexLayout.Title)
	if err != nil {
		name = ""
	}

	return &Tab{
		Name:     name,
		Contents: flexLayout,
	}
}
