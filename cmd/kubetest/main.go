package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	fakeObjectPrefix = "kubetest"
)

func main() {
	cmd := newCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running kubetest: %v\n", err)
		os.Exit(1)
	}
}

type options struct {
	// The path of kubeconfig configuration file.
	kubeconfig string
	context    string
}

func newOption() *options {
	defaultKubeconfig := os.Getenv("KUBECONFIG")
	if defaultKubeconfig == "" {
		defaultKubeconfig = filepath.Join(os.Getenv("HOME"), ".kube/config")
	}
	return &options{
		kubeconfig: defaultKubeconfig,
	}
}

func (o *options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.kubeconfig, "kubeconfig", o.kubeconfig, "Path to the kubeconfig file to use for CLI requests.")
	fs.StringVar(&o.context, "context", o.context, "The name of the kubeconfig context to use")
}

func newCommand() *cobra.Command {
	o := newOption()
	cmd := &cobra.Command{
		Use: "kubetest [--kubeconfig PATH]",
	}
	createCommand := newCreateCommand(o)
	flushCommand := newFlushCommand(o)
	cmd.AddCommand(createCommand)
	cmd.AddCommand(flushCommand)
	flags := cmd.PersistentFlags()
	o.AddFlags(flags)
	return cmd
}

func getClient(o *options) (kubernetes.Interface, error) {
	kubeConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: o.kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: o.context}).ClientConfig()
	if err != nil {
		return nil, err
	}
	kubeConfig.QPS = 1000
	kubeConfig.Burst = 1000

	k8sClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}
