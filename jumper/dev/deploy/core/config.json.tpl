{
    "common": {
{{#common}}
        "version": "{{version}}",

        "area": {
            "id":   {{area.id}},
            "name": "{{area.name}}"
        },

        "db_share":  "mongodb://{{format db_share  (add ../proj_name 's')}}",
        "db_center": "mongodb://{{format db_center (add ../proj_name 'c')}}",
        "db_stats":  "mongodb://{{format db_stats  (add ../proj_name 'st')}}",
        "db_log":    "mongodb://{{format db_log    (add ../proj_name 'lg')}}",
        "db_bill":   "mongodb://{{format db_bill   (add ../proj_name 'b')}}",
        "db_cross":  "mongodb://{{format db_cross  (add ../proj_name 'cr')}}",

        "db_user": {
{{#db_user}}
            "u{{id}}": "mongodb://{{format db (add ../../proj_name 'u' id)}}"{{#unless @last}},{{/unless}}
{{/db_user}}
        },

        "kfk": {
            "urls": [
{{#kfk.urls}}
                "{{.}}"{{#unless @last}},{{/unless}}
{{/kfk.urls}}
            ]
        },

        "redis": {
            "pool_size": {{redis.pool_size}},

            "urls": [
{{#redis.urls}}
                "{{.}}"{{#unless @last}},{{/unless}}
{{/redis.urls}}
            ]
        },

        "behind_proxy": {{behind_proxy}},

        "dev_mode": false
{{/common}}
    },

    "auth": {
{{#auth}}
        "ip":   "{{ip}}",
        "port": 5000,

        "threads": {{threads}}
{{/auth}}
    },

    "bill": {
{{#bill}}
        "ip":   "{{ip}}",
        "port": 6000,

        "threads": {{threads}}
{{/bill}}
    },

    "switcher": {
{{#switcher}}
        "ip":   "{{ip}}",
        "port": 7000,

        "threads": {{threads}},

        "token": "{{token}}"
{{/switcher}}
    },

    "reporter": {
{{#reporter}}
        "ip":   "{{ip}}",
        "port": 9000,

        "threads": {{threads}},

        "token": "{{token}}"
{{/reporter}}
    },

    "agent": {
{{#agent}}
        "ip":   "{{ip}}",
        "port": 9001
{{/agent}}
    },

    "admin": {
{{#admin}}
        "ip":   "{{ip}}",
        "port": 8000,

        "pwdfill": "{{pwdfill}}",

        "ccy": "{{ccy}}"
{{/admin}}
    },

    "routers": {
{{#routers}}
        "router{{id}}": {"addr4c": "{{ip}}:{{add 16000 id}}"}{{#unless @last}},{{/unless}}
{{/routers}}
    },

    "bats": {
{{#bats}}
        "bat{{id}}": {"ip": "{{ip}}", "port": {{add 17000 id}}, "threads": {{threads}} }{{#unless @last}},{{/unless}}
{{/bats}}
    },

    "games": {
{{#games}}
        "game{{id}}": {
            "addr4gw": "{{ip}}:{{add 10000 id}}",
            "svc":     "{{ip}}:{{add 13000 id}}",
            "db_game": "mongodb://{{format db (add ../proj_name '_' id)}}"
        }{{#unless @last}},{{/unless}}
{{/games}}
    },

    "gates": {
{{#gates}}
        "gate{{id}}": {
            "ip_wan": "{{ip_wan}}",
            "port":   {{add 20000 id}},
            "wsport": {{add 23000 id}}
        }{{#unless @last}},{{/unless}}
{{/gates}}
    }
}
