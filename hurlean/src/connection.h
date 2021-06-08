#pragma once

#include <asio.hpp>

#include "message.h"
#include "queue.h"


namespace hl
{
	template <class T>
	class Connection
	{
	private:
		asio::ip::tcp::socket socket;
		Queue<Message<T>> in;
		Queue<Message<T>> out;

	public:
		Connection(asio::io_context& io_context)
			: socket(io_context)
		{}

		~Connection() = default;

	public:
		// wait for connection
		void wait(asio::ip::tcp::acceptor& acceptor)
		{
			acceptor.accept(socket);
		}

		void send(std::string message)
		{
			asio::error_code ignored_error;
			asio::write(socket, asio::buffer(message), ignored_error);
		}
	};
}