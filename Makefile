# run statsd stub that listen UPD port (for monitoring that proxy sends metrics to UPD)
run-statsd-stub:
	nc -kluv localhost 8125

run-bench:
	go test -bench . -benchmem -benchtime=5s