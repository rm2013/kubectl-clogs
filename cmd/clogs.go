package cmd

import (
	"fmt"
	"log"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/rm2013/kubectl-clogs/clogs"
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	clogsLong = `
cLogs is an interactive pod and container selector for 'kubectl logs'

Arg[1] will act as a filter, any pods that match will be returned in a list
that the user can select from.


`
	clogsExample = `
	
	# select from all pods matching [busybox] then run: 'kubectl logs [pod_name] '
	%[1]s clogs busybox

	# select from all pods matching [multi_container_pod]
	# then select from all containers in pod matching [container]
	# then run: 'kubectl logs [pod_name] [container_name] '
	%[1]s clogs multi_container_pod -c container
`
)

// CLogsOptions
type CLogsOptions struct {
	configFlags *genericclioptions.ConfigFlags
	clientCfg   *rest.Config

	configOverrides clientcmd.ConfigOverrides
	allNamespaces   bool
	containerFilter string
	namespace       string

	genericclioptions.IOStreams
}

// NewCLogsOptions provides an instance of CLogsOptions with default values
func NewCLogsOptions(streams genericclioptions.IOStreams) *CLogsOptions {
	return &CLogsOptions{
		configFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

// NewCmdClogs provides a cobra command wrapping CLogsOptions
func NewCmdClogs(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewCLogsOptions(streams)

	cmd := &cobra.Command{
		Use:          "clogs [pod filter] [flags]",
		Short:        "Logs from a Kubernetes Pod container",
		Args:         cobra.MinimumNArgs(1),
		Example:      fmt.Sprintf(clogsExample, "kubectl"),
		Long:         clogsLong,
		SilenceUsage: false,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Run(args); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&o.allNamespaces, "all-namespaces", "A", o.allNamespaces, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	cmd.PersistentFlags().StringVarP(&o.containerFilter, "container", "c", "", "Container to search")
	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete clogs
func (o *CLogsOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	o.clientCfg, err = o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	c := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&o.configOverrides)

	o.namespace, _, err = c.Namespace()
	if err != nil {
		return err
	}

	if *o.configFlags.Namespace != "" {
		o.namespace = *o.configFlags.Namespace
	}

	if o.allNamespaces {
		o.namespace = ""
	}

	return nil
}

// Run clogs
func (o *CLogsOptions) Run(args []string) error {

	podFilter := args[0]

	config := &clogs.Config{
		Namespace:       o.namespace,
		PodFilter:       podFilter,
		ContainerFilter: o.containerFilter,
	}

	r := clogs.NewClogs(o.clientCfg, config)

	if err := r.Do(); err != nil {
		log.Fatal(err)
	}

	return nil
}
