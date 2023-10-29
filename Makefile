NAME = mox
CMD = cmd/$(NAME)/main.go
BIN = bin/$(NAME)
TARGET = /usr/local/bin/$(NAME)
INSTALL_LOCATION = /usr/local/Cellar/$(NAME)

.PHONY: build clean install uninstall

build:
	@echo "Building..."
	@go build -o $(BIN) $(CMD)
	@chmod +x $(BIN)
	@echo "Done!"

clean:
	@echo "Cleaning..."
	@rm -f $(BIN)
	@echo "Done!"

install: 
	@echo "Installing Moxie..."
	@make clean >> /dev/null
	@make build >> /dev/null
	@make uninstall >> /dev/null
	@sudo mkdir -p $(INSTALL_LOCATION)
	@sudo cp -r . $(INSTALL_LOCATION)
	@sudo ln -s $(INSTALL_LOCATION)/$(BIN) $(TARGET)
	@echo "Done!"

uninstall:
	@echo "Uninstalling..."
	@sudo rm -rf $(INSTALL_LOCATION)
	@@sudo rm -f $(TARGET)
	@echo "Done!"
