
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

$(function ()
{
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
	}

	// button: update
	$('#btn_update').on('click', function () {
		let $this = $(this);

        $.ajax({
			method: "post",
			url: '/server/status',
            data: $('form select[name=area]')
                .add($this.closest('.panel').find('select, input, textarea'))
                .serialize(),
            dataType: "json",
		}).done(function (r) {
			err_show(r, $this);
		}).fail(function () {
			err_show(ErrNet, $this);
		});
    });

});
