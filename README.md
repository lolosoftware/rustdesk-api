# RustDesk API

[English Doc](README_EN.md)

Ce projet implémente l’API RustDesk en Go et inclut un panneau d’administration web ainsi qu’un client web.

<div align=center>
<img src="https://img.shields.io/badge/golang-1.22-blue"/>
<img src="https://img.shields.io/badge/gin-v1.9.0-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-v1.25.7-green"/>
<img src="https://img.shields.io/badge/swag-v1.16.3-yellow"/>
<img src="https://goreportcard.com/badge/github.com/lejianwen/rustdesk-api/v2"/>
<img src="https://github.com/lejianwen/rustdesk-api/actions/workflows/build.yml/badge.svg"/>
</div>

## Préférable avec [lejianwen/rustdesk-server].
> [lejianwen/rustdesk-server] est un fork du dépôt officiel RustDesk Server.
> 1. Corrige le problème de délai d’attente de connexion API.
> 2. Permet de forcer la connexion avant d’initier une session.
> 3. Prend en charge le websocket côté client.

# Fonctionnalités

- API côté PC
    - API personnelle
    - Connexion
    - Carnet d’adresses
    - Groupes
    - Authentification autorisée
      - Prise en charge de la connexion `GitHub`, `Google` et `OIDC`
      - Prise en charge de l’authentification via le `back office web`
      - Prise en charge de `LDAP` (AD et OpenLDAP testés), si le serveur API est configuré pour LDAP
    - i18n
- Web Admin
    - Gestion des utilisateurs
    - Gestion des appareils
    - Gestion du carnet d’adresses
    - Gestion des étiquettes
    - Gestion des groupes
    - Gestion OAuth
    - Configuration LDAP via fichier ou variables d’environnement
    - Journaux de connexion
    - Journaux de connexion/liaison
    - Journaux de transfert de fichiers
    - Accès rapide au client web
    - i18n
    - Partage vers des invités via le client web
    - Contrôle serveur (quelques commandes officielles simples : [WIKI](https://github.com/lejianwen/rustdesk-api/wiki/Rustdesk-Command))
- Web Client
    - Récupération automatique du serveur API
    - Récupération automatique du serveur ID et de la clé
    - Récupération automatique du carnet d’adresses
    - Les invités peuvent accéder à distance au poste via un lien de partage temporaire
- CLI
    - Réinitialisation du mot de passe administrateur

## Fonctionnalités

### Service API
L’API de base du client PC est implémentée. La version personnelle est prise en charge et peut être activée via le fichier `rustdesk.personal` ou la variable d’environnement `RUSTDESK_API_RUSTDESK_PERSONAL`.

<table>
    <tr>
      <td width="50%" align="center" colspan="2"><b>Connexion</b></td>
    </tr>
    <tr>
        <td width="50%" align="center" colspan="2"><img src="docs/pc_login.png"></td>
    </tr>
     <tr>
      <td width="50%" align="center"><b>Carnet d’adresses</b></td>
      <td width="50%" align="center"><b>Groupes</b></td>
    </tr>
    <tr>
        <td width="50%" align="center"><img src="docs/pc_ab.png"></td>
        <td width="50%" align="center"><img src="docs/pc_gr.png"></td>
    </tr>
</table>

### Web Admin :

* Le projet utilise une séparation front/back pour offrir une interface d’administration conviviale, principalement pour la gestion et l’affichage. Le code du front est dans [rustdesk-api-web](https://github.com/lejianwen/rustdesk-api-web).

* L’adresse d’accès au panneau d’administration est `http://<votre serveur>[:port]/_admin/`.
* Lors de la première installation, l’administrateur est nommé `admin` et le mot de passe est affiché dans la console. Il peut être modifié via la [ligne de commande](#CLI).

  ![img.png](./docs/init_admin_pwd.png)

1. Interface d’administration
   ![web_admin](docs/web_admin.png)
2. Interface utilisateur standard
   ![web_user](docs/web_admin_user.png)

3. Chaque utilisateur peut avoir plusieurs carnets d’adresses et les partager avec d’autres utilisateurs.
4. Les groupes sont personnalisables. Deux types sont actuellement supportés : `groupe partagé` et `groupe normal`.
5. Il est possible d’ouvrir directement le client web, ou de le partager à des invités qui peuvent accéder à distance aux postes via le client web.
6. OAuth : prise en charge de `GitHub`, `Google` et `OIDC`. Il faut créer une `OAuth App` puis la configurer dans le panneau d’administration.
    - Pour `Google` et `GitHub`, `Issuer` et `Scopes` ne sont pas requis.
    - Pour `OIDC`, `Issuer` est obligatoire. `Scopes` est optionnel et par défaut `openid,profile,email`. Vérifiez que l’application a bien accès à `sub`, `email` et `preferred_username`.
    - La `GitHub OAuth App` est créée dans `Settings` -> `Developer settings` -> `OAuth Apps` -> `New OAuth App`.
      Voir : [https://github.com/settings/developers](https://github.com/settings/developers)
    - Définir l’URL de callback : `http://<votre serveur[:port]>/api/oidc/callback`
      par exemple `http://127.0.0.1:21114/api/oidc/callback`
7. Journaux de connexion
8. Journaux de liaison/connexion
9. Journaux de transfert de fichiers
10. Contrôle du serveur

  - `Mode simple` : certaines commandes sont rendues conviviales et peuvent être exécutées directement depuis le panneau d’administration.
    ![rustdesk_command_simple](./docs/rustdesk_command_simple.png)

  - `Mode avancé` : les commandes peuvent être exécutées directement depuis le back-office.
      * Commandes officielles prises en charge
      * Ajout de commandes personnalisées
      * Exécution de commandes personnalisées

11. **Support LDAP** : si LDAP est configuré sur le serveur API (AD et LDAP testés), les utilisateurs peuvent se connecter via leurs données LDAP. Voir : https://github.com/lejianwen/rustdesk-api/issues/114. En cas d’échec LDAP, l’authentification revient à l’utilisateur local.

### Client Web :

1. Si vous êtes déjà connecté au panneau d’administration, le client web se connecte automatiquement.
2. Si vous n’êtes pas connecté, cliquez sur le bouton de connexion en haut à droite ; le serveur API est déjà configuré automatiquement.
3. Après la connexion, le serveur ID et la clé sont synchronisés automatiquement.
4. Après la connexion, le carnet d’adresses est enregistré automatiquement dans le client web pour faciliter l’utilisation.

### Documentation automatisée :
L’API est documentée via Swag pour faciliter la compréhension et l’utilisation de l’API par les développeurs.

1. Documentation du back-office : `<votre serveur[:port]>/admin/swagger/index.html`
2. Documentation côté PC : `<votre serveur[:port]>/swagger/index.html`
   ![api_swag](docs/api_swag.png)

### CLI

```bash
# afficher l’aide
./apimain -h
```

#### Réinitialiser le mot de passe administrateur
```bash
./apimain reset-admin-pwd <pwd>
```

## Installation et exécution

### Configuration associée

* [Fichier de configuration](./conf/config.yaml)
* Modifier la configuration dans `conf/config.yaml`.
* Si `gorm.type` vaut `sqlite`, la configuration MySQL n’est pas nécessaire.
* Si la langue n’est pas définie, la valeur par défaut est `zh-CN`.

### Variables d’environnement
Les variables d’environnement correspondent une par une aux configurations du fichier `conf/config.yaml`. Le préfixe est `RUSTDESK_API`.
Le tableau ci-dessous n’est pas exhaustif. Consultez le fichier `conf/config.yaml` pour le détail complet.

| Variable | Description | Exemple |
|---|---|---|
| TZ | Fuseau horaire | Asia/Shanghai |
| RUSTDESK_API_LANG | Langue | `en`, `zh-CN` |
| RUSTDESK_API_APP_WEB_CLIENT | Activer/désactiver le client web ; 1 = activé, 0 = désactivé ; actif par défaut | 1 |
| RUSTDESK_API_APP_REGISTER | Inscription activée ; `true`, `false` ; défaut `false` | `false` |
| RUSTDESK_API_APP_SHOW_SWAGGER | Afficher Swagger ; `1` = oui, `0` = non ; défaut `0` | `1` |
| RUSTDESK_API_APP_TOKEN_EXPIRE | Durée de validité du token | `168h` |
| RUSTDESK_API_APP_DISABLE_PWD_LOGIN | Désactiver la connexion par mot de passe ; `true`, `false` ; défaut `false` | `false` |
| RUSTDESK_API_APP_REGISTER_STATUS | Statut par défaut des utilisateurs inscrits ; 1 = activé, 2 = désactivé, défaut 1 | `1` |
| RUSTDESK_API_APP_CAPTCHA_THRESHOLD | Seuil captcha ; -1 = désactivé, 0 = toujours, >0 = activé après N échecs ; défaut `3` | `3` |
| RUSTDESK_API_APP_BAN_THRESHOLD | Seuil de bannissement IP ; 0 = désactivé, >0 = bannir après N échecs ; défaut `0` | `0` |
| ----- CONFIG ADMIN ----- | --- | --- |
| RUSTDESK_API_ADMIN_TITLE | Titre du panneau d’administration | `RustDesk Api Admin` |
| RUSTDESK_API_ADMIN_HELLO | Message d’accueil admin, accepte du HTML |  |
| RUSTDESK_API_ADMIN_HELLO_FILE | Fichier du message d’accueil, le fichier remplace la variable `RUSTDESK_API_ADMIN_HELLO` | `./conf/admin/hello.html` |
| ----- CONFIG GIN ----- | --- | --- |
| RUSTDESK_API_GIN_TRUST_PROXY | IP proxy de confiance, séparées par des virgules ; par défaut accepte tout | 192.168.1.2,192.168.1.3 |
| ----- CONFIG GORM ----- | --- | --- |
| RUSTDESK_API_GORM_TYPE | Type de base de données (`sqlite` ou `mysql`), défaut `sqlite` | sqlite |
| RUSTDESK_API_GORM_MAX_IDLE_CONNS | Nombre maximum de connexions inactives | 10 |
| RUSTDESK_API_GORM_MAX_OPEN_CONNS | Nombre maximum de connexions ouvertes | 100 |
| RUSTDESK_API_RUSTDESK_PERSONAL | Activer l’API personnelle ; 1 = activée, 0 = désactivée ; défaut activée | 1 |
| ----- CONFIG MYSQL ----- | --- | --- |
| RUSTDESK_API_MYSQL_USERNAME | Nom d’utilisateur MySQL | root |
| RUSTDESK_API_MYSQL_PASSWORD | Mot de passe MySQL | 111111 |
| RUSTDESK_API_MYSQL_ADDR | Adresse MySQL | 192.168.1.66:3306 |
| RUSTDESK_API_MYSQL_DBNAME | Nom de la base MySQL | rustdesk |
| RUSTDESK_API_MYSQL_TLS | Activer TLS ; valeurs possibles : `true`, `false`, `skip-verify`, `custom` | `false` |
| ----- CONFIG RUSTDESK ----- | --- | --- |
| RUSTDESK_API_RUSTDESK_ID_SERVER | Adresse du serveur ID RustDesk | 192.168.1.66:21116 |
| RUSTDESK_API_RUSTDESK_RELAY_SERVER | Adresse du relay RustDesk | 192.168.1.66:21117 |
| RUSTDESK_API_RUSTDESK_API_SERVER | Adresse du serveur API RustDesk | http://192.168.1.66:21114 |
| RUSTDESK_API_RUSTDESK_KEY | Clé RustDesk | 123456789 |
| RUSTDESK_API_RUSTDESK_KEY_FILE | Fichier contenant la clé RustDesk | `./conf/data/id_ed25519.pub` |
| RUSTDESK_API_RUSTDESK_WEBCLIENT<br/>_MAGIC_QUERYONLINE | Activer la nouvelle méthode de requête de statut en ligne dans le client web v2 ; `1` = activée, `0` = désactivée, défaut `0` | `0` |
| RUSTDESK_API_RUSTDESK_WS_HOST | Base HTTPS/WSS optionnelle du proxy WebSocket. Le client ajoute `/id` pour le port 21118 et `/relay` pour le port 21119 | `https://rdapi.example.com/webclient-ws` |
| ---- PROXY ----- | --- | --- |
| RUSTDESK_API_PROXY_ENABLE | Activer le proxy : `false`, `true` | `false` |
| RUSTDESK_API_PROXY_HOST | Adresse du proxy | `http://127.0.0.1:1080` |
| ---- JWT ---- | --- | --- |
| RUSTDESK_API_JWT_KEY | Clé JWT personnalisée ; vide = JWT désactivé. Si vous n’utilisez pas `MUST_LOGIN` du serveur `lejianwen/rustdesk-server`, il est conseillé de laisser vide |  |
| RUSTDESK_API_JWT_EXPIRE_DURATION | Durée de validité JWT | `168h` |

### Exécution

#### Docker

1. Exécution directe avec Docker. La configuration peut être modifiée en montant le fichier `/app/conf/config.yaml`, ou en remplaçant les valeurs via des variables d’environnement.

    ```bash
    docker run -d --name rustdesk-api -p 21114:21114 \
    -v /data/rustdesk/api:/app/data \
    -e TZ=Asia/Shanghai \
    -e RUSTDESK_API_LANG=fr \
    -e RUSTDESK_API_RUSTDESK_ID_SERVER=192.168.1.66:21116 \
    -e RUSTDESK_API_RUSTDESK_RELAY_SERVER=192.168.1.66:21117 \
    -e RUSTDESK_API_RUSTDESK_API_SERVER=http://192.168.1.66:21114 \
    -e RUSTDESK_API_RUSTDESK_KEY=<key> \
    lejianwen/rustdesk-api
    ```

2. Utiliser `docker compose` : voir [WIKI](https://github.com/lejianwen/rustdesk-api/wiki)

#### Lancer depuis le binaire release

[Télécharger la release](https://github.com/lejianwen/rustdesk-api/releases)

#### Installation depuis le code source

1. Cloner le dépôt
   ```bash
   git clone https://github.com/lejianwen/rustdesk-api.git
   cd rustdesk-api
   ```

2. Installer les dépendances

    ```bash
    go mod tidy
    # installer swag si vous voulez générer la documentation ; sinon, vous pouvez l’ignorer
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

3. Compiler le front d’administration. Le code du front est dans [rustdesk-api-web](https://github.com/lejianwen/rustdesk-api-web).
   ```bash
   cd resources
   mkdir -p admin
   git clone https://github.com/lejianwen/rustdesk-api-web
   cd rustdesk-api-web
   npm install
   npm run build
   cp -ar dist/* ../admin/
   ```
4. Exécution
    ```bash
    # exécution directe
    go run cmd/apimain.go
    # ou générer l’API puis exécuter
    go generate generate_api.go
    ```
   > Note : lors de l’utilisation de `go run` ou du binaire compilé, les dossiers `conf` et `resources` doivent être présents dans le répertoire courant.
   > Si vous exécutez ailleurs, vous pouvez spécifier un chemin absolu avec `-c` et la variable d’environnement `RUSTDESK_API_GIN_RESOURCES_PATH`, par exemple :
   > ```bash
   > RUSTDESK_API_GIN_RESOURCES_PATH=/opt/rustdesk-api/resources ./apimain -c /opt/rustdesk-api/conf/config.yaml
   > ```
5. Compiler : si vous souhaitez compiler vous-même, placez-vous à la racine du projet. Sous Windows, exécutez `build.bat`; sous Linux, exécutez `build.sh`. L’exécutable généré se trouve dans le dossier `release`.

6. Ouvrez le navigateur sur `http://<votre serveur[:port]>/_admin/`. Le nom d’utilisateur et le mot de passe par défaut sont `admin`. Changez-le rapidement.

#### Exécution avec l’image `lejianwen/server-s6`

- Problème de délai d’attente de connexion corrigé
- Peut forcer la connexion avant de lancer une session
- GitHub : https://github.com/lejianwen/rustdesk-server

```yaml
 networks:
   rustdesk-net:
     external: false
 services:
   rustdesk:
     ports:
       - 21114:21114
       - 21115:21115
       - 21116:21116
       - 21116:21116/udp
       - 21117:21117
       - 21118:21118
       - 21119:21119
     image: lejianwen/rustdesk-server-s6:latest
     environment:
       - RELAY=<relay_server[:port]>
       - ENCRYPTED_ONLY=1
       - MUST_LOGIN=N
       - TZ=Asia/Shanghai
       - RUSTDESK_API_RUSTDESK_ID_SERVER=<id_server[:21116]>
       - RUSTDESK_API_RUSTDESK_RELAY_SERVER=<relay_server[:21117]>
       - RUSTDESK_API_RUSTDESK_API_SERVER=http://<api_server[:21114]>
       - RUSTDESK_API_KEY_FILE=/data/id_ed25519.pub
       - RUSTDESK_API_JWT_KEY=xxxxxx # clé JWT
     volumes:
       - /data/rustdesk/server:/data
       - /data/rustdesk/api:/app/data # monter la base de données
     networks:
       - rustdesk-net
     restart: unless-stopped
```

## Divers

- [WIKI](https://github.com/lejianwen/rustdesk-api/wiki)
- [Problème de délai d’attente de connexion](https://github.com/lejianwen/rustdesk-api/issues/92)
- [Modifier l’ID du client](https://github.com/abdullah-erturk/RustDesk-ID-Changer)
- [Source du client web](https://hub.docker.com/r/keyurbhole/flutter_web_desk)

## Remerciements

Merci à tous ceux qui ont contribué !

<a href="https://github.com/lejianwen/rustdesk-api/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=lejianwen/rustdesk-api" />
</a>

## Merci pour votre soutien ! Si ce projet vous est utile, donnez une étoile ⭐️ pour encourager le projet.

[lejianwen/rustdesk-server]: https://github.com/lejianwen/rustdesk-server
