
hm.get:
	git -C ./host-manager checkout tags/1.1.0 || git clone git@github.com:art-sitedesign/host-manager.git ./host-manager && git -C ./host-manager checkout tags/1.1.0

hm.build:
	make hm.get
	make -C ./host-manager build

hm.add:
	sudo ./host-manager/bin/$(E) -d $(D) -add

hm.rm:
	sudo ./host-manager/bin/$(E) -d $(D) -rm

app.build:
	bash ./builder.sh
	make hm.build

app.run.mac:
	./bin/darwin-amd64

app.run.linux:
	./bin/linux-amd64
