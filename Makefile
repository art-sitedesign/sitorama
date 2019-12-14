.PHONY: hm.get

hm.get:
	git -C ./host-manager pull || git clone git@github.com:art-sitedesign/host-manager.git ./host-manager

hm.add:
	make hm.get
	sudo ./host-manager/bin/$(E) -d $(D) -add

hm.rm:
	make hm.get
	sudo ./host-manager/bin/$(E) -d $(D) -rm

app.build:
	bash ./builder.sh

app.run.mac:
	./bin/darwin-amd64

app.run.linux:
	./bin/linux-amd64
