$(document).ready(function(){
    $.get('http://localhost/', {}, function (data) {
        if (!Array.isArray(data.data)) {
            return true;
        }

        $.each(data.data, function(index, item) {
            var row = $('<tr>');
            row.append($('<td>').text(item.name));
            row.append($('<td>').text(item.pay_date));
            row.append($('<td>').text(item.pay_price + ' ₽'));
            row.append($('<td>').text(item.sale_date));
            row.append($('<td>').text(item.sale_price + ' ₽'));
            row.append($('<td>').text(item.days));
            row.append($('<td>').text(item.pay_day + ' ₽'));
            $('#home_table_id tbody').append(row);
        });
    });
});