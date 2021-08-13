
var edit_opt;
var edit_html_tab, edit_html_row;

var edit_del_rows  = null;
var edit_errmsg    = null;
var edit_last_page = true;

// ============================================================================

function edit_set_error(obj) {
    obj.addClass("err");
    setTimeout(function() {
        obj.removeClass("err");
    }, 1000);

    if (!edit_errmsg)
        edit_errmsg = "输入存在错误，请检查";
}

function edit_check_number(obj, e) {
    let v = obj.val().trim();

    if (v == "" && !e.noedit) {
        edit_set_error(obj);
        return 0;
    }

    v = Number(v);
    if (Number.isNaN(v)) {
        edit_set_error(obj);
        return 0;
    }

    return v;
}

function edit_check_string(obj, e) {
    let v = obj.val().trim();

    if (v == "" && !e.noedit && !e.optional) {
        edit_set_error(obj);
        return "";
    }

    if (e.format && v != "" && !v.match(new RegExp(e.format))) {
        edit_set_error(obj);
        return "";
    }

    return v;
}

// ============================================================================

function edit_gen_ui(obj) {
    let table = $('#pnl_edit table');
    table.html(edit_html_tab);

    let last_tr = table.find('tbody tr:last-child');
    obj.forEach(r => {
        let tr = edit_row_insert(last_tr);

        edit_opt.uicols.forEach(e => {
            if (e.unbind) return;

            edit_get_ctrl(tr, e).val(r[e.col]);
        });

        edit_calc_formula(tr);

        if (edit_opt.autoid)
            edit_set_val(tr, '_id', r._id, true);
    });

    edit_op_init();
}

function edit_gen_obj() {
    let obj = [];

    // add, update
    edit_each_row(tr => {
        let op = edit_op_get(tr);
        if (op == "add" || op == "up") {
            let rec = {};

            edit_opt.uicols.forEach(e => {
                if (e.unbind) return;

                let ctrl = edit_get_ctrl(tr, e);

                switch (e.dtype) {
                    case "number":
                        rec[e.col] = edit_check_number(ctrl, e);
                        break;

                    case "string":
                        rec[e.col] = edit_check_string(ctrl, e);
                        break;
                }
            });

            edit_opt.hdcols.forEach(e => {
                if (e.unbind) return;

                let v = tr.attr(`edit-${e.col}`);

                switch (e.dtype) {
                    case "number":
                        rec[e.col] = Number(v);
                        break;

                    case "string":
                        rec[e.col] = v;
                        break;
                }
            });

            rec.__edit_op__ = op;
            obj.push(rec);
        }
    });

    // delete
    edit_op_deleted().forEach(pk => {
        obj.push({__edit_op__: "del", _id: pk});
    });

    return obj;
}

function edit_row_insert(tr) {
    tr.before(edit_html_row);

    let tr_new = tr.prev();

    edit_opt.uicols.forEach(e => {
        if (e.noedit)
            tr_new.find(`td:nth-child(${e.idx + 1}) ${e.ctrl}`).prop('disabled', true);
    });

    edit_calc_formula(tr_new);
    edit_op_set(tr_new, 'add');

    return tr_new;
}

function edit_row_del(tr) {
    tr.remove();
}

function edit_each_row(f) {
    $('#pnl_edit table tbody tr').each(function () {
        let tr = $(this);
        if (tr.find('td:first-child button').text() == "+") return;

        f(tr);
    });
}

function edit_get_ctrl(tr, col) {
    let e = col;
    if (typeof e == 'string') {
        e = edit_opt.col_index[e];
        if (!e) return;
    }
    if (e.ctrl == 'hidden') return;

    return tr.find(`td:nth-child(${e.idx + 1}) ${e.ctrl}`);
}

function edit_get_td_idx(col) {
    let e = edit_opt.col_index[col];
    if (!e || e.ctrl == 'hidden') return;

    return e.idx + 1;
}

function edit_get_val(tr, col) {
    let e = edit_opt.col_index[col];
    if (!e) return;

    let v = e.ctrl == 'hidden'
        ? tr.attr(`edit-${e.col}`)
        : edit_get_ctrl(tr, e).val().trim();

    switch (e.dtype) {
        case "number":
            return Number(v);

        case "string":
            return v;
    }
}

function edit_set_val(tr, col, v, rawset) {
    let e = edit_opt.col_index[col];
    if (!e) return;
    if (e.formula && !rawset) return;

    if (e.ctrl == 'hidden') {
        tr.attr(`edit-${e.col}`, v);
    } else {
        obj = edit_get_ctrl(tr, e);
        obj.val(v);
        if (!rawset)
            obj.trigger('change');
    }
}

function edit_calc_formula(tr) {
    edit_opt.fmcols.forEach(e => {
        try {
            let r = Function('return ' +
                e.formula.replace(/\$([\w_]+)/g, (m, c) => {
                    let v = edit_get_val(tr, c);
                    return v === undefined ? m : (edit_opt.col_index[c].dtype == 'string' ? `'${v}'` : v);
                })
            )();

            edit_set_val(tr, e.col, r, true);

        } catch(e) {}
    });
}

function edit_commit_change() {
    edit_op_init();
    $('.chg').removeClass('chg');

    edit_event_fire('save');
}

function edit_op_init() {
    edit_each_row(tr => {
        tr.attr('op', null);

        edit_opt.uicols.forEach(e => {
            if (e.col == '_id' || e.immutable)
                edit_get_ctrl(tr, e).prop('disabled', true);
        });
    });

    edit_del_rows = [];
}

function edit_op_set(tr, op) {
    switch (op) {
        case "add":
            tr.attr('op', op);
            break;

        case "del":
            if (tr.attr('op') != "add")
                edit_del_rows.push(edit_get_val(tr, '_id'));
            break;

        case "up":
            if (!tr.attr('op'))
                tr.attr('op', op);
            break;

        default:
            return;
    }
}

function edit_op_get(tr) {
    return tr.attr('op');
}

function edit_op_deleted() {
    return edit_del_rows;
}

function edit_on_change(tr) {
    edit_calc_formula(tr);
    edit_op_set(tr, 'up');
}

// ============================================================================

function edit_load(page) {
    edit_last_page = true;

    $.ajax({
        method: "post",
        url:    `${location.pathname}?op=load&page=${page}`,
        data:   $('#pnl_edit form').find('input, select').serialize(),
        dataType: "json",
    }).done(function (r) {
        if (err_ok(r)) {
            edit_gen_ui(r);

            $('.pager .pageno').text(page);
            edit_last_page = r.length < edit_opt.rpp;

            edit_event_fire('load');
        }
    }).fail(function () {
        err_show(ErrNet);
    });
}

function edit_save(obj, cb) {
    if (edit_opt.autoid) {
        edit_fill_autoid(obj, (r) => {
            if (err_ok(r))
                edit_real_save(r, cb);
            else
                cb(r);
        });
    } else {
        edit_real_save(obj, cb);
    }
}

function edit_fill_autoid(obj, cb) {
    let n = 0;
    obj.forEach(v => {if (v.__edit_op__ == "add") n++});
    if (n == 0) {
        cb(obj);
        return;
    }

    $.ajax({
        method: "post",
        url: `${location.pathname}?op=autoid`,
        data: `n=${n}`,
        dataType: "json",
    }).done(function (r) {
        if (!err_ok(r)) {
            cb(r);
            return;
        }

        // fill
        edit_each_row(tr => {
            let op = edit_op_get(tr);
            if (op == "add") {
                edit_set_val(tr, '_id', r.from++, true);
                edit_calc_formula(tr);
            }
        });

        // re-gen obj
        obj = edit_gen_obj();
        cb(obj);

    }).fail(function () {
        cb(ErrNet);
    });
}

function edit_real_save(obj, cb) {
    $.ajax({
        method: "post",
        url: `${location.pathname}?op=save`,
        contentType: 'application/json; charset=UTF-8',
        data: JSON.stringify(obj),
        dataType: "json",
    }).done(function (r) {
        if (err_ok(r))
            edit_commit_change();

        cb(r);

    }).fail(function () {
        cb(ErrNet);
    });
}

// ============================================================================

var edit_events = {};

function edit_event_on(evt, f) {
    let arr = edit_events[evt];
    if (!arr) {
        arr = [];
        edit_events[evt] = arr;
    }

    arr.push(f);
}

function edit_event_fire(evt, ...args) {
    let arr = edit_events[evt];
    if (!arr) return;

    let ret;
    arr.forEach(f => ret = f(...args));
    return ret;
}

// ============================================================================

function edit_create(opt) {
    // prepare
    opt.col_index = {};
    let i = 0;

    opt.cols.forEach(e => {
        opt.col_index[e.col] = e;

        // set ui col index
        if (e.ctrl != "hidden")
            e.idx = ++i;

        // html ctrl is 'unbind'
        if (e.ctrl == "html")
            e.unbind = true;
    });

    opt.uicols = opt.cols.filter(e => e.ctrl != "hidden");
    opt.hdcols = opt.cols.filter(e => e.ctrl == "hidden");
    opt.fmcols = opt.cols.filter(e => e.formula);

    // formula cols are NOT editable
    opt.fmcols.forEach(e => {
        e.noedit = true;
    });
    if (opt.autoid) opt.col_index._id.noedit = true;

    // set
    edit_opt = opt;

    // html
    edit_html_tab = `
        <thead><tr>
        <th style='min-width: 90px; width: 90px'>删除</th>
        ${edit_opt.uicols.map(e=>`<th${e.width ? ` style='min-width: ${e.width}'` : ''}>${e.name}</th>`).join('')}
        </tr></thead>
        <tbody><tr><td><button class='btn btn-default'>+</button></td>${edit_opt.uicols.map(e=>`<td></td>`).join('')}</tr></tbody>
    `;

    edit_html_row = `
        <tr>
        <td><button class='btn btn-default'>-</button></td>
        ${edit_opt.uicols.map(e => {
            switch (e.ctrl) {
                case "input":
                    return `<td><input class='form-control' type='text'></td>`;

                case "select":
                    return `<td><select class='form-control'>
                           ${e.src.map(v=>`<option value='${v[0]}'>${v instanceof Array ? v[1] : v[0]}</option>`).join('')}
                           </select></td>`;

                case "html":
                    return `<td>${e.src}</td>`;
            }
        }).join('')}
        </tr>
    `;

    // initial refresh
    edit_load(1);
}

// ============================================================================

$(function () {
    // add .form-control
    {
        $('select, input, textarea').addClass("form-control");
    }

    // add/remove
    $('#pnl_edit table').on('click', 'td:first-child button', function () {
        let $this = $(this);
        let tr = $this.closest('tr');

        if ($this.text() == "+") {
            // add row
            let tr_new = edit_row_insert(tr);
            edit_event_fire('new_row', tr_new);
        } else {
            // del row
            let b = edit_event_fire('del_row', tr);
            if (!b) return;

            bs_confirm(
                '确定删除该行 ?',
            ).ok(() => {
                edit_op_set(tr, 'del');
                edit_row_del(tr);
            }).before(() => {
                tr.addClass('warn');
            }).after(() => {
                tr.removeClass('warn');
            });
        }
    });

    // data change
    $('#pnl_edit table').on('change', 'input, select', function () {
        let tr = $(this).closest('tr');

        edit_on_change(tr);

        $(this).addClass('chg');
    });

    // save
    $('#btn_save').on('click', function () {
        edit_errmsg = null;
        let obj = edit_gen_obj();
        if (obj.length == 0) return;
        if (edit_errmsg) {
            err_show({err: edit_errmsg});
            return;
        }

        let $this = $(this);

        $this.prop('disabled', true)

        edit_save(obj, (r) => {
            $this.prop('disabled', false)
            err_show(r);
        });
    });

    // filter
    $('#pnl_edit form').on('submit', function (e) {
        edit_load(1);
        e.preventDefault();
    });

    // paging
    $('#page_prev').on('click', function () {
        let page = Number($('.pager .pageno').text());

        if (page <= 1) return;

        page--;
        edit_load(page);
    });

    $('#page_next').on('click', function () {
        let page = Number($('.pager .pageno').text());

        if (edit_last_page) return;

        page++;
        edit_load(page);
    });

});
