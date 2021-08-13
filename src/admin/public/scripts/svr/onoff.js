
$(function ()
{
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
	}

	// init auto refresh
	{
		setInterval(() => {
			let ctx = $('.panel-body');
			let target = ctx.find('input[name=target]').val();

			if (!target) return;

			$.ajax({
				method: "post",
				url: '/server/onoff/refresh',
				data: ctx.find('select, input, textarea').serialize(),
				dataType: "json",
			}).done(function (r) {
				if (!err_ok(r)) return;

				$('#started').prev().text(r.started.length);
				$('#started').text(r.started.join(', '));

				$('#stopped').prev().text(r.stopped.length);
				$('#stopped').text(r.stopped.join(', '));

				$('#starting').prev().text(r.starting.length);
				$('#starting').text(r.starting.join(', '));

				$('#stopping').prev().text(r.stopping.length);
				$('#stopping').text(r.stopping.join(', '));
			});
		}, 3500);
	}

	// start/stop button
	$('.head button').on('click', function () {
		let $this = $(this);
		let name = $this.attr('name');

		let ctx = $('.panel-body');
		let target = ctx.find('input[name=target]').val();

		if (!target) {
			err_show({err: 'no target specified'});
			return;
		}

		bs_confirm(
			`Are you sure to ${name} the servers ?`,
		).ok(() => {
			$.ajax({
				method: "post",
				url: `/server/onoff/cmd?op=${name}`,
				data: ctx.find('select, input, textarea').serialize(),
				dataType: "json",
			}).done(function (r) {
				err_show(r);
			}).fail(function () {
				err_show(ErrNet);
			});
		});
	});

});
