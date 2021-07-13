# target directories
workingDir = /opt/x-idp
systemdDir = /etc/systemd/system
configDir = /etc/idp

build:
	mkdir -p bin
	go build -o ./bin/idp

install:
	mkdir -p $(workingDir)
	cp ./bin/idp $(workingDir)
	cp -r ./static $(workingDir)
	cp -r ./templates $(workingDir)

setup:
	useradd x-idp -d $(workingDir)
	mkdir -p $(workingDir)/system
	chown -R x-idp:x-idp $(workingDir)
	touch $(workingDir)/system/secret.key
	chmod 600 $(workingDir)/system/secret.key
	chown x-idp:x-idp $(workingDir)/system/secret.key
	tr -dc 'a-z0-9!@#$%^&*(-_=+)' < /dev/urandom | head -c50 > $(workingDir)/system/secret.key
	cp ./config/x-idp.service $(systemdDir)
	mkdir -p $(configDir)
	cp	./config/config.yaml $(configDir)
	systemctl daemon-reload
