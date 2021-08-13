
var config  = require('../../../config.json');
var dbpool  = require('../../lib/dbpool');
var utils   = require('../../lib/utils');
var gtab    = require('../../models/gtab');
var session = require('../../models/session');

// ============================================================================

const C_MaxRecords = 10000;

// ============================================================================

function make_default_query(conf) {
    let q = conf.q;

    // filters
    if (!q.fkey) {
        q.fkey = [];
        q.fval = [];
    } else if (!(q.fkey instanceof Array)) {
        q.fkey = [q.fkey];
        q.fval = [q.fval];
    }

    if (!conf.hide.fkey) {
        q.filters = [];
        for (let i = 0; i < q.fkey.length; i++) {
            if (q.fkey[i] == "--" || q.fval[i] == "") continue;
            q.filters.push([q.fkey[i], q.fval[i]]);
        }
    }

    // dr
    if (!conf.hide.dr) {
        if (!q.dr) {
            let now = new Date();
            q.dr = `${now.addDay(-10).toDateString()} ~ ${now.toDateString()}`;
        }
    }

    // grpstat
    q.grpstat = q.grpstat || '';

    // sort
    q.sort = q.sort || '';

    // ok
    return;
}

function make_col_index(conf) {
    let i = 0;
    conf.cols.forEach(obj => obj.index = i++);
}

function make_pipeline(conf) {
    let p = [];

    // p0
    if (conf.p0)
        p = p.concat(conf.p0);

    // filter: fixed
    p = p.concat(make_filter_fixed(conf));

    // filter0
    if (conf.filter0)
        p = p.concat(make_filter_n(conf, conf.filter0));

    // p1 (normalization)
    if (conf.p1)
        p = p.concat(conf.p1);

    // filter1
    if (conf.filter1)
        p = p.concat(make_filter_n(conf, conf.filter1));

    // group
    p = p.concat(make_group(conf));

    // apply sort
    p = p.concat(make_sort(conf));

    // limit total number of records
    p = p.concat([{$limit: C_MaxRecords}]);

    // return
    return p;
}

function make_filter_fixed(conf) {
    let q = conf.q;
    let p = [{$match:{}}];
    let b = false;

    // fixed filters
    if (!conf.hide.fkey) {
        utils.filter(p[0].$match, q.filters);
        b = true;
    }

    // dr
    if (!conf.hide.dr) {
        let arr = q.dr.split(" ~ ");
        let d_start = Date.fromString(arr[0]).startOfDay();
        let d_end   = Date.fromString(arr[1]).endOfDay();

        p[0].$match['ts'] = {$gte: d_start, $lt: d_end};
        b = true
    }

    return b ? p : [];
}

function make_filter_n(conf, filter) {
    let q = conf.q;
    let p = [{$match:{}}];

    let b = false;
    filter.forEach(obj => {
        let v = q[obj.key];
        if (!v) return;

        // check combo value
        if (obj.combo) {
            let r = v.match(/- ([^-]+)$/);
            if (r) v = r[1];
        }

        // tonumber
        if (obj.dt == 'n')
            v = Number(v);

        // add match
        p[0].$match[obj.key] = v;
        b = true;
    });

    return b ? p : [];
}

function make_group(conf) {
    // keep track of ordered result fields
    let fields = [];
    conf.fields = fields;

    // decode grpstat
    let gs = conf.q.grpstat.split('|');
    let sel_grps  = gs.length > 0 ? gs[0].split(',') : [];
    let sel_stats = gs.length > 1 ? gs[1].split(',') : [];

    // pipeline
    let p = [{$group:{}}, {$project:{}}];

    // group
    let _id = {};
    let tz = new Date().getTimezoneOffset() * 60 * 1000;

    conf.cols.forEach(obj => {
        if (obj.type != 'group') return;
        if (obj.selectable && !sel_grps.some(v => v == obj.key)) return;

        if (obj.isdate) {
            _id[obj.key] = {$dateToString: {format: "%Y-%m-%d", date: {$subtract: ["$"+obj.key, tz]}}};
        } else {
            _id[obj.key] = "$"+obj.key;
        }

        fields.push(obj);
    });

    p[0].$group['_id'] = _id;

    // stats
    conf.cols.forEach(obj => {
        if (obj.type != 'stat') return;
        if (obj.selectable && !sel_stats.some(v => v == obj.key)) return;

        p[0].$group[obj.key] = obj.op;

        fields.push(obj);
    });

    // project _id fields
    Object.keys(p[0].$group).forEach(k => {
        if (k != '_id') p[1].$project[k] = 1;
    });

    Object.keys(_id).forEach(k => {
        p[1].$project[k] = "$_id."+k;
    });

    // project set fields to size
    conf.cols.forEach(obj => {
        if (obj.type != 'stat') return;

        if (p[0].$group[obj.key] && obj.tosize)
            p[1].$project[obj.key] = {$size: "$"+obj.key};
    });

    // return
    return p;
}

function make_sort(conf) {
    if (conf.hide.sort)
        return [];

    // check result fields
    if (conf.fields.length == 0)
        return [];

    // pipeline
    let p = [{$sort:{}}];

    // decode sort
    let so = conf.q.sort.split(',');
    if (so.length != 2)
        so = [conf.fields[0].key, 1];

    p[0].$sort[so[0]] = Number(so[1]);

    // return
    return p;
}

async function make_tab(conf, p) {
    // check result fields
    if (conf.fields.length == 0) return;

    // get connection string
    let cnnstr = null;
    switch (conf.db) {
        case 'share':
            cnnstr = config.common.db_share;
            break;

        case 'stats':
            cnnstr = config.common.db_stats;
            break;

        case 'log':
            cnnstr = config.common.db_log;
            break;

        case 'bill':
            cnnstr = config.common.db_bill;
            break;
    }

    if (!cnnstr) return;

    // connect
    let docs = await dbpool.exec(cnnstr, async db => {
        return await db.collection(conf.coll).aggregate(p).toArray();
    });

    // make calc fields
    make_calc_field(conf);

    // tab
    conf.tab.header = conf.fields.map(v => v.text);
    conf.tab.body   = docs.map(v => conf.fields.map(obj => v[obj.key]));

    // make calc data
    make_calc_data(conf);

    // make dict data
    make_dict_data(conf);
}

function make_calc_field(conf) {
    // add calc fields
    conf.cols.forEach(obj => {
        if (obj.type != 'calc') return;

        let arr = obj.formula.match(/\$[\w_][\w\d_]*/g);
        if (!arr || arr.every(v => conf.fields.some(v2 => v2.key == v.slice(1))))
            conf.fields.push(obj);
    });

    // sort result fields
    conf.fields.sort((a, b) => a.index - b.index);
}

function make_calc_data(conf) {
    // prepare formulas
    let arr = [];
    for (let i = 0; i < conf.fields.length; i++) {
        let obj = conf.fields[i];
        if (obj.type != 'calc') continue;

        let formula = obj.formula.replace(/\$([\w_][\w\d_]*)/g, (_, key) => {
            let n = conf.fields.findIndex(v => v.key == key);
            return `row[${n}]`;
        });

        arr.push([i, formula]);
    }

    // fill calc data
    if (arr.length > 0)
        conf.tab.body.forEach(row => arr.forEach(v => row[v[0]] = Function('row', `return ${v[1]}`)(row)));
}

function make_dict_data(conf) {
    for (let i = 0; i < conf.fields.length; i++) {
        let obj = conf.fields[i];

        let col = conf.cols.find(v => v.key == obj.key);
        if (!col || !col.dict) continue;

        conf.tab.body.forEach(row => row[i] = gtab[`dict_${col.dict}`](row[i]));
    }
}

// ============================================================================

const C_filters = [
    ["--",      "--"],
    ["area:n",  "区域"],
    ["svr:s",   "服务器名"],
    ["sdk:re",  "Sdk"],
];

function render(req, res, conf) {
    let sess = session.data(req);
    let q    = req.query;

    // make conf
    if (!conf) conf = {
        /*
            ------------ default or generated -------------

            sess
            q
            combo:   [{key: 'field', combo: 'source'}, ...]
            tab:     {name, header, body}
            fields:  [the final result fields]
            scripts: []

            ------------ user defined -------------

            name:    will be set as tab.name

            db:      'share/stats/log/bill'
            coll:    ''

            p0:      pipeline0. []
            filter0: [{key: '', text: '', dt: 'n/s', combo: 'source'}]

            p1:      pipeline1. []
            filter1: [{key: '', text: '', dt: 'n/s', combo: 'source'}]

            cols:    [
                {type: 'group', key: 'field', text: 'display', dict: 'name', selectable: true, isdate: true},
                {type: 'stat',  key: 'field', text: 'display', dict: 'name', selectable: true, op: object, tosize: true},
                {type: 'calc',  key: 'field', text: 'display', dict: 'name', formula: '$a + $b'},
            ],

            form: prepended user-defined forms
            hide: ['fkey', 'dr', 'sort']

            on_pipeline(p): event fired when pipeline is constructed
        */
    };

    // session & q
    conf.sess = sess;
    conf.q    = q;

    // prepare hide option
    let hide = conf.hide || [];
    conf.hide = {};
    hide.forEach(v => conf.hide[v] = true);
    if (conf.hide.fkey) conf.hide.fval = true;

    // make default query
    make_default_query(conf);

    // form
    conf.form = conf.form || [];
    conf.form = conf.form.concat([
        {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[0], icon:  "filter"},
        {type: "input",     name: "fval", placeholder: "值", value: q.fval[0], style: "width: 120px"},
        {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[1], icon:  "filter"},
        {type: "input",     name: "fval", placeholder: "值", value: q.fval[1], style: "width: 120px"},
        {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[2], icon:  "filter"},
        {type: "input",     name: "fval", placeholder: "值", value: q.fval[2], style: "width: 120px"},
        {type: "daterange", name: "dr",                      value: q.dr,      icon:  "calendar-alt"},

        {type: "hidden",    name: "grpstat",                 value: q.grpstat},
        {type: "hidden",    name: "sort",                    value: q.sort},
    ]);

    conf.form = conf.form.filter(v => !conf.hide[v.name]);

    // combo
    conf.combo = [];

    // scripts
    conf.scripts = ['/scripts/fx/fx.js'];

    // name
    conf.tab = {name: conf.name};

    // filter
    [conf.filter0, conf.filter1].forEach(filter => {
        if (!filter) return;

        let combo = [];

        let form = filter.map(v => {
            let obj = {type: 'input', name: v.key, placeholder: v.text, value: q[v.key]};
            if (v.combo)
                combo.push({key: v.key, combo: v.combo});

            return obj;
        });

        conf.form  = conf.form.concat(form);
        conf.combo = conf.combo.concat(combo);
    });

    // make col index
    make_col_index(conf);

    // make pipeline
    let p = make_pipeline(conf);

    // event
    if (conf.on_pipeline)
        conf.on_pipeline(p);

    // make tab & render
    make_tab(conf, p).then(r => res.render('fx/fx', conf)).catch(err => {
        res.status(500).end('fx error. ask admin for help');
        console.error(err);
    });
}

// ============================================================================

module.exports = {
    render: render,
}
