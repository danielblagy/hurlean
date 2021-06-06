#include <iostream>
#include <asio.hpp>


int main() {
	try {
		asio::io_context io_context;

		// ip version 4, port 13
		asio::ip::tcp::endpoint endpoint(asio::ip::tcp::v4(), 4545);

		asio::ip::tcp::acceptor acceptor(io_context, endpoint);

		while (true) {
			asio::ip::tcp::socket socket(io_context);
			acceptor.accept(socket);	// wait for connection

			// we got connection
			std::string message = "hello";

			asio::error_code ignored_error;
			asio::write(socket, asio::buffer(message), ignored_error);
		}
	}
	catch (std::exception& e) {
		std::cerr << e.what() << std::endl;
	}
	
	return 0;
}