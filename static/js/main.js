
(function($) {
  if ($('.x_alert_msg').length){
    setTimeout(function(){
        $('.x_alert_msg').hide();
        $('.x_alert_msg').remove();
    },3000)
  }

  $('#create_account_btn').click(function(){
    window.location = "/user/create";
    // $.ajax({
    //   url: "/createUser",
    //   success: function(r){
    //     console.log('success');
    //   }
    // });
  });

  $('.x_day').on('click',function(){
    var day = $(this).data('day')
    $(this).toggleClass('selected');
    var val = '';
    $('.x_day.selected').each(function(){
        var curr_day = $(this).data('day');
        val += curr_day+',';
    });
    if (val){
      $('#days').val(val);
    }
  });
})(jQuery);
