package armeria

import "go.uber.org/zap"

func searchForDanglingInstances() {
	var mobsAndItems []interface{}

	mi := Armeria.registry.GetAllFromType(RegistryTypeMobInstance)
	mobsAndItems = append(mobsAndItems, mi...)
	ii := Armeria.registry.GetAllFromType(RegistryTypeItemInstance)
	mobsAndItems = append(mobsAndItems, ii...)

	for _, o := range mobsAndItems {
		obj := o.(ContainerObject)
		container := Armeria.registry.GetObjectContainer(obj.Id())
		if container == nil {
			Armeria.log.Info(
				"found dangling object instance",
				zap.String("uuid", obj.Id()),
			)

			if obj.Type() == ContainerObjectTypeMob {
				obj.(*MobInstance).Parent.DeleteInstance(obj.(*MobInstance))
			} else if obj.Type() == ContainerObjectTypeItem {
				obj.(*ItemInstance).Parent.DeleteInstance(obj.(*ItemInstance))
			}

			Armeria.log.Info(
				"dangling instance deleted",
				zap.String("uuid", obj.Id()),
			)
		}
	}
}
