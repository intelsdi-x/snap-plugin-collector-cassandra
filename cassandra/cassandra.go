/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cassandra

import (
	"reflect"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
)

// const defines constant varaibles
const (
	// Name of plugin
	Name = "cassandra"
	// Version of plugin
	Version = 4
	// Type of plugin
	PluginType = plugin.CollectorPluginType

	// Timeout duration
	DefaultTimeout = 5 * time.Second

	CassURL    = "url"
	Port       = "port"
	Hostname   = "hostname"
	InvalidURL = "Invalid URL in Global configuration"
	NoHostname = "No hostname define in Global configuration"
)

// Meta returns the snap plug.PluginMeta type
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(Name, Version, PluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// NewCassandraCollector returns a new instance of Cassandra struct
func NewCassandraCollector() *Cassandra {
	return &Cassandra{}
}

// Cassandra struct
type Cassandra struct {
	client *CassClient
}

// CollectMetrics collects metrics from Cassandra through JMX
func (p *Cassandra) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	metrics := []plugin.MetricType{}

	if p.client == nil {
		err := p.loadMetricAPI(mts[0].Config())
		if err != nil {
			return nil, err
		}
	}

	for _, m := range mts {
		results := []nodeData{}
		search := strings.Split(replaceUnderscoreToDot(strings.TrimLeft(m.Namespace().String(), "/")), "/")
		if len(search) > 3 {
			p.client.Root.Get(p.client.client.GetUrl(), search[4:], 0, &results)
		}

		for _, result := range results {
			ns := append([]string{"intel", "cassandra", "node", p.client.host}, strings.Split(result.Path, Slash)...)
			nss, tags := processTagNamespace(ns, m.Tags())

			metrics = append(metrics, plugin.MetricType{
				Namespace_: core.NewNamespace(nss...),
				Timestamp_: time.Now(),
				Data_:      result.Data,
				Unit_:      reflect.TypeOf(result.Data).String(),
				Tags_:      tags,
			})
		}
	}
	return metrics, nil
}

// processTagNamespace creates tags from the giving namespace
// and returns a new namespace with tags removed.
func processTagNamespace(ns []string, mp map[string]string) ([]string, map[string]string) {
	m := map[string]string{}

	// Create the node and the node name tag.
	m[ns[2]] = ns[3]

	// Create other tags except the last one and the type(ns[5],ns[6]).
	for i := 7; i < len(ns)-1; i += 2 {
		m[ns[i]] = ns[i+1]
	}

	// Create a new namespace with tags removed.
	nss := []string{}
	nss = append(nss, ns[0], ns[1], ns[5], ns[6], ns[len(ns)-1])

	// Copy over Snap tags
	if mp != nil {
		for k, v := range mp {
			m[k] = v
		}
	}
	return nss, m
}

// GetMetricTypes returns the metric types exposed by Cassandra
func (p *Cassandra) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	return NewEmptyCassClient().getMetricType(cfg)
}

// GetConfigPolicy returns a ConfigPolicy
func (p *Cassandra) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

// loadMetricAPI returns the root node
func (p *Cassandra) loadMetricAPI(config *cdata.ConfigDataNode) error {
	var err error
	// inits CassClient
	p.client, err = initClient(plugin.ConfigType{ConfigDataNode: config})
	if err != nil {
		return err
	}

	// reads the root metric node from the memory
	nod, err := readMetricAPI()
	if err != nil {
		err = p.client.buidMetricAPI()
		if err != nil {
			return err
		}
	} else {
		p.client.Root = nod
	}
	return nil
}
