start_test_icinga2:
	docker run -d --name icinga2 -p 8080:80 -p 8443:443 -p 5665:5665 -it jordan/icinga2:latest
	sleep 60
	docker exec icinga2 bash -c 'echo -e "object ApiUser \"icinga-test\" {\n  password = \"icinga\"\n  permissions = [ \"*\" ]\n}" > /etc/icinga2/conf.d/test-api-user.conf'
	docker exec icinga2 supervisorctl restart icinga2
