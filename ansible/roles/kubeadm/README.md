# kubeadm

Роль устанавливает `kubeadm` из официального репозитория Kubernetes.
После установки пакет фиксируется (`apt-mark hold`), чтобы случайный `apt upgrade` не сломал кластер.

## Зависимости

Используются только встроенные модули `ansible.builtin`.

## Переменные

| Переменная | Тип | По умолчанию | Описание |
|------------|-----|-------------|----------|
| `kubernetes_version` | string | `"1.32"` | Версия Kubernetes |
| `kubernetes_gpg_key_url` | string | формируется из `kubernetes_version` | URL GPG-ключа |
| `kubernetes_apt_repo_url` | string | формируется из `kubernetes_version` | URL APT-репозитория |

## Пример использования

```yaml
- name: Установка kubeadm
  hosts: all
  become: true
  roles:
    - role: kubeadm
```
