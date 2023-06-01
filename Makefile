all:
	git submodule update --init --recursive && cd go-llama.cpp && make libbinding.a
