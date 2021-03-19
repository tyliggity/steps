package get

import (
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	base2 "github.com/stackpulse/public-steps/kubectl/base"
)

const (
	LabelsJsonPath      = "metadata.labels"
	NodeRoleLabelPrefix = "node-role.kubernetes.io/"
)

const NodeReadyJsonKey = "ready"

var parsingConfiguration = map[string]*base2.JsonParseConfig{
	"name":             {ParseFunc: base2.JsonPathStringParser, Args: []string{"metadata.name"}},
	NodeReadyJsonKey:   {ParseFunc: base2.JsonPathSearchInObjectArray, Args: []string{"status.conditions", "type", "Ready", "status", "Unknown"}},
	"roles":            {ParseFunc: ExtractRolesParser},
	"age":              {ParseFunc: base2.JsonPathDurationFromDate, Args: []string{"metadata.creationTimestamp", time.RFC3339}},
	"version":          {ParseFunc: base2.JsonPathStringParser, Args: []string{"status.nodeInfo.kubeletVersion"}},
	"os_image":         {ParseFunc: base2.JsonPathStringParser, Args: []string{"status.nodeInfo.osImage"}},
	"kernelVersion":    {ParseFunc: base2.JsonPathStringParser, Args: []string{"status.nodeInfo.kernelVersion"}},
	"containerRuntime": {ParseFunc: base2.JsonPathStringParser, Args: []string{"status.nodeInfo.containerRuntimeVersion"}},
	"architecture":     {ParseFunc: base2.JsonPathStringParser, Args: []string{"status.nodeInfo.architecture"}},
	"internalIP":       {ParseFunc: base2.JsonPathSearchInObjectArray, Args: []string{"status.addresses", "type", "InternalIP", "address", "Unknown"}},
	"externalIP":       {ParseFunc: base2.JsonPathSearchInObjectArray, Args: []string{"status.addresses", "type", "ExternalIP", "address", "Unknown"}},
	"labels":           {ParseFunc: base2.JsonPathObjectKeys, Args: []string{LabelsJsonPath}},
}

func ExtractRolesParser(args []string, item *gabs.Container) (interface{}, error) {
	roles := make([]string, 0)

	for label, _ := range item.Path(LabelsJsonPath).ChildrenMap() {
		if strings.HasPrefix(label, NodeRoleLabelPrefix) {
			if len(label) == len(NodeRoleLabelPrefix) {
				// Dealing with edge case which the label is exactly the prefix
				continue
			}
			roles = append(roles, label[len(NodeRoleLabelPrefix):])
		}
	}
	return roles, nil
}
