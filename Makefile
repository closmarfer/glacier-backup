build:
	go build -o bin/glacier-backup cmd/main.go

install: create-app-folder install-config build
	cp bin/glacier-backup /usr/local/bin/glacier-backup

update: build
	cp bin/glacier-backup /usr/local/bin/glacier-backup

install-config:
	cp ./var/config.yaml "$(HOME)/.glacier-backup/config.yaml"

create-app-folder:
	mv "$(HOME)/.glacier-backup" "$(HOME)/.glacier-backup_bk"
	mkdir -p "$(HOME)/.glacier-backup"

uninstall:
	rm /usr/local/bin/glacier-backup