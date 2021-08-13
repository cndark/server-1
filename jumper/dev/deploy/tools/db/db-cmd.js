#!/usr/bin/env node

// ============================================================================

var dbs     = require('../../1-dbs');
var hosts   = require('../../2-hosts');
var servers = require('../../3-servers');

var PROJ_NAME = process.env.PROJ_NAME;

// ============================================================================

function print(...args) {
    console.log(...args);
}

// ============================================================================

function db_entries() {
    // entries
    let entries = [
        {db: servers.common.db_share,  name: `${PROJ_NAME}s`},
        {db: servers.common.db_center, name: `${PROJ_NAME}c`},
        {db: servers.common.db_stats,  name: `${PROJ_NAME}st`},
        {db: servers.common.db_cross,  name: `${PROJ_NAME}cr`},
        {db: servers.common.db_bill,   name: `${PROJ_NAME}b`},
    ];

    servers.common.db_user.forEach(v => entries.push({db: v.db, name: `${PROJ_NAME}u${v.id}`}));
    servers.games.forEach(v => entries.push({db: v.db, name: `${PROJ_NAME}_${v.id}`}));

    // bind db object
    let db_index = {};
    dbs.forEach(v => db_index[v.name] = v);

    entries.forEach(v => {
        let obj = db_index[v.db];
        if (!obj) throw `db NOT found: ${v.db}`;

        v.db = obj;
    });

    // ok
    return entries;
}

// ============================================================================

function get_cmd_dump() {
    let arr = [];

    db_entries().forEach(v => {
        arr.push(`mongodump -h ${v.db.ip} --port ${v.db.port} ` +
                 `${v.db.user ? `-u ${v.db.user} -p ${v.db.pwd} --authenticationDatabase=admin ` : ''}` +
                 `--db=${v.name} --excludeCollectionsWithPrefix=replay --gzip --archive=${v.name}.db`);
    });

    return arr;
}

function get_cmd_restore() {
    let arr = [];

    db_entries().forEach(v => {
        arr.push(`mongorestore -h ${v.db.ip} --port ${v.db.port} ` +
        `${v.db.user ? `-u ${v.db.user} -p ${v.db.pwd} --authenticationDatabase=admin ` : ''}` +
                 `--drop --gzip --archive=${v.name}.db`);
    });

    return arr;
}

// ============================================================================

function avail_hosts() {
    return hosts.filter(v => v.name.match(/m_gs/)).map(v => v.ip_lan);
}

function alloc_host(arr) {
    let ips = avail_hosts();
    let L   = ips.length;

    if (L == 0) throw `no backup machines are found`;

    return arr.map((v, i) => `${ips[i % L]} ${v}`);
}

function help() {
    print(`Usage: ${process.argv[1]} options`);
    print("options:");
    print("    --nhost      available host count");
    print("    -d           dump commands");
    print("    -r           restore commands");

    process.exit(1);
}

// ============================================================================

(async () => {
    let args = process.argv.slice(2);

    if (args.length < 1) help();

    if (args[0] == "--nhost") {
        print(avail_hosts().length);

    } else if (args[0] == "-d") {
        let lines = alloc_host(get_cmd_dump()).join('\n');
        print(lines);

    } else if (args[0] == "-r") {
        let lines = alloc_host(get_cmd_restore()).join('\n');
        print(lines);

    } else {
        help();
    }
})().catch(console.error);
