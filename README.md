# KeeneticToMqtt
KeeneticToMqtt is a bridge beetwen keeneticOS and mqtt. It allows you to do some api calls through mqtt commands.

Available features are:
- choosing internet policy (for example turn on wireguard) for keenetic clients.
- permit or disallow internet access for keenetic clients.

## <a name="home_assistant_addon"></a>Home Assistant addon
### <a name="home_assistant_addon_installation"></a> Installation

Supervisor > Add-on Store > ![image](https://user-images.githubusercontent.com/45158965/126977982-fc0a743c-68d9-4034-99aa-28011a3431ab.png) > Repositories

Add https://github.com/BlenderistDev/homeassistant-addons to your addons

Install addon keeneticToMqtt as usual.

More information: https://www.home-assistant.io/common-tasks/os#installing-third-party-add-ons

## Config options

### Config example
```
keenetic:
  host: 192.168.0.1
  login: login
  password: password
mqtt:
  host: mqtt://localhost:1883
  login: login
  password: password
  clientId: keeneticToMqtt
  baseTopic: keeneticToMqtt
homeassistant:
  deviceId: keeneticToMqtt
  updateInterval: 10s
  whitelist: ['00:00:00:00:00:00']
```
### keenetic
- host - keenetic host. Usually like http://192.168.0.1.
- login - keenetic user with api access. [more info](https://help.keenetic.com/hc/en-us/articles/360015786580-How-to-regain-access-to-the-web-interface).
- password - password for keenetic user.
  
### mqtt
- host - mqtt server host.
- login - mqtt user username.
- password - mqtt user password.
- clientId - mqtt client id.
- baseTopic - keeneticToMqtt mqtt base topic, if empty "keeneticToMqtt" will be used.

### homeassistant
- deviceId - home assistant device id
- updateInterval - home assistant entities update interval. You need to add unit, for example:
  - `10s` for 10 seconds.
  - `1m` for 1 minute.
- whitelist - list of mac addresses to handle.
