
``OTECOL`` in ``cri-o``
=======================

Start stop ``cri-o``
====================

.. code :: bash 

  sudo systemctl daemon-reload
  sudo systemctl enable crio
  sudo systemctl start crio

Running Kubernetes with ``kubeadm``
===================================

* ``init`` a Kubernetes cluster

.. code :: bash 

  sudo kubeadm init --config=$KUBEADM_CONFIG --cri-socket /var/run/crio/crio.sock 

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

* Install CRDS and Calico networking

.. code-block:: bash
  

  kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml

  cat >>EOF kubectl create -f -

  # This section includes base Calico installation configuration.
  # For more information, see: https://docs.projectcalico.org/v3.18/reference/installation/api#operator.tigera.io/v1.Installation
  apiVersion: operator.tigera.io/v1
  kind: Installation
  metadata:
    name: default
  spec:
    # Configures Calico networking.
    calicoNetwork:
      # Note: The ipPools section cannot be modified post-install.
      ipPools:
      - blockSize: 26
        cidr: 10.85.0.0/16    #Change this to your CIDR value
        encapsulation: VXLANCrossSubnet
        natOutgoing: Enabled
        nodeSelector: all()


.. code-block:: bash

  kubectl taint nodes --all node-role.kubernetes.io/master-

* Verify the pods

.. code-block:: bash

  kubectl get pods --all-namespaces

  NAMESPACE         NAME                                                   READY   STATUS    RESTARTS   AGE
  kube-system       coredns-558bd4d5db-fwpnd                               1/1     Running   0          16h
  kube-system       coredns-558bd4d5db-szpjs                               1/1     Running   0          16h
  kube-system       etcd-ip-172-31-80-66.ec2.internal                      1/1     Running   0          16h
  kube-system       kube-apiserver-ip-172-31-80-66.ec2.internal            1/1     Running   1          16h
  kube-system       kube-controller-manager-ip-172-31-80-66.ec2.internal   1/1     Running   0          16h
  kube-system       kube-proxy-mrrkm                                       1/1     Running   0          16h
  kube-system       kube-scheduler-ip-172-31-80-66.ec2.internal            1/1     Running   0          16h


* Verify node

.. code-block:: bash

  kubectl get nodes -o wide

  NAME                           STATUS   ROLES                  AGE   VERSION   INTERNAL-IP    EXTERNAL-IP   OS-IMAGE                    KERNEL-VERSION           CONTAINER-RUNTIME
  ip-172-31-80-66.ec2.internal   Ready    control-plane,master   17h   v1.21.0   172.31.80.66   <none>        Fedora 31 (Cloud Edition)   5.8.18-100.fc31.x86_64   cri-o://1.21.0




# * Verify the container
# sudo crictl ps