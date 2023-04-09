[![microboxdesktop logo](http://microbox.rocks/assets/readme-headers/microboxdesktop.png)](http://microbox.cloud/open-source#microboxdesktop)
[![Build Status](https://github.com/mu-box/microbox/actions/workflows/ci.yaml/badge.svg)](https://github.com/mu-box/microbox/actions)

## Microbox

[Microbox](https://microbox.cloud/) automates the creation of isolated, repeatable environments for local and production applications. When developing locally, Microbox provisions your app's infrastructure inside of a virtual machine (VM) and mounts your local codebase into the VM. Any changes made to your codebase are reflected inside the virtual environment.

Once code is built and tested locally, Microbox provisions and deploys an identical infrastructure on a production platform.


## How It Works

Microbox uses [Virtual Box](http://virtualbox.org) and [Docker](https://www.docker.com/) to create virtual development environments on your local machine. App configuration is handled in the [boxfile.yml](https://docs.microbox.cloud/boxfile/), a small yaml config file used to provision and configure your apps' environments both locally and in production.


## Why Microbox?

Microbox allows you to stop configuring environments and just code. It guarantees that any project you start will work the same for anyone else collaborating on the project. When it's time to launch the project, you'll know that your production app will work, because it already works locally.


### Installation

By using the [Microbox installer](https://microbox.cloud/download). *(Recommended)* .The installer includes all required dependencies (Virtual Box & Docker).


### Usage

```
Usage:
  microbox [flags]
  microbox [command]

Available Commands:
  configure     Configure Microbox.
  run           Start your local development environment.
  build-runtime Build your app's runtime.
  compile-app   Compile your application.
  deploy        Deploy your application to a live remote or a dry-run environment.
  console       Open an interactive console inside a component.
  remote        Manage application remotes.
  status        Display the status of your Microbox VM & apps.
  login         Authenticate your microbox client with your microbox.cloud account.
  logout        Remove your microbox.cloud api token from your local microbox client.
  clean         Clean out any apps that no longer exist.
  info          Show information about the specified environment.
  tunnel        Create a secure tunnel between your local machine & a live component.
  implode       Remove all Microbox-created containers, files, & data.
  destroy       Destroy the current project and remove it from Microbox.
  start         Start the Microbox virtual machine.
  stop          Stop the Microbox virtual machine.
  update-images Updates docker images.
  evar          Manage environment variables.
  dns           Manage dns aliases for local applications.
  log           Streams application logs.
  version       Show the current Microbox version.
  server        Start a dedicated microbox server

Flags:
      --debug     In the event of a failure, drop into debug context
  -h, --help      help for microbox
  -t, --trace     Increases display output and sets level to trace
  -v, --verbose   Increases display output and sets level to debug

Use "microbox [command] --help" for more information about a command.
```


### Documentation

- Microbox documentation is available at [docs.microbox.cloud](https://docs.microbox.cloud/).
- Guides for popular languages, frameworks and services are available at [guides.microbox.cloud](http://guides.microbox.cloud).


## Contributing

Contributing to Microbox is easy. Just follow these [contribution guidelines](https://docs.microbox.cloud/contributing/).
Microbox uses [govendor](https://github.com/kardianos/govendor#the-vendor-tool-for-go) to vendor dependencies. Use `govendor sync` to restore dependencies.


### Contact

For help using Microbox or if you have any questions/suggestions, please reach out to help@microbox.cloud or find us on [slack](https://slack.microbox.rocks/). You can also [create a new issue on this project](https://github.com/mu-box/microbox/issues/new).

[![microbox logo](http://microbox.rocks/assets/open-src/microbox-open-src.png)](https://microbox.cloud/open-source/)
