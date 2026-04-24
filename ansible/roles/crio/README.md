# crio

Роль устанавливает CRI-O.

Помимо самого CRI-O роль загружает модуль ядра `overlay` — он нужен CRI-O для работы с файловой системой контейнеров (overlayfs).

## Зависимости

`community.general` используется для модуля `modprobe` для загрузки модулей ядра

Установить:
```bash
ansible-galaxy collection install -r requirements.yml
```

## Переменные

| Переменная | Тип | По умолчанию | Описание |
|------------|-----|-------------|----------|
| `crio_version` | string | `"1.32"` | Версия CRI-O |
| `crio_gpg_key_url` | string | формируется из `crio_version` | URL GPG-ключа |
| `crio_apt_repo_url` | string | формируется из `crio_version` | URL APT-репозитория |

## Пример использования

```yaml
- name: Установка CRI-O
  hosts: all
  become: true
  roles:
    - role: crio
```
