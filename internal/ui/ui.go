package ui

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aryanmaurya1/deployterm/internal"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var currentPage string
var prevPage string

func switchErrorPage(app *tview.Application, pages *tview.Pages, err error) {
	currentPage = fmt.Sprintf("%d", time.Now().UnixNano())
	errorModal := tview.NewModal().
		SetText(fmt.Sprintf("Error\n%+v", err.Error())).
		AddButtons([]string{"Back", "Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				// Remove dynamically created page
				pages.RemovePage(currentPage)

				currentPage = prevPage
				pages.SwitchToPage(currentPage)
			} else {
				app.Stop()
			}
		}).
		SetBackgroundColor(tcell.ColorGreen)
	pages.AddAndSwitchToPage(currentPage, errorModal, false)
}

func namespaceListPage(app *tview.Application, pages *tview.Pages, opsClient internal.IOperation) tview.Primitive {
	currentPage = NAME_NAMESPACE_PAGE

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Namespaces")
	list.SetTitleColor(tcell.ColorAntiqueWhite)
	list.SetFocusFunc(
		func() {
			list.Clear()
			namespaces, err := opsClient.ListNamespaces(context.Background())
			if err != nil {
				fmt.Printf("error in listing namespace: %+v\n", err)
				prevPage = currentPage
				switchErrorPage(app, pages, err)
			}

			for idx, ns := range namespaces {
				primary := ns.Name
				secondary := fmt.Sprintf("create_at: %+v", ns.CreationTimestamp)
				shortcut := rune(int('a') + idx)
				if shortcut == 'q' {
					shortcut = '1'
				}

				list.AddItem(primary, secondary, shortcut, nil)
			}

			list.AddItem("Quit", "Press to exit", rune('q'), func() {
				app.Stop()
			})

			list.AddItem("AQ", "Press to exit", rune('x'), func() {
				prevPage = currentPage
				switchErrorPage(app, pages, errors.New("type"))
			})
		})
	return list
}

func GetRootPage(app *tview.Application, pages *tview.Pages, opsClient internal.IOperation) (tview.Primitive, string) {
	return namespaceListPage(app, pages, opsClient), NAME_NAMESPACE_PAGE
}
