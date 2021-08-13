
var dbs     = require('../1-dbs');
var hosts   = require('../2-hosts');
var servers = require('../3-servers');

// ============================================================================

var svr_cmd = {
    auth: {
        files: ["auth", "node_modules", "gamedata"],
        cmd:   "./auth/www",
    },

    bill: {
        files: ["bill", "node_modules", "gamedata"],
        cmd:   "./bill/www",
    },

    switcher: {
        files: ["switcher", "node_modules"],
        cmd:   "./switcher/www",
    },

    reporter: {
        files: ["reporter", "node_modules"],
        cmd:   "./reporter/www",
    },

    agent: {
        files: ["agent", "node_modules"],
        cmd:   "./agent/www",
    },

    admin: {
        files: ["admin", "node_modules", "gamedata"],
        cmd:   "./admin/www",
    },

    router: {
        files: ["router"],
        cmd:   "./router -config config.json -server router$id -log router$id.log",
    },

    bat: {
        files: ["bat", "node_modules"],
        cmd:   "./bat/www -server bat$id",
    },

    game: {
        files: ["game", "gamedata"],
        cmd:   "./game -config config.json -server game$id -log game$id.log",
    },

    gate: {
        files: ["gate"],
        cmd:   "./gate -config config.json -server gate$id -log gate$id.log",
    },
}

// ============================================================================

function make_index() {
    let index = {
        dbs:   {},
        hosts: {},
    };

    // make
    dbs.forEach(v => index.dbs[v.name] = v);
    hosts.forEach(v => index.hosts[v.name] = v);

    index.auth     = servers.auth;
    index.bill     = servers.bill;
    index.switcher = servers.switcher;
    index.reporter = servers.reporter;
    index.agent    = servers.agent;
    index.admin    = servers.admin;

    servers.routers.forEach(v => index[`router${v.id}`] = v);
    servers.bats.forEach(v => index[`bat${v.id}`] = v);
    servers.games.forEach(v => index[`game${v.id}`] = v);
    servers.gates.forEach(v => index[`gate${v.id}`] = v);

    // set
    servers.index = index;
}

function make_db(obj, col) {
    let name = obj[col];
    if (!name) return;

    let db = servers.index.dbs[name];
    if (!db) throw `db NOT found: ${name}`;

    if (db.user && db.user != "") {
        obj[col] = `${db.user}:${db.pwd}@${db.ip}:${db.port}/%s?authSource=admin`;
    } else {
        obj[col] = `${db.ip}:${db.port}/%s`;
    }
}

function make_host(obj) {
    let name = obj.host;
    if (!name) return;

    let host = servers.index.hosts[name];
    if (!host) throw `host NOT found: ${name}`;

    obj.ip     = host.ip_lan;
    obj.ip_wan = host.ip_wan;
}

function make_cmd(obj, name) {
    if (obj.id) {
        name = name.match(/^\D+/);
    }

    let cmd   = svr_cmd[name].cmd;
    let files = svr_cmd[name].files;

    if (obj.id) {
        cmd = cmd.replace(/\$id/g, obj.id);
    }

    obj.cmd   = cmd;
    obj.files = files;
}

function make_server(name) {
    let obj = servers.index[name];

    make_db(obj, "db")
    make_host(obj)
    make_cmd(obj, name)
}

// ============================================================================

// index
make_index()

// common
make_db(servers.common, "db_share");
make_db(servers.common, "db_center");
make_db(servers.common, "db_stats");
make_db(servers.common, "db_log");
make_db(servers.common, "db_bill");
make_db(servers.common, "db_cross");
servers.common.db_user.forEach(v => make_db(v, "db"));

// servers
make_server("auth")
make_server("bill")
make_server("switcher")
make_server("reporter")
make_server("agent")
make_server("admin")

servers.routers.forEach(v => make_server(`router${v.id}`));
servers.bats.forEach(v => make_server(`bat${v.id}`));
servers.games.forEach(v => make_server(`game${v.id}`));
servers.gates.forEach(v => make_server(`gate${v.id}`));

// inject project name
servers.proj_name = process.env.PROJ_NAME;

// ============================================================================

module.exports = servers;
