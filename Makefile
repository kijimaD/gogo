.PHONY: test
test:
	go test ./... -v && \
	./test.sh

.PHONY: clean
clean:
	rm -f *.out *.s
