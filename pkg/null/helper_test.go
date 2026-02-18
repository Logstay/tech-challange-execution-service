package null

var (
	intJSON         = []byte(`12345`)
	nullIntJSON     = []byte(`{"Int64":12345,"Valid":true}`)
	boolJSON        = []byte(`true`)
	nullBoolJSON    = []byte(`{"Bool":true,"Valid":true}`)
	floatJSON       = []byte(`1.2345`)
	nullFloatJSON   = []byte(`{"Float64":1.2345,"Valid":true}`)
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"String":"test","Valid":true}`)
	nullJSON        = []byte(`null`)
)

func maybePanic(e error) {
	if e != nil {
		panic(e)
	}
}
