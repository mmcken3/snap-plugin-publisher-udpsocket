# Snap publisher plugin : UDP Socket

This plugin supports pushing metrics to a socket at a client IP over UDP.

It's used with the [Snap framework](http://github.com/intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Examples](#examples)
3. [Community Support](#community-support)
4. [Contributing and Roadmap](#contributing-and-roadmap)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Under Construction
Travis CI was working at first, but this project has not been used or touched much since the fall of 2017. There may be build items and updates needed to get this working properly again.

## Getting Started

### System Requirements

* [golang 1.6+](https://golang.org/dl/) (needed only for building)

### Operating Systems
Currently supported for Linux and Darwin systems.
* Linux/amd64
* Darwin/amd64

### Installation

#### Building the plugin from binary:
Fork https://github.com/mmcken3/snap-plugin-publisher-udpsocket

Clone repo into `$GOPATH/src/github.com/mmcken3/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-publisher-udpsocket
```

Run make within the cloned repo to build the plugin:
```
$ make
```
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

## Documentation
To use this plugin you have to specify a config with the IP Address and Port you want to write to.

```
# JSON
"config": {
    "IP_Port": "127.0.0.1:5100"
}

# YAML
config: 
    IP_Port: "127.0.0.1:5100"

```

The plugin will write all metrics serialized as JSON to the specified IP and Port. An example of this output is below:

```json
[
  {
    "timestamp": "2016-07-25T11:27:59.795513984+02:00",
    "namespace": "/intel/mock/host0/baz",
    "data": 86,
    "unit": "",
    "tags": {
      "plugin_running_on": "my-machine"
    },
    "version": 0,
    "last_advertised_time": "0001-01-01T00:00:00Z"
  },
  {
    "timestamp": "2016-07-25T11:27:59.795514856+02:00",
    "namespace": "/intel/mock/host1/baz",
    "data": 70,
    "unit": "",
    "tags": {
      "plugin_running_on": "my-machine"
    },
    "version": 0,
    "last_advertised_time": "0001-01-01T00:00:00Z"
  },
  {
    "timestamp": "2016-07-25T11:27:59.795548989+02:00",
    "namespace": "/intel/mock/bar",
    "data": 82,
    "unit": "",
    "tags": {
      "plugin_running_on": "my-machine"
    },
    "version": 0,
    "last_advertised_time": "2016-07-25T11:27:21.852064032+02:00"
  },
  {
    "timestamp": "2016-07-25T11:27:59.795549268+02:00",
    "namespace": "/intel/mock/foo",
    "data": 72,
    "unit": "",
    "tags": {
      "plugin_running_on": "my-machine"
    },
    "version": 0,
    "last_advertised_time": "2016-07-25T11:27:21.852063228+02:00"
  }
]
```

### Examples

Example of running [psutil collector plugin](https://github.com/intelsdi-x/snap-plugin-collector-psutil), and writing data as a JSON to a UDP Socket at a given Port and IP:

Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

Ensure [snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `service snap-telemetry start`
* systemd: `systemctl start snap-telemetry`
* command line: `sudo snapteld -l 1 -t 0 &`


Download and load Snap plugins:
Ensure that the udpsocket publisher is intalled like above.
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-psutil/latest/linux/x86_64/snap-plugin-collector-psutil
$ snaptel plugin load snap-plugin-publisher-udpsocket
$ snaptel plugin load snap-plugin-collector-psutil
```

Create a [task manifest](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md) (see [exemplary tasks](examples/tasks/)),
Or like this example `psutil-udp-publish.json` with following content:
```json
{
  "version": 1,
  "schedule": {
    "type": "simple",
    "interval": "1s"
  },
  "max-failures": 10,
  "workflow": {
    "collect": {
      "metrics": {
        "/intel/psutil/load/load1": {},
        "/intel/psutil/load/load15": {},
        "/intel/psutil/load/load5": {},
        "/intel/psutil/vm/available": {},
        "/intel/psutil/vm/free": {},
        "/intel/psutil/vm/used": {}
      },
      "publish": [
        {
          "plugin_name": "udpsocket",
          "config": {
            "IP_Port": "127.0.0.1:5100"
          }
        }
      ]
    }
  }
}

```

Create a task:
```
$ snaptel task create -t psutil-udp-publish.json
```

See JSON data by listening on the given port at the given IP address.

Example go udp-server that will listen and print results until killed.

####  udp-server.go
```
package main
import (
    "fmt"
    "net"
)

func main() {
    p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 5100,
        IP: net.ParseIP("127.0.0.1"),
    }
    ser, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Connection Error  %v\n", err)
        return
    }
    for {
        _,remoteaddr,err := ser.ReadFromUDP(p)
        fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
        if err !=  nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
    }
}
```

To stop previously created task:
```
$ snaptel task stop <task_id>
```

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing and Roadmap
There isn't a current roadmap for this plugin. As we launch this plugin, we do not have any outstanding requirements for the next release. Any contribution is appreciated. If you have a feature request or contribution, please add it as an issue and/or submit a pull request. 

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
Original code from the [Snap](http://github.com/intelsdi-x/snap) repo.
There is lots of original code from the [File Publisher](http://github.com/intelsdi-x/snap-plugin-publisher-file) repo. This plugin formats data in the same way as file publisher, it does everything the same way but publishes to the network instead of writing a file.
