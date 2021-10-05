package lookup

import (
	"github.com/ory/kratos/text"
	"github.com/ory/kratos/ui/node"
)

func NewRevealLookupNode() *node.Node {
	return node.NewInputField(node.LookupReveal, "true", node.LookupGroup, node.InputAttributeTypeSubmit,
		node.WithInputAttributes(func(a *node.InputAttributes) {
			a.InitialValue = "false"
		})).
		WithMetaLabel(text.NewInfoSelfServiceSettingsRevealLookup())
}

func NewRegenerateLookupNode() *node.Node {
	return node.NewInputField(
		node.LookupRegenerate, "true", node.LookupGroup, node.InputAttributeTypeSubmit, node.WithInputAttributes(func(a *node.InputAttributes) {
			a.InitialValue = "false"
		})).
		WithMetaLabel(text.NewInfoSelfServiceSettingsRegenerateLookup())
}

func NewDisableLookupNode() *node.Node {
	return node.NewInputField(node.LookupDisable, "true", node.LookupGroup, node.InputAttributeTypeSubmit, node.WithInputAttributes(func(a *node.InputAttributes) {
		a.InitialValue = "false"
	})).
		WithMetaLabel(text.NewInfoSelfServiceSettingsDisableLookup())
}

func NewConfirmLookupNode() *node.Node {
	return node.NewInputField(node.LookupConfirm, "true", node.LookupGroup, node.InputAttributeTypeSubmit, node.WithInputAttributes(func(a *node.InputAttributes) {
		a.InitialValue = "false"
	})).
		WithMetaLabel(text.NewInfoSelfServiceSettingsLookupConfirm())
}