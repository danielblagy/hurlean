#pragma once

#include <thread>

#include "connection.h"


namespace hl
{
	template <class T>
	class ClientConnection
	{
	private:
		Connection connection;
		std::thread connection_thread;

	public:
		ClientConnection(const asio::io_context& io_context);
		~ClientConnection() = default;

	private:
		void listen();	// used as an update function in connection thread for listening for messages from the client
	};
}