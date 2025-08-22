package main

func userPostsHandler(r *Request) *Response {
	s := r.PathParams["id"]
	response := NewResponse(StatusOK, PlainTextHeaders(), "Nazdar vitaj user id - "+s)
	return response
}

func userHandler(*Request) *Response {
	response := NewResponse(StatusOK, PlainTextHeaders(), "Nazdar vitaj")
	return response
}
