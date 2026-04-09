package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestConfigFlagsFromClientGetterUnwrapsMatchVersionFlags(t *testing.T) {
	configFlags := genericclioptions.NewConfigFlags(true)

	require.Same(t, configFlags, configFlagsFromClientGetter(NewMatchVersionFlags(configFlags)))
}

func TestFactoryFileSourceHandlersUseWrappedConfigFlags(t *testing.T) {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.HandleSecretFromFileSources = func(secret *corev1.Secret, fileSources []string) error {
		secret.Data["keystore.jks"] = []byte("memfs")
		return nil
	}
	configFlags.HandleConfigMapFromFileSources = func(configMap *corev1.ConfigMap, fileSources []string) error {
		configMap.Data["config.yml"] = "memfs"
		return nil
	}
	configFlags.HandleConfigMapFromEnvFileSources = func(configMap *corev1.ConfigMap, envFileSources []string) error {
		configMap.Data["ENV"] = "VALUE"
		return nil
	}

	factory := NewFactory(NewMatchVersionFlags(configFlags))

	secret := &corev1.Secret{Data: map[string][]byte{}}
	require.NoError(t, factory.SecretFromFileSources()(secret, []string{"keystore.jks"}))
	require.Equal(t, []byte("memfs"), secret.Data["keystore.jks"])

	configMap := &corev1.ConfigMap{Data: map[string]string{}, BinaryData: map[string][]byte{}}
	require.NoError(t, factory.ConfigMapFromFileSources()(configMap, []string{"config.yml"}))
	require.NoError(t, factory.ConfigMapFromEnvFileSources()(configMap, []string{"settings.env"}))
	require.Equal(t, "memfs", configMap.Data["config.yml"])
	require.Equal(t, "VALUE", configMap.Data["ENV"])
}
