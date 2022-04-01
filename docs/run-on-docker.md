# Deploy with Datadog on Docker

These instructions will walk through building a Java application with Datadog's trace agent and running it on Docker.

## Prerequisites

The following are prerequisites to building applications with Paketo Buildpacks.

1. Install Docker by following [this guide](https://docs.docker.com/get-docker/).
2. Install the [Docker Compose plugin](https://docs.docker.com/compose/cli-command/#installing-compose-v2).
3. Install the pack CLI by following [this guide](https://buildpacks.io/docs/install-pack/).
4. Download the [Paketo Samples](https://github.com/paketo-buildpacks/samples) by running `git clone https://github.com/paketo-buildpacks/samples`. Alternatively, bring your own Java application, however, it's recommended that for your first build you use a sample application.
5. Run `cd samples/java/maven` as we'll use the Maven sample. The process should work with the `samples/java/gradle` application as well, if you prefer using Gradle.

## Build an OCI Image

The first step is to build an OCI image from your application. This is the standard build process, with the one additional flag required by this buildpack `BP_DATADOG_ENABLED`. Setting this to true will cause the Datadog Trace agent to be included in the image and setting it to false will skip adding the agent.

```bash
pack build apps/maven -e BP_DATADOG_ENABLED=true
``` 

## Prepare to Run

Reporting metrics and traces to Datadog from your Java application requires a side-car agent. We will use the [standard Datadog instructions for Docker](https://docs.datadoghq.com/tracing/setup_overview/setup/java/?tab=containers#configure-the-datadog-agent-for-apm).

The process requires a few additional files to be created. We'll do that now.

### JMX Config

First, we'll create a JMX configuration. This tells the Agent how to talk to our Java application over JMX for fetching metrics.

1. From the root of our project, run `mkdir conf.d/jmx.d/`. This format should be done exactly as listed here, the Datadog agent expects this way. 
2. Then `cd conf.d/jmx.d` and create the file `conf.yaml`. In it, put the following contents.

    ```yaml
    init_config:
      is_jmx: true
      collect_default_metrics: true

    instances:
      - host: localhost
        port: 8000
        name: jmx_instance
    ```

### Docker Compose

We'll use Docker Compose to deploy the container image. In the root of the project, create a file `docker-compose.yml`. In it, put the following contents.

```yaml
version: "3.9"
services:
  maven:
    image: apps/maven
    ports:
      - "8080:8080"
    environment:
      BPL_JMX_ENABLED: true
      BPL_JMX_PORT: 8000
      DD_ENV: paketo-prod
      DD_SERVICE: paketo-maven
      DD_VERSION: 1.0
      DD_PROFILING_ENABLED: true
      DD_TRACE_SAMPLE_RATE: 0.1
  datadog:
    image: gcr.io/datadoghq/agent:latest-jmx
    environment:
      DD_API_KEY: <insert-dd-api-key>
      DD_APM_ENABLED: true
      DD_APM_NON_LOCAL_TRAFFIC: true
    network_mode: service:maven
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro 
      - /proc/:/host/proc/:ro
      - /sys/fs/cgroup/:/host/sys/fs/cgroup:ro
      - ./conf.d:/conf.d
```

Notes on this file:

1. Insert your Datadog API key where it says `<insert-dd-api-key>`. 
2. Adjust the `DD_*` env variables as you see fit. Refer to the [Datadog documentation](https://docs.datadoghq.com/tracing/setup_overview/setup/java/?tab=containers#configuration) for a list of available settins and their meanings.

Besides configuration, the important thing being done in the Docker compose file is that we're instructing Docker to put the Datadog agent into the same network namespace as our application's container, which allows the agent to talk to the JVM through JMX over localhost. This is critical or we'd have to make JMX remotely available which you should not do.

## Run the Application

Running the application is simple. Execute `docker compose up` from the `samples/java/maven` directory (or the root of your project). 

It should start two containers, `maven-datadog-1` and `maven-maven-1`. The former is the Datadog agent and the latter is our App.

## Validation

You can use the following instructions to validate that everything is working as expected.

1. You should see this line in the logs, which confirms that the Datadog agent is able to connect over JMX.

    ```
    maven-datadog-1 | 2022-03-10 18:27:26 GMT | JMX | INFO | Instance | Connected to JMX Server at jmx_instance
    ```

2. You should see this line which confirms that the Datadog agent was able to pull metrics over JMX.

    ```
    maven-datadog-1 | 2022-03-10 18:27:26 GMT | JMX | INFO | Reporter | Instance jmx_instance is sending 27 metrics to the metrics reporter during collection #1
    ```

3. You should see these two lines which indicate the app is up and running and listening for traffic:

    ```
    maven-maven-1 | 2022-03-10 18:27:23.402 INFO 1 --- [ main] o.s.b.web.embedded.netty.NettyWebServer : Netty started on port 8080
    maven-maven-1 | 2022-03-10 18:27:23.425 INFO 1 --- [ main] io.paketo.demo.DemoApplication : Started DemoApplication in 3.644 seconds (JVM running for 5.706)
    ```

4. You should be able to `curl http://localhost:8080/actuator/health` and it should respond with `{"status":"UP"}`. (If you have a remote docker daemon, you'll need to change localhost to the IP of your docker daemon.

5. Look in the Datadog dashboard. You should see metrics being reported.

6. Run `watch -n 1 curl -s http://localhost:8080/actuator/health` and let it run for a few minutes. Refresh your Datadog dashboard and look at APM -> Traces. You should see traces for the requests being generated.
