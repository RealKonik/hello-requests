module github.com/RealKonik/hello-requests

go 1.15

require (
	github.com/dsnet/compress v0.0.1
	github.com/gwatts/rootcerts v0.0.0-20230201191557-c2e6d643fa97
	github.com/hunterbdm/hello-requests v0.0.0-00010101000000-000000000000
	github.com/tam7t/hpkp v0.0.0-20160821193359-2b70b4024ed5
	gitlab.com/yawning/bsaes.git v0.0.0-20190805113838-0a714cd429ec
	golang.org/x/crypto v0.8.0
	golang.org/x/net v0.9.0
)

replace github.com/hunterbdm/hello-requests v0.0.0-00010101000000-000000000000 => ./
