# Stream output example
This example shows how to use the `BuildStream` function to generate the document into a stream. You can use this stream to send the document to any destination that accepts a stream, like an HTTP response or a file.

# Annotations
The example also shows that now we can add routes without parsing the code looking for annotations. This feature can be helpful in several use cases, like generating the documentation from a framework, or some definition or manifest because you don't have access to code to write annotations.

## Update index.html
This example serves the document in the `/docs/oas`, no file is generated, and the renderer in `/docs/api`. To correctly render the document you must uncomment line 40 and make sure lines 38 and 39 are commented in `internal/dist/index.html`

```html
    ...
        window.ui = SwaggerUIBundle({
            // url: "openapi.yaml",
            // url: "https://petstore.swagger.io/v2/swagger.json",
            url: "/docs/oas",
    ...
```
