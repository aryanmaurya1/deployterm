package infra

import (
	"context"

	"github.com/aryanmaurya1/deployterm/internal/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type clientsetWrapper struct {
	clientset *kubernetes.Clientset
}

func NewClientsetWrapper(clientset *kubernetes.Clientset) *clientsetWrapper {
	return &clientsetWrapper{clientset: clientset}
}

func (csw *clientsetWrapper) ListNamespaces(ctx context.Context) ([]*corev1.Namespace, error) {
	nsList, err := csw.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return utils.ConvertToPtrList[corev1.Namespace](nsList.Items), nil
}

func (rcw *runclientWrapper) ListDeployments(ctx context.Context, namespace string) ([]*appsv1.Deployment, error) {
	return nil, nil
}

func (rcw *runclientWrapper) GetDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	return nil, nil
}

func (rcw *runclientWrapper) CreateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return nil, nil
}

func (rcw *runclientWrapper) UpdateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return nil, nil
}

func (rcw *runclientWrapper) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	return nil, nil
}
