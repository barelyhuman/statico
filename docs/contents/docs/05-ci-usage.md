---
title: CI / Deployments / Hosting
published: true
---

Since the generated `out` folder from the given config is just simple HTML files you can host it on anything that can serve static web format files.

And if you are working with something that has a build step then you can let the service clone your code and use the below command to generate the `out` at their build stage.

```sh
curl -sf https://gobinaries.com/barelyhuman/statico | sh; statico
```

The above would download the latest version of statico and then execute it in the given environment.

**Note:** If you are using a differently named `config.yml` , please make sure to add that into the command above.

```sh
# eg:
curl -sf https://gobinaries.com/barelyhuman/statico | sh; statico -c doc.statico.yml
```
