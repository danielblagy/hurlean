#pragma once

#include <thread>
#include <memory>

#include "connection.h"


namespace hl
{
	template <class T>
	class ClientSession
	{
	private:
		std::shared_ptr<Connection<T>> connection;
		std::thread connection_thread;

	public:
		ClientSession(std::shared_ptr<Connection<T>> _connection)
			: connection(std::move(_connection))
		{}

		~ClientSession() = default;

	private:
		void listen();	// used as an update function in connection thread for listening for messages from the client
	};
}