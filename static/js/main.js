
(function($) {
  console.log('GOT HERE');
  if ($('.x_alert_msg').length){
    setTimeout(function(){
        $('.x_alert_msg').hide();
        $('.x_alert_msg').remove();
    },3000)
  }
  console.log($('#create_account_btn'));
  $('#create_account_btn').click(function(){
    window.location = "/user/create";
    // $.ajax({
    //   url: "/createUser",
    //   success: function(r){
    //     console.log('success');
    //   }
    // });
  });
})(jQuery);
