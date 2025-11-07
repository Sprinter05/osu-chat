CC=go
BUILD=build

.PHONY: clean
	
$(BUILD)/oauthd: $(BUILD)
	$(CC) build -o $(BUILD)/oauthd . 

$(BUILD):
	if ! [ -d "./$(BUILD)" ], then mkdir $(BUILD); fi

clean: $(BUILD)
	rm -r $(BUILD)