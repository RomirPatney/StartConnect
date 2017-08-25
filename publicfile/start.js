$(window).scroll(function() {
  	if($(document).scrollTop() > 200) {
    	$('.navbar').addClass('shrink');
    }
    else {
    $('.navbar').removeClass('shrink');
    }
  });


    $(window).scroll(function() {
      if($(document).scrollTop() > 200) {
        $('nav img').css('height', '50px', 'padding-bottom', '10px');
      }
      else {
        $('nav img')
      }
    });
