
<p align="center">
  <img src="docs/header-small.jpg" alt="header image for tbm, a picture of a gopher looking out of a tunnel">
</p>

# Tunnel Boring Machine (tbm)

A simple way to start multiple services like [cloud-sql-proxy](https://github.com/GoogleCloudPlatform/cloud-sql-proxy), [IAP proxy](https://cloud.google.com/iap), [kubefwd](https://kubefwd.com/) or any other command really at once and keep them running in the background.

This is heavily inspired by [foreman](https://github.com/ddollar/foreman) and especially by [mattn/goreman](https://github.com/mattn/goreman) and adds features I missed from these projects:

✅ Enable / disable services  
✅ Use remote template configuration    
✅ Organize by environment (prod/stage/develop)  

## Why & Motivation

If you are using proxies and port forwarding commands for development these ports are often hardcoded in apps like database clients, Makefiles or environment variables.

These informal conventions open the door for accidents where person A maps "database-prod" to port 1234, while person B maps "database-stage" to port 1234. Then sharing a Makefile or curl command could very easily hit the wrong database without being obvious at first.

Having a common and portable configuration in a team where people can enable or disable services they don't need (or don't have permissions to) solves this issue.

## Usage

On first usage run `tbm init` to create a new configuration file at the default location (`~/.tbm.yaml`). You can also init based on a remote configuration file with `tbm init --config-url https://example.com/company-default.yaml`, this will download the remote file instead of creating a default configuration.

After that run `tbm start` to start the services defined by your configuration file to see how everything works in practice.

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

The configuration is stored in a simple YAML file in the user's home directory by default. The default name is `~/.tbm.yaml`.

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
services:
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


## Acknowledgments

A big part of the code, especially around managing processes (start, stop, find, terminate) is mostly taken from
Yasuhiro Matsumoto's [goreman](https://github.com/mattn/goreman) project. The license file for that is included in the `log` package.