package ui

import (
	"fmt"

	"github.com/rivo/tview"
	corev1 "k8s.io/api/core/v1"
)

func NamespaceList(app *tview.Application, namespaces []*corev1.Namespace) *tview.List {
	list := tview.NewList()
	for idx, ns := range namespaces {
		primary := ns.Name
		secondary := fmt.Sprintf("create_at: %+v", ns.CreationTimestamp)
		shortcut := rune(int('a') + idx)
		if shortcut == 'q' {
			shortcut = '1'
		}

		list.AddItem(primary, secondary, shortcut, nil)
	}

	list.AddItem("Quit", "Press to exit", rune('q'), func() { app.Stop() })
	return list
}
