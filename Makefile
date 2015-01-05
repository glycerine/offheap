all:
	go build
	go install

clean:
	rm -f *~

kv:
	cd keyval; bambam -p="keyval" -o="." account.go && mv schema.capnp keyval.capnp && capnpc -ogo keyval.capnp
