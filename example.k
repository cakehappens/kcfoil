import kcl_plugin.helm

helm.template("example", "./charts/example", {
  values: {
    nameOverride: "foo"
  }
}) | {
  deployment_example_foo.metadata.labels.kcl = "true"
}
