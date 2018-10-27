package UDPSocket

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	PluginName    = "UDPSocket"
	PluginVersion = 1
	PluginType    = plugin.PublisherPluginType
)

type udpPublisher struct {
}

type MetricToPublish struct {
	// Metric creation timestamp
	Timestamp time.Time         `json:"timestamp"`
	Namespace string            `json:"namespace"`
	Data      interface{}       `json:"data"`
	Unit      string            `json:"unit"`
	Tags      map[string]string `json:"tags"`
	Version_  int               `json:"version"`
	// Last advertised time is the last time the snap agent was told about a metric
	LastAdvertisedTime time.Time `json:"last_advertised_time"`
}

// NewUdpPublisher returns an instance of udpPublisher
func NewUdpPublisher() *udpPublisher {
	return &udpPublisher{}
}

func (u *udpPublisher) Publish(contentType string, content []byte, config map[string]ctypes.ConfigValue) error {
	logger := log.New()
	logger.Println("Publisher started")
	var metrics []plugin.MetricType

	switch contentType {
	case plugin.SnapGOBContentType:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.Printf("Error decoding: error=%v content=%v", err, content)
			return fmt.Errorf("Error decoding %v", err)
		}
	default:
		return fmt.Errorf("Unknown content type '%s'", contentType)
	}

	logger.Printf("publishing %v metrics to %v", len(metrics), config)

	conn, err := net.Dial("udp", config["IP_Port"].(ctypes.ConfigValueStr).Value)
	if err != nil {
		return fmt.Errorf("Some connection error %v", err)
	}
	mts := formatMetricTypes(metrics)
	jsonOut, err := json.Marshal(mts)
	if err != nil {
		return fmt.Errorf("Error while marshaling metrics to JSON: %v", err)
	}

	fmt.Fprintf(conn, string(jsonOut))
	conn.Close()

	return nil
}

//Meta returns metadata about the plugin
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(PluginName, PluginVersion, PluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

func (u *udpPublisher) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	r1, err := cpolicy.NewStringRule("IP_Port", true)
	handleErr(err)
	r1.Description = "IP address to send metrics too"

	config.Add(r1)
	cp.Add([]string{""}, config)
	return cp, nil
}

// formatMetricTypes returns metrics in format to be publish as a JSON based on incoming metrics types;
// i.a. namespace is formatted as a single string
func formatMetricTypes(mts []plugin.MetricType) []MetricToPublish {
	var metrics []MetricToPublish
	for _, mt := range mts {
		metrics = append(metrics, MetricToPublish{
			Timestamp:          mt.Timestamp(),
			Namespace:          mt.Namespace().String(),
			Data:               mt.Data(),
			Unit:               mt.Unit(),
			Tags:               mt.Tags(),
			Version_:           mt.Version(),
			LastAdvertisedTime: mt.LastAdvertisedTime(),
		})
	}
	return metrics
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
