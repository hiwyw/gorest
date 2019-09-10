package field

type Embed struct {
	Id  string `rest:"default=xxxx"`
	Age int64  `rest:"default=20"`
}

type IncludeStruct struct {
	Int8WithRange     int8   `json:"int8WithRange" rest:"min=1,max=20"`
	Uint16WithDefault uint16 `json:"uint16WithDefault,omitempty" rest:"default=11"`
}

type TestStruct struct {
	Embed `json:",inline"`

	Name               string `json:"name" rest:"required=true"`
	StringWithOption   string `json:"stringWithOption,omitempty" rest:"required=true,options=lvm|ceph"`
	StringWithDefault  string `json:"stringWithDefault,omitempty" rest:"default=boy"`
	StringWithLenLimit string `json:"stringWithLenLimit" rest:"minLen=2,maxLen=10"`
	IntWithDefault     int    `json:"intWithDefault,omitempty" rest:"default=11"`
	IntWithRange       uint32 `json:"intWithRange" rest:"min=1,max=1000"`
	BoolWithDefault    bool   `json:"boolWithDefault,omitempty" rest:"default=true"`

	Composition         []IncludeStruct          `json:"composition" rest:"required=true"`
	StringMapCompostion map[string]IncludeStruct `json:"stringMapComposition" rest:"required=true"`
	IntMapCompostion    map[int32]IncludeStruct  `json:"intMapComposition" rest:"required=true"`
}
