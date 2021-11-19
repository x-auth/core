# target directories
workingDir = /opt/x-idp
systemdDir = /etc/systemd/system
configDir = /etc/idp

# hydra stuff
hydra_installer = https://raw.githubusercontent.com/ory/hydra/v1.9.0/install.sh
hydra_version = v1.9.0

build:
	mkdir -p bin
	go build -o ./bin/idp

install:
	mkdir -p $(workingDir)
	cp ./bin/idp $(workingDir)
	cp -r ./static $(workingDir)
	cp -r ./templates $(workingDir)

setup:
	# set up the portal software
	useradd x-idp -d $(workingDir)
	mkdir -p $(workingDir)/system
	chown -R x-idp:x-idp $(workingDir)
	touch $(workingDir)/system/secret.key
	chmod 600 $(workingDir)/system/secret.key
	chown x-idp:x-idp $(workingDir)/system/secret.key
	tr -dc 'a-z0-9!@#$%^&*(-_=+)' < /dev/urandom | head -c50 > $(workingDir)/system/secret.key
	cp ./config/portal.service $(systemdDir)
	cp ./config/hydra.service $(systemdDir)
	cp ./config/x-idp.service $(systemdDir)
	mkdir -p $(configDir)
	cp	./config/config.yaml $(configDir)
	systemctl daemon-reload

	# download and install ory hydra
	mkdir -p $(workingDir)/services
	bash <(curl $(hydra_installer)) -b $(workingDir)/services $(hydra_version)

.PHONY: build install setup