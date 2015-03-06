SOURCES = $(shell ls *.c)
EXECUTABLE = ttk4235demoelevator
COMPILER = clang
CFLAGS = -Wall -g -std=c11 -fsanitize=address
LDFLAGS = -fsanitize=address -lcomedi -lm

all: $(EXECUTABLE)

rebuild:	clean all
	
clean:
	rm $(EXECUTABLE)

# unable to replace with obj, because of some ASAN stuff
$(EXECUTABLE): $(SOURCES)
	$(COMPILER) $(CFLAGS) $^ -o $@ $(LDFLAGS)	
	
	
.PHONY: all rebuild clean

