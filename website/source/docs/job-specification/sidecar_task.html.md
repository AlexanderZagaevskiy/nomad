---
layout: "docs"
page_title: "sidecar_service Stanza - Job Specification"
sidebar_current: "docs-job-specification-sidecar-task"
description: |-
  The "sidecar_task" stanza allows specifying options for configuring
  the task of the sidecar proxies used in Consul Connect integration
---

# `sidecar_task` Stanza

<table class="table table-bordered table-striped">
  <tr>
    <th width="120">Placement</th>
    <td>
      <code>job -> group -> service -> connect -> **sidecar_task** </code>
    </td>
  </tr>
</table>

The `sidecar_task` stanza allows configuring various options for the proxy
sidecar managed by Nomad for [Consul
Connect](/guides/integrations/consul-connect/index.html) integration such as
resource requirements, kill timeouts and more as defined below. It is valid
only within the context of a [`connect`][connect] stanza.

```hcl
 job "countdash" {
   datacenters = ["dc1"]
   group "api" {
     network {
       mode = "bridge"
     }

     service {
       name = "count-api"
       port = "9001"

       connect {
         sidecar_service {}
         sidecar_task {
            resources {
               cpu = 500
               memory = 1024
            }
         }
       }
     }
     task "web" {
         driver = "docker"
         config {
           image = "test/test:v1"
         }
     }
   }
 }

```

## Default Envoy proxy sidecar

Nomad automatically includes a default Envoy proxy sidecar task whenever a
group service has a [`sidecar_service`][sidecar_service] stanza.

The default sidecar task is equivalent to:

```hcl
   sidecar_task {
     name   = "connect-proxy-<service>"

     driver = "docker"
     config {
       image = "${meta.connect.sidecar_image}
       args  = [
         "-c",
         "${NOMAD_SECRETS_DIR}/envoy_bootstrap.json",
         "-l",
         "${meta.connect.log_level}"
       ]
     }

     logs {
       max_files     = 2
       max_file_size = 2 # MB
     }

     resources {
       cpu    = 250 # MHz
       memory = 128 # MB
     }

     shutdown_delay = "5s"
   }
```

The `meta.connect.sidecar_image` and `meta.connect.log_level` are [*client*
configurable][nodemeta] variables with the following defaults:

- `sidecar_image` - `"envoyproxy/envoy:v1.11.1"` - The official upstream Envoy Docker image.
- `sidecar_log_level` - `"info"` - Debug is useful for debugging Connect related issues.

## `sidecar_task` Parameters
- `name` `(string: "connect-proxy-<service>")` - Name of the task. Defaults to
  including the name of the service it is a proxy for.

- `driver` `(string: "docker")` - Driver used for the sidecar task.

- `user` `(string: nil)` - Determines which user is used to run the task, defaults
   to the same user the Nomad client is being run as.

- `config` `(map: nil)` - Configuration provided to the driver for initialization.

- `env` `(map: nil)` - Map of environment variables used by the driver.

- `resources` <code>([Resources][resources])</code> - Resources needed by the sidecar task.

- `meta` `(map: nil)` - Arbitrary metadata associated with this task that's opaque to Nomad.

- `logs` <code>([Logs][]: nil)</code> - Specifies logging configuration for the
  `stdout` and `stderr` of the task.

- `kill_timeout` `(string: "5s")` - Time between signalling a task that will be
  killed and killing it.

- `shutdown_delay` `(string: "5s")` - Delay between deregistering the task from
  Consul and sending it a signal to shutdown.

- `kill_signal` `(string:SIGINT)` - Kill signal to use for the task, defaults to SIGINT.


## `sidecar_task` Examples
The following example configures resources for the sidecar task and other configuration.

```hcl
   sidecar_task {
     resources {
       cpu = 500
       memory = 1024
     }

     env {
       FOO = "abc"
     }

     shutdown_delay = "5s"
   }

 ```

[connect]: /docs/job-specification/connect.html "Nomad connect Job Specification"
[job]: /docs/job-specification/job.html "Nomad job Job Specification"
[group]: /docs/job-specification/group.html "Nomad group Job Specification"
[task]: /docs/job-specification/task.html "Nomad task Job Specification"
[interpolation]: /docs/runtime/interpolation.html "Nomad interpolation"
[sidecar_service]: /docs/job-specification/sidecar_service.html "Nomad sidecar service Specification"
[resources]: /docs/job-specification/resources.html "Nomad resources Job Specification"
[logs]: /docs/job-specification/logs.html "Nomad logs Job Specification"
[nodemeta]: /docs/configuration/client.html#meta
