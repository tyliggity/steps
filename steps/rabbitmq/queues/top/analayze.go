package main

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/public-steps/rabbitmq/queues/base"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
)

var CheckValuesDescription = map[string]string{
	"memory":         "Top memory used queues",
	"messages":       "Top most messages queues",
	"message_bytes":  "Top messages bytes used queues",
	"messages_ready": "Top queues with unconsumed messages",
}

type MaxQueuesEntry struct {
	QueueName string
	Value     float64
}

type TopQueues struct {
	MaxQueues         map[string][]MaxQueuesEntry
	FullDetailsQueues map[string]interface{} // Map of queue name to it's full JSON details
}

type DescriptiveMaxQueues struct {
	Description string           `json:"description"`
	MaxQueues   []MaxQueuesEntry `json:"max_queues"`
}

type Args struct {
	TotalTopQueues int    `env:"TOTAL_TOP_QUEUES" envDefault:"5"`
	ExcludeRe      string `env:"EXCLUDE_RE"`
}

type RMQQueuesTop struct {
	*base.RMQQueues
	args      *Args
	excludeRe *regexp.Regexp
}

func (r *RMQQueuesTop) Init() error {
	rmqQueue, err := base.NewRmqQueues()
	if err != nil {
		return err
	}
	r.RMQQueues = rmqQueue

	args := &Args{}
	if err := envconf.Parse(args); err != nil {
		return err
	}

	r.args = args
	if r.args.ExcludeRe != "" {
		r.excludeRe, err = regexp.Compile(r.args.ExcludeRe)
		if err != nil {
			return fmt.Errorf("compile exclude re: %w", err)
		}
	}

	return nil
}

func (r *RMQQueuesTop) sortEntries(entries []MaxQueuesEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})
}

func (r *RMQQueuesTop) insertToMaxQueue(maxQueues map[string][]MaxQueuesEntry, val interface{}, key, queueName string) bool {
	// Currently supporting only float comparison
	valFloat, ok := val.(float64)
	if !ok {
		log.Debugln("value in key %s is not float (it's %T) in queue: %s", key, val, queueName)
		return false
	}

	entries, ok := maxQueues[key]
	if !ok {
		entries = make([]MaxQueuesEntry, 0)
	}

	// If we are not in the max top queues limit or if the current value is greater from the smallest number in the top queues,
	// add it to the top queues
	if len(entries) < r.args.TotalTopQueues || valFloat > entries[len(entries)-1].Value {
		entry := MaxQueuesEntry{
			QueueName: queueName,
			Value:     valFloat,
		}
		if len(entries) == r.args.TotalTopQueues {
			// If we are already reached the max total queues, remove the last one
			entries = entries[:len(entries)-1]
		}
		entries = append(entries, entry)

		r.sortEntries(entries)
		maxQueues[key] = entries

		return true
	}

	return false
}

func (r *RMQQueuesTop) getTopQueues(queues []*gabs.Container) TopQueues {
	maxQueues := make(map[string][]MaxQueuesEntry)
	fullDetailsQueues := make(map[string]interface{})

	for _, queue := range queues {
		name, ok := queue.S("name").Data().(string)
		if !ok || name == "" {
			log.Debugln("Found entry with empty name: %s", queue)
			continue
		}
		if r.excludeRe != nil && r.excludeRe.Match([]byte(name)) {
			log.Debugln("Skipping %s as it matching the exclude regex", r.excludeRe)
			continue
		}

		for key := range CheckValuesDescription {
			if !queue.ExistsP(key) {
				log.Debugln("key %s not exists in queue %s", queue)
				continue
			}

			val := queue.Path(key).Data()
			if r.insertToMaxQueue(maxQueues, val, key, name) {
				fullDetailsQueues[name] = queue
			}
		}
	}

	return TopQueues{
		MaxQueues:         maxQueues,
		FullDetailsQueues: fullDetailsQueues,
	}
}

func (r *RMQQueuesTop) Run() (int, []byte, error) {
	output, err := r.Get()
	if err != nil {
		return 1, nil, err
	}

	gc, err := gabs.ParseJSON(output)
	if err != nil {
		return 1, output, fmt.Errorf("parse output as JSON: %w", err)
	}

	if len(gc.Children()) == 0 {
		return 0, output, nil
	}

	topQueues := r.getTopQueues(gc.Children())
	// Delete the full details queues that is not present in the final result
	finalFullDetailsQueues := make(map[string]interface{})
	descriptiveMaxQueues := make(map[string]DescriptiveMaxQueues)
	for key, description := range CheckValuesDescription {
		descriptiveMaxQueues[key] = DescriptiveMaxQueues{
			Description: description,
			MaxQueues:   topQueues.MaxQueues[key],
		}

		// Save the queue details in the final queues details dict
		for _, entry := range topQueues.MaxQueues[key] {
			finalFullDetailsQueues[entry.QueueName] = topQueues.FullDetailsQueues[entry.QueueName]
		}
	}

	retGc := gabs.New()
	retGc.Set(finalFullDetailsQueues, "full_details_queues")
	retGc.Set(descriptiveMaxQueues, "max_queues")

	return 0, retGc.Bytes(), nil
}

func main() {
	step.Run(&RMQQueuesTop{})
}
