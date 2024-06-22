package server

type ConfigResponse struct {
	Result Result `json:"result"`
}

type Result struct {
	JobList   []Job `json:"jobList"`
	BotStatus int   `json:"status"`
}

type Job struct {
	Id                string      `json:"id"`
	Name              string      `json:"name"`
	SourceConnectType int         `json:"sourceConnectType"`
	TargetConnectType int         `json:"targetConnectType"`
	JobStatus         int         `json:"status"`
	Cron              string      `json:"cron"`
	Type              int         `json:"type"`
	TriggerType       string      `json:"triggerType"`
	JobConfigList     []JobConfig `json:"jobConfigList"`
	_JobConfigMap     map[string]map[string]JobConfig
}

type JobConfig struct {
	Id      string `json:"id"`
	Group   string `json:"group"`
	Name    string `json:"name"`
	Label   string `json:"label"`
	Content string `json:"content"`
}

// 根据 name 和 group 获取对应 content
func (job *Job) getValueByName(name string, group string) JobConfig {
	var jobConfig JobConfig

	if job._JobConfigMap == nil {
		job._JobConfigMap = make(map[string]map[string]JobConfig)
		for _, v := range job.JobConfigList {
			if job._JobConfigMap[v.Group] == nil {
				job._JobConfigMap[v.Group] = make(map[string]JobConfig)
			}
			job._JobConfigMap[v.Group][v.Name] = v
		}
	}

	if groupMap, ok := job._JobConfigMap[group]; ok {
		if config, exists := groupMap[name]; exists {
			return config
		}
	}

	return jobConfig
}
