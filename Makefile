CFLAGS=-Wall -std=gnu99

.PHONY: gogo
gogo: gogo.o

.PHONY: test
test:
	go test ./... -v && \
	./test.sh

.PHONY: clean
clean:
	rm -f gogo *.out *.s *.o
