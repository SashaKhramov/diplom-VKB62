module example.com/dpi

go 1.23.4

require (
	github.com/google/gopacket v1.1.19
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
)

require github.com/mushorg/go-dpi v0.0.0-20221106151256-6cae3029b928

require golang.org/x/sys v0.0.0-20190412213103-97732733099d // indirect

replace example.com/dpi => ../dpi
