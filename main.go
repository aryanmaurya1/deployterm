package main

import (
	"context"
	"flag"
	"fmt"
	"reflect"

	"github.com/aryanmaurya1/deployterm/internal"
	"github.com/aryanmaurya1/deployterm/internal/infra"
	"github.com/aryanmaurya1/deployterm/internal/ui"
	"github.com/rivo/tview"
)

func main() {
	var kconfigPath string
	var useControllerRuntime bool

	flag.StringVar(&kconfigPath, "kubeconfig", "", "path to kubeconfig file")
	flag.BoolVar(&useControllerRuntime, "use-controller-runtime", false, "use controller-runtime library instead of client-go")
	flag.Parse()

	k8sClient, err := internal.NewK8sClient(kconfigPath)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to k8s cluster: %+v\n", err.Error()))
	}

	var operationClient internal.IOperation

	operationClient = infra.NewClientsetWrapper(k8sClient.GetClientset())
	if useControllerRuntime {
		operationClient = infra.NewRunclientWrapper(k8sClient.GetRunclient())
	}

	fmt.Printf("using client: %+v\n", reflect.TypeOf(operationClient))

	nsList, err := operationClient.ListNamespaces(context.Background())
	if err != nil {
		panic("")
	}

	app := tview.NewApplication()
	namespaceList := ui.NamespaceList(app, nsList)
	if err := app.SetRoot(namespaceList, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
