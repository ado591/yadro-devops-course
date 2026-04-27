## ДЗ k8s. Как это было

Сначала определилась с последовательностью шагов: 
1. Отключить swap
2. Установить br_netfilter и overlay
3. Включить ipv4 forward.
4. Установить cri-o, kubelet, kubeadm, kubectl
5. Стартануть cri-o
6. ```kubeadm init```
7. Если запуск прошел успешно, то появится ```Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 10.184.0.131:6443 ...```. Команду скопировать, но bootstrap token можно при необходимости переиздать
8. Переместить конфиг в ```~/.kube/config```
9. Включить calico по [инструкции](https://docs.tigera.io/calico/latest/getting-started/kubernetes/quickstart)
10. Повторить действия 1-5 с worker нодой
11. Включить kubeadm join из пункта 7

С установкой пакетов возникли трудности, тк сильно ограничена скорость загрузки на ВМ. Самым простым способом перекинула себе через scp. Без пункта 8 попыталась дернуть ```kubectl get nodes``` получила ```"Unhandled Error" err="couldn't get current server API group list```. После пункта 8 повторила и получила: 
```
NAME     STATUS     ROLES           AGE     VERSION
master   NotReady   control-plane   8m47s   v1.32.13
```
Поняла, что Calico не понимает какую подсеть использовать. Сделала ```kubeadm reset```, а затем ```kubeadm init --pod-network-cidr=192.168.0.0/16```. После того как все Calico образы скачались пошла работать с worker нодой. Получила
```
NAME     STATUS     ROLES           AGE   VERSION
master   Ready      control-plane   17m   v1.32.13
worker   NotReady   <none>          10s   v1.32.13
```
спустя время Calico распространился на worker ноду и Status стал Ready.

### Busybox 

```kubectl run busybox --image=busybox -- ping 8.8.8.8``` + ```kubectl describe pod busybox``` выдал следующую картину:

```
Events:
  Type     Reason            Age   From               Message
  ----     ------            ----  ----               -------
  Warning  FailedScheduling  11s   default-scheduler  0/2 nodes are available: 1 node(s) had untolerated taint {[node-role.kubernetes.io/control-plane](http://node-role.kubernetes.io/control-plane): }, 1 node(s) had untolerated taint {node.kubernetes.io/disk-pressure: }. preemption: 0/2 nodes are available: 2 Preemption is not helpful for scheduling.
```

2 вывода: 1. Control Plane не принимает обычные поды; 2. Disk pressure на worker ноде. Идем на worker ноду разбираться

```df -h``` и получаем 89%. Попыталась почистить всякий хлам и снизила загрузку до 80%, но желаемого эффекта не добилась. Поэтому залетаем в ```/var/lib/kubelet/config.yaml``` и ставим крайне низкие лимиты:
```
evictionHard:
  nodefs.available: "5%"
  nodefs.inodesFree: "5%"
  imagefs.available: "5%"
```

Хочется поисследовать что же там так сильно заняло место, но для тестирования подойдет. После чего наш busybox поселился на worker. Успех!

### Ansible

В прошлой домашке уже сделала несколько ролей. Хотелось более-менее эффективно их использовать. Пришла к следующей схеме: 
1. Имеем k8s_cluster, где будут жить все наши ВМ
2. Отдельно выделим master для Control Plane, отедльно workers для остальных. Исторически сложилось, что в прошлой ДЗ я и так сделала разделение на staging и production, его и оставила.
3. В common вынесла шаги 1-5
4. Для control plane сделала сохранение токена в переменную, workers ее использовали 

Вышло не особо секьюрно, хочется доработать до Ansible Vault. Но в свое оправдание скажу, что токен живет сутки и маловероятно, что его успеют угнать.


