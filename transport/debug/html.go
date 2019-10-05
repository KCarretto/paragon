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
    		$('<li>').append(
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
                $("#scriptTextArea").val("");
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

<h1>Paragon Script</h1>

<div class="form-group">
	<textarea class="form-control" id="scriptTextArea" rows="3">def main():
	    print("Hello World")
	</textarea>
</div>
<button type="button" id="submit"class="btn btn-primary">Submit</button>
<br/>
<br/>
<br/>
<h1>Results</h1>
<ul id="resultsList">

</ul>
</body>
</html>
`
