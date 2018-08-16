# Blue Jay Fork - Blueprint

> [Original Blue Jay Repo](https://github.com/blue-jay/blueprint)

## How to run

```bash
npm install
npm run init
npm run watch
```

## How to update JavaScript dependencies (packages.json)?

```bash
sudo npm install -g npm-check-updates
npm-check-updates -u -a
```

## How to add secrets and configmap to openshift

```bash
oc create -f .openshift/secrets.yaml
oc create -f .openshift/configmap.yaml
```
