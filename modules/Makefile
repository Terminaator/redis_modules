SRC_DIR = ./src/

OUT_DIR = ./out/
OUT_O_DIR = ./out/o/
OUT_SO_DIR = ./out/so/

CC = gcc  # C compiler
LD = ld  # ld compiler

CFLAGS = -fPIC -std=gnu99 -ljson-c -Iheader -c -o # Cc flags
LDFLAGS = -shared -Bsymbolic -lc -o  # ld flags

FLAG = -lcurl -linit

C_FILES = $(notdir $(shell find $(SRC_DIR) -name *.c))

all: make_dirs make_o make_so
#-lcurl -linit
make_dirs:
	mkdir -p $(OUT_DIR) $(OUT_O_DIR) $(OUT_SO_DIR)
make_o:
	$(foreach var,$(C_FILES), $(CC) $(CFLAGS) $(addprefix $(OUT_O_DIR), $(var:.c=.o)) $(addprefix $(SRC_DIR), $(var));)

make_so:
	$(foreach var,$(C_FILES), $(LD) $(LDFLAGS)  $(addprefix $(OUT_SO_DIR), $(var:.c=.so)) $(addprefix $(OUT_O_DIR), $(var:.c=.o));)