$(document).ready( function() {

	var user_name;
	var conexion_final;

	$("#form_registro").on("submit", function(e){
		user_name = $("#user_name").val();
		e.preventDefault();

		$.ajax({
			type: "POST",
			url: "http://localhost:8000/validate", //En este caso como estamos sobre el mismo servidor
			data: {
				"user_name" : user_name
			},
			success: function(data){
				result_validate(data)
			}
		})
	})

	function result_validate(data){
		obj = JSON.parse(data);
		console.log(obj)
		if (obj.valid === true){
			create_conection()
		}else{
			location.reload()
		}
	}

	function create_conection(){
		var conexion = new WebSocket("ws://localhost:8000/chat/" + user_name)
		conexion_final = conexion

		conexion.onopen = function(response){
			conexion.onmessage = function(response){
				console.log(response)
				//Lo ultimo
				val = $("#chat_area").val();
		   	$("#chat_area").val(val + "\n" + response.data); 
			}
		}
		$("#registro").hide();
   	$("#container_chat").show();
	}

	$("#form_message").on("submit", function(e) {
		e.preventDefault()
		message = $("#msg").val();
		conexion_final.send(message)
		$("#msg").val("")
	});
})