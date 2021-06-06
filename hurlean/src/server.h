#pragma once

#include <string>
#include <vector>
#include <memory>

#include <asio.hpp>

#include "client_connection.h"


namespace hl
{
	class Server {

	private:
		asio::io_context io_context;
		//asio::ip::tcp::endpoint endpoint(asio::ip::tcp::v4(), 4545);
		//asio::ip::tcp::acceptor acceptor(io_context, endpoint);
		
		std::vector<std::shared_ptr<ClientConnection>> connections;

	public:
		Server();	// require endpoint info here
		~Server() = default;

		void run();
		void update();
	
	public:
		virtual void on_client_connect() = 0;
		virtual void on_client_disconnect() = 0;
		virtual void on_client_message() = 0;

	protected:
		void write();

	private:
		void listen();

	};
}