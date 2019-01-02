package memership

func (node *MemManager) GetViewsInfo(viewMap map[string]*ViewNode) []string {

	var views []string
	for _, item := range viewMap {
		views = append(views, item.String())
	}
	return views
}
func (node *MemManager) RemoveViewsInfo(viewMap map[string]*ViewNode) int {
	lenV := len(node.PartialView)
	for _, item := range viewMap {
		node.removeFromView(item, node.PartialView)
	}

	return lenV
}
