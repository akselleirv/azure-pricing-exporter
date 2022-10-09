An opinionated REST API layout for Go.


- oapi-codegen - Follows design first principles.
- Echo - High performance, extensible, minimalist Go web framework.
- Tilt - A toolkit for fixing the pains of microservice development.


# Start

Flag arguments:
```
-module-name string
    the name of the new module
-output string
    the path to output the files
-service-name string
    the name of the new service
```

Example:

```
go run main.go -output $HOME/my-service -module-name="github.com/akselleirv/my-service" -service-name my-service
```