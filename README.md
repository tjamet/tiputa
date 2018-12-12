# Tiputa

Tiputa is a PoC to use [pass](https://www.passwordstore.org/) to encrypt user authentication of kuberntes clients.
It implements Kubernetes [client-go credential plugins](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins)
available in beta since kubernetes 1.11.

# Usage

install using

```bash
go get -u github.com/tjamet/tiputa
go build -o /usr/local/bin/tiputa github.com/tjamet/tiputa
```

Then, export your kubernetes credentials to pass and update your kubernetes configuration:

```yaml
<...>
- context:
    cluster: your-cluster
    user: password-encrypted-user
  name: your-cluster
<...>
- name: password-encrypted-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args:
      - -pass-token
      - where/you/saved/your/token
      command: /usr/local/bin/tiputa
```

Next time you run `kubectl --context your-cluster`, your the token will be retrieved from pass