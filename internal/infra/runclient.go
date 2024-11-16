package infra

import (
	"context"

	"github.com/aryanmaurya1/deployterm/internal/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type runclientWrapper struct {
	runclient client.Client
}

func NewRunclientWrapper(runclient client.Client) *runclientWrapper {
	return &runclientWrapper{runclient: runclient}
}

func (rcw *runclientWrapper) ListNamespaces(ctx context.Context) ([]*corev1.Namespace, error) {
	var nsList corev1.NamespaceList
	err := rcw.runclient.List(ctx, &nsList)
	if err != nil {
		return nil, err
	}

	return utils.ConvertToPtrList[corev1.Namespace](nsList.Items), nil
}

func (rcw *clientsetWrapper) ListDeployments(ctx context.Context, namespace string) ([]*appsv1.Deployment, error) {
	return nil, nil
}

func (rcw *clientsetWrapper) GetDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	return nil, nil
}

func (rcw *clientsetWrapper) CreateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return nil, nil
}

func (rcw *clientsetWrapper) UpdateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return nil, nil
}

func (rcw *clientsetWrapper) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	return nil, nil
}
