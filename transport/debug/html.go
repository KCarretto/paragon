package debug

const debugHTML = `<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/toastr.js/latest/css/toastr.min.css"/>
<script
  src="https://code.jquery.com/jquery-3.4.1.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js"></script>

<script src="https://cdnjs.cloudflare.com/ajax/libs/toastr.js/latest/js/toastr.min.js"></script>

<script>
	var OFFSET = 0;
    var MAX_RESULTS_SHOWN = 15;
    toastr.options = {
      "closeButton": true,
      "debug": false,
      "newestOnTop": false,
      "progressBar": true,
      "positionClass": "toast-bottom-full-width",
      "preventDuplicates": false,
      "onclick": null,
      "showDuration": "300",
      "hideDuration": "1000",
      "timeOut": "5000",
      "extendedTimeOut": "1000",
      "showEasing": "swing",
      "hideEasing": "linear",
      "showMethod": "fadeIn",
      "hideMethod": "fadeOut"
	}
    
	function addToResults(result) {
    	$('#resultsList').prepend(
    		$('<li class="list-group-item">').append(
            	result
            )
         );
         if ($('#resultsList').children().length > MAX_RESULTS_SHOWN) {
         	$('#resultsList').children().last().remove();
         }
    }
    function queueScript(){
    	let script = {
            "content": btoa($("#scriptTextArea").val())
        }
        $.ajax({
            url: '/queue',
          	type: 'POST',
    		contentType: 'application/json',
    		data: JSON.stringify(script),
            success: function(data){
                toastr.success('Your Task was successfully Queued!', 'Task Queued')
            },
           	error: function(){
                toastr.error('We failed send the script to the server.', 'Task Queue Error')
          	},
		});
    }
    function pollResults(){
    	let queryParams = {
        	"offset": OFFSET
        }
    	$.get("/messages", queryParams, function(responses) {
            responses.map(function(response){
              	OFFSET += 1
            	if (response.results != null){
                    response.results.map(function(result){
                        if (result.output != null){
                            addToResults(atob(result.output));
                        }
                    });
                }
            });
        });
    }
	$(document).ready(function() {
		$('#scriptTextArea').keydown(function(e) {
			var keyCode = e.keyCode || e.which;

			if (keyCode == 9) {
			  e.preventDefault();
			  var start = $(this).get(0).selectionStart;
			  var end = $(this).get(0).selectionEnd;

			  // set textarea value to: text before caret + tab + text after caret
			  spaces = "    "
			  $(this).val($(this).val().substring(0, start)
						  + spaces
						  + $(this).val().substring(end));

			  // put caret at right position again
			  $(this).get(0).selectionStart =
			  $(this).get(0).selectionEnd = start + 4;
			}
		});
        $("#submit").click(function() {
        	queueScript()
		});
        setInterval(
        	pollResults, 
           	3000
        );
	});
</script>
</head>
<body>

<div class="jumbotron" style="padding-top:10px;padding-bottom:10px;">
  <h3 class="display-6" style="padding-top:5px;">Welcome to Paragon</h3>
  <p class="lead"></p>
  <hr class="my-6">
  <p>The agent is currently running in debug mode. Try running paragon scripts below.</p>
  <a class="btn btn-primary btn-md" target="_blank" href="https://github.com/KCarretto/paragon" role="button">Learn more</a>
</div>

<div class="row">

	<div class="col-sm-6">
		<div class="card">
			<h4 class="card-header">Code</h4>
			<div class="card-body">
				<div class="form-group">
<textarea class="form-control" id="scriptTextArea" rows="8">
def main():
    print("Hello World")
</textarea>
				</div>
			<button type="button" id="submit"class="btn btn-primary float-right">Submit</button>
			</div>
		</div>
	</div>


	<div class="col-sm-6">
		<div class="card">
			<h4 class="card-header">Results</h4>
			<div class="card-body">
				<ul id="resultsList" class="list-group list-group-flush"></ul>
			</div>
		</div>
	</div>
</div>
</body>
</html>
`
