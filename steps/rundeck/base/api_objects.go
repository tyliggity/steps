package base

import "time"

type execution struct {
	Description     string   `json:"description,omitempty"`
	Href            string   `json:"href,omitempty"`
	Permalink       string   `json:"permalink,omitempty"`
	Status          string   `json:"status,omitempty"`
	Project         string   `json:"project,omitempty"`
	ExecutionType   string   `json:"execution_type,omitempty"`
	User            string   `json:"user,omitempty"`
	DateStarted     *Date    `json:"date_started,omitempty"`
	DateEnded       *Date    `json:"date_ended,omitempty"`
	Job             *Job     `json:"job,omitempty"`
	Argstring       string   `json:"argstring,omitempty"`
	SuccessfulNodes []string `json:"successful_nodes,omitempty"`
	FailedNodes     []string `json:"failed_nodes,omitempty"`
	ServerUUID      string   `json:"server_uuid,omitempty"`
	CustomStatus    string   `json:"custom_status,omitempty"`
}

type ExecutionIDFromString struct {
	ID int `json:"id,string"`
	execution
}

type ExecutionIDFromInt struct {
	ID int `json:"id"`
	execution
}

type Date struct {
	Unixtime int64     `json:"unixtime,omitempty"`
	Date     time.Time `json:"date,omitempty"`
}

type Job struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name,omitempty"`
	Group           string                  `json:"group,omitempty"`
	Project         string                  `json:"project,omitempty"`
	Description     string                  `json:"description,omitempty"`
	Options         *map[string]interface{} `json:"options,omitempty"`
	Href            string                  `json:"href,omitempty"`
	Permalink       string                  `json:"permalink,omitempty"`
	AverageDuration int64                   `json:"average_duration,omitempty"`
	Scheduled       bool                    `json:"scheduled,omitempty"`
	ScheduleEnabled bool                    `json:"schedule_enabled,omitempty"`
	Enabled         bool                    `json:"enabled,omitempty"`
}

type Abort struct {
	Status string `json:"status,omitempty"`
	Reason string `json:"reason,omitempty"`
}
