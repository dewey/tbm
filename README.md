
<p align="center">
  <img height="500" src="docs/header.jpg" alt="header image for tbm, a picture of a gopher looking out of a tunnel">
</p>

# Tunnel Boring Machine (tbm)

A simple way to run multiple
proxies like [cloud-sql-proxy](https://github.com/GoogleCloudPlatform/cloud-sql-proxy), [IAP proxy](https://cloud.google.com/iap), [kubefwd](https://kubefwd.com/) in the background.

It's easy to enable or disable services without messing with the port
configuration. This makes it easy to have a standardized and portable configuration file in a team where people can
enable / disable
services
they don't need (or don't have permissions to) while still keeping a common port configuration. Having a standardized
port configuration is useful to
prevent accidents where person A maps "database-prod" to port 1234, while person B maps "database-stage" to port 1234
and
them sharing a curl command / script which will then hit the wrong database.

This is heavily inspired
by [foreman](https://github.com/ddollar/foreman) and especially by [mattn/goreman](https://github.com/mattn/goreman).

## Usage

On first usage you should run `tbm init` to create a new configuration file at the default location (`~/.tbm.yaml`).

After that you can run `tbm start` to start all the default services to see how everything works in practice. The next
step is to adapt the configuration file and add all your services.

Run `tbm help` to get an overview over the available commands.

![Screenshot of a terminal with tbm running two ping commands concurrently](/docs/screenshot.png "Example of tbm running two ping commands")

## Install
You can either install from source or use Homebrew on macOS.

```
brew tap dewey/tbm https://github.com/dewey/tbm
brew install tbm
```

### Configuration

#### Location

The configuration is stored in a simple YAML file in the user's home directory by default. The default name
is `~/.tbm.yaml`.

A custom configuration location can be defined with the config flag `tbm start --config ~/myconfigs/tbm.yaml`. More
information about this command can be found with `tbm start --help`.

#### Configuration file

The configuration file can contain the following keys.

- Service name: The top level key is the name of the service. In this example it's `cloudsql-db`
    - Command: The command that should be executed when tbm starts
    - Environment: This is a free form string which can be used to differentiate services that are named the same across
      environments
    - Enable: You can enable or disable a service, this is also useful if you have a company-wide configuration file and
      you only want to enable the services you have access to
    - Variables: A list of variables, these will be injected into the `command` if the name of the field maps to the
      placeholder name in the command.

Example file with two services defined:

```yaml
cloudsql-db:
  command: cloud_sql_proxy -instances=europe-west1:prod-db=tcp:0.0.0.0:{{.port}}
  environment: prod
  enable: true
  variables:
    - port: 10001
cloudsql-db-replica:
  command: cloud_sql_proxy -instances=europe-west1:prod-db-replica=tcp:0.0.0.0:{{.port}}
  environment: prod
  enable: true
  variables:
    - port: 10002
```

## Motivation

In my daily life as a developer there's a number of services that I have to keep running (CloudSQL proxy to access
databases running in Google Cloud Platform), a company specific "devproxy" which tunnels specific traffic through an IAP
proxy that provides authentication in front of services running in Kubernetes and
various [port forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/)
commands to access services in Kubernetes directly.

There's other projects that can take care of this too, but I haven't found one where I can easily define a bunch of
services, enable, disable them and set a pre-defined port where they are mapped to localhost. This project covers for
this specific use case.

## Acknowledgments

A big part of the code, especially around managing processes (start, stop, find, terminate) is mostly taken from
Yasuhiro MATSUMOTO's [goreman](https://github.com/mattn/goreman) project. The license file for that is included in the `log` package.