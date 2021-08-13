
function load() {
    $.ajax({
        method: "post",
        url: '/server/settings/load',
        data: $('select[name=area]').serialize(),
        dataType: "json",
    }).done(function (r) {
        r.closereg_hwater = r.closereg_hwater || '';
        r.closereg_limit = r.closereg_limit || '';
        r.opennew_mode = r.opennew_mode || 'manual';
        r.opennew_itv = r.opennew_itv || '';

        let ctx = $('.panel-body');

        ctx.find('[name=closereg_hwater]').val(r.closereg_hwater);
        ctx.find('[name=closereg_limit]').val(r.closereg_limit);
        ctx.find('[name=opennew_mode]').val(r.opennew_mode);
        ctx.find('[name=opennew_itv]').val(r.opennew_itv);
    }).fail(function () {
        err_show(ErrNet);
    });
}

function save_check() {
    let ctx = $('.panel-body');

    let chk_int = name => {
        let v = ctx.find(`[name=${name}]`).val();
        if (v == '') return true;

        v = Number(v);
        return !isNaN(v) && v > 0;
    };

    if (!chk_int('closereg_hwater')) return false;
    if (!chk_int('closereg_limit')) return false;
    if (!chk_int('opennew_itv')) return false;

    return true;
}

// ============================================================================

$(function ()
{
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
    }
    
    // init load
    {
        load();
    }

    // area change
    $('select[name=area]').on('change', function () {
		load();
	});

	// button: common
	$('.head button').on('click', function () {
        let $this = $(this);

        if (!save_check()) {
            err_show({err: 'error params'});
            return;
        }

        $.ajax({
			method: "post",
			url: '/server/settings/save',
            data: $this.closest('.panel').find('select, input, textarea').serialize(),
            dataType: "json",
		}).done(function (r) {
			err_show(r);
		}).fail(function () {
			err_show(ErrNet);
		});
    });

});
