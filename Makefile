run:
	@go run main.go

# sending a zip file called zip_file as a parameter named code
# the '@' symbol in curl is used to specify that the value of the parameter should be read from a file
submit:
	curl -X POST -F "code=@(zip_file)" http://localhost:8000/api/submit

# this contains a query parameter names fn that is set to the value of $(fn) var
# this is for the execution of the code that was previously sent and stored within the fn parameter
execute:
	curl http://localhost:8080/api/execute\?fn\=$(fn)
