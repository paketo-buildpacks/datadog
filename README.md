# `gcr.io/paketo-buildpacks/datadog`

The Paketo Buildpack for Datadog is a Cloud Native Buildpack that contributes and configures the Datadog Agent.

## Behavior

This buildpack will participate if all the following conditions are met

* The `$BP_DATADOG_ENABLED` is set to a truthy value (i.e. `true`, `t`, `1` ignoring case)

The buildpack will do the following for Java applications:

* Contributes the Datadog Java agent to a layer and configures `$JAVA_TOOL_OPTIONS` to use it

The buildpack will do the following for Node.js applications:

* Contributes the Datadog Node.js trace agent to a layer
* Require the trace agent, if it's not present


## Configuration
| Environment Variable | Description
| -------------------- | -----------
| `$BP_DATADOG_ENABLED` | whether to contribute the Datadog trace agent
| `$BPL_DATADOG_DISABLED` | whether to disable the Datadog trace agent (Java only!)

## Usage

Instructions for using the buildpack can be found at the links below:

 - [Running on Docker](docs/run-on-docker.md)

## Bindings

The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`

| Key                   | Value   | Description                                                                                       |
| --------------------- | ------- | ------------------------------------------------------------------------------------------------- |
| `<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>` |

## License

This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
