package storage

import "time"

type FilterByFn func(task Task) bool

func FilterSlice(tasks []Task, fns ...FilterByFn) ([]Task, []Task) {
	fn := func(task Task) bool {
		for _, f := range fns {
			if !f(task) {
				return false
			}
		}
		return true
	}
	filtered := make([]Task, 0, len(tasks))
	removed := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if fn(task) {
			filtered = append(filtered, task)
		} else {
			removed = append(removed, task)
		}
	}
	return filtered, removed
}

func FilterByPriorityFn(priority Priority) FilterByFn {
	return func(task Task) bool {
		return task.Priority == priority
	}
}

func FilterByStatusFn(status Status) FilterByFn {
	return func(task Task) bool {
		return task.Status == status
	}
}

func FilterByDueFn(date time.Time) FilterByFn {
	return func(task Task) bool {
		return task.Due != nil && task.Due.Before(date)
	}
}
