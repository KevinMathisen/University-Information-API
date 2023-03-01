# Assignment 1

Assignment 1 for PROG2005 Spring 2023

Owner: Kevin Nikolai Mathisen
Email: Kevinnm@stud.ntnu.no

# Running the assignment
You can run the application either by downloading the repository and running the application locally, or visting the render deployment.

### Render
The render deployment is located at [render](https://prog2005-assignment-1.onrender.com/)

### Locally
If you downloaded the repository you can run it in one of two ways, as long as you have golang installed on your computer.
- In the repository run the command `go run ./cmds/main.go`
- In the repository, first run the command `go build ./cmds`, which will create a file calles *cmds.exe*. Running this file will start the application.

When running the application locally it will be available on `localhost:8080`, or a different port if specified in environment. 

# How to use the API
There are four enpoints accessible on the application:
```
/
/unisearcher/v1/uniinfo/
/unisearcher/v1/neighbourunis/
/unisearcher/v1/diag/
```
## Default path
Only for accessibility, no functionality. Lists and links to the other endpoints provided.
```
Method: GET
Path: /
```
### Response
```
Content-type: text/html
Status code: Appropriate http status code, 200 if everything ok
```
## Uniinfo
For retriving information about all universities which match the given name. This information includes:
- Name
- Country located in
- IsoCode of country
- Webpages for the university
- Langauges spoken in country
- Openstreetmap of the country

### Request
```
Method: GET
Path: unisearcher/v1/uniinfo/{:partial_or_complete_university_name}/
```
Name of university can be complete or partial. If partial is used, the name given has to either be a single word, or if multiple words are given, have to be found in the beginning of the real university name. This is because of a limitation in the API used in the backend. Leading and trailing whitespace around university name provided is ignored.

Example request `unisearcher/v1/uniinfo/norwegian%20university%20of%20science/`

### Response
```
Content-type: application/json
Status code: Appropriate http status code, 200 if everything ok
```
Body example:
```
[
  {
      "name": "Norwegian University of Science and Technology", 
      "country": "Norway",
      "isocode": "NO",
      "webpages": ["http://www.ntnu.no/"],
      "languages": {"nno": "Norwegian Nynorsk",
                    "nob": "Norwegian Bokmål",
                    "smi": "Sami"},
      "map": "https://www.openstreetmap.org/relation/2978650"
  },
  ...
]
```
## Neighboursunis
For retriving information about all universities which match the given name in the neighbouring countries of the country provided. 
The information returned will be identical in format to the uniinfo endpoint. 

### Request
```
Method: GET
Path: unisearcher/v1/neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}
```
`{:country_name}` is the English name of the country that is used to find neighbouring contries, which is then used to find universities.

`{:partial_or_complete_university_name}` is the partial or complete university name, which is used to find universities in neighbouring countries

`{?limit={:number}}` is an optional parameter that limits the number of universities in neighbouring countries that are returned. If the limit is not set, or if it is set to 0, there will be no limit. Else the limit has to be a positive integer.

Example request `unisearcher/v1/neighbourunis/china/University%20of%20Science?limit=1`

### Response
Even though the request is different, the response will be identical to the response for the uniinfo. See Response uniinfo

## Diag
For viewing the status of the application, and the services it depends on. Returns the status code of the universities and countries api, and the version and uptime in seconds of the application running.

### Request
```
Method: GET
Path: unisearcher/v1/diag/
```

### Response
```
Content-type: application/json
Status code: Appropriate http status code, 200 if everything ok
```
Body example:
```
{
   "universitiesapi": "200",
   "countriesapi": "200",
   "version": "v1",
   "uptime": 150.5
}

```