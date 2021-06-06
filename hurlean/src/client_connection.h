#pragma once

#include <thread>
#include <deque>

#include <asio.hpp>

#include "message.h"


namespace hl
{
	class ClientConnection {
	
	private:
		asio::ip::tcp::socket socket;
		std::thread connection_thread;
		std::deque<Message> in;
		std::deque<Message> out;

	public:
		ClientConnection(const asio::io_context& io_context);
		~ClientConnection() = default;

	private:
		void listen();	// used as an update function in connection thread for listening for messages from the server
	};
}