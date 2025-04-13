```toml
name = 'app05 APIs'
id = 'cd64ce57-2890-4f36-b3b9-85c89d77523e'

[[environmentGroups]]
name = 'Default'
environments = ['development', 'production']
```

#### Variables

```json5
{
  development: {
    "baseUrl": "http://localhost:8081/api/v1",
    "token": "qrpm29EaFPQRvfOxdqGIz9kVcDn6n7rSBHxQVg1SPy8="
  },
  production: {
    "baseUrl": "https://backend.somolabs.com/api/v1",
    "token": "3qimBmhNU82xt_FPaTacJR5gKGjIN-4dnOGvRbfOcDk="
  }

}
```
