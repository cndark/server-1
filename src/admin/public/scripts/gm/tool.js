
function fill_select(name, r) {
	let opt = [];
	for (let row of r) {
		opt.push('<option value="' + row[0] + '">' + row[1] + '</option>');
	}
	$(`.panel-body select[name=${name}]`).html(opt.join());
}

// ============================================================================

$(function () {
	// add .form-control
	{
		$('select, input, textarea').addClass("form-control");
	}

	// fetch init data: res
	{
		$.ajax({
			method: "post",
			url: "/gm/gtab",
			data: "key=res",
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

	// fetch init data: hero
	{
		$.ajax({
			method: "post",
			url: "/gm/gtab",
			data: "key=hero",
			dataType: "json",
		}).done(function (r) {
			let combo = $('.panel-body input[name=hero_k]');

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

	// button: others
	$('button:not([id])').on('click', function () {
		let $this = $(this);

		$.ajax({
			method: "post",
			url: "/gm/tool",
			data: $('#plrinfo').parent().next().find('input[name="plrid"], input[name="plrname"]')
				.add($this.parent().next().find('select, input, textarea'))
				.serialize(),
			dataType: "json",
		}).done(function (r) {
			err_show(r, $this);
		}).fail(function () {
			err_show(ErrNet, $this);
		})
	});

	// button: plrinfo
	$('#plrinfo').on('click', function () {
		let $this = $(this);

		$.ajax({
			method: "post",
			url: "/gm/tool",
			data: $this.parent().next().find('select, input, textarea').serialize(),
			dataType: "json",
		}).done(function (r) {
			$tab = $this.parent().next().find('table');
			$tab.empty();

			if (!err_show(r, $this)) return;
			if (r.length < 1) return;

			let arr = [];

			// header
			arr.push('<thead>');
			arr.push('<tr>');
			for (let col of r[0]) {
				arr.push('<th>' + col + '</th>');
			}
			arr.push('</tr>');
			arr.push('</thead>');

			// remove header and sort
			r = r.slice(1)
			if (r.length > 0 && r[0][0] != '*') {
				r.sort();
			}

			// body
			arr.push('<tbody>');
			for (let row of r) {
				arr.push('<tr>');
				for (let col of row) {
					arr.push('<td>' + col + '</td>');
				}
				arr.push('</tr>');
			}
			arr.push('</tbody>');

			// set table html
			$tab.html(arr.join());
		}).fail(function () {
			err_show(ErrNet, $this);
		})
	});

});
