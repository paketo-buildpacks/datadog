# `gcr.io/paketo-buildpacks/datadog`

The Paketo Datadog Buildpack is a Cloud Native Buildpack that contributes and configures the Datadog Agent.

## Behavior

This buildpack will participate all the following conditions are met

* A binding exists with `type` of `Datadog`

The buildpack will do the following for Java applications:

* Contributes the Datadog Java agent to a layer and configures `$JAVA_TOOL_OPTIONS` to use it
* Transforms the contents of the binding secret to environment variables with the pattern `DD_<KEY>=<VALUE>`

## Bindings

The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`

| Key                   | Value   | Description                                                                                       |
| --------------------- | ------- | ------------------------------------------------------------------------------------------------- |
| `<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>` |

### Type: `Datadog`

This binding supports arbitrary key/value mappings. If you set a secret with the key `dd-foo=bar` or `dd.foo=bar` then the secret is mapped to an env variable of `DD_FOO=bar`. These environment variables are read by the Datadog Java Agent. The [full list of what's supported by the Java Agent is found here](https://docs.datadoghq.com/tracing/setup_overview/setup/java/?tab=containers#configuration).

| Key      | Value   | Description                                                    |
| -------- | ------- | -------------------------------------------------------------- |
| `dd-foo` | `value` | The env variable `DD_FOO` will be set to the value of `value`. |


## License

This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
