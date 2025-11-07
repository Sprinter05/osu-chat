CC=go
BUILD=build

oauthd: $(BUILD)/oauthd
.PHONY: clean
	
$(BUILD)/oauthd: $(BUILD)
	$(CC) build -o $(BUILD)/oauthd ./callback 

$(BUILD):
	if ! [ -d "./$(BUILD)" ]; then mkdir $(BUILD); fi

clean: $(BUILD)
	rm -r $(BUILD)