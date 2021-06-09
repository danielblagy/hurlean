#pragma once

#include <asio.hpp>

#include "message.h"
#include "queue.h"


namespace hl
{
	template <class T>
	struct Connection
	{
		asio::ip::tcp::socket socket;
		Queue<Message<T>> in;
		Queue<Message<T>> out;
		Message<T> temp_message_in;
		bool opened;

		Connection(asio::io_context& io_context)
			: socket(io_context)
		{}

		~Connection()
		{
			if (opened)
				close();
		}

		// used on server side
		void wait_for_client(asio::ip::tcp::acceptor& acceptor)
		{
			acceptor.accept(socket);
			opened = true;
		}

		// used on client side
		bool connect_to_server(const asio::ip::tcp::resolver::results_type& endpoints)
		{
			asio::error_code error;
			asio::connect(socket, endpoints, error);

			if (error)
			{
				// TODO : log 'error on connect to server' & error
			}

			return !error;
		}

		bool write(const Message<T>& message)
		{
			asio::error_code error;
			asio::write(socket, asio::buffer(&message.header, sizeof(MessageHeader<T>)), error);

			if (error)
			{
				// TODO : log 'error on writing a message header' & error
			}

			asio::write(socket, asio::buffer(message.body.data(), message.body.size()), error);

			return !error;
		}

		// true if read was successful, false if not
		bool read()
		{
			asio::error_code error;
			asio::read(socket, asio::buffer(&temp_message_in.header, sizeof(MessageHeader<T>)), error);

			// TODO : log 'error on reading message header' & error
			if (error)
			{
				return false;
			}
			
			temp_message_in.body.resize(temp_message_in.header.size);
			asio::read(socket, asio::buffer(temp_message_in.body.data(), temp_message_in.body.size()), error);

			if (error)
			{
				// TODO : log 'error on reading message body' & error
			}

			return !error;
		}

		void close()
		{
			opened = false;
			socket.close();
		}
	};
}