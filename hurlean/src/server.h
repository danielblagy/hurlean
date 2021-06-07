#pragma once

#include <string>
#include <vector>
#include <memory>

#include <asio.hpp>

#include "client_connection.h"


namespace hl
{
	template <class T>
	class Server
	{
	private:
		asio::io_context io_context;
		asio::ip::tcp::acceptor acceptor;
		
		std::vector<std::shared_ptr<ClientConnection<T>>> connections;

	public:
		// port on which server is going to be running
		Server(unsigned short port)
			: acceptor(io_context, asio::ip::tcp::v4(), port)
		{}
		
		~Server() = default;

		void start();
		void update(size_t max_messages);
	
	public:
		virtual void on_client_connect() = 0;
		virtual void on_client_disconnect() = 0;
		virtual void on_client_message() = 0;

	private:
		void write();
	};
}