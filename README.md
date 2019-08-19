# aws-rds
OpenShift/Kubernetes operator to manage creating/destroying RDS databases on AWS

## Set-up and config

This operator requires AWS credentials in order to be able to work with AWS.
To allow the operator to access AWS make sure you have one of the following prior the operator installation:

* [Configure AWS CLI](#configure-aws-cli), or
* [Update secret file manually and then apply](#update-the-secret-file-manually), or
* [Use `oc` CLI tool](#use-oc-cli-tool)

### Configure AWS CLI

This is the most convenient. If you have [AWS CLI configured](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) simply
use the `make install-operator-secrets` command as is and the AWS credentials are picked up from the `~/.aws/credentials` file and used to
replace `REPLACE_*` placeholders in the [templates/aws.secret.yaml](./templates/aws.secret.yaml) file.

### Update the secret file manually

To install the AWS secret, update the [./templates/aws.secret.yaml)](./templates/aws.secret.yaml) file manually,

```yaml
apiVersion: v1
kind: Secret
metadata:
    name: aws-rds-operator
    namespace: openshift-operators
    type: Opaque
data:
    AWS_ACCESS_KEY_ID: ...
    AWS_SECRET_ACCESS_KEY: ...
    AWS_REGION: dXMtZWFzdC0y #(BASE64:us-east-2)
```

and then run the same `make install-operator-secrests` command.

### Use `oc` CLI tool

```sh
oc create secret generic aws-rds-operator --from-literal=AWS_ACCESS_KEY_ID=... --from-literal=AWS_SECRET_ACCESS_KEY=... --from-literal=AWS_REGION=us-east-2 -n openshift-operators
```

## Working with operator locally

To build the operator locally run

```sh
make build
```

To run the operator locally run

```sh
make run-locally
```

## Deploying the operator to OpenShift

### Using OperatorHub

Create an `OperatorSource`

```yaml
apiVersion: operators.coreos.com/v1
kind: OperatorSource
metadata:
  name: aws-rds-operator
  namespace: openshift-marketplace
spec:
  type: appregistry
  endpoint: https://quay.io/cnr
  registryNamespace: pmacik
```

Now go to OperatorHub in OpenShift console and install the AWS RDS Operator.

### Directly

Coming Soon...

## Creating a RDS database

Create a secret with the desired DB username and password:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mydb
  namespace: default
  labels:
    app: mydb
type: Opaque
data:
  DB_USERNAME: cG9zdGdyZXM= #(BASE64:postgres)
  DB_PASSWORD: cGFzc3dvcmRvcnNvbWV0aGluZw== #(BASE64:passwordorsomething)
```

Create a `RDSDatabase` custom resource:

```yaml
apiVersion: aws.pmacik.dev/v1alpha1
kind: RDSDatabase
metadata:
  name: mydb
  namespace: default
  labels:
    app: mydb
spec:
  class: db.t2.micro
  engine: postgres
  dbName: mydb
  name: mydb
  password:
    key: DB_PASSWORD
    name: mydb # the name of the secret created above
  username: postgres
  publiclyAccessible: true
  size: 10
```

The creation of the DB takes approximately 5 minutes. A progress can be watched in the `.status.state` or `.status.message` attributes of the `RDSDatabase` custom resource:

```yaml
...
status:
  dbConnectionConfig: mydb
  dbCredentials: mydb
  message: ConfigMap Created
  state: Completed
...
```

Once the state is `Complete` a `ConfigMap` referenced by `.status.dbConnectionConfig` attribute is created and it contains the connection information:

```sh
oc get cm mydb -n default -o yaml
```

```yaml
apiVersion: v1
data:
  DB_HOST: <AWS DB URL>
  DB_PORT: "9432"
kind: ConfigMap
metadata:
...
```

while the secret referenced by the `.status.dbCredentials` attribute contains the DB username and password.
