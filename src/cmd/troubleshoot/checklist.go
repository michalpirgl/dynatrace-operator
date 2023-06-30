package troubleshoot

import (
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

type Task struct {
	check  *Check
	next   []*Task
	pre    []*Task
	result error
}

func newTask(c *Check) *Task {
	return &Task{
		check:  c,
		next:   make([]*Task, 0),
		pre:    make([]*Task, 0),
		result: errors.Errorf("Pending"),
	}
}

func (t *Task) Run(log logr.Logger) {
	t.result = (*t.check).Do(log)
}

type Checklist struct {
	rootTasks []*Task
	tasks     map[*Check]*Task
}

func (cl *Checklist) Add(c *Check, prerequisites ...*Check) {
	task := newTask(c)

	if len(prerequisites) == 0 {
		cl.rootTasks = append(cl.rootTasks, task)
		cl.tasks[c] = task
		return
	}

	for _, p := range prerequisites {
		preTask := cl.tasks[p]

		// register at pre-requisite task
		preTask.next = append(preTask.next, task)

		// store pre-requisite task
		task.pre = append(task.pre, preTask)
	}
}

func (cl *Checklist) Run() {
	var log logr.Logger
	// start with root tasks
	execList := make([]*Task, 0)
	copy(execList, cl.rootTasks)

	i := 0
	for {
		task := execList[i]
		task.Run(log)

		if task.result == nil {

			// TODO: check for each task if pre-requisistes are already met bevor adding it to execution list
			execList = append(execList, task.next...)
		}
	}
}
