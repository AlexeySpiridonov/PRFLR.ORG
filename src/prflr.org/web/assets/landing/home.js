$(function() {

  // Screenshot browser

  $('.screenshots').click(function(){
    gallery = $(this);
    lastChild = gallery.children("li:last-child");

    lastChild.remove();

    gallery.prepend( lastChild );

  });


  $(".pricing-link").click(function(e) {
    $('html, body').animate({
        scrollTop: $("#pricing").offset().top
    }, 700);

    e.preventDefault();
  });
  

});
