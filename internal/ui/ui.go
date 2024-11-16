package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/aryanmaurya1/deployterm/internal"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	NAME_NAMESPACE_PAGE = "NAMESPACE_PAGE"
)

type stack []string

func (stk *stack) push(element string) {
	*stk = append(*stk, element)
}

func (stk *stack) pop() string {
	if len(*stk) == 0 {
		return ""
	}

	v := (*stk)[len(*stk)-1]
	*stk = (*stk)[:len(*stk)-1]
	return v
}

var stk stack

func switchToErrorPage(app *tview.Application, pages *tview.Pages, err error) {
	currentPageName := fmt.Sprintf("%d", time.Now().UnixNano())
	errorModal := tview.NewModal().
		SetText(fmt.Sprintf("Error\n%+v", err.Error())).
		AddButtons([]string{"Back", "Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				// Remove dynamically added page
				pages.RemovePage(currentPageName)

				pages.SwitchToPage(stk.pop())
			} else {
				app.Stop()
			}
		}).
		SetBackgroundColor(tcell.ColorGreen)
	pages.AddAndSwitchToPage(currentPageName, errorModal, false)
}

func switchToOptionsPage(app *tview.Application, pages *tview.Pages, namespace string, deploymentName string, opsClient internal.IK8sOperation) {
	currentPageName := fmt.Sprintf("%d", time.Now().UnixNano())

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(fmt.Sprintf("Namespace - <%s> | Deployment - <%s>", namespace, deploymentName))
	list.SetTitleColor(tcell.ColorAntiqueWhite)
	list.SetFocusFunc(
		func() {
			list.Clear()

			list.AddItem("Back", "Press to go back", rune('x'), func() {
				// Remove dynamically added page
				pages.RemovePage(currentPageName)

				pages.SwitchToPage(stk.pop())
			})

			list.AddItem("Quit", "Press to exit", rune('q'), func() {
				app.Stop()
			})

		})

	pages.AddAndSwitchToPage(currentPageName, list, true)
}

func switchToDeploymentListPage(app *tview.Application, pages *tview.Pages, namespace string, opsClient internal.IK8sOperation) {
	currentPageName := fmt.Sprintf("%d", time.Now().UnixNano())

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(fmt.Sprintf("Namespace - <%s>", namespace))
	list.SetTitleColor(tcell.ColorAntiqueWhite)
	list.SetFocusFunc(
		func() {
			list.Clear()
			deployments, err := opsClient.ListDeployments(context.Background(), namespace)
			if err != nil {
				fmt.Printf("error in listing deployments: %+v\n", err)
				stk.push(currentPageName)
				switchToErrorPage(app, pages, err)
			}

			for idx, deployment := range deployments {
				primary := deployment.Name
				secondary := fmt.Sprintf("desired_replicas: %+v | current_replica: %+v", *deployment.Spec.Replicas, deployment.Status.Replicas)
				shortcut := rune(int('a') + idx)

				// Avoid shortcut collision with Quit/Back operations
				if shortcut == 'q' {
					shortcut = '1'
				} else if shortcut == 'x' {
					shortcut = '2'
				}

				list.AddItem(primary, secondary, shortcut, func() {
					stk.push(currentPageName)
					switchToOptionsPage(app, pages, namespace, deployment.Name, opsClient)
				})
			}

			list.AddItem("Back", "Press to go back", rune('x'), func() {
				// Remove dynamically added page
				pages.RemovePage(currentPageName)

				pages.SwitchToPage(stk.pop())
			})

			list.AddItem("Quit", "Press to exit", rune('q'), func() {
				app.Stop()
			})

		})

	pages.AddAndSwitchToPage(currentPageName, list, true)
}

func rootPage(app *tview.Application, pages *tview.Pages, opsClient internal.IK8sOperation) tview.Primitive {
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
				stk.push(NAME_NAMESPACE_PAGE)
				switchToErrorPage(app, pages, err)
			}

			for idx, ns := range namespaces {
				primary := ns.Name
				secondary := fmt.Sprintf("create_at: %+v", ns.CreationTimestamp)
				shortcut := rune(int('a') + idx)

				// Avoid shortcut collision with Quit/Back operations
				if shortcut == 'q' {
					shortcut = '1'
				}

				list.AddItem(primary, secondary, shortcut, func() {
					stk.push(NAME_NAMESPACE_PAGE)
					switchToDeploymentListPage(app, pages, ns.Name, opsClient)
				})
			}

			list.AddItem("Quit", "Press to exit", rune('q'), func() {
				app.Stop()
			})

		})
	return list
}

func GetRootPage(app *tview.Application, pages *tview.Pages, opsClient internal.IK8sOperation) (tview.Primitive, string) {
	return rootPage(app, pages, opsClient), NAME_NAMESPACE_PAGE
}
