/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package flexlayout

import (
	"github.com/kubenext/kubeon/pkg/action"
	"github.com/kubenext/kubeon/pkg/view/component"
)

type FlexLayout struct {
	sections    []*Section
	buttonGroup *component.ButtonGroup
}

func NewFlexLayout() *FlexLayout {
	return &FlexLayout{
		buttonGroup: component.NewButtonGroup(),
	}
}

// AddSection adds a new section to the flex layout.
func (fl *FlexLayout) AddSection() *Section {
	section := NewSection()
	fl.sections = append(fl.sections, section)
	return section
}

// AddButton adds a button the button group for a flex layout.
func (fl *FlexLayout) AddButton(name string, payload action.Payload, buttonOptions ...component.ButtonOption) {
	button := component.NewButton(name, payload, buttonOptions...)
	fl.buttonGroup.AddButton(button)
}

// ToComponent converts the FlexLayout to a FlexLayout.
func (fl *FlexLayout) ToComponent(title string) *component.FlexLayout {
	var sections []component.FlexLayoutSection

	for _, section := range fl.sections {
		layoutSection := component.FlexLayoutSection{}

		for _, member := range section.Members {
			item := component.FlexLayoutItem{
				Width: member.Width,
				View:  member.View,
			}

			layoutSection = append(layoutSection, item)
		}

		sections = append(sections, layoutSection)
	}

	if title == "" {
		title = "Summary"
	}

	view := component.NewFlexLayout(title)
	view.AddSections(sections...)
	view.SetButtonGroup(fl.buttonGroup)

	return view
}
