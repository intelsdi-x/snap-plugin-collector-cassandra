# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**

# Snap collector plugin - Cassandra

This plugin collects Cassandra cluster statistics with the [Snap Framework] (http://github.com/intelsdi-x/snap).


1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#openrating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

In order to use this plugin you need the Cassandra node or cluster that you can collect metrics from.

### System Requirements

* [Snap](http://github.com/intelsdi-x/snap)
* Cassandra node/cluster
* [golang 1.6+](https://golang.org/dl/)
* [snap-plugin-utilities](http://github.com/intelsdi-x/snap-plugin-utilities)

Note that Go and plugin utilities are needed only if building the plugin from source.

### Operating systems
All OSs currently supported by Snap:
* Linux/amd64
* Darwin/amd64

### Installation

#### Download plugin binary:

You can also download prebuilt binaries for OS X and Linux (64-bit) at the [releases](https://github.com/intelsdi-x/snap-plugin-collector-cassandra/releases) page

#### To build the plugin binary:
Get the source by running a `go get` to fetch the code:
```
$ go get -d github.com/intelsdi-x/snap-plugin-collector-cassandra
```

Build the plugin by running make within the cloned repo:
```
$ cd $GOPATH/src/github.com/intelsdi-x/snap-plugin-collector-cassandra && make
```
This builds the plugin in `/build/`.


### Configuration and Usage

* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).

* Add Cassandra configuration information to the Task Manifest (see [examples](#examples))

* Load the plugin and create the task

## Documentation 

### Collected Metrics
This plugin has the ability to gather all metrics within the Cassandra org.apache.cassandra.metrics package. View [Metric Types](./METRICS.md) for the full list. They are in the following catalog:

**Cassandra Metric Catalog**
* Cache (Counter, Row, Key)
* Client Request (CASRead, CASWrite, Read, Write)
* ColumnFamily (system_auth, system, ...)
* CommitLog
* Compaction 
* CQL
* Dropped Message
* File Cache
* Keyspace
* Storage
* Thread Pool

The dynamic metric queries are supported. You may view the [sample dynamic metrics](./DYNAMIC_METRICS.md).

### Examples
Example running snap-plugin-collector-cassandra and writing data to a file. 

Run Cassandra from docker image:
```
docker run --detach --name snap-cassandra -p 9042:9042 -p 7199:7199 -p 8082:8082 -p 9160:9160 -d candysmurfhub/cassandra
```

Ensure [Snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
```
$ snapteld -l 1 -t 0 &
```

Download and load Snap plugins:
```
$ snaptel plugin load snap-plugin-collector-cassandra
$ snaptel plugin load snap-plugin-publisher-file
```

See available metrics for your system (this is just part of the list)
```
$snaptel metric list                                
NAMESPACE 												 VERSIONS
/intel/cassandra/node/*/org_apache_cassandra_metrics/type/*/scope/*/name/*/999thPercentile 		 3
/intel/cassandra/node/*/org_apache_cassandra_metrics/type/*/scope/*/name/*/99thPercentile 		 3
/intel/cassandra/node/*/org_apache_cassandra_metrics/type/*/scope/*/name/*/Count 			 3
/intel/cassandra/node/*/org_apache_cassandra_metrics/type/*/scope/*/name/*/FifteenMinuteRate 		 3
/intel/cassandra/node/*/org_apache_cassandra_metrics/type/*/scope/*/name/*/FiveMinuteRate 		 3
/intel/cassandra/node/*/org_apache_cassandra_metrics/type/*/scope/*/name/*/Max
/intel/cassandra/node/*/org_apache_cassandra_metrics/type/*/keyspace/*/scope/*/name/*/50thPercentile 		 3
/intel/cassandra/node/*/org_apache_cassandra_metrics/type/*/keyspace/*/scope/*/name/*/75thPercentile 		 3
```

$snaptel metric list --verbose
```
/intel/cassandra/node/[Node Name]/org_apache_cassandra_metrics/type/[typeValue]/scope/[scopeValue]/name/[nameValue]/FiveMinuteRate 		 		 3 		 double
/intel/cassandra/node/[Node Name]/org_apache_cassandra_metrics/type/[typeValue]/scope/[scopeValue]/name/[nameValue]/Max 			 		 3 		 double
/intel/cassandra/node/[Node Name]/org_apache_cassandra_metrics/type/[typeValue]/scope/[scopeValue]/name/[nameValue]/Mean 			 		 3 		 double
/intel/cassandra/node/[Node Name]/org_apache_cassandra_metrics/type/[typeValue]/scope/[scopeValue]/name/[nameValue]/MeanRate 		 			 3 		 double
/intel/cassandra/node/[Node Name]/org_apache_cassandra_metrics/type/[typeValue]/scope/[scopeValue]/name/[nameValue]/Min 			 		 3 		 double
```

Create a task
```
$ snaptel task create -t cassandra-file.yml
Using task manifest to create task
Task created
ID: 37cd9903-daf6-4e53-b15d-b9082666a830
Name: Task-37cd9903-daf6-4e53-b15d-b9082666a830
State: Running
```
You may view [example tasks](https://github.com/intelsdi-x/snap-plugin-collector-cassandra/blob/master/examples/tasks/).

See the file output (this is just part of the file):
```
$ tail -f collected_cassandra.log
{
  "timestamp": "2016-12-09T17:21:51.537472191-08:00",
  "namespace": "/intel/cassandra/node/192.168.99.100/org.apache.cassandra.metrics/type/Table/keyspace/system_auth/scope/resource_role_permissons_index/name/CoordinatorReadLatency/Max",
  "data": 0,
  "unit": "float64",
  "tags": {
    "plugin_running_on": "egu-mac01.lan"
  },
  "version": 0,
  "last_advertised_time": "2016-12-09T17:21:51.616846048-08:00"
},
{
  "timestamp": "2016-12-09T17:21:51.537473191-08:00",
  "namespace": "/intel/cassandra/node/192.168.99.100/org.apache.cassandra.metrics/type/Table/keyspace/system_auth/scope/resource_role_permissons_index/name/ViewReadTime/Max",
  "data": 0,
  "unit": "float64",
  "tags": {
    "plugin_running_on": "egu-mac01.lan"
  },
  "version": 0,
  "last_advertised_time": "2016-12-09T17:21:51.616846191-08:00"
}
```

![Docker example](https://cloud.githubusercontent.com/assets/13841563/21120807/0427882a-c07e-11e6-944c-9b8cb46c9844.gif)

### Roadmap
This plugin is still in active development. As we launch this plugin, we have a few items in mind for the next few releases:
- [ ] Additional error handling
- [ ] Cluster load and scalability testing

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-cassandra/issues) 
and feel free to submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-cassandra/pulls) that resolves it.

## Community Support
This repository is one of **many** plugins in **Snap**, the open telemetry framework. See the full project at http://github.com/intelsdi-x/snap. To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).


## Acknowledgements

* Author: [@candysmurf](https://github.com/candysmurf/)

