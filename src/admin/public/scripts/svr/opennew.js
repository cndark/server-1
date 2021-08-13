
let refresh_ts = 0;

function refresh() {
	let ctx = $('.panel-body');

	$.ajax({
		method: "post",
		url: '/server/opennew/refresh',
		data: ctx.find('select, input, textarea').serialize(),
		dataType: "json",
	}).done(function (r) {
		if (!err_ok(r)) return;

		$('#last_id').text(r.last.id);
		$('#last_ts').text(r.last.ts);
		$('#opening').text(r.opening || '');

		$('#wait_ids').prev().text(r.wait_ids.length);
		$('#wait_ids').text(r.wait_ids.join(', '));
	});

	refresh_ts = Date.now();
}

// ============================================================================

$(function ()
{
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
	}

	// auto refresh
	{
		refresh();
		setInterval(() => {
			if ($('#opening').text().length > 0) {
				refresh();
			} else if (Date.now() - refresh_ts > 60 * 1000) {
				refresh();
			}
		}, 5000);
	}

	// button: ok
	$('.head button').on('click', function () {
		let $this = $(this);

		let ids = $('#wait_ids').text();
		let m = ids.match(/^\s*(\d+)/);
		if (!m) {
			err_show({err: 'no wait server'});
			return;
		}
		let id = Number(m[1]);

		bs_confirm(
			`Are you sure to manually open new server: game${id} ?`,
		).ok(() => {
			$.ajax({
				method: "post",
				url: '/server/opennew/cmd',
				data: $this.closest('.panel').find('select, input, textarea').serialize() + `&id=${id}`,
				dataType: "json",
			}).done(function (r) {
				err_show(r);
				if (err_ok(r)) refresh();
			}).fail(function () {
				err_show(ErrNet);
			});
		});
	});

});
