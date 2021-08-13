#!/usr/bin/env node

// ============================================================================

var fs    = require('fs');
var path  = require('path');
var hbars = require('handlebars');

var svrinfo = require('./svrinfo');

// ============================================================================

function print(...args) {
    console.log(...args);
}

// ============================================================================

function help() {
    print(`Usage: ./deploy.js opt [args...]`);
    print();
    print("opt:");
    print("    -c          generate config.json");
    print("    -s  svr     generate SERVERS");
    print("    -f  svr     generate FILES");
    print("    -h  svr     query host info");
    print("    -d  db      query db info");
    print("    -n [svr]    query names");
    print("        '':        all servers");
    print("        webs:      all web servers");
    print("        routers:   all routers");
    print("        bats:      all battles");
    print("        games:     all games");
    print("        gates:     all gates");

    process.exit(1);
}

function get_svrinfo(name) {
    let si = svrinfo.index[name];
    if (!si) throw `server NOT found: ${name}`;

    return si
}

// ============================================================================

async function gen_config() {
    hbars.registerHelper('format', function(s, ...p) {
        p.pop();
        let i = 0;
        return s.replace(/%s/g, () => p[i++]);
    });
    hbars.registerHelper('add', function(...a) {
        a.pop();
        if (a.length == 0) return;

        let tp = typeof a[0];
        if (tp == 'string') {
            return a.join('');
        } else if (tp == 'number') {
            return a.reduce((a, v) => a + Number(v));
        }
    });

    let tpl = await fs.promises.readFile(path.join(__dirname, './config.json.tpl'), 'utf8');

    if (process.env.DEV_MODE == 'true') {
        tpl = tpl.replace('"dev_mode": false', '"dev_mode": true');
    }

    let r = hbars.compile(tpl, {noEscape: true})(svrinfo);
    print(r);
}

function gen_servers(name) {
    let si = get_svrinfo(name);
    if (!si.cmd) throw `server has NO cmd: ${name}`;

    print(`SERVERS=(`);
    print(`"${name.padEnd(12)}${si.cmd}"`);
    print(`)`);
}

function gen_files(name) {
    let si = get_svrinfo(name);
    if (!si.files) throw `server has NO files: ${name}`;

    print(`FILES="${si.files.join(' ')}"`);
}

function query_host(name) {
    let si = get_svrinfo(name);

    if (si.ip) print(si.ip);
}

function query_db(name) {
    let db = svrinfo.index.dbs[name];
    if (!db) throw `db NOT found: ${name}`;

    if (db.user && db.user != "") {
        print(`mongodb://${db.user}:${db.pwd}@${db.ip}:${db.port}/test?authSource=admin`);
    } else {
        print(`mongodb://${db.ip}:${db.port}/test`);
    }
}

function query_names(name) {
    let arr = [];

    if (!name || name == "webs") {
        arr.push("auth");
        arr.push("bill");
        arr.push("switcher");
        arr.push("reporter");
        arr.push("agent");
        arr.push("admin");
    }

    if (!name || name == "routers") {
        svrinfo.routers.forEach(v => arr.push(`router${v.id}`));
    }
    if (!name || name == "bats") {
        svrinfo.bats.forEach(v => arr.push(`bat${v.id}`));
    }
    if (!name || name == "games") {
        svrinfo.games.forEach(v => arr.push(`game${v.id}`));
    }
    if (!name || name == "gates") {
        svrinfo.gates.forEach(v => arr.push(`gate${v.id}`));
    }

    print(arr.join(' '));
}

// ============================================================================

(async () => {
    let args = process.argv.slice(2);

    let opt = args[0]
    let p1  = args[1]

    if (!opt) help();

    if (opt == "-c") {
        await gen_config();
    
    } else if (opt == "-s") {
        if (!p1) help();
        gen_servers(p1);
    
    } else if (opt == "-f") { 
        if (!p1) help();
        gen_files(p1);
    
    } else if (opt == "-h") { 
        if (!p1) help();
        query_host(p1);
    
    } else if (opt == "-d") { 
        if (!p1) help();
        query_db(p1);
    
    } else if (opt == "-n") { 
        query_names(p1);
    
    } else {
        help();
    }
})().catch(e => {
    console.error(typeof e == 'string' ? e : e.message);
    process.exit(1);
});
