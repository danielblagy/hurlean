# hurlean

## To Do
- [x] console interface cleanup
	- [x] server: 'hurlean' prefix in prints and errors
	- [x] client: 'hurlean' prefix in prints and errors
	- [x] server: print 'server starting' on success
	- [x] server: make debug prints conditional
	- [x] client: make debug prints conditional

- [ ] code cleanup (formatting, delete commented code)

- [ ] simple example: time querying

Bugs:
- client connection may close unexpectedly
- if it is closed like that, on the server side, the client's go routines stop and the connection is not deleted from the Clients map