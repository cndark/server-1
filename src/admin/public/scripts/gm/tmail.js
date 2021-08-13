
var paged_opt;

function paged_create(opt) {
    paged_opt = opt;
}

// ============================================================================

$(function () {
    // add page input
    $('form').append("<input type='hidden' name='page'>");

    // paging
    $('#page_prev').on('click', function () {
        let page = Number($('.pager .pageno').text());

        if (page <= 1) return;

        page--;
        $('form input[name=page]').val(page);
        $('form').submit();
    });

    $('#page_next').on('click', function () {
        let page = Number($('.pager .pageno').text());

        if (paged_opt.rec_n < paged_opt.rpp) return;

        page++;
        $('form input[name=page]').val(page);
        $('form').submit();
    });
});

// ============================================================================

$(function () {
    // add .form-control
    {
        $('select, input, textarea').addClass("form-control");
    }

    // field color
    {
        // status
        $('table tbody td:nth-child(4)').each(function () {
            if ($(this).text() == "failed")
                $(this).css('color', '#ff0000');
        });

        // target
        $('table tbody td:nth-child(2)').each(function () {
            if ($(this).text() == "ALL")
                $(this).css('color', '#ff00ff');
        });
    }

    // fetch init data: res
    {
        $.ajax({
            method: "post",
            url:    "/gm/gtab",
            data:   "key=res",
            dataType: "json",
        }).done(function (r) {
            let combo = $('.panel-body input[name=res_k]');

            combo.typeahead({
                source: r.map(e => `${e[1]} - ${e[0]}`),
            });

            combo.on('change', function () {
				let $this = $(this);
                let e = $this.typeahead('getActive');
                if (!e || e != $this.val())
                    $this.val("");
			});
        });
    }

    // link-button
    {
        $('table tbody tr td:nth-child(7)').html('<a href="javascript:;">详细</a>');
        $('table tbody tr td:nth-child(8)').html('<a href="javascript:;">删除</a>');

        // detail
        $('table tbody tr td:nth-child(7) a').on('click', function () {
            let tr = $(this).closest('tr');
            let btn_update = $('#btn_update');
            let btn_audit = $('#btn_audit');
            let panel = btn_update.closest('.panel').find('.panel-body');

            let id = tr.attr('_id');

            $.ajax({
                method: "post",
                url: "/gm/tmail/detail",
                data: {_id: id},
                dataType: "json",
            }).done(function (r) {
                err_show(r);

                if (err_ok(r)) {
                    btn_update.text('更新');
                    btn_update.prop('disabled', r.status != "wait");

                    if (r.status == 'wait' && r.audit == '-' && btn_audit.attr('priv') == 'true') {
                        btn_audit.show();
                    } else {
                        btn_audit.hide();
                    }

                    panel.find('[name=_id]').val(r._id);
                    panel.find('[name=send_ts]').val(r.send_ts);
                    panel.find('[name=area]').val(r.area);
                    panel.find('[name=target]').val(r.target);
                    panel.find('[name=title]').val(r.title);
                    panel.find('[name=text]').val(r.text);

                    let arr_k = panel.find('[name=res_k]');
                    let arr_v = panel.find('[name=res_v]');

                    for (let i = 0; i < 10; i++) {
                        $(arr_k.get(i)).val(r.res_k[i]);
                        $(arr_v.get(i)).val(r.res_v[i]);
                    }
                }
            }).fail(function () {
                err_show(ErrNet);
            });
        });

        // delete
        $('table tbody tr td:nth-child(8) a').on('click', function () {
            let tr = $(this).closest('tr');

            let id = tr.attr('_id');
            let status = tr.find('td:nth-child(4)').text();

            if (status != "wait" && status != "audit") {
                err_show({err: "已发送邮件不能被删除"});
                return;
            }

            bs_confirm(
                '确定删除记录 ?',
            ).ok(() => {
                $.ajax({
                    method: "post",
                    url: "/gm/tmail/delete",
                    data: {_id: id},
                    dataType: "json",
                }).done(function (r) {
                    err_show(r);

                    if (err_ok(r))
                        location.reload();

                }).fail(function () {
                    err_show(ErrNet);
                });
            }).before(() => {
                tr.addClass('warn');
            }).after(() => {
                tr.removeClass('warn');
            });
        });
    }

    // button: add
    $('#btn_add').on('click', function () {
        let btn_update = $('#btn_update');
        let panel = btn_update.closest('.panel').find('.panel-body');

        btn_update.text('新增');
        btn_update.prop('disabled', false);

        panel.find('[name=_id]').val(`${Date.now()}-${Math.floor(Math.random() * 100000)}`);

        err_show('在这里新增');
    });

    // button: update
    $('#btn_update').on('click', function () {
        let btn_update = $(this);
        let panel = btn_update.closest('.panel').find('.panel-body');

        let isnew = btn_update.text() == '新增';
        let id = panel.find('input[name=_id]').val();
        let area = panel.find('select[name=area]').val();
        let target = panel.find('input[name=target]').val();

        if (area == -1 && target != '') {
            err_show({err: '所有区域时不需要目标, 这是全区域邮件'});
            return;
        }

        bs_confirm(
            area == -1 || target == 'all',
            'Global-Mail: are you sure to add/update this global mail ?',
        ).ok(() => {
            $.ajax({
                method: "post",
                url: `/gm/tmail/update?isnew=${isnew}`,
                data: `${panel.find('select, input, textarea').serialize()}&_id=${id}`,
                dataType: "json",
            }).done(function (r) {
                err_show(r);

                if (err_ok(r))
                    location.reload();

            }).fail(function () {
                err_show(ErrNet);
            });
        });
    });

    // button: audit
    $('#btn_audit').on('click', function () {
        let btn_audit = $(this);
        let panel = btn_audit.closest('.panel').find('.panel-body');

        let id = panel.find('input[name=_id]').val();

        $.ajax({
            method: "post",
            url: '/gm/tmail/audit',
            data: `_id=${id}`,
            dataType: "json",
        }).done(function (r) {
            err_show(r);

            if (err_ok(r))
                location.reload();

        }).fail(function () {
            err_show(ErrNet);
        });
    });
});
