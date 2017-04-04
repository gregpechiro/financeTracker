var today = new Date();
var range = 1;
var beg = new Date(today.getFullYear(), today.getMonth(), 1);
var end = new Date(today.getFullYear(), today.getMonth() + range, 0);
displayRange();

function logRange() {
	console.log(beg);
	console.log(end);
	console.log(' ');
}

$('a#reset').click(function() {
	beg = new Date(today.getFullYear(), today.getMonth(), 1);
	end = new Date(today.getFullYear(), today.getMonth() + 1, 0);
	$('select#rangeType').val('1');
	displayRange();
	table.draw();
});

$('a#rangePrev').click(function() {
	if ($('select#rangeType').val() == 'year') {
		beg.setFullYear(end.getFullYear() - 1, 0, 1);
		end.setFullYear(end.getFullYear() - 1, 12, 0);
	} else {
		beg.setMonth(beg.getMonth() - 1);
		end = new Date(beg.getFullYear(), beg.getMonth() + range, 0);
	}
	displayRange();
	table.draw();

});

$('a#rangeNext').click(function() {
	if ($('select#rangeType').val() == 'year') {
		beg.setFullYear(end.getFullYear() + 1, 0, 1);
		end.setFullYear(end.getFullYear() + 1, 12, 0);
	} else {
		beg.setMonth(beg.getMonth() + 1);
		end = new Date(beg.getFullYear(), beg.getMonth() + range, 0);
	}

	displayRange();
	table.draw();
});

function displayRange() {
	logRange();
	var rangeType = $('select#rangeType').val();
	if (rangeType === '1') {
		$('a#range').html(beg.toLocaleString({}, {
			month: "long"
		}) + ' ' + beg.getFullYear());
		return;
	}
	if (rangeType === '3' || rangeType === '6') {
		$('a#range').html(
			beg.toLocaleString({}, {
				month: "long"
			}) + ' ' + beg.getFullYear() +
			' - ' +
			end.toLocaleString({}, {
				month: "long"
			}) + ' ' + end.getFullYear()
		);
		return
	}
	if (rangeType === 'year') {
		$('a#range').html(beg.getFullYear());
	}
}

$('#begDate').change(function() {
	beg = $(this).datepicker('getDate');
	table.draw();
	logRange();
});

$('#endDate').change(function() {
	end = $(this).datepicker('getDate');
	table.draw();
	logRange();
});

$('select#rangeType').change(function() {
	var rangeType = $(this).val();
	switch (rangeType) {
		case 'pick':
			$('#dateRange').addClass('hide');
			$('#datePicker').removeClass('hide');
			$('#begDate').datepicker('update', beg);
			$('#endDate').datepicker('update', end);
			displayRange();
			table.draw();
			return;
		case 'year':
			beg.setFullYear(end.getFullYear(), 0, 1);
			end.setFullYear(end.getFullYear(), 12, 0);
			range = 12;
			break;
		default:
			range = +rangeType
			beg.setDate(1);
			end = new Date(beg.getFullYear(), beg.getMonth() + range, 0);
	}
	$('#datePicker').addClass('hide');
	$('#dateRange').removeClass('hide');
	displayRange();
	table.draw();
});

var category2 = [];

$('select#category1').change(function() {
	category2 = [];
});

$.fn.dataTable.ext.search.push(
	function(settings, data, dataIndex) {

		var date = new Date(data[0]);
		if (date == 'Invalid Date') {
			return false;
		}
		return (beg <= date && date <= end);
	}
);

$.fn.dataTable.ext.search.push(
	function(settings, data, dataIndex) {

		var aType = $('select#aType-filter').val();

		if (aType == 'expense') {
			var v = intVal(data[2])
			if (v < 0) {
				return true;
			}
			return false;
		}

		if (aType == 'income') {
			var v = intVal(data[3]);
			if (v > 0) {
				return true;
			}
			return false;
		}

		return true;
	}
);

$.fn.dataTable.ext.search.push(
	function(settings, data, dataIndex) {

		var category = $('select#category1').val();
		if (category == '') {
			return true;
		}

		if (data[5] == category) {
			if (data[6] != '') {
				category2.push(data[6]);
			}
			return true;
		}
		if (data[6] == category) {
			if (data[5] != '') {
				category2.push(data[5]);
			}
			return true
		}
		return false;
	}
);

$.fn.dataTable.ext.search.push(
	function(settings, data, dataIndex) {

		var category = $('select#category2').val();
		if (category == '') {
			return true;
		}

		if (data[5] == category) {
			return true;
		}
		if (data[6] == category) {
			return true
		}
		return false;
	}
);

$.fn.dataTable.ext.search.push(
	function(settings, data, dataIndex) {

		var who = $('select#who').val();
		if (who == '') {
			return true;
		}

		if (data[4] == who) {
			return true;
		}
		return false;
	}
);

$('button#resetAll').click(function() {
	$('select#aType-filter').val('');
	$('select#who').val('');
	$('select#category1').val('');
	$('select#category2').val('');
	$('a#reset').click();
});

$('select.filter').change(function() {
	table.draw();
});

var table = $('#transactions').DataTable({
	"lengthMenu": [5, 10, 15, 20],
	"pageLength": 5,
	"order": [
		[0, "desc"]
	],
	"columnDefs": [{
			"name": "date",
			"type": "date",
			"targets": 0
		},
		{
			"name": "title",
			"targets": 1
		},
		{
			"name": "expense",
			"targets": 2
		},
		{
			"name": "income",
			"targets": 3
		},
		{
			"name": "who",
			"targets": 4
		},
		{
			"name": "category",
			"targets": 5
		},
		{
			"name": "secondaryCategory",
			"targets": 6
		}
	],
	"footerCallback": function(row, data, start, end, display) {
		api = this.api(), data;

		expenseTotal = api
			.column(2, {
				search: 'applied'
			})
			.data()
			.reduce(function(a, b) {
				return intVal(a) + intVal(b);
			}, 0);

		// Total over this page
		expenseSubTotal = api
			.column(2, {
				page: 'current'
			})
			.data()
			.reduce(function(a, b) {
				return intVal(a) + intVal(b);
			}, 0);

		// Update footer
		$(api.column(2).footer()).html(
			((api.page.info().pages > 1) ?
				((expenseSubTotal < 0) ? '<span class="pull-left">Sub Total: </span>' + currencyFormatter.format(expenseSubTotal) : '') +
				'<br>' :
				'') +

			((expenseTotal < 0) ? '<span class="pull-left">Total: </span>' + currencyFormatter.format(expenseTotal) : '')
		);

		incomeTotal = api
			.column(3, {
				search: 'applied'
			})
			.data()
			.reduce(function(a, b) {
				return intVal(a) + intVal(b);
			}, 0);

		// Total over this page
		incomeSubTotal = api
			.column(3, {
				page: 'current'
			})
			.data()
			.reduce(function(a, b) {
				return intVal(a) + intVal(b);
			}, 0);

		// Update footer
		$(api.column(3).footer()).html(
			((api.page.info().pages > 1) ?
				((incomeSubTotal > 0) ? '<span class="pull-left">Sub Total: </span>' + currencyFormatter.format(incomeSubTotal) : '') +
				'<br>' :
				'') +

			((incomeTotal > 0) ? '<span class="pull-left">Total: </span>' + currencyFormatter.format(incomeTotal) : '')
		);
	}
});

table.on('search.dt', function() {
	if ($('select#category2').val() !== '') {
		return;
	}

	if ($('select#category1').val() !== '') {
		$.uniqueSort(category2).sort();
	} else {
		category2 = [];
	}
	renderCategory2();
});

function intVal(i) {
	return typeof i === 'string' ?
		i.replace(/[\$,]/g, '') * 1 :
		typeof i === 'number' ?
		i : 0;
};

function renderCategory2() {
	if ($('select#category1').val() == '') {
		$('select#category2').html('<option></option>');
		return;
	}
	var options = '<option value="">All</option>';
	for (var i = 0; i < category2.length; i++) {
		options += '<option value="' + category2[i] + '">' + category2[i] + '</option>';
	}
	$('select#category2').html(options);
}
