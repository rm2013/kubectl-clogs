package clogs

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

var podTemplate = &promptui.SelectTemplates{
	Active:   fmt.Sprintf("Namespace: {{ .Namespace | blue }} | Pod: %s {{ .Name | cyan }}", promptui.IconSelect),
	Inactive: "Namespace: {{ .Namespace | blue }} | Pod: {{ .Name | magenta }}",
	Selected: fmt.Sprintf("Namespace: {{ .Namespace | blue }} | Pod: %s {{ .Name | cyan }}", promptui.IconGood),
}
var containerTemplates = &promptui.SelectTemplates{
	Active:   fmt.Sprintf("Container: %s {{ .Name | cyan }}", promptui.IconSelect),
	Inactive: "Container: {{ .Name | magenta }}",
	Selected: fmt.Sprintf("Container: %s {{ .Name | cyan }}", promptui.IconGood),
}
