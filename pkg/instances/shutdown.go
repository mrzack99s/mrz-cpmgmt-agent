package instances

func ShutdownInstance(instance_id string) (*Machine, string) {

	if CheckExistInstance(instance_id) {
		machine := InstanceLists[instance_id]
		if machine.Status == "up" {
			machine.Status = "down"
			if err := machine.MachineState.Shutdown(machine.MachineState.Logger().Context); err != nil {
			}
			return machine, "success"
		}
		return nil, "machine is down"
	} else {
		return nil, "not have instance"
	}
}
