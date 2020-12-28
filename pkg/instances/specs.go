package instances

type InstanceSpecs struct {
	Vcpu       int64
	HtEnabled  bool
	MemSizeMib int64
}

type InstanceTemplateStruct struct {
	Standard_T1_Nano   InstanceSpecs
	Standard_T1_Micro  InstanceSpecs
	Standard_T1_Medium InstanceSpecs
	Standard_T1_Large  InstanceSpecs
	Standard_T1_XLarge InstanceSpecs
	Standard_T2_Nano   InstanceSpecs
	Standard_T2_Micro  InstanceSpecs
	Standard_T2_Medium InstanceSpecs
	Standard_T2_Large  InstanceSpecs
	Standard_T2_XLarge InstanceSpecs
}

var (
	InstanceTemplates = InstanceTemplateStruct{
		Standard_T1_Nano: InstanceSpecs{
			Vcpu: 1, HtEnabled: false, MemSizeMib: 528,
		},
		Standard_T1_Micro: InstanceSpecs{
			Vcpu: 1, HtEnabled: false, MemSizeMib: 1059,
		},
		Standard_T1_Medium: InstanceSpecs{
			Vcpu: 2, HtEnabled: false, MemSizeMib: 1588,
		},
		Standard_T1_Large: InstanceSpecs{
			Vcpu: 2, HtEnabled: false, MemSizeMib: 2092,
		},
		Standard_T1_XLarge: InstanceSpecs{
			Vcpu: 4, HtEnabled: false, MemSizeMib: 4200,
		},
		Standard_T2_Nano: InstanceSpecs{
			Vcpu: 1, HtEnabled: true, MemSizeMib: 528,
		},
		Standard_T2_Micro: InstanceSpecs{
			Vcpu: 1, HtEnabled: true, MemSizeMib: 1048,
		},
		Standard_T2_Medium: InstanceSpecs{
			Vcpu: 2, HtEnabled: true, MemSizeMib: 1570,
		},
		Standard_T2_Large: InstanceSpecs{
			Vcpu: 2, HtEnabled: true, MemSizeMib: 2092,
		},
		Standard_T2_XLarge: InstanceSpecs{
			Vcpu: 4, HtEnabled: true, MemSizeMib: 4200,
		},
	}
)
