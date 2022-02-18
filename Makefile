install: create-app-folder install-config
	go build -o bin/glacier-backup cmd/main.go
	cp bin/glacier-backup /usr/local/bin/glacier-backup

install-config:
	cp ./var/config.yaml "$(HOME)/.glacier-backup/config.yaml"

create-app-folder:
	cp "$(HOME)/.glacier-backup" "$(HOME)/.glacier-backup_bk"
	mkdir -p "$(HOME)/.glacier-backup"

uninstall:
	rm /usr/local/bin/glacier-backup