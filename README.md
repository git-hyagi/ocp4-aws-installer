### ABOUT
The process of install the `OpenShift 4` is a lot easier compared to the `3.x` version (which was based on running Ansible playbooks).  
To make it even easier, the `ocp4-aws-installer` just help with some steps:
* gather the pull-secret
* download and install the `openshift-install` program
* download and install the `oc` (openshift client)

### PRE-REQUISITES
* **A ssh key** that can be used to SSH into the master nodes as the user `core`. When you deploy the cluster, the key is added to the `core` user's `~/.ssh/authorized_keys` list. Example of command to create a new ssh key:
~~~
$ ssh-keygen -t rsa -b 4096 -N '' -f ~/.ssh/ocp4
~~~

* **Follow the steps from [documentation](https://docs.openshift.com/container-platform/4.5/installing/installing_aws/installing-aws-account.html) to configure the AWS account:**
  * [Configuring Route53](https://docs.openshift.com/container-platform/4.5/installing/installing_aws/installing-aws-account.html#installation-aws-route53_installing-aws-account)
  * [AWS account limits](https://docs.openshift.com/container-platform/4.5/installing/installing_aws/installing-aws-account.html#installation-aws-limits_installing-aws-account)
  * [Required AWS permissions](https://docs.openshift.com/container-platform/4.5/installing/installing_aws/installing-aws-account.html#installation-aws-permissions_installing-aws-account)
  * [Creating an IAM user](https://docs.openshift.com/container-platform/4.5/installing/installing_aws/installing-aws-account.html#installation-aws-iam-user_installing-aws-account)
  * (optional) Configure the `~/.aws/credentials`:
~~~
[default]
aws_access_key_id = <AWSKEY>
aws_secret_access_key = <AWSSECRETKEY>
region = sa-east-1
~~~

* **A Red Hat [account](https://access.redhat.com)**

### INSTALLING THE ocp4-aws-installer
There is no need to install the `ocp4-aws-installer`. Just clone or download this repo:
~~~
$ git clone https://github.com/git-hyagi/ocp4-aws-installer.git
~~~

### RUNNING
After completing all the [pre-requisites](#pre\-requisites) and [downloading](#installing-the-ocp4-aws-installer) the program, run:
~~~
$ go run ocp4-aws-installer.go
~~~

During the program execution it will be asked to pass an `username` and `password`. They'll be used to get the `pull-secret` from [https://cloud.redhat.com/openshift/install](https://cloud.redhat.com/openshift/install) and create the `/tmp/pull-secret.txt` file.  
If the `ocp4-aws-installer` finished successfuly, begin the installation with:
~~~
$ openshift-install create cluster --dir <install dir>
~~~


### MORE INFORMATION
Official `OpenShift 4` installation documentation:  
[https://docs.openshift.com/container-platform/4.5/installing/installing_aws/installing-aws-default.html](https://docs.openshift.com/container-platform/4.5/installing/installing_aws/installing-aws-default.html)
