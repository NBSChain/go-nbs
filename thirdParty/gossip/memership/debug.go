package memership

func (node *MemManager) outputViewTask(task *gossipTask) error {
	var views []string
	for _, item := range node.PartialView {
		views = append(views, item.String())
	}

	task.result <- views
	return nil
}

func (node *MemManager) inputViewTask(task *gossipTask) error {

	var views []string
	for _, item := range node.InputView {
		views = append(views, item.String())
	}

	task.result <- views
	return nil
}

func (node *MemManager) GetViewsInfo(typ int) []string {

	task := &gossipTask{
		taskType: typ,
	}
	task.result = make(chan interface{})
	node.taskQueue <- task
	views := (<-task.result).([]string)
	return views
}

func (node *MemManager) removeOV(task *gossipTask) error {
	lenV := len(node.PartialView)
	for _, item := range node.PartialView {
		node.removeFromView(item, node.PartialView)
	}
	task.result <- lenV
	return nil
}

func (node *MemManager) removeIV(task *gossipTask) error {
	lenV := len(node.InputView)
	for _, item := range node.InputView {
		node.removeFromView(item, node.InputView)
	}
	task.result <- lenV
	return nil
}

func (node *MemManager) RemoveViewsInfo(typ int) int {
	task := &gossipTask{
		taskType: typ,
	}
	task.result = make(chan interface{})
	node.taskQueue <- task
	no := (<-task.result).(int)

	return no
}
