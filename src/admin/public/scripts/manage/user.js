
$(function ()
{
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
    }

    // del row event
    {
        edit_event_on('del_row', tr => {
            if (edit_op_get(tr) == 'add') return true;

            let rank = Number($('#userdata').attr('rank'));
            let _rank = edit_get_val(tr, 'rank');

            if (_rank <= rank) {
                err_show({err: 'no right to delete users of high-rank'});
                return false;
            }

            return true;
        });
    }

    // row ui rendering
    {
        let f = () => {
            edit_each_row(tr => {
                let user = $('#userdata').attr('user');
                let rank = Number($('#userdata').attr('rank'));

                let _id = edit_get_val(tr, '_id');
                let _rank = edit_get_val(tr, 'rank');

                if (_rank <= rank) {
                    if  (_id != user) edit_get_ctrl(tr, 'pwd').prop('disabled', true);
                    edit_get_ctrl(tr, 'rank').prop('disabled', true);
                    edit_get_ctrl(tr, 'memo').prop('disabled', true);
                } else {
                    tr.find(`td:nth-child(${edit_get_td_idx('grant')})`).html('<a href=javascript:;>点我授权</a>');
                }
            });
        };

        edit_event_on('load', f);
        edit_event_on('save', f);
    }

    // grant link click
    $('#pnl_edit table').on('click', `tbody td:nth-child(${edit_get_td_idx('grant')}) a`, function () {
        let tr = $(this).closest('tr');

        $('#dlg_priv').modal({
            backdrop: 'static',
        }, tr);
    });

    // priv dialog init
    $('#dlg_priv').on('show.bs.modal', function (e) {
        let tr = e.relatedTarget;
        let _id = edit_get_val(tr, '_id');

        let $this = $(this);
        $this.find('.modal-header .modal-title span').text(_id);

        // get priv data
        $.ajax({
			method: "post",
			url:    "/manage/user/priv_get",
			data:   `_id=${_id}`,
			dataType: "json",
		}).done(function (r) {
            if (!err_ok(r)) {
                $this.modal('hide');
                err_show({err: 'get priv failed'});
                return;
            }

            // make dict
            let dict = {};
            r.forEach(key => dict[key] = true);

            // highlight priv-btns
            $this.find('.modal-body button[type=item]').each(function () {
                let btn = $(this);

                if (dict[btn.attr('key')])
                    btn.addClass('sel');
                else
                    btn.removeClass('sel');
            });
		});
    });

    // priv dialog priv-btn click
    $('#dlg_priv .modal-body button').on('click', function () {
        $(this).toggleClass('sel');
    });

    // priv dialog OK
    $('#dlg_priv .modal-footer .btn-primary').on('click', function () {
        let dlg = $('#dlg_priv');

        // _id
        let _id = dlg.find('.modal-header .modal-title span').text();

        // parr
        let parr = [];
        dlg.find('.modal-body button.sel').each(function () {
            parr.push($(this).attr('key'));
        });

        // set
        $.ajax({
			method: "post",
			url:    "/manage/user/priv_set",
			data:   {_id: _id, pstr: parr.join(',')},
			dataType: "json",
		}).done(function (r) {
            err_show(r);
            $('#dlg_priv').modal('hide');
		});
    });
});
