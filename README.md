 SSL Certificate Checker

SSL Certificate Checker is a web service written in the Go programming language, designed for monitoring and managing information about security certificates (SSL/TLS) for various hosts and their associated services. The service provides a RESTful API to retrieve information about hosts and their certificates, as well as import and export data related to services.

 Key Features:

- Retrieve List of Services and Associated Hosts:
  The /services endpoint returns a JSON object containing information about available services and their associated hosts.

- Retrieve Certificate Information by Host:
  The /certificate endpoint accepts a host parameter in the request and returns detailed information about the certificate, including:
  - Host name
  - Certificate issuer (Issuer)
  - Validity start and end dates of the certificate
  - Serial number of the certificate
  - System time
  - System uptime
  - CPU usage
  - RAM usage

- Import Service Data from a JSON File:
  The /import endpoint allows uploading service data and their associated hosts in JSON format. After successful import, the data is stored in the service's memory.

- Export Service Data to a JSON File:
  The /export endpoint returns current service data and their associated hosts in JSON format.

- Progress Animation:
  The simulateProgress function demonstrates task progress visually using the mpb library. This function runs in the background when the server starts.

 Technologies and Libraries Used:

- Go: The programming language used to write the service.
- net/http: Library for creating an HTTP server and handling requests.
- encoding/json: Library for encoding and decoding JSON data.

 Running and Using the Service:

To run the service, compile and execute the binary file. The service will be accessible at http://localhost:8080. Interaction with the API can be done using HTTP requests such as GET and POST, depending on the chosen endpoint.

 Example Usage:

- Retrieve list of services:

  curl http://localhost:8080/services
