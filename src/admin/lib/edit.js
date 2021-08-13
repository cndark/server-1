
require('./ext');

var dbpool  = require('./dbpool');
var utils   = require('./utils');
var session = require('../models/session');

// ============================================================================

function make_q(q) {
    // filter
    if (!q.fkey) {
        q.fkey = [];
        q.fval = [];
    } else if (!(q.fkey instanceof Array)) {
        q.fkey = [q.fkey];
        q.fval = [q.fval];
    }

    q.filters = [];
    for (let i = 0; i < q.fkey.length; i++) {
        if (q.fkey[i] == "--" || q.fval[i] == "") continue;
        q.filters.push([q.fkey[i], q.fval[i]]);
    }

    // page
    q.page = tonumber(q.page);
    if (q.page < 1) q.page = 1;

    // ok
    return q;
}

// ============================================================================

function show(req, res, opt) {
    /*
        view:    extended view of edit.pug if specified
        data:    custom data for rendering

        db:      db name in dbpool
        coll:    collection name
        autoid:  uses a sequence for _id generation
        rpp:     records per page
        sort:    sort-object
        scripts: array of script-files to be included

        -------------------------------

        on_load(docs):
            doc tranforming is allowed

        on_before_save(docs):
            doc tranforming is allowed. throw exceptions to reject the save.
            __edit_cond__ of each row: query condition for 'up' and 'del'.

        on_after_save():
            do post actions

        -------------------------------

        cols: [
            {col, name, dtype, ctrl, src, format, formula, immutable, optional, unbind, width},
        ]

            '_id' is the primary key and MUST be defined

            dtype:     number, string
            ctrl:      input, select, hidden, html
            src:       data source for 'select, html' ctrl
            format:    regexp
            formula:   e.g.  $col1 + $col2
            immutable: can NOT be modified
            unbind:    dont save to db
            width:     e.g.  20px, 10%

        -------------------------------

        filters: [
            "col:n",    --> number
            "col:s",    --> string
            "col:re",   --> regexp
        ]

        -------------------------------

        form: [
            {type, name, src, value, desc, icon, placeholder, style},
            or plain-string,
        ]
    */

    switch (req.method) {
        case "GET":
            render(req, res, opt);
            break;

        case "POST":
            let op = req.query.op;
            if (op == "load")
                load(req, res, opt).catch(console.error);
            else if (op == "save")
                save(req, res, opt).catch(console.error);
            else if (op == "autoid")
                autoid(req, res, opt).catch(console.error);
            else
                res.status(404).end();

            break;

        default:
            res.status(404).end();
    }
}

function render(req, res, opt) {
    let sess = session.data(req);
    let q    = make_q(req.query);

    // make col index
    let col_index = {};
    opt.cols.forEach(v => col_index[v.col] = v);

    // check
    if (!col_index['_id']) {
        res.end("fatal: primary key '_id' is NOT defined");
        return;
    }

    if (opt.autoid && col_index._id.formula) {
        res.end("fatal: formula of '_id' should NOT be set in autoid mode");
        return;
    }

    if (opt.autoid && col_index._id.dtype != "number") {
        res.end("fatal: '_id' MUST be of number type in autoid mode");
        return;
    }

    for (let i = 0; i < opt.cols.length; i++) {
        let c = opt.cols[i];

        if (c.ctrl == 'hidden' && !c.formula) {
            if (!(c.col == '_id' && opt.autoid)) {
                res.end(`fatal: 'formula' MUST be set for hidden col ${c.col}`);
                return;
            }
        }
    }

    // make filters
    let filters = (opt.filters || []).map(v => {
        let arr = v.split(':');
        let c = col_index[arr[0]];
        return c ? [v, c.name] : ["--", "--"];
    })

    // make form
    let form = opt.form || [];

    if (filters.length > 0) {
        filters = [["--", "--"]].concat(filters);

        let n = filters.length - 1;
        if (n > 3) n = 3;

        for (let i = 0; i < n; i++) {
            form.push({type: "select", name: "fkey", src: filters,      value: q.fkey[i], icon: "filter"});
            form.push({type: "input",  name: "fval", placeholder: "值", value: q.fval[i], style: "width: 150px"});
        }
    }

    // make scripts
    let scripts = opt.scripts || [];

    // render opt
    let ropt = {
        autoid: opt.autoid,
        rpp:    opt.rpp,
        cols:   opt.cols,
    };

    // render
    res.render(opt.view || 'lib/table/edit', {
        sess:    sess,
        form:    form,
        filters: filters.length > 0,
        opt:     ropt,
        scripts: scripts,
        data:    opt.data,
    });
}

async function load(req, res, opt) {
    let q = make_q(req.body);

    let cond = utils.filter({}, q.filters);
    let rpp  = opt.rpp;
    let skip = (q.page - 1) * rpp;

    try {
        let db = dbpool.get(opt.db);

        let docs = await db.collection(opt.coll).find(cond).sort(opt.sort).skip(skip).limit(rpp).toArray();

        // event
        if (opt.on_load) {
            opt.on_load(docs);
        }

        res.json(docs).end();
    } catch(e) {
        res.json({err: e.message}).end();
    }
}

async function save(req, res, opt) {
    let q = req.body;

    try {
        // event
        if (opt.on_before_save) {
            opt.on_before_save(q);
        }

        // create bulk
        let db = dbpool.get(opt.db);

        let bulk = db.collection(opt.coll).initializeOrderedBulkOp();
        q.forEach(v => {
            // get op
            let op = v.__edit_op__;
            delete v.__edit_op__;

            // get cond (server side)
            let cond = v.__edit_cond__ || {};
            delete v.__edit_cond__;
            cond._id = v._id;

            // add to bulk
            switch (op) {
                case "add":
                    bulk.insert(v);
                    break;

                case "up":
                    bulk.find(cond).updateOne({$set: v});
                    break;

                case "del":
                    bulk.find(cond).deleteOne();
                    break;
            }
        });

        // check
        if (bulk.length == 0) throw 'nothing to update';

        // save
        try {
            await bulk.execute();
        } catch (e) {
            if (e.message.match(/^E11000/))
                throw 'primary key already exists';
            else
                throw e;
        }

        // event
        if (opt.on_save) {
            opt.on_save();
        }

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}

async function autoid(req, res, opt) {
    let q = req.body;

    let n = tonumber(q.n);
    if (n < 1 || n > 1000) {
        res.json({err: 'invalid args'}).end();
        return;
    }

    try {
        let db = dbpool.get(opt.db);

        let r = await db.collection('edit_autoid').findOneAndUpdate(
            {_id: opt.coll},
            {$inc: {seq: n}},
            {projection: {seq: 1}, upsert: true, returnOriginal: false},
        );

        res.json({from: r.value.seq - n + 1, to: r.value.seq}).end();
    } catch(e) {
        res.json({err: 'fill autoid failed'}).end();
    }
}

// ============================================================================

module.exports = {
    show: show,
}
