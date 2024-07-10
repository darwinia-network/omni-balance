package daemons

import (
	"context"
	"omni-balance/utils"
	"omni-balance/utils/configs"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type TaskFunc func(ctx context.Context, conf configs.Config) error

var (
	tasks = make(map[string]Task)
	m     sync.Mutex
)

type Task struct {
	Name            string
	Description     string
	TaskFunc        TaskFunc
	DefaultInterval time.Duration
	RunOnStart      bool
}

func Help() string {
	var result = strings.Builder{}

	result.WriteString("Available tasks:\n")
	var names []string
	for index := range tasks {
		result.WriteString(" Name: ")
		result.WriteString(tasks[index].Name)
		result.WriteString("\n")
		result.WriteString(" Description: ")
		result.WriteString(tasks[index].Description)
		result.WriteString("\n")
		result.WriteString(" Default Run Interval: ")
		result.WriteString(tasks[index].DefaultInterval.String())
		names = append(names, tasks[index].Name)
		result.WriteString("\n\n")
	}
	result.WriteString("\n\nYou can override the Run Interval time for these tasks in your configuration. For example:\n")
	result.WriteString("task_interval:\n")
	times := time.Second * 10
	for _, v := range names {
		result.WriteString("    ")
		result.WriteString(v)
		result.WriteString(": ")
		result.WriteString(times.String())
		result.WriteString("\n")
		times = times * 55
	}
	return result.String()
}

func GetTaskConfig() []Task {
	var result = make([]Task, 0)
	for _, v := range tasks {
		result = append(result, v)
	}
	return result
}

func RegisterIntervalTask(task Task) {
	m.Lock()
	defer m.Unlock()
	tasks[task.Name] = task
}

func runForever(ctx context.Context, conf configs.Config, task Task) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("task %s failed, err: %v, will retry after 2s", task.Name, err)
			debug.PrintStack()
			time.Sleep(time.Second * 2)
			runForever(ctx, conf, task)
		}
	}()

	interval := conf.GetTaskInterval(task.Name, task.DefaultInterval)
	var t = time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			logrus.Infof("task %s stopped", task.Name)
			return
		case <-t.C:
			if !utils.IsFinishedInit() {
				logrus.Infof("task %s waiting for init finished", task.Name)
				continue
			}
			if err := task.TaskFunc(ctx, conf); err != nil {
				logrus.Errorf("task %s failed, err: %v", task.Name, err)
			}
			t.Reset(interval)
		}
	}
}

func Run(ctx context.Context, conf configs.Config) error {
	for index := range tasks {
		if tasks[index].RunOnStart {
			logrus.Debugf("task %s run on start, wait for the task finished", tasks[index].Name)
			if err := tasks[index].TaskFunc(ctx, conf); err != nil {
				logrus.Errorf("task %s failed, err: %v", tasks[index].Name, err)
				continue
			}
			logrus.Debugf("task %s run on start finished", tasks[index].Name)
			continue
		}
	}
	for index := range tasks {
		logrus.Debugf("task %s run in background", tasks[index].Name)
		go runForever(ctx, conf, tasks[index])
	}
	return nil
}
