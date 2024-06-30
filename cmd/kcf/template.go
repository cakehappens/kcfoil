package main

import (
	"fmt"

	"github.com/go-clix/cli"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/native"

	//"github.com/grafana/tanka/pkg/tanka"
	//"github.com/posener/complete"
	"kcl-lang.io/kcl-go/pkg/plugin"

	"github.com/cakehappens/kcfoil/pkg/helm"
	"github.com/cakehappens/kcfoil/pkg/kubernetes/manifest"
)

func init() {
	plugin.RegisterPlugin(plugin.Plugin{
		Name: "helm",
		MethodMap: map[string]plugin.MethodSpec{
			"template": {
				Body: func(args *plugin.MethodArgs) (*plugin.MethodResult, error) {
					name := args.StrArg(0)
					chartpath := args.StrArg(1)
					values := args.Arg(2)

					h := helm.ExecHelm{}

					chart, err := h.ChartExists(chartpath)
					if err != nil {
						return nil, fmt.Errorf("helmTemplate: Failed to find a chart at '%s': %s. See https://tanka.dev/helm#failed-to-find-chart", chart, err)
					}

					//// check if resources exist in cache
					//helmKey, err := templateKey(name, chartpath, opts.TemplateOpts)
					//if err != nil {
					//	return nil, err
					//}
					//if entry, ok := helmTemplateCache.Load(helmKey); ok {
					//	log.Debug().Msgf("Using cached template for %s", name)
					//	return entry, nil
					//}

					opts := helm.TemplateOpts{
						Values: values.(map[string]interface{}),
					}

					// render resources
					list, err := h.Template(name, chart, opts)
					if err != nil {
						return nil, err
					}

					// convert list to map
					out, err := manifest.ListAsMap(list, "")
					if err != nil {
						return nil, err
					}

					//helmTemplateCache.Store(helmKey, out)
					//return out, nil

					return &plugin.MethodResult{V: out}, nil
				},
			},
		},
	})
}

const code = `
import kcl_plugin.helm

_three = helm.template("example", "./charts/example", {"nameOverride": "foo"})
_three | {
   deployment_example_foo.metadata.labels.kcl = "true"
}
`

func templateCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "template <path>",
		Short: "",
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		result, err := native.Run("main.k", kcl.WithCode(code))
		if err != nil {
			return err
		}
		fmt.Println(result.GetRawYamlResult())
		return nil
	}
	return cmd
}
