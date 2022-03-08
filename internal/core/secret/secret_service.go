// Package secret internal is the domain layer for creating secrets
// uses the k8 client to CRUD the node
package secret

import (
	"context"
	"fmt"
	"github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/logger"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type secretService struct{}

type secretServiceInterface interface {
	Get(name string) (*corev1.Secret, *errors.RestErr)
	Create(*SecretDto) (*corev1.Secret, *errors.RestErr)
	List() (*corev1.SecretList, *errors.RestErr)
	Delete(secret *corev1.Secret) *errors.RestErr
	Count() (*int, *errors.RestErr)
}

var (
	SecretService secretServiceInterface
	k8Client      = k8s.K8ClientService
)

func init() { SecretService = &secretService{} }

// Get returns a single secret  by name
func (service secretService) Get(name string) (*corev1.Secret, *errors.RestErr) {
	secret := &corev1.Secret{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8Client.Get(context.Background(), key, secret); err != nil {
		if apiErrors.IsNotFound(err) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("secret by name %s doesn't exist", name))
		}
		go logger.Error(service.Get, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("can't get secret by name %s", name))
	}

	return secret, nil
}

// Create creates a secret from the given spec
func (service secretService) Create(dto *SecretDto) (*corev1.Secret, *errors.RestErr) {
	t := true
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dto.Name,
			Namespace: "default",
			Labels: map[string]string{
				"kotal.io/key-type":            dto.Type,
				"app.kubernetes.io/created-by": "kotal-api",
			},
		},
		StringData: dto.Data,
		Immutable:  &t,
	}

	if err := k8Client.Create(context.Background(), secret); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			return nil, errors.NewBadRequestError(fmt.Sprintf("secret by name %s already exist", dto.Name))
		}
		go logger.Error(service.Create, err)
		return nil, errors.NewInternalServerError("error creating secret")
	}
	return secret, nil
}

// List returns all secrets
func (service secretService) List() (*corev1.SecretList, *errors.RestErr) {
	secrets := &corev1.SecretList{}

	if err := k8Client.List(context.Background(), secrets, client.InNamespace("default"), client.HasLabels{"app.kubernetes.io/created-by"}); err != nil {
		go logger.Error(service.List, err)
		return nil, errors.NewInternalServerError("failed to get all secrets")
	}

	return secrets, nil
}

// Delete a single secret node by name
func (service secretService) Delete(secret *corev1.Secret) *errors.RestErr {
	if err := k8Client.Delete(context.Background(), secret); err != nil {
		go logger.Error(service.Delete, err)
		return errors.NewInternalServerError(fmt.Sprintf("can't delete secret by name %s", secret.Name))
	}

	return nil
}

// Delete a list of secrets
func (service secretService) Count() (*int, *errors.RestErr) {
	secrets := &corev1.SecretList{}
	if err := k8Client.List(context.Background(), secrets, client.InNamespace("default"), client.HasLabels{"kotal.io/key-type"}); err != nil {
		go logger.Error(service.Count, err)
		return nil, errors.NewInternalServerError("failed to get all secrets")
	}
	length := len(secrets.Items)
	return &length, nil
}