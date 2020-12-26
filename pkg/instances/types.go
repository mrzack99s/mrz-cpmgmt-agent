package instances

var InstanceLists = make(map[string]*Machine)

func CheckExistInstance(instance_id string) bool {
	if _, ok := InstanceLists[instance_id]; ok {
		return true
	}

	return false
}
