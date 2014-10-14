{{define "js"}}
$(function(){
  $('.page').on('click', '.answer', function(){
    var index = $(this).index('.answer');
    $.post('/', {answer: index}, function(data) {
      $('.question-body').replaceWith(data.html);
      if (data.correct) {
        $('body').animate({'background-color': '#00AA00'}, 300, 'swing', function() {
          $('body').animate({'background-color': '#060606'}, 300, 'swing');
        });
      } else {
        $('body').animate({'background-color': '#AA0000'}, 300, 'swing', function() {
          $('body').animate({'background-color': '#060606'}, 300, 'swing');
        });
      }
    },'json');
  });
});
{{end}}
